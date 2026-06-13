/**
 * moonBASIC IDE — full-featured editor for .mb files
 */
import { SBIcons as I } from './icons.js';
import { PanelResize } from './panel-resize.js';
import { initPanels, panelMenuItems, setLeftPanelVisible, setRightPanelVisible, toggleLeftPanel, toggleRightPanel } from './panels.js';
import { EXAMPLES, EXAMPLE_CATEGORIES, SNIPPETS } from './examples.js';
import {
  loadLangData, buildMoonBasicMode, getLangData,
  formatCommandHelp, lookupAtCursor, completionItems
} from './lang.js';
import {
  isDesktopApp, getToolchain, checkFile, compileFile, runFile,
  openNativeFile, saveNativeFile, saveNativeFileAs,
  getLSPHover, getLSPCompletion, openProjectFolder
} from './toolchain.js';
import { initModals, showNewProjectModal, showGoToModal, showShortcutsModal } from './modals.js';
import { showSettingsModal } from './settings.js';
import { initAppearance, setAppearanceCallback } from './themes.js';
import { buildProject, PROJECT_TYPES } from './projects.js';
import { pickProjectFolder, writeProjectToDisk, readProjectFromFolder, folderPickerAvailable } from './project-folder.js';
import { buildMenubar } from './menus.js';
import { createExplorer } from './explorer.js';
import { createFindBar } from './find.js';
import { initDocs, loadDoc, markdownToHtml, getDocIndex, showDocPanel, openDocForCommand } from './docs.js';

const RECENT_KEY = 'moonbasic-ide-recent';
const DRAFT_KEY = 'moonbasic-ide-draft';

