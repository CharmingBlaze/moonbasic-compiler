/**
 * IDE settings — appearance, themes, and toolchain
 */
import { openModal, closeModal, setModalCloseHook } from './modals.js';
import {
  isDesktopApp, appApi, getToolchain,
  getIDESettings, setIDESettings, browseToolchain, testToolchain
} from './toolchain.js';
import {
  DEFAULT_APPEARANCE, COLOR_GROUPS, MONO_FONTS,
  getAppearance, applyAppearance, saveAppearance, getResolvedVars, hexForCssColor, listThemes
} from './themes.js';

const TOOLCHAIN_KEY = 'moonbasic-ide-settings';

let draftAppearance = null;
let savedAppearanceSnapshot = null;

export async function loadToolchainSettings() {
  if (isDesktopApp()) {
    try {
      const s = await getIDESettings();
      localStorage.setItem(TOOLCHAIN_KEY, JSON.stringify(s));
      return s;
    } catch (_) {}
  }
  try {
    return JSON.parse(localStorage.getItem(TOOLCHAIN_KEY) || '{}');
  } catch (_) {
    return {};
  }
}

export async function saveToolchainSettings(s) {
  const payload = {
    moonbasicPath: (s.moonbasicPath || '').trim(),
    moonrunPath: (s.moonrunPath || '').trim()
  };
  if (isDesktopApp()) {
    try {
      const saved = await setIDESettings(payload);
      localStorage.setItem(TOOLCHAIN_KEY, JSON.stringify(saved));
      return saved;
    } catch (_) {}
  }
  localStorage.setItem(TOOLCHAIN_KEY, JSON.stringify(payload));
  return payload;
}

function populateFontSelect(modal) {
  const fontEl = modal.querySelector('#set-font-mono');
  if (!fontEl || fontEl.options.length) return;
  MONO_FONTS.forEach(f => {
    const opt = document.createElement('option');
    opt.value = f;
    opt.textContent = f;
    fontEl.appendChild(opt);
  });
}

function revertAppearance() {
  if (savedAppearanceSnapshot) applyAppearance(savedAppearanceSnapshot);
}

function readDraftFromModal(modal) {
  const appearance = {
    themeId: modal.querySelector('.theme-card.active')?.dataset.themeId || draftAppearance.themeId,
    editorFontSize: Number(modal.querySelector('#set-editor-size')?.value) || 15,
    editorLineHeight: Number(modal.querySelector('#set-editor-line')?.value) || 1.65,
    uiFontSize: Number(modal.querySelector('#set-ui-size')?.value) || 13,
    fontMono: modal.querySelector('#set-font-mono')?.value || 'JetBrains Mono',
    colorOverrides: { ...draftAppearance.colorOverrides }
  };
  modal.querySelectorAll('[data-color-key]').forEach(input => {
    const key = input.dataset.colorKey;
    const themeVal = getResolvedVars({ ...appearance, colorOverrides: {} })[key];
    const hex = input.value;
    if (hex && hexForCssColor(themeVal) !== hex) {
      appearance.colorOverrides[key] = hex;
    } else {
      delete appearance.colorOverrides[key];
    }
  });
  return appearance;
}

function bindAppearanceControls(modal) {
  const sizeEl = modal.querySelector('#set-editor-size');
  const sizeVal = modal.querySelector('#set-editor-size-val');
  const lineEl = modal.querySelector('#set-editor-line');
  const lineVal = modal.querySelector('#set-editor-line-val');
  const uiEl = modal.querySelector('#set-ui-size');
  const uiVal = modal.querySelector('#set-ui-size-val');
  const fontEl = modal.querySelector('#set-font-mono');

  const preview = () => {
    draftAppearance = readDraftFromModal(modal);
    applyAppearance(draftAppearance);
  };

  sizeEl?.addEventListener('input', () => {
    if (sizeVal) sizeVal.textContent = sizeEl.value;
    preview();
  });
  lineEl?.addEventListener('input', () => {
    if (lineVal) lineVal.textContent = Number(lineEl.value).toFixed(2);
    preview();
  });
  uiEl?.addEventListener('input', () => {
    if (uiVal) uiVal.textContent = uiEl.value;
    preview();
  });
  fontEl?.addEventListener('change', preview);
  modal.querySelectorAll('[data-color-key]').forEach(el => {
    el.addEventListener('input', preview);
  });
}

