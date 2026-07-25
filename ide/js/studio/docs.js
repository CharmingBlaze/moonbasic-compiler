/**
 * moonBASIC IDE — documentation browser
 */
import { isDesktopApp, appApi } from './toolchain.js';
import { setRightPanelVisible } from './panels.js';

let docIndex = [];
let docSearchEl;
let docListEl;
let docViewerEl;
let currentDocPath = '';
let onDocNavigate = null;

export async function initDocs({ searchId, listId, viewerId, onSelect }) {
  docSearchEl = document.getElementById(searchId);
  docListEl = document.getElementById(listId);
  docViewerEl = document.getElementById(viewerId);
  onDocNavigate = onSelect;

  try {
    const res = await fetch('js/studio/docs-index.json');
    if (res.ok) docIndex = await res.json();
  } catch (_) {}

  docSearchEl?.addEventListener('input', () => renderDocList(docSearchEl.value));
  docViewerEl?.addEventListener('click', handleDocClick);

  renderDocList('');

  if (docIndex.length) {
    const first = docIndex.find(d => d.id === 'begin') || docIndex[0];
    await loadDoc(first.path, onSelect);
  }
}

function handleDocClick(e) {
  const anchor = e.target.closest('a.doc-anchor, a.doc-heading-link');
  if (anchor) {
    const href = anchor.getAttribute('href') || '';
    if (href.startsWith('#')) {
      e.preventDefault();
      scrollToAnchor(href.slice(1));
    }
    return;
  }
  const internal = e.target.closest('a.doc-internal');
  if (internal) {
    e.preventDefault();
    const path = internal.dataset.docPath;
    const hash = internal.dataset.docHash || '';
    loadDoc(path, onDocNavigate).then(() => {
      if (hash) scrollToAnchor(hash);
    });
    showDocPanel();
  }
}

export function showDocPanel() {
  setRightPanelVisible(true);
  document.querySelector('[data-activity="docs"]')?.click();
  document.querySelectorAll('[data-doc-tab]').forEach(t => {
    t.classList.toggle('active', t.dataset.docTab === 'guide');
  });
  document.getElementById('doc-viewer')?.classList.remove('hidden');
  document.getElementById('help-cursor-panel')?.classList.add('hidden');
  document.getElementById('preview-overlay')?.classList.add('hidden');
}

function scrollToAnchor(id) {
  const el = docViewerEl?.querySelector(`#${CSS.escape(id)}`);
  if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' });
}

function renderDocList(filter) {
  if (!docListEl) return;
  const q = (filter || '').trim().toLowerCase();
  docListEl.innerHTML = '';
  const cats = new Map();
  for (const d of docIndex) {
    const hay = `${d.title} ${d.category} ${d.path}`.toLowerCase();
    if (q && !hay.includes(q)) continue;
    if (!cats.has(d.category)) cats.set(d.category, []);
    cats.get(d.category).push(d);
  }
  const order = ['Start', 'Language', 'Guides', 'Systems', 'Reference', 'API Reference', 'Command Set', 'Tutorials', 'Developers', 'Architecture', 'General'];
  const sorted = [...cats.entries()].sort((a, b) => {
    const ai = order.indexOf(a[0]);
    const bi = order.indexOf(b[0]);
    if (ai !== -1 || bi !== -1) return (ai === -1 ? 999 : ai) - (bi === -1 ? 999 : bi);
    return a[0].localeCompare(b[0]);
  });
  for (const [cat, items] of sorted) {
    const h = document.createElement('div');
    h.className = 'sidebar-section';
    h.textContent = cat;
    docListEl.appendChild(h);
    for (const d of items) {
      const btn = document.createElement('button');
      btn.type = 'button';
      btn.className = 'sidebar-item doc-item';
      btn.dataset.path = d.path;
      btn.innerHTML = `<span class="doc-item-title">${escapeHtml(d.title)}</span>` +
        (q ? `<span class="doc-item-path">${escapeHtml(d.path)}</span>` : '');
      btn.classList.toggle('active', d.path === currentDocPath);
      btn.addEventListener('click', () => loadDoc(d.path, onDocNavigate));
      docListEl.appendChild(btn);
    }
  }
  if (!cats.size && q) {
    const p = document.createElement('p');
    p.className = 'sidebar-hint';
    p.textContent = 'No docs match your search.';
    docListEl.appendChild(p);
  }
}

export async function loadDoc(path, callback) {
  currentDocPath = path.replace(/\\/g, '/');
  let content = '';
  if (isDesktopApp() && appApi()?.ReadBundledDoc) {
    const res = await appApi().ReadBundledDoc(currentDocPath);
    if (res.success) content = res.content;
  }
  if (!content) {
    try {
      const res = await fetch('bundled-docs/' + currentDocPath);
      if (res.ok) content = await res.text();
    } catch (_) {}
  }
  if (!content) {
    content = `*Documentation file not found: \`${currentDocPath}\`*\n\nRun \`npm run langdata\` or \`go run ./tools/docsexport\` from the repo root.`;
  }
  const html = markdownToHtml(content, currentDocPath);
  if (docViewerEl) docViewerEl.innerHTML = html;
  document.querySelectorAll('.doc-item').forEach(el => {
    el.classList.toggle('active', el.dataset.path === currentDocPath);
  });
  document.getElementById('preview-overlay')?.classList.add('hidden');
  if (callback) callback(currentDocPath, html);
  return content;
}