export async function startStudio() {
  initAppearance();
  const langData = await loadLangData();
  const mbLang = buildMoonBasicMode(langData);

  let editor;
  let activeFile = 0;
  let openFiles = [{ name: 'untitled.mb', content: EXAMPLES.hello.code, dirty: false, diskPath: '' }];
  let projectPath = '';
  let findBar;

  const $ = id => document.getElementById(id);
  const consoleEl = $('console-output');
  const helpEl = $('help-cursor-panel');
  const refList = $('ref-list');
  const refSearch = $('ref-search');

  const explorer = createExplorer({
    treeEl: $('file-tree'),
    onOpenFile: f => openBuffer(f.name, f.content, f.diskPath)
  });

  function log(kind, msg) {
    const line = document.createElement('div');
    line.className = 'log-' + (kind || 'info');
    line.textContent = msg;
    consoleEl.appendChild(line);
    consoleEl.scrollTop = consoleEl.scrollHeight;
  }

  function setStatus(state, text) {
    $('status-dot').className = 'statusbar-dot ' + (state || '');
    $('status-text').textContent = text || 'Ready';
  }

  async function refreshToolchainLog() {
    const tc = await getToolchain();
    if (tc.found) {
      log('success', `Toolchain ready: moonbasic${tc.moonrun ? ' + moonrun' : ''}`);
      setStatus('ok', 'Toolchain ready');
    } else if (isDesktopApp()) {
      log('warn', 'Toolchain not found — run npm run toolchain:build or File → Settings');
      setStatus('warn', 'No toolchain');
    }
  }

  function syncEditor() {
    if (editor) openFiles[activeFile].content = editor.getValue();
  }

  function updateTabs() {
    const bar = $('tab-bar');
    bar.innerHTML = '';
    openFiles.forEach((f, i) => {
      const tab = document.createElement('button');
      tab.type = 'button';
      tab.className = 'editor-tab' + (i === activeFile ? ' active' : '') + (f.dirty ? ' dirty' : '');
      const label = document.createElement('span');
      label.textContent = f.name + (f.dirty ? ' •' : '');
      tab.appendChild(label);
      if (openFiles.length > 1) {
        const close = document.createElement('span');
        close.className = 'tab-close';
        close.textContent = '×';
        close.addEventListener('click', e => { e.stopPropagation(); closeTab(i); });
        tab.appendChild(close);
      }
      tab.addEventListener('click', () => switchTab(i));
      bar.appendChild(tab);
    });
    $('bc-file').textContent = openFiles[activeFile]?.name || 'untitled.mb';
  }

  function closeTab(i) {
    if (openFiles.length <= 1) return;
    syncEditor();
    openFiles.splice(i, 1);
    activeFile = Math.min(activeFile, openFiles.length - 1);
    editor.setValue(openFiles[activeFile].content);
    updateTabs();
    scheduleCheck();
  }

  function switchTab(i) {
    syncEditor();
    activeFile = i;
    editor.setValue(openFiles[i].content);
    updateTabs();
    clearErrors();
    scheduleCheck();
    focusEditor();
    if (activeDocTab === 'cursor') scheduleHelpAtCursor();
  }

  function openBuffer(name, content, diskPath = '') {
    syncEditor();
    const existing = openFiles.findIndex(f => f.diskPath && diskPath && f.diskPath === diskPath);
    const byName = openFiles.findIndex(f => f.name === name && !diskPath);
    const idx = diskPath ? (existing >= 0 ? existing : -1) : byName;
    if (idx >= 0) {
      activeFile = idx;
      openFiles[idx].content = content;
      if (diskPath) openFiles[idx].diskPath = diskPath;
    } else {
      openFiles.push({ name, content, dirty: false, diskPath });
      activeFile = openFiles.length - 1;
    }
    editor.setValue(content);
    updateTabs();
    addRecent(diskPath || name);
    $('preview-overlay')?.classList.add('hidden');
    scheduleCheck();
    focusEditor();
    if (activeDocTab === 'cursor') scheduleHelpAtCursor();
  }

  function addRecent(path) {
    if (!path) return;
    try {
      let r = JSON.parse(localStorage.getItem(RECENT_KEY) || '[]');
      r = [path, ...r.filter(x => x !== path)].slice(0, 8);
      localStorage.setItem(RECENT_KEY, JSON.stringify(r));
    } catch (_) {}
  }

  function getRecent() {
    try { return JSON.parse(localStorage.getItem(RECENT_KEY) || '[]'); } catch { return []; }
  }

  function parseProblems(stderr) {
    const lines = (stderr || '').split(/\r?\n/);
    const out = [];
    const re = /^\[moonBASIC\].+ in (.+) line (\d+) col (\d+):$/;
    for (let i = 0; i < lines.length; i++) {
      const m = re.exec(lines[i]);
      if (!m) continue;
      let msg = 'Error';
      if (i + 1 < lines.length && lines[i + 1].trim()) msg = lines[i + 1].trim();
      out.push({ line: parseInt(m[2], 10), message: msg });
    }
    return out;
  }

  function clearErrors() {
    errorMarks.forEach(m => m.clear());
    errorMarks = [];
    $('pane-problems').innerHTML = '';
    updateProblemsBadge(0);
  }

  let errorMarks = [];
  function addProblem(line, message) {
    if (editor) errorMarks.push(editor.addLineClass(line - 1, 'background', 'cm-error-line'));
    const pane = $('pane-problems');
    const row = document.createElement('div');
    row.className = 'problem-row';
    row.innerHTML = `<span class="problem-loc">Ln ${line}</span><span class="problem-msg">${message}</span>`;
    row.addEventListener('click', () => { editor.setCursor(line - 1, 0); showPanel('problems'); });
    pane.appendChild(row);
  }

  function updateProblemsBadge(n) {
    const tab = document.querySelector('.panel-tab[data-panel="problems"]');
    if (tab) tab.textContent = n ? `Problems (${n})` : 'Problems';
  }

  function showPanel(id) {
    document.querySelectorAll('.panel-tab').forEach(t => t.classList.toggle('active', t.dataset.panel === id));
    document.querySelectorAll('.panel-pane').forEach(p => p.classList.toggle('active', p.id === 'pane-' + id));
  }

  let checkTimer = null;
  function scheduleCheck() {
    clearTimeout(checkTimer);
    checkTimer = setTimeout(() => void doCheck(false), 1200);
  }

  async function ensureSaved() {
    syncEditor();
    const f = openFiles[activeFile];
    if (f.diskPath) {
      if (isDesktopApp()) {
        const res = await saveNativeFile(f.content, f.name);
        if (!res.success) throw new Error(res.error || 'Save failed');
        f.dirty = false;
        updateTabs();
        return f.diskPath;
      }
      return f.diskPath;
    }
    if (isDesktopApp()) {
      const res = await saveNativeFileAs(f.content, f.name);
      if (!res.success) throw new Error(res.error || 'Save cancelled');
      f.diskPath = res.path;
      f.name = res.filename || f.name;
      f.dirty = false;
      updateTabs();
      addRecent(f.diskPath);
      return f.diskPath;
    }
    throw new Error('Save the file first — use File → Save As (desktop app required for run/compile)');
  }

  async function doCheck(showOk = true) {
    if (!editor) return;
    try {
      const path = await ensureSaved();
      const res = await checkFile(path);
      clearErrors();
      if (res.success) {
        if (showOk) { log('success', 'Check OK'); setStatus('ready', 'Check OK'); }
        return true;
      }
      const errs = parseProblems(res.stderr);
      if (errs.length) {
        errs.forEach(e => addProblem(e.line, e.message));
        updateProblemsBadge(errs.length);
        showPanel('problems');
        setStatus('error', errs.length + ' problem(s)');
      } else {
        log('error', res.stderr || res.error || 'Check failed');
        setStatus('error', 'Check failed');
      }
      return false;
    } catch (e) {
      if (showOk) log('warn', e.message);
      return false;
    }
  }

  async function doCompile() {
    try {
      const path = await ensureSaved();
      log('accent', 'Compiling ' + path);
      const res = await compileFile(path);
      if (res.success) {
        log('success', res.message || 'Compiled to .mbc');
        setStatus('ready', res.message || 'Compiled');
      } else {
        parseProblems(res.stderr).forEach(e => addProblem(e.line, e.message));
        showPanel('problems');
        log('error', res.stderr || res.error || 'Compile failed');
        setStatus('error', 'Compile failed');
      }
    } catch (e) { log('error', e.message); }
  }

  async function doRun() {
    try {
      const path = await ensureSaved();
      log('accent', 'Running ' + path);
      const res = await runFile(path);
      if (res.success) {
        log('success', res.message || 'moonrun started — game window opening');
        setStatus('running', 'Running');
        $('preview-overlay')?.classList.add('hidden');
      } else {
        log('error', res.error || res.stderr || 'Run failed');
        setStatus('error', 'Run failed');
      }
    } catch (e) { log('error', e.message); }
  }

  let activeDocTab = 'guide';
  let helpDebounce = null;

  function editorDocPath(f) {
    if (f.diskPath) return f.diskPath;
    return `_ide_workspace_/${f.name.replace(/[/\\:]/g, '_')}`;
  }

  function renderCursorHelp(html, key) {
    helpEl.innerHTML = html;
    if (key) {
      const link = document.createElement('p');
      link.innerHTML = `<button type="button" class="doc-open-ref" data-cmd="${key}">Open full reference →</button>`;
      helpEl.appendChild(link);
      link.querySelector('.doc-open-ref')?.addEventListener('click', () => openDocForCommand(key));
    }
  }

  async function updateHelpAtCursor(opts = {}) {
    if (!editor || !helpEl) return;
    const { showTab = false } = opts;
    if (showTab) {
      setRightPanelVisible(true);
      showDocTab('cursor');
    }
    if (activeDocTab !== 'cursor' && !showTab) return;

    const pos = editor.getCursor();
    const line = editor.getLine(pos.line);
    const f = openFiles[activeFile];
    syncEditor();

    $('preview-overlay')?.classList.add('hidden');

    if (isDesktopApp()) {
      const md = await getLSPHover(editorDocPath(f), f.content, pos.line, pos.ch);
      if (md?.trim()) {
        renderCursorHelp(markdownToHtml(md));
        return;
      }
    }

    const hit = lookupAtCursor(line, pos.ch, getLangData());
    if (hit?.entry) {
      renderCursorHelp(markdownToHtml(formatCommandHelp(hit.key, hit.entry)), hit.key);
      return;
    }
    helpEl.innerHTML = '<p class="help-placeholder">Place the cursor on <code>APP.OPEN</code>, <code>PRINT</code>, or any command for docs.</p>';
  }

  function scheduleHelpAtCursor() {
    clearTimeout(helpDebounce);
    helpDebounce = setTimeout(() => updateHelpAtCursor(), 120);
  }

  function showDocTab(tab) {
    activeDocTab = tab;
    document.querySelectorAll('[data-doc-tab]').forEach(t => t.classList.toggle('active', t.dataset.docTab === tab));
    $('doc-viewer')?.classList.toggle('hidden', tab !== 'guide');
    helpEl?.classList.toggle('hidden', tab !== 'cursor');
    if (tab === 'cursor' || tab === 'guide') {
      $('preview-overlay')?.classList.add('hidden');
    }
  }

  function buildReferenceList(filter = '') {
    const data = getLangData();
    const q = filter.trim().toUpperCase();
    refList.innerHTML = '';
    const keys = Object.keys(data.commands || {}).sort();
    let count = 0;
    for (const key of keys) {
      if (q && !key.includes(q)) continue;
      const entry = data.commands[key];
      const btn = document.createElement('button');
      btn.type = 'button';
      btn.className = 'ref-item';
      btn.innerHTML = `<span class="ref-key">${key}</span><span class="ref-hint">${(entry.description || '').slice(0, 72)}</span>`;
      btn.addEventListener('click', async () => {
        $('preview-overlay')?.classList.add('hidden');
        const opened = await openDocForCommand(key);
        if (!opened) {
          showDocTab('cursor');
          setRightPanelVisible(true);
          renderCursorHelp(markdownToHtml(formatCommandHelp(key, entry)), key);
        }
      });
      refList.appendChild(btn);
      if (++count >= 400 && !q) {
        const more = document.createElement('p');
        more.className = 'sidebar-hint';
        more.textContent = `Type to search all ${keys.length} commands…`;
        refList.appendChild(more);
        break;
      }
    }
  }

  function focusEditor() {
    if (!editor) return;
    editor.focus();
    editor.refresh();
  }

  function initCodeMirror() {
    CodeMirror.defineMode('moonbasic', () => ({ token(stream) { return mbLang.token(stream); } }));

    CodeMirror.registerHelper('hint', 'moonbasic', (cm, options) => {
      const cur = cm.getCursor();
      const line = cm.getLine(cur.line);
      const f = openFiles[activeFile];
      syncEditor();
      const word = (line.slice(0, cur.ch).match(/([A-Za-z_.][A-Za-z0-9_.]*)$/) || ['', ''])[1];
      const from = CodeMirror.Pos(cur.line, cur.ch - word.length);
      const finish = items => {
        const list = (items || []).map(it => ({
          text: it.insertText || it.text || it.label,
          displayText: (it.displayText || it.label || it.text) + (it.hint || it.detail ? ' — ' + (it.hint || it.detail) : ''),
          from, to: cur
        }));
        if (typeof options === 'function') { options(list.length ? { list, from, to: cur } : null); return; }
        return list.length ? { list, from, to: cur } : null;
      };
      if (isDesktopApp() && f.diskPath) {
        getLSPCompletion(f.diskPath, f.content, cur.line, cur.ch).then(finish).catch(() => finish(completionItems(line, cur.ch, getLangData())));
        return;
      }
      return finish(completionItems(line, cur.ch, getLangData()));
    });

    editor = CodeMirror($('editor-mount'), {
      value: openFiles[0].content,
      mode: 'moonbasic',
      theme: 'mb-dark',
      lineNumbers: true,
      matchBrackets: true,
      autoCloseBrackets: true,
      indentUnit: 4,
      tabSize: 4,
      indentWithTabs: false,
      smartIndent: true,
      electricChars: true,
      styleActiveLine: true,
      styleActiveSelected: true,
      highlightSelectionMatches: { minChars: 2, showToken: true },
      showCursorWhenSelecting: true,
      cursorBlinkRate: 530,
      cursorHeight: 1,
      undoDepth: 300,
      scrollbarStyle: 'native',
      autoRefresh: true,
      extraKeys: {
        Tab: cm => {
          if (cm.somethingSelected()) cm.indentSelection('add');
          else cm.replaceSelection('    ', 'end');
        },
        'Shift-Tab': cm => cm.indentSelection('subtract'),
        'Ctrl-Space': cm => CodeMirror.showHint(cm, CodeMirror.hint.moonbasic, { async: true }),
        'Ctrl-S': () => saveCurrent(),
        'Ctrl-N': () => openBuffer('untitled.mb', EXAMPLES.hello.code),
        'Ctrl-O': () => openFileDialog(),
        'Ctrl-F': () => findBar?.show(),
        'Ctrl-D': cm => {
          cm.execCommand('selectNextOccurrence');
          return true;
        },
        'Ctrl-/': cm => {
          cm.operation(() => {
            for (let line = cm.getCursor('from').line; line <= cm.getCursor('to').line; line++) {
              const text = cm.getLine(line);
              if (/^\s*;/.test(text)) {
                cm.replaceRange(text.replace(/^(\s*);\s?/, '$1'), { line, ch: 0 }, { line, ch: text.length });
              } else {
                const lead = (text.match(/^(\s*)/) || ['', ''])[1];
                cm.replaceRange(`${lead}; ${text.slice(lead.length)}`, { line, ch: 0 }, { line, ch: text.length });
              }
            }
          });
        },
        F5: () => doRun(),
        'Ctrl-Shift-C': () => doCheck(),
        'Ctrl-Shift-B': () => doCompile(),
        'Alt-H': () => updateHelpAtCursor({ showTab: true }),
        F1: () => { document.querySelector('[data-activity="reference"]')?.click(); refSearch?.focus(); },
        'Ctrl-B': () => toggleLeftPanel(),
        'Ctrl-\\': () => toggleRightPanel()
      }
    });
    editor.on('cursorActivity', () => {
      const pos = editor.getCursor();
      $('cursor-pos').textContent = `Ln ${pos.line + 1}, Col ${pos.ch + 1}`;
      if (activeDocTab === 'cursor') scheduleHelpAtCursor();
    });
    editor.on('change', () => {
      openFiles[activeFile].dirty = true;
      updateTabs();
      scheduleCheck();
      scheduleDraftSave();
    });
    editor.on('focus', () => {
      editor.getWrapperElement().classList.add('is-focused');
    });
    editor.on('blur', () => {
      editor.getWrapperElement().classList.remove('is-focused');
    });
    $('editor-mount')?.addEventListener('mousedown', () => {
      requestAnimationFrame(() => focusEditor());
    });
    setTimeout(focusEditor, 0);
  }

  let draftTimer = null;
  function scheduleDraftSave() {
    clearTimeout(draftTimer);
    draftTimer = setTimeout(() => {
      syncEditor();
      try { localStorage.setItem(DRAFT_KEY, JSON.stringify({ openFiles, activeFile, ts: Date.now() })); } catch (_) {}
    }, 2000);
  }

  async function saveCurrent() {
    syncEditor();
    const f = openFiles[activeFile];
    if (isDesktopApp()) {
      const res = f.diskPath ? await saveNativeFile(f.content, f.name) : await saveNativeFileAs(f.content, f.name);
      if (res.success) {
        f.diskPath = res.path || f.diskPath;
        f.name = res.filename || f.name;
        f.dirty = false;
        updateTabs();
        addRecent(f.diskPath);
        log('success', 'Saved ' + f.name);
      } else log('error', res.error || 'Save failed');
    } else {
      downloadText(f.name, f.content);
      f.dirty = false;
      updateTabs();
      log('success', 'Downloaded ' + f.name);
    }
  }

  function downloadText(name, content) {
    const a = document.createElement('a');
    a.href = URL.createObjectURL(new Blob([content], { type: 'text/plain' }));
    a.download = name;
    a.click();
    URL.revokeObjectURL(a.href);
  }

  function icon(name, cls) {
    const s = document.createElement('span');
    s.className = cls || '';
    s.innerHTML = I[name] || '';
    return s;
  }

  function docSubmenu() {
    const featured = new Set([
      'BEGIN_HERE.md', 'GETTING_STARTED.md', 'FIRST_HOUR.md', 'LANGUAGE.md', 'PROGRAMMING.md',
      'COMMANDS.md', 'systems/GUIDES.md', 'systems/00-START.md', 'systems/01-CORE.md',
      'reference/ENTITY.md', 'reference/WINDOW.md', 'reference/DRAW.md', 'BUILDING.md', 'DEVELOPER.md'
    ]);
    const items = getDocIndex()
      .filter(d => featured.has(d.path))
      .map(d => ({ label: d.title, action: () => { showDocPanel(); loadDoc(d.path); } }));
    items.push({ sep: true });
    items.push({
      label: 'Browse all docs…',
      action: () => {
        showDocPanel();
        document.getElementById('doc-search')?.focus();
      }
    });
    return items;
  }

  function recentSubmenu() {
    const items = getRecent().map(p => ({
      label: p.split(/[/\\]/).pop() || p,
      action: () => { if (isDesktopApp()) openBuffer(p.split(/[/\\]/).pop(), '', p); }
    }));
    return items.length ? items : [{ label: '(empty)', disabled: true }];
  }

  function showSidebarPanel(activityId) {
    setLeftPanelVisible(true);
    const btn = document.querySelector(`[data-activity="${activityId}"]`);
    btn?.click();
  }

  function openSettings(tab) {
    showSettingsModal({ tab, onSaved: () => { refreshToolchainLog(); editor?.refresh(); } });
  }

  function setupMenus() {
    buildMenubar('menubar', [
      {
        label: 'File',
        items: [
          { label: 'New File', shortcut: 'Ctrl+N', action: () => openBuffer('untitled.mb', EXAMPLES.hello.code) },
          { label: 'Open File…', shortcut: 'Ctrl+O', action: () => openFileDialog() },
          { label: 'Open Folder…', action: () => openFolderDialog() },
          { label: 'Save', shortcut: 'Ctrl+S', action: () => saveCurrent() },
          { label: 'Save As…', action: () => saveAsDialog() },
          { sep: true },
          { label: 'New Project…', action: () => newProject() },
          { label: 'Recent Files', submenu: recentSubmenu() },
          { sep: true },
          { label: 'Settings…', action: () => openSettings('themes') },
          { sep: true },
          { label: 'Close File', action: () => closeTab(activeFile) }
        ]
      },
      {
        label: 'Edit',
        items: [
          { label: 'Undo', shortcut: 'Ctrl+Z', action: () => editor.undo() },
          { label: 'Redo', shortcut: 'Ctrl+Y', action: () => editor.redo() },
          { sep: true },
          { label: 'Find…', shortcut: 'Ctrl+F', action: () => findBar?.show() },
          { label: 'Go to Line…', shortcut: 'Ctrl+G', action: () => showGoToModal({ maxLine: editor.lineCount(), onGo: n => editor.setCursor(n - 1, 0) }) }
        ]
      },
      {
        label: 'Run',
        items: [
          { label: 'Run (moonrun)', shortcut: 'F5', action: () => doRun() },
          { label: 'Check Syntax', shortcut: 'Ctrl+Shift+C', action: () => doCheck() },
          { label: 'Compile to .mbc', shortcut: 'Ctrl+Shift+B', action: () => doCompile() }
        ]
      },
      {
        label: 'View',
        items: [
          { label: 'Explorer', action: () => showSidebarPanel('explorer') },
          { label: 'Examples', action: () => showSidebarPanel('examples') },
          { label: 'API Reference', action: () => showSidebarPanel('reference') },
          { label: 'Documentation', action: () => showSidebarPanel('docs') },
          { sep: true },
          {
            label: 'Panels',
            submenu: [
              ...panelMenuItems(),
              { sep: true },
              { label: 'Help at Cursor', shortcut: 'Alt+H', action: () => updateHelpAtCursor({ showTab: true }) },
              { sep: true },
              { label: 'Terminal', action: () => showPanel('terminal') },
              { label: 'Problems', action: () => showPanel('problems') }
            ]
          },
          { sep: true },
          {
            label: 'Settings',
            submenu: [
              { label: 'Themes…', action: () => openSettings('themes') },
              { label: 'Editor & Fonts…', action: () => openSettings('editor') },
              { label: 'Custom Colors…', action: () => openSettings('colors') },
              { label: 'Toolchain…', action: () => openSettings('toolchain') }
            ]
          }
        ]
      },
      {
        label: 'Help',
        items: [
        { label: 'Begin Here', action: () => { showDocPanel(); loadDoc('BEGIN_HERE.md'); } },
          { label: 'Language Reference', action: () => { showDocPanel(); loadDoc('LANGUAGE.md'); } },
          { label: 'Programming Guide', action: () => { showDocPanel(); loadDoc('PROGRAMMING.md'); } },
          { label: 'Command Index', action: () => { showDocPanel(); loadDoc('COMMANDS.md'); } },
          { sep: true },
          { label: 'All Documentation', submenu: docSubmenu() },
          { sep: true },
          { label: 'Keyboard Shortcuts', action: () => showShortcutsModal() },
          { label: 'Search API (F1)', action: () => { document.querySelector('[data-activity="reference"]')?.click(); refSearch?.focus(); } }
        ]
      }
    ]);
  }

  function buildActivityBar() {
    const items = [
      { id: 'explorer', icon: 'files', title: 'Explorer' },
      { id: 'examples', icon: 'library', title: 'Examples' },
      { id: 'snippets', icon: 'snippet', title: 'Snippets' },
      { id: 'docs', icon: 'book', title: 'Documentation' },
      { id: 'reference', icon: 'search', title: 'API Reference' }
    ];
    const bar = $('activity-bar');
    bar.innerHTML = '';

    const hideLeft = document.createElement('button');
    hideLeft.type = 'button';
    hideLeft.className = 'activity-btn panel-toggle-btn';
    hideLeft.id = 'btn-hide-left';
    hideLeft.title = 'Hide left sidebar (Ctrl+B)';
    hideLeft.textContent = '‹';
    bar.appendChild(hideLeft);

    items.forEach((it, i) => {
      const btn = document.createElement('button');
      btn.type = 'button';
      btn.className = 'activity-btn' + (i === 0 ? ' active' : '');
      btn.dataset.activity = it.id;
      btn.title = it.title;
      btn.appendChild(icon(it.icon, 'icon'));
      btn.addEventListener('click', () => {
        setLeftPanelVisible(true);
        if (it.id === 'docs') setRightPanelVisible(true);
        bar.querySelectorAll('.activity-btn[data-activity]').forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
        document.querySelectorAll('.sidebar-panel').forEach(p => p.classList.toggle('active', p.dataset.panel === it.id));
      });
      bar.appendChild(btn);
    });
  }

  function buildExamples() {
    const list = $('example-list');
    list.innerHTML = '';
    EXAMPLE_CATEGORIES.forEach(cat => {
      const h = document.createElement('div');
      h.className = 'sidebar-section';
      h.textContent = cat;
      list.appendChild(h);
      Object.entries(EXAMPLES).filter(([, ex]) => ex.category === cat).forEach(([, ex]) => {
        const btn = document.createElement('button');
        btn.type = 'button';
        btn.className = 'sidebar-item';
        btn.textContent = ex.title;
        btn.addEventListener('click', () => openBuffer(ex.title.replace(/\s+/g, '_').toLowerCase() + '.mb', ex.code));
        list.appendChild(btn);
      });
    });
  }

  function buildSnippets() {
    const list = $('snippet-list');
    list.innerHTML = '';
    Object.entries(SNIPPETS).forEach(([, sn]) => {
      const btn = document.createElement('button');
      btn.type = 'button';
      btn.className = 'sidebar-item';
      btn.textContent = sn.title;
      btn.addEventListener('click', () => editor.replaceSelection(sn.code));
      list.appendChild(btn);
    });
  }

  async function openFileDialog() {
    if (isDesktopApp()) {
      const res = await openNativeFile();
      if (res.success) openBuffer(res.filename, res.content, res.path);
      return;
    }
    $('file-input').click();
  }

  async function openFolderDialog() {
    if (isDesktopApp()) {
      const res = await openProjectFolder();
      if (res.success && res.files?.length) {
        projectPath = res.path;
        const files = res.files.map(f => ({
          name: f.name,
          content: f.content,
          diskPath: projectPath + '/' + f.name.replace(/\\/g, '/')
        }));
        explorer.setProject(projectPath, files);
        openBuffer(files[0].name, files[0].content, files[0].diskPath);
        log('success', `Opened folder: ${projectPath} (${files.length} .mb files)`);
      }
      return;
    }
    log('warn', 'Open Folder requires the desktop app');
  }

  async function saveAsDialog() {
    syncEditor();
    const f = openFiles[activeFile];
    if (isDesktopApp()) {
      const res = await saveNativeFileAs(f.content, f.name);
      if (res.success) { f.diskPath = res.path; f.name = res.filename || f.name; f.dirty = false; updateTabs(); addRecent(f.diskPath); }
    } else downloadText(f.name, f.content);
  }

  function newProject() {
    if (!folderPickerAvailable()) { log('error', 'New project needs desktop app or Chrome/Edge'); return; }
    showNewProjectModal({
      folderLabel: '', types: PROJECT_TYPES,
      onCreate: async type => {
        const proj = buildProject(type);
        const folder = await pickProjectFolder();
        if (!folder?.path) return;
        await writeProjectToDisk(folder, proj.files);
        const loaded = await readProjectFromFolder(folder.path);
        if (loaded?.length) {
          projectPath = folder.path;
          const files = loaded.map(f => ({ name: f.name, content: f.content, diskPath: folder.path + '/' + f.name }));
          explorer.setProject(projectPath, files);
          openFiles = files.map(f => ({ ...f, dirty: false }));
          activeFile = 0;
          editor.setValue(openFiles[0].content);
          updateTabs();
          log('success', 'Project created: ' + folder.path);
        }
      }
    });
  }

  function wireUi() {
    const addIcon = (id, name) => { const el = $(id); if (el) el.appendChild(icon(name, 'icon')); };
    addIcon('btn-open-folder', 'folder');
    addIcon('btn-new', 'newFile');
    addIcon('btn-open', 'folder');
    addIcon('btn-save', 'save');
    $('btn-open-folder')?.addEventListener('click', () => openFolderDialog());
    $('btn-new')?.addEventListener('click', () => openBuffer('untitled.mb', EXAMPLES.hello.code));
    $('btn-open')?.addEventListener('click', () => openFileDialog());
    $('btn-save')?.addEventListener('click', () => saveCurrent());
    addIcon('btn-run', 'play');
    $('btn-run').title = 'Run (F5)';
    $('btn-run').addEventListener('click', () => doRun());
    addIcon('btn-check', 'check');
    $('btn-check').title = 'Check';
    $('btn-check').addEventListener('click', () => doCheck());
    addIcon('btn-compile', 'debug');
    $('btn-compile').title = 'Compile';
    $('btn-compile').addEventListener('click', () => doCompile());

    const settingsBtn = $('btn-settings');
    const settingsDrop = $('settings-dropdown');
    const settingsRoot = $('settings-menu-root');
    if (settingsBtn) {
      settingsBtn.appendChild(icon('settings', 'icon'));
      settingsBtn.addEventListener('click', e => {
        e.stopPropagation();
        const open = !settingsDrop.hidden;
        settingsDrop.hidden = open;
        settingsRoot?.classList.toggle('open', !open);
      });
      settingsDrop?.querySelectorAll('[data-settings-open]').forEach(btn => {
        btn.addEventListener('click', e => {
          e.stopPropagation();
          settingsDrop.hidden = true;
          settingsRoot?.classList.remove('open');
          openSettings(btn.dataset.settingsOpen);
        });
      });
      document.addEventListener('click', () => {
        if (settingsDrop) settingsDrop.hidden = true;
        settingsRoot?.classList.remove('open');
      });
    }

    $('logo').appendChild(icon('moon', 'icon-lg'));
    $('logo').appendChild(document.createTextNode(' moonBASIC IDE'));

    document.querySelectorAll('.panel-tab').forEach(tab => {
      tab.addEventListener('click', () => showPanel(tab.dataset.panel));
    });

    document.querySelectorAll('[data-doc-tab]').forEach(tab => {
      tab.addEventListener('click', () => {
        showDocTab(tab.dataset.docTab);
        if (tab.dataset.docTab === 'cursor') {
          setRightPanelVisible(true);
          updateHelpAtCursor({ showTab: true });
        }
      });
    });

    refSearch?.addEventListener('input', () => buildReferenceList(refSearch.value));
    $('file-input')?.addEventListener('change', e => {
      const file = e.target.files?.[0];
      if (!file) return;
      const reader = new FileReader();
      reader.onload = () => openBuffer(file.name, reader.result);
      reader.readAsText(file);
    });
    document.querySelectorAll('[data-example]').forEach(btn => {
      btn.addEventListener('click', () => openBuffer(btn.dataset.example + '.mb', EXAMPLES[btn.dataset.example]?.code || ''));
    });
    $('btn-new-project')?.addEventListener('click', () => newProject());
    $('overlay-logo')?.appendChild(icon('moon', 'icon-lg'));
    $('find-prev')?.addEventListener('click', () => findBar?.findNext(true));
    $('find-next')?.addEventListener('click', () => findBar?.findNext(false));
  }

  const style = document.createElement('style');
  style.textContent = `
    .cm-s-mb-dark .cm-error-line{background:rgba(251,113,133,0.15)}
    .CodeMirror-hints{font-family:JetBrains Mono,Consolas,monospace;background:#12151c;border:1px solid rgba(168,85,247,0.35)}
    .help-panel-body,.doc-viewer{padding:12px;overflow:auto;height:100%;font-size:12px;line-height:1.55}
    .ref-item,.sidebar-item,.tree-file{display:flex;flex-direction:column;align-items:flex-start;width:100%;text-align:left;padding:6px 10px;background:none;border:none;color:inherit;cursor:pointer;border-radius:4px}
    .ref-item:hover,.sidebar-item:hover,.tree-file:hover{background:rgba(168,85,247,0.12)}
    .ref-key,.tree-file{font-family:var(--sb-mono);color:#67e8f9;font-size:11px}
    .ref-hint{color:#7c8aa5;font-size:11px;margin-top:2px}
    #ref-search,#doc-search{width:100%;margin-bottom:8px;padding:6px 8px;background:#0a0b10;border:1px solid rgba(120,130,180,0.2);color:inherit;border-radius:4px}
    .preview-overlay.hidden{display:none} .hidden{display:none!important}
    .problem-row{padding:4px 8px;cursor:pointer;border-bottom:1px solid rgba(120,130,180,0.08)}
    .problem-row:hover{background:rgba(251,113,133,0.1)} .problem-loc{color:#fb7185;margin-right:8px;font-family:var(--sb-mono)}
    .tab-close{margin-left:8px;opacity:0.6;font-size:14px} .tab-close:hover{opacity:1;color:#fb7185}
    .tree-root{font-weight:600;padding:6px 10px;color:#c084fc} .tree-files{padding-left:4px}
    .tree-file{flex-direction:row;gap:6px;align-items:center} .doc-item.active{background:rgba(168,85,247,0.2)}
    .menu-submenu{position:relative} .menu-has-sub .menu-sub-arrow{margin-left:auto}
    .menu-sub-dropdown{display:none;position:absolute;left:100%;top:0;min-width:220px;background:rgba(18,20,28,0.98);border:1px solid rgba(168,85,247,0.35);border-radius:8px;padding:6px 0}
    .menu-submenu:hover .menu-sub-dropdown{display:block}
    .preview-body-stack{position:relative;min-height:100%} .preview-overlay{position:absolute;inset:0;z-index:2}
  `;
  document.head.appendChild(style);

  document.addEventListener('keydown', e => {
    if (e.ctrlKey && !e.shiftKey && e.key.toLowerCase() === 'b') {
      e.preventDefault();
      toggleLeftPanel();
    }
    if (e.ctrlKey && e.key === '\\') {
      e.preventDefault();
      toggleRightPanel();
    }
  });

  initModals();
  buildActivityBar();
  initPanels({
    onLayoutChange: () => {
      setupMenus();
      editor?.refresh();
    }
  });
  buildExamples();
  buildSnippets();
  buildReferenceList();
  initCodeMirror();
  setAppearanceCallback(() => editor?.refresh());
  findBar = createFindBar({ barId: 'find-bar', inputId: 'find-input', closeId: 'find-close', getEditor: () => editor });
  wireUi();
  setupMenus();
  updateTabs();
  PanelResize.init({ onEditorRefresh: () => editor?.refresh() });
  window.addEventListener('resize', () => editor?.refresh());

  await initDocs({
    searchId: 'doc-search', listId: 'doc-list', viewerId: 'doc-viewer',
    onSelect: () => showDocPanel()
  });

  updateHelpAtCursor({ showTab: false });
  await refreshToolchainLog();
  if (!isDesktopApp()) log('dim', 'Browser mode — use desktop app for run/compile');
  log('accent', 'moonBASIC IDE — File menus, docs, F5 run, Ctrl+Shift+B compile');
}
