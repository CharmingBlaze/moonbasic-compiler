/**
 * Show/hide left sidebar and right help panel
 */
const STORAGE_KEY = 'moonbasic-ide-panels';

let state = {
  leftVisible: true,
  rightVisible: true
};

let onLayoutChange = null;

function loadState() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) Object.assign(state, JSON.parse(raw));
  } catch (_) {}
}

function saveState() {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
  } catch (_) {}
}

function applyLayout() {
  document.body.classList.toggle('left-panel-hidden', !state.leftVisible);
  document.body.classList.toggle('right-panel-hidden', !state.rightVisible);
  const leftBtn = document.getElementById('btn-hide-left');
  if (leftBtn) {
    leftBtn.textContent = state.leftVisible ? '‹' : '›';
    leftBtn.title = state.leftVisible ? 'Hide left sidebar (Ctrl+B)' : 'Show left sidebar (Ctrl+B)';
  }
  onLayoutChange?.();
  saveState();
}

export function initPanels(opts = {}) {
  onLayoutChange = opts.onLayoutChange || null;
  loadState();
  applyLayout();

  document.getElementById('btn-hide-left')?.addEventListener('click', e => {
    e.stopPropagation();
    toggleLeftPanel();
  });
  document.getElementById('btn-hide-right')?.addEventListener('click', e => {
    e.stopPropagation();
    setRightPanelVisible(false);
  });
}

export function isLeftPanelVisible() {
  return state.leftVisible;
}

export function isRightPanelVisible() {
  return state.rightVisible;
}

export function setLeftPanelVisible(visible) {
  state.leftVisible = !!visible;
  applyLayout();
}

export function setRightPanelVisible(visible) {
  state.rightVisible = !!visible;
  applyLayout();
}

export function toggleLeftPanel() {
  setLeftPanelVisible(!state.leftVisible);
}

export function toggleRightPanel() {
  setRightPanelVisible(!state.rightVisible);
}

export function panelMenuItems() {
  return [
    {
      label: state.leftVisible ? 'Hide Left Sidebar' : 'Show Left Sidebar',
      shortcut: 'Ctrl+B',
      action: () => toggleLeftPanel()
    },
    {
      label: state.rightVisible ? 'Hide Help Panel' : 'Show Help Panel',
      shortcut: 'Ctrl+\\',
      action: () => toggleRightPanel()
    }
  ];
}