export function resolveDocPath(currentDoc, href) {
  if (!href || /^https?:\/\//i.test(href) || href.startsWith('mailto:')) return null;
  const hashIdx = href.indexOf('#');
  const pathPart = hashIdx >= 0 ? href.slice(0, hashIdx) : href;
  const hash = hashIdx >= 0 ? href.slice(hashIdx + 1) : '';
  if (!pathPart) return hash ? { path: currentDoc, hash } : null;
  if (!/\.md$/i.test(pathPart)) return null;

  const baseDir = currentDoc.includes('/') ? currentDoc.replace(/\/[^/]+$/, '') : '';
  const joined = pathPart.startsWith('/') ? pathPart.slice(1) : (
    baseDir ? `${baseDir}/${pathPart}` : pathPart
  );
  const parts = joined.split('/');
  const out = [];
  for (const p of parts) {
    if (p === '..') out.pop();
    else if (p && p !== '.') out.push(p);
  }
  return { path: out.join('/'), hash };
}

export function markdownToHtml(md, basePath = '') {
  let s = md.replace(/\r\n/g, '\n');

  s = s.replace(/```(\w*)\n([\s\S]*?)```/g, (_, lang, code) =>
    `<pre class="doc-code"><code>${escapeHtml(code.trim())}</code></pre>`);

  s = s.replace(/^(\|.+\|)\n(\|[-:| ]+\|)\n((?:\|.+\|\n?)+)/gm, (_, header, _sep, body) => {
    const ths = header.split('|').filter(Boolean).map(c => `<th>${inlineMd(c.trim())}</th>`).join('');
    const rows = body.trim().split('\n').map(row => {
      const tds = row.split('|').filter(Boolean).map(c => `<td>${inlineMd(c.trim())}</td>`).join('');
      return `<tr>${tds}</tr>`;
    }).join('');
    return `<table class="doc-table"><thead><tr>${ths}</tr></thead><tbody>${rows}</tbody></table>`;
  });

  const slugCounts = {};
  const heading = (level, text) => {
    let slug = slugify(text);
    slugCounts[slug] = (slugCounts[slug] || 0) + 1;
    if (slugCounts[slug] > 1) slug += `-${slugCounts[slug]}`;
    return `<h${level} id="${slug}"><a class="doc-heading-link" href="#${slug}">${inlineMd(text)}</a></h${level}>`;
  };
  s = s.replace(/^#### (.+)$/gm, (_, t) => heading(4, t));
  s = s.replace(/^### (.+)$/gm, (_, t) => heading(3, t));
  s = s.replace(/^## (.+)$/gm, (_, t) => heading(2, t));
  s = s.replace(/^# (.+)$/gm, (_, t) => heading(1, t));

  s = s.replace(/^---$/gm, '<hr class="doc-hr">');
  s = s.replace(/^> (.+)$/gm, '<blockquote>$1</blockquote>');

  s = s.replace(/\[([^\]]+)\]\(([^)]+)\)/g, (match, text, href) => {
    if (/^https?:\/\//i.test(href)) {
      return `<a href="${escapeHtml(href)}" target="_blank" rel="noopener">${inlineMd(text)}</a>`;
    }
    if (href.startsWith('#')) {
      return `<a href="${escapeHtml(href)}" class="doc-anchor">${inlineMd(text)}</a>`;
    }
    const resolved = resolveDocPath(basePath, href);
    if (resolved?.path) {
      const attrs = `class="doc-internal" data-doc-path="${escapeHtml(resolved.path)}"` +
        (resolved.hash ? ` data-doc-hash="${escapeHtml(resolved.hash)}"` : '');
      return `<a href="#" ${attrs}>${inlineMd(text)}</a>`;
    }
    return `<span class="doc-missing" title="${escapeHtml(href)}">${inlineMd(text)}</span>`;
  });

  s = s.replace(/^\d+\. (.+)$/gm, '<li class="doc-ol">$1</li>');
  s = s.replace(/(<li class="doc-ol">.*<\/li>\n?)+/g, m => `<ol>${m.replace(/ class="doc-ol"/g, '')}</ol>`);
  s = s.replace(/^- (.+)$/gm, '<li>$1</li>');
  s = s.replace(/(<li>.*<\/li>\n?)+/g, m => `<ul>${m}</ul>`);

  s = s.split(/\n\n+/).map(p => {
    p = p.trim();
    if (!p) return '';
    if (/^<(h[1-4]|pre|ul|ol|table|blockquote|hr)/.test(p)) return p;
    return `<p>${p.replace(/\n/g, '<br>')}</p>`;
  }).join('\n');

  return `<article class="doc-article">${s}</article>`;
}

function inlineMd(text) {
  return escapeHtml(text)
    .replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
    .replace(/`([^`]+)`/g, '<code>$1</code>');
}

function slugify(text) {
  return text.toLowerCase()
    .replace(/[^\w\s-]/g, '')
    .trim()
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-');
}

function escapeHtml(t) {
  return String(t)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

export function getDocIndex() {
  return docIndex;
}

export function openDocForCommand(key) {
  const ns = key.includes('.') ? key.split('.')[0] : key;
  const candidates = [
    `reference/${ns}.md`,
    `reference/moonbasic-command-set/${ns.toLowerCase()}.md`
  ];
  for (const path of candidates) {
    if (docIndex.some(d => d.path.toLowerCase() === path.toLowerCase())) {
      showDocPanel();
      return loadDoc(path);
    }
  }
  return null;
}