function buildThemeGrid(modal, themeId) {
  const grid = modal.querySelector('#settings-theme-grid');
  if (!grid) return;
  grid.innerHTML = '';
  listThemes().forEach(t => {
    const card = document.createElement('button');
    card.type = 'button';
    card.className = 'theme-card' + (t.id === themeId ? ' active' : '');
    card.dataset.themeId = t.id;
    card.innerHTML =
      `<div class="theme-swatch">${t.swatch.map(c => `<span style="background:${c}"></span>`).join('')}</div>` +
      `<div class="theme-card-name">${t.name}</div>` +
      `<div class="theme-card-desc">${t.desc}</div>`;
    card.addEventListener('click', () => {
      grid.querySelectorAll('.theme-card').forEach(c => c.classList.remove('active'));
      card.classList.add('active');
      draftAppearance.themeId = t.id;
      draftAppearance.colorOverrides = {};
      rebuildColorPickers(modal);
      applyAppearance(draftAppearance);
      syncSliders(modal);
    });
    grid.appendChild(card);
  });
}

function rebuildColorPickers(modal) {
  const host = modal.querySelector('#settings-color-groups');
  if (!host) return;
  host.innerHTML = '';
  const resolved = getResolvedVars(draftAppearance);
  COLOR_GROUPS.forEach(group => {
    const h = document.createElement('h4');
    h.className = 'settings-section-title';
    h.textContent = group.label;
    host.appendChild(h);
    const grid = document.createElement('div');
    grid.className = 'color-grid';
    group.keys.forEach(({ id, label }) => {
      const row = document.createElement('div');
      row.className = 'color-picker-row';
      const val = draftAppearance.colorOverrides[id] || resolved[id];
      row.innerHTML =
        `<label>${label}</label>` +
        `<input type="color" data-color-key="${id}" value="${hexForCssColor(val)}" title="${id}">`;
      grid.appendChild(row);
    });
    host.appendChild(grid);
  });
  modal.querySelectorAll('[data-color-key]').forEach(el => {
    el.addEventListener('input', () => {
      draftAppearance = readDraftFromModal(modal);
      applyAppearance(draftAppearance);
    });
  });
}

function syncSliders(modal) {
  const a = draftAppearance;
  const sizeEl = modal.querySelector('#set-editor-size');
  const lineEl = modal.querySelector('#set-editor-line');
  const uiEl = modal.querySelector('#set-ui-size');
  const fontEl = modal.querySelector('#set-font-mono');
  if (sizeEl) sizeEl.value = a.editorFontSize;
  if (lineEl) lineEl.value = a.editorLineHeight;
  if (uiEl) uiEl.value = a.uiFontSize;
  if (fontEl) fontEl.value = a.fontMono;
  modal.querySelector('#set-editor-size-val').textContent = a.editorFontSize;
  modal.querySelector('#set-editor-line-val').textContent = Number(a.editorLineHeight).toFixed(2);
  modal.querySelector('#set-ui-size-val').textContent = a.uiFontSize;
}

function showSettingsTab(modal, tab) {
  modal.querySelectorAll('.settings-nav-btn').forEach(b => {
    b.classList.toggle('active', b.dataset.settingsTab === tab);
  });
  modal.querySelectorAll('.settings-pane').forEach(p => {
    p.classList.toggle('active', p.dataset.settingsPane === tab);
  });
}

async function updateToolchainStatus(el) {
  if (!el) return;
  el.textContent = 'Checking…';
  el.className = 'settings-status';
  try {
    const tc = isDesktopApp() ? await testToolchain() : await getToolchain();
    if (tc.found) {
      el.textContent = `moonbasic: ${tc.moonbasic || 'found'}${tc.moonrun ? ' · moonrun: ' + tc.moonrun : ''}`;
      el.classList.add('ok');
    } else {
      el.textContent = 'Toolchain not found — set paths below or leave blank to auto-detect';
      el.classList.add('warn');
    }
  } catch (e) {
    el.textContent = String(e?.message || e);
    el.classList.add('err');
  }
}

