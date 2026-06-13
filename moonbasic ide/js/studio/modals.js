/**
 * Simple modal dialogs (replaces browser prompt())
 */

let activeModal = null;
/** @type {((modal: HTMLElement) => void) | null} */
let onModalCloseHook = null;

export function setModalCloseHook(fn) {
  onModalCloseHook = fn || null;
}

export function closeModal() {
  if (!activeModal) return;
  onModalCloseHook?.(activeModal);
  onModalCloseHook = null;
  activeModal.hidden = true;
  activeModal = null;
}

export function openModal(id) {
  closeModal();
  const el = document.getElementById(id);
  if (!el) return null;
  el.hidden = false;
  activeModal = el;
  return el;
}

export function initModals() {
  document.querySelectorAll('[data-modal-close]').forEach(btn => {
    btn.addEventListener('click', closeModal);
  });
  document.querySelectorAll('.modal-backdrop').forEach(bg => {
    bg.addEventListener('click', e => {
      if (e.target === bg) closeModal();
    });
  });
  document.addEventListener('keydown', e => {
    if (e.key === 'Escape') closeModal();
  });
}

/**
 * @param {{ title: string, options: { id: string, label: string, hint?: string }[], onPick: (id: string) => void }} opts
 */
export function showChoiceModal({ title, options, onPick }) {
  const modal = openModal('modal-choice');
  if (!modal) return;
  modal.querySelector('.modal-title').textContent = title;
  const list = modal.querySelector('.modal-options');
  list.innerHTML = '';
  options.forEach(opt => {
    const btn = document.createElement('button');
    btn.type = 'button';
    btn.className = 'modal-option';
    btn.innerHTML = '<span class="modal-option-label">' + opt.label + '</span>' +
      (opt.hint ? '<span class="modal-option-hint">' + opt.hint + '</span>' : '');
    btn.addEventListener('click', () => {
      closeModal();
      onPick(opt.id);
    });
    list.appendChild(btn);
  });
}

export const EXPORT_FORMATS = [
  { id: 'html', label: 'Single HTML file', hint: 'One self-contained .html game' },
  { id: 'html5', label: 'HTML5 ZIP', hint: 'Unzip and open index.html in a browser' },
  { id: 'desktop', label: 'Desktop ZIP', hint: 'Ready-to-run .exe / launcher — no Go or build tools' }
];

/**
 * @param {{ folderLabel: string, types: { id: string, label: string, hint?: string }[], onCreate: (type: string) => void }} opts
 */
export function showNewProjectModal({ folderLabel, types, onCreate }) {
  let pickedType = null;

  const modal = openModal('modal-project');
  if (!modal) return;

  const titleEl = modal.querySelector('.modal-title');
  const stepEl = modal.querySelector('.project-step-label');
  const folderInput = modal.querySelector('#project-folder-display');
  const list = modal.querySelector('.modal-options');
  const createBtn = modal.querySelector('[data-project-create]');

  titleEl.textContent = 'New Project';
  stepEl.textContent = 'Choose what you are building — grow it in code with #Include, modules, and libraries.';
  if (folderInput) folderInput.value = folderLabel || '';

  list.innerHTML = '';
  (types || []).forEach(opt => {
    const btn = document.createElement('button');
    btn.type = 'button';
    btn.className = 'modal-option';
    btn.innerHTML = '<span class="modal-option-label">' + opt.label + '</span>' +
      (opt.hint ? '<span class="modal-option-hint">' + opt.hint + '</span>' : '');
    btn.addEventListener('click', () => {
      pickedType = opt.id;
      list.querySelectorAll('.modal-option').forEach(b => b.classList.remove('selected'));
      btn.classList.add('selected');
    });
    list.appendChild(btn);
  });

  createBtn.onclick = () => {
    if (!pickedType) return;
    closeModal();
    onCreate(pickedType);
  };
}

export const PLATFORM_OPTIONS = [
  { id: 'windows', label: 'Windows (64-bit)', hint: 'Includes .exe + Start Game.bat' },
  { id: 'linux', label: 'Linux (64-bit)', hint: 'Native binary + start.sh' },
  { id: 'mac-intel', label: 'macOS (Intel)', hint: 'Universal desktop player' },
  { id: 'mac-arm', label: 'macOS (Apple Silicon)', hint: 'Apple Silicon desktop player' }
];

export const KEYBOARD_SHORTCUTS = [
  ['Run (moonrun)', 'F5'],
  ['Toggle left sidebar', 'Ctrl+B'],
  ['Toggle help panel', 'Ctrl+\\'],
  ['Check syntax', 'Ctrl+Shift+C'],
  ['Compile to .mbc', 'Ctrl+Shift+B'],
  ['Save', 'Ctrl+S'],
  ['Open file', 'Ctrl+O'],
  ['Find', 'Ctrl+F'],
  ['Go to line', 'Ctrl+G'],
  ['Autocomplete', 'Ctrl+Space'],
  ['Help at cursor', 'Alt+H'],
  ['API search', 'F1']
];

/**
 * @param {{ maxLine: number, onGo: (line: number) => void }} opts
 */
export function showGoToModal({ maxLine, onGo }) {
  const modal = openModal('modal-goto');
  if (!modal) return;
  const input = modal.querySelector('#goto-line-input');
  input.value = '';
  input.max = String(maxLine);
  const go = () => {
    const n = parseInt(input.value, 10);
    if (n >= 1 && n <= maxLine) {
      closeModal();
      onGo(n);
    }
  };
  modal.querySelector('#goto-line-go').onclick = go;
  input.onkeydown = e => { if (e.key === 'Enter') { e.preventDefault(); go(); } };
  input.focus();
}

export function showShortcutsModal() {
  const modal = openModal('modal-shortcuts');
  if (!modal) return;
  const list = modal.querySelector('#shortcuts-list');
  list.innerHTML = '';
  KEYBOARD_SHORTCUTS.forEach(([label, key]) => {
    const row = document.createElement('div');
    row.className = 'shortcut-row';
    row.innerHTML = '<span>' + label + '</span><kbd>' + key + '</kbd>';
    list.appendChild(row);
  });
}