/**
 * @param {{ onSaved?: () => void, tab?: string }} [opts]
 */
export async function showSettingsModal(opts = {}) {
  const modal = openModal('modal-settings');
  if (!modal) return;

  draftAppearance = getAppearance();
  savedAppearanceSnapshot = JSON.parse(JSON.stringify(draftAppearance));

  populateFontSelect(modal);
  buildThemeGrid(modal, draftAppearance.themeId);
  rebuildColorPickers(modal);
  syncSliders(modal);
  bindAppearanceControls(modal);

  setModalCloseHook(m => {
    if (m.id === 'modal-settings') revertAppearance();
  });

  modal.querySelector('#settings-close-x')?.addEventListener('click', () => {
    revertAppearance();
    closeModal();
  }, { once: true });

  const startTab = opts.tab || 'themes';
  showSettingsTab(modal, startTab);

  modal.querySelectorAll('.settings-nav-btn').forEach(btn => {
    btn.onclick = () => showSettingsTab(modal, btn.dataset.settingsTab);
  });

  const mbInput = modal.querySelector('#settings-moonbasic');
  const mrInput = modal.querySelector('#settings-moonrun');
  const statusEl = modal.querySelector('#settings-status');
  const hintEl = modal.querySelector('#settings-hint');

  if (hintEl) {
    const localDir = isDesktopApp() && appApi()?.GetLocalToolchainDir
      ? await appApi().GetLocalToolchainDir()
      : 'toolchain/';
    hintEl.textContent = isDesktopApp()
      ? `Local folder: ${localDir || 'toolchain/'} — or set explicit paths below.`
      : 'Compiler paths apply in the desktop app only.';
  }

  const tc = await loadToolchainSettings();
  mbInput.value = tc.moonbasicPath || '';
  mrInput.value = tc.moonrunPath || '';
  await updateToolchainStatus(statusEl);

  modal.querySelector('#settings-browse-moonbasic').onclick = async () => {
    if (!isDesktopApp()) return;
    const p = await browseToolchain('moonbasic');
    if (p) mbInput.value = p;
    await updateToolchainStatus(statusEl);
  };
  modal.querySelector('#settings-browse-moonrun').onclick = async () => {
    if (!isDesktopApp()) return;
    const p = await browseToolchain('moonrun');
    if (p) mrInput.value = p;
    await updateToolchainStatus(statusEl);
  };
  modal.querySelector('#settings-test-toolchain').onclick = async () => {
    await saveToolchainSettings({ moonbasicPath: mbInput.value, moonrunPath: mrInput.value });
    await updateToolchainStatus(statusEl);
  };
  modal.querySelector('#settings-clear-toolchain').onclick = async () => {
    mbInput.value = '';
    mrInput.value = '';
    await saveToolchainSettings({ moonbasicPath: '', moonrunPath: '' });
    await updateToolchainStatus(statusEl);
  };

  modal.querySelector('#settings-reset-appearance').onclick = () => {
    draftAppearance = {
      ...DEFAULT_APPEARANCE,
      themeId: 'moonbasic',
      colorOverrides: {}
    };
    buildThemeGrid(modal, draftAppearance.themeId);
    rebuildColorPickers(modal);
    syncSliders(modal);
    applyAppearance(draftAppearance);
  };

  modal.querySelector('#settings-reset-colors').onclick = () => {
    draftAppearance.colorOverrides = {};
    rebuildColorPickers(modal);
    applyAppearance(draftAppearance);
  };

  modal.querySelector('#settings-cancel').onclick = () => {
    revertAppearance();
    closeModal();
  };

  modal.querySelector('#settings-save').onclick = async () => {
    draftAppearance = readDraftFromModal(modal);
    saveAppearance(draftAppearance);
    savedAppearanceSnapshot = JSON.parse(JSON.stringify(draftAppearance));
    setModalCloseHook(null);
    await saveToolchainSettings({
      moonbasicPath: mbInput.value,
      moonrunPath: mrInput.value
    });
    await updateToolchainStatus(statusEl);
    closeModal();
    opts.onSaved?.();
  };
}

// Back-compat alias for toolchain-only loaders
export const loadSettings = loadToolchainSettings;
export const saveSettings = saveToolchainSettings;
