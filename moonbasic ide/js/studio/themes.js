/**
 * IDE appearance themes and CSS variable application
 */
const APPEARANCE_KEY = 'moonbasic-ide-appearance';

export const THEME_PRESETS = {
  moonbasic: {
    id: 'moonbasic', name: 'moonBASIC', desc: 'Purple & cyan — default',
    swatch: ['#0c0e14', '#a855f7', '#22d3ee'],
    vars: {
      '--sb-bg': '#0c0e14', '--sb-bg-elevated': '#12151c', '--sb-sidebar': 'rgba(16, 18, 26, 0.92)',
      '--sb-activity': '#0a0b10', '--sb-panel': '#0e1018', '--sb-border': 'rgba(120, 130, 180, 0.14)',
      '--sb-border-bright': 'rgba(168, 85, 247, 0.35)', '--sb-tab-inactive': '#10131a', '--sb-tab-active': '#0c0e14',
      '--sb-accent': '#a855f7', '--sb-accent-2': '#22d3ee', '--sb-accent-hover': '#c084fc',
      '--sb-glow': 'rgba(168, 85, 247, 0.25)', '--sb-text': '#e2e8f0', '--sb-text-dim': '#7c8aa5',
      '--sb-text-bright': '#ffffff', '--sb-selection': 'rgba(168, 85, 247, 0.22)',
      '--sb-titlebar': 'linear-gradient(180deg, #141820 0%, #0f1118 100%)',
      '--sb-statusbar': 'linear-gradient(90deg, #6d28d9 0%, #0891b2 100%)',
      '--ed-bg': '#0a0c12', '--ed-fg': '#e8edf7', '--ed-gutter': '#080a10',
      '--ed-gutter-border': 'rgba(120, 130, 180, 0.18)', '--ed-line-num': '#4a5568',
      '--ed-cursor': '#22d3ee', '--ed-cursor-glow': 'rgba(34, 211, 238, 0.75)',
      '--ed-active-line': 'rgba(34, 211, 238, 0.07)', '--ed-active-line-bar': 'rgba(34, 211, 238, 0.55)',
      '--ed-selection': 'rgba(56, 189, 248, 0.22)',
      '--syn-keyword': '#e879f9', '--syn-builtin': '#67e8f9', '--syn-string': '#86efac',
      '--syn-comment': '#6b7a94', '--syn-number': '#fde68a', '--syn-property': '#fdba74',
      '--syn-atom': '#c084fc', '--syn-variable': '#e8edf7',
      '--doc-h1': '#c084fc', '--doc-h2': '#67e8f9', '--doc-h3': '#e879f9', '--doc-link': '#38bdf8'
    }
  },
  midnight: {
    id: 'midnight', name: 'Midnight', desc: 'Deep blue focus',
    swatch: ['#0a1628', '#3b82f6', '#60a5fa'],
    vars: {
      '--sb-bg': '#0a1628', '--sb-bg-elevated': '#0f1d32', '--sb-sidebar': 'rgba(10, 22, 40, 0.95)',
      '--sb-activity': '#081220', '--sb-panel': '#0c1830', '--sb-border': 'rgba(96, 165, 250, 0.15)',
      '--sb-border-bright': 'rgba(59, 130, 246, 0.4)', '--sb-tab-inactive': '#0c1830', '--sb-tab-active': '#0a1628',
      '--sb-accent': '#3b82f6', '--sb-accent-2': '#60a5fa', '--sb-accent-hover': '#93c5fd',
      '--sb-glow': 'rgba(59, 130, 246, 0.25)', '--sb-text': '#dbeafe', '--sb-text-dim': '#7da4d4',
      '--sb-text-bright': '#f0f9ff', '--sb-selection': 'rgba(59, 130, 246, 0.25)',
      '--sb-titlebar': 'linear-gradient(180deg, #0f1d32 0%, #0a1628 100%)',
      '--sb-statusbar': 'linear-gradient(90deg, #1d4ed8 0%, #2563eb 100%)',
      '--ed-bg': '#081220', '--ed-fg': '#e0f2fe', '--ed-gutter': '#061018',
      '--ed-gutter-border': 'rgba(96, 165, 250, 0.2)', '--ed-line-num': '#4b6a8a',
      '--ed-cursor': '#60a5fa', '--ed-cursor-glow': 'rgba(96, 165, 250, 0.7)',
      '--ed-active-line': 'rgba(59, 130, 246, 0.1)', '--ed-active-line-bar': 'rgba(96, 165, 250, 0.5)',
      '--ed-selection': 'rgba(59, 130, 246, 0.28)',
      '--syn-keyword': '#93c5fd', '--syn-builtin': '#60a5fa', '--syn-string': '#86efac',
      '--syn-comment': '#5a7a9a', '--syn-number': '#fcd34d', '--syn-property': '#f9a8d4',
      '--syn-atom': '#c4b5fd', '--syn-variable': '#e0f2fe',
      '--doc-h1': '#93c5fd', '--doc-h2': '#60a5fa', '--doc-h3': '#bfdbfe', '--doc-link': '#60a5fa'
    }
  },
  forest: {
    id: 'forest', name: 'Forest', desc: 'Calm green tones',
    swatch: ['#0c1410', '#22c55e', '#4ade80'],
    vars: {
      '--sb-bg': '#0c1410', '--sb-bg-elevated': '#121a14', '--sb-sidebar': 'rgba(12, 20, 16, 0.95)',
      '--sb-activity': '#0a100c', '--sb-panel': '#0e1612', '--sb-border': 'rgba(74, 222, 128, 0.12)',
      '--sb-border-bright': 'rgba(34, 197, 94, 0.35)', '--sb-tab-inactive': '#101810', '--sb-tab-active': '#0c1410',
      '--sb-accent': '#22c55e', '--sb-accent-2': '#4ade80', '--sb-accent-hover': '#86efac',
      '--sb-glow': 'rgba(34, 197, 94, 0.2)', '--sb-text': '#dcfce7', '--sb-text-dim': '#7a9a82',
      '--sb-text-bright': '#f0fdf4', '--sb-selection': 'rgba(34, 197, 94, 0.2)',
      '--sb-titlebar': 'linear-gradient(180deg, #141c16 0%, #0c1410 100%)',
      '--sb-statusbar': 'linear-gradient(90deg, #15803d 0%, #16a34a 100%)',
      '--ed-bg': '#0a100c', '--ed-fg': '#ecfdf5', '--ed-gutter': '#080e0a',
      '--ed-gutter-border': 'rgba(74, 222, 128, 0.15)', '--ed-line-num': '#4a6a52',
      '--ed-cursor': '#4ade80', '--ed-cursor-glow': 'rgba(74, 222, 128, 0.65)',
      '--ed-active-line': 'rgba(34, 197, 94, 0.08)', '--ed-active-line-bar': 'rgba(74, 222, 128, 0.45)',
      '--ed-selection': 'rgba(34, 197, 94, 0.22)',
      '--syn-keyword': '#86efac', '--syn-builtin': '#4ade80', '--syn-string': '#fde68a',
      '--syn-comment': '#6b8a72', '--syn-number': '#fbbf24', '--syn-property': '#a7f3d0',
      '--syn-atom': '#bbf7d0', '--syn-variable': '#ecfdf5',
      '--doc-h1': '#86efac', '--doc-h2': '#4ade80', '--doc-h3': '#bbf7d0', '--doc-link': '#4ade80'
    }
  },
  sunset: {
    id: 'sunset', name: 'Sunset', desc: 'Warm orange & rose',
    swatch: ['#14100c', '#f97316', '#fb7185'],
    vars: {
      '--sb-bg': '#14100c', '--sb-bg-elevated': '#1c1610', '--sb-sidebar': 'rgba(20, 16, 12, 0.95)',
      '--sb-activity': '#100c08', '--sb-panel': '#181410', '--sb-border': 'rgba(251, 146, 60, 0.15)',
      '--sb-border-bright': 'rgba(249, 115, 22, 0.35)', '--sb-tab-inactive': '#1a1410', '--sb-tab-active': '#14100c',
      '--sb-accent': '#f97316', '--sb-accent-2': '#fb7185', '--sb-accent-hover': '#fdba74',
      '--sb-glow': 'rgba(249, 115, 22, 0.22)', '--sb-text': '#fef3c7', '--sb-text-dim': '#a08a70',
      '--sb-text-bright': '#fffbeb', '--sb-selection': 'rgba(249, 115, 22, 0.22)',
      '--sb-titlebar': 'linear-gradient(180deg, #1c1610 0%, #14100c 100%)',
      '--sb-statusbar': 'linear-gradient(90deg, #c2410c 0%, #e11d48 100%)',
      '--ed-bg': '#100c08', '--ed-fg': '#fef3c7', '--ed-gutter': '#0c0806',
      '--ed-gutter-border': 'rgba(251, 146, 60, 0.18)', '--ed-line-num': '#7a6550',
      '--ed-cursor': '#fb923c', '--ed-cursor-glow': 'rgba(251, 146, 60, 0.7)',
      '--ed-active-line': 'rgba(249, 115, 22, 0.09)', '--ed-active-line-bar': 'rgba(251, 146, 60, 0.5)',
      '--ed-selection': 'rgba(249, 115, 22, 0.25)',
      '--syn-keyword': '#fdba74', '--syn-builtin': '#fb923c', '--syn-string': '#fde68a',
      '--syn-comment': '#8a7560', '--syn-number': '#fbbf24', '--syn-property': '#fda4af',
      '--syn-atom': '#fecdd3', '--syn-variable': '#fef3c7',
      '--doc-h1': '#fdba74', '--doc-h2': '#fb923c', '--doc-h3': '#fda4af', '--doc-link': '#fb923c'
    }
  },
  solarized: {
    id: 'solarized', name: 'Solarized Dark', desc: 'Classic readable dark',
    swatch: ['#002b36', '#2aa198', '#b58900'],
    vars: {
      '--sb-bg': '#002b36', '--sb-bg-elevated': '#073642', '--sb-sidebar': 'rgba(0, 43, 54, 0.96)',
      '--sb-activity': '#001f27', '--sb-panel': '#073642', '--sb-border': 'rgba(131, 148, 150, 0.2)',
      '--sb-border-bright': 'rgba(42, 161, 152, 0.4)', '--sb-tab-inactive': '#073642', '--sb-tab-active': '#002b36',
      '--sb-accent': '#2aa198', '--sb-accent-2': '#268bd2', '--sb-accent-hover': '#5dc4ba',
      '--sb-glow': 'rgba(42, 161, 152, 0.2)', '--sb-text': '#93a1a1', '--sb-text-dim': '#657b83',
      '--sb-text-bright': '#fdf6e3', '--sb-selection': 'rgba(42, 161, 152, 0.25)',
      '--sb-titlebar': 'linear-gradient(180deg, #073642 0%, #002b36 100%)',
      '--sb-statusbar': 'linear-gradient(90deg, #2aa198 0%, #268bd2 100%)',
      '--ed-bg': '#00212b', '--ed-fg': '#839496', '--ed-gutter': '#001a22',
      '--ed-gutter-border': 'rgba(131, 148, 150, 0.25)', '--ed-line-num': '#586e75',
      '--ed-cursor': '#2aa198', '--ed-cursor-glow': 'rgba(42, 161, 152, 0.6)',
      '--ed-active-line': 'rgba(42, 161, 152, 0.1)', '--ed-active-line-bar': 'rgba(42, 161, 152, 0.45)',
      '--ed-selection': 'rgba(42, 161, 152, 0.28)',
      '--syn-keyword': '#cb4b16', '--syn-builtin': '#268bd2', '--syn-string': '#2aa198',
      '--syn-comment': '#586e75', '--syn-number': '#d33682', '--syn-property': '#b58900',
      '--syn-atom': '#6c71c4', '--syn-variable': '#839496',
      '--doc-h1': '#2aa198', '--doc-h2': '#268bd2', '--doc-h3': '#b58900', '--doc-link': '#268bd2'
    }
  },
  light: {
    id: 'light', name: 'Paper Light', desc: 'Bright daytime theme',
    swatch: ['#f8fafc', '#7c3aed', '#0891b2'],
    vars: {
      '--sb-bg': '#f1f5f9', '--sb-bg-elevated': '#ffffff', '--sb-sidebar': 'rgba(248, 250, 252, 0.98)',
      '--sb-activity': '#e2e8f0', '--sb-panel': '#ffffff', '--sb-border': 'rgba(100, 116, 139, 0.2)',
      '--sb-border-bright': 'rgba(124, 58, 237, 0.35)', '--sb-tab-inactive': '#e2e8f0', '--sb-tab-active': '#ffffff',
      '--sb-accent': '#7c3aed', '--sb-accent-2': '#0891b2', '--sb-accent-hover': '#8b5cf6',
      '--sb-glow': 'rgba(124, 58, 237, 0.15)', '--sb-text': '#1e293b', '--sb-text-dim': '#64748b',
      '--sb-text-bright': '#0f172a', '--sb-selection': 'rgba(124, 58, 237, 0.15)',
      '--sb-titlebar': 'linear-gradient(180deg, #ffffff 0%, #f1f5f9 100%)',
      '--sb-statusbar': 'linear-gradient(90deg, #7c3aed 0%, #0891b2 100%)',
      '--ed-bg': '#ffffff', '--ed-fg': '#1e293b', '--ed-gutter': '#f1f5f9',
      '--ed-gutter-border': 'rgba(100, 116, 139, 0.2)', '--ed-line-num': '#94a3b8',
      '--ed-cursor': '#0891b2', '--ed-cursor-glow': 'rgba(8, 145, 178, 0.4)',
      '--ed-active-line': 'rgba(124, 58, 237, 0.06)', '--ed-active-line-bar': 'rgba(124, 58, 237, 0.35)',
      '--ed-selection': 'rgba(8, 145, 178, 0.18)',
      '--syn-keyword': '#7c3aed', '--syn-builtin': '#0891b2', '--syn-string': '#15803d',
      '--syn-comment': '#94a3b8', '--syn-number': '#b45309', '--syn-property': '#c2410c',
      '--syn-atom': '#6d28d9', '--syn-variable': '#1e293b',
      '--doc-h1': '#7c3aed', '--doc-h2': '#0891b2', '--doc-h3': '#6d28d9', '--doc-link': '#0891b2'
    }
  },
  highcontrast: {
    id: 'highcontrast', name: 'High Contrast', desc: 'Maximum readability',
    swatch: ['#000000', '#ffff00', '#00ffff'],
    vars: {
      '--sb-bg': '#000000', '--sb-bg-elevated': '#0a0a0a', '--sb-sidebar': '#0a0a0a',
      '--sb-activity': '#000000', '--sb-panel': '#111111', '--sb-border': 'rgba(255, 255, 255, 0.25)',
      '--sb-border-bright': 'rgba(255, 255, 0, 0.5)', '--sb-tab-inactive': '#111111', '--sb-tab-active': '#000000',
      '--sb-accent': '#ffff00', '--sb-accent-2': '#00ffff', '--sb-accent-hover': '#ffff66',
      '--sb-glow': 'rgba(255, 255, 0, 0.3)', '--sb-text': '#ffffff', '--sb-text-dim': '#cccccc',
      '--sb-text-bright': '#ffffff', '--sb-selection': 'rgba(255, 255, 0, 0.3)',
      '--sb-titlebar': 'linear-gradient(180deg, #1a1a1a 0%, #000000 100%)',
      '--sb-statusbar': 'linear-gradient(90deg, #666600 0%, #006666 100%)',
      '--ed-bg': '#000000', '--ed-fg': '#ffffff', '--ed-gutter': '#0a0a0a',
      '--ed-gutter-border': 'rgba(255,255,255,0.3)', '--ed-line-num': '#aaaaaa',
      '--ed-cursor': '#00ffff', '--ed-cursor-glow': 'rgba(0, 255, 255, 0.8)',
      '--ed-active-line': 'rgba(255, 255, 0, 0.12)', '--ed-active-line-bar': '#ffff00',
      '--ed-selection': 'rgba(255, 255, 0, 0.35)',
      '--syn-keyword': '#ffff00', '--syn-builtin': '#00ffff', '--syn-string': '#00ff00',
      '--syn-comment': '#888888', '--syn-number': '#ff9900', '--syn-property': '#ff66ff',
      '--syn-atom': '#cc99ff', '--syn-variable': '#ffffff',
      '--doc-h1': '#ffff00', '--doc-h2': '#00ffff', '--doc-h3': '#ff9900', '--doc-link': '#00ffff'
    }
  }
};

export const DEFAULT_APPEARANCE = {
  themeId: 'moonbasic',
  editorFontSize: 15,
  editorLineHeight: 1.65,
  uiFontSize: 13,
  fontMono: 'JetBrains Mono',
  colorOverrides: {}
};

export const COLOR_GROUPS = [
  { label: 'Interface', keys: [
    { id: '--sb-bg', label: 'Background' }, { id: '--sb-sidebar', label: 'Sidebar' },
    { id: '--sb-accent', label: 'Accent' }, { id: '--sb-accent-2', label: 'Accent 2' },
    { id: '--sb-text', label: 'Text' }, { id: '--sb-text-dim', label: 'Muted text' },
    { id: '--sb-border', label: 'Borders' }
  ]},
  { label: 'Editor', keys: [
    { id: '--ed-bg', label: 'Editor background' }, { id: '--ed-fg', label: 'Editor text' },
    { id: '--ed-gutter', label: 'Line numbers bg' }, { id: '--ed-cursor', label: 'Cursor' },
    { id: '--ed-active-line', label: 'Active line' }, { id: '--ed-selection', label: 'Selection' }
  ]},
  { label: 'Syntax', keys: [
    { id: '--syn-keyword', label: 'Keywords' }, { id: '--syn-builtin', label: 'Built-ins' },
    { id: '--syn-string', label: 'Strings' }, { id: '--syn-comment', label: 'Comments' },
    { id: '--syn-number', label: 'Numbers' }, { id: '--syn-property', label: 'Properties' }
  ]}
];

export const MONO_FONTS = [
  'JetBrains Mono', 'Cascadia Code', 'Fira Code', 'Consolas', 'Courier New', 'monospace'
];

let currentAppearance = { ...DEFAULT_APPEARANCE, colorOverrides: {} };
let onApplyCallback = null;

function mergeVars(appearance) {
  const theme = THEME_PRESETS[appearance.themeId] || THEME_PRESETS.moonbasic;
  return { ...theme.vars, ...appearance.colorOverrides };
}

export function getResolvedVars(appearance = currentAppearance) {
  return mergeVars(appearance);
}

export function applyAppearance(appearance = currentAppearance) {
  currentAppearance = {
    ...DEFAULT_APPEARANCE,
    ...appearance,
    colorOverrides: { ...(appearance.colorOverrides || {}) }
  };
  const vars = mergeVars(currentAppearance);
  const root = document.documentElement;
  for (const [k, v] of Object.entries(vars)) {
    root.style.setProperty(k, v);
  }
  root.style.setProperty('--ed-font-size', `${currentAppearance.editorFontSize}px`);
  root.style.setProperty('--ed-line-height', String(currentAppearance.editorLineHeight));
  root.style.setProperty('--ui-font-size', `${currentAppearance.uiFontSize}px`);
  const mono = currentAppearance.fontMono || 'JetBrains Mono';
  root.style.setProperty('--sb-mono', `"${mono}", Consolas, monospace`);
  document.body.style.fontSize = `${currentAppearance.uiFontSize}px`;
  document.body.dataset.theme = currentAppearance.themeId;
  onApplyCallback?.();
}

export function loadAppearance() {
  try {
    const raw = localStorage.getItem(APPEARANCE_KEY);
    if (raw) {
      const parsed = JSON.parse(raw);
      currentAppearance = {
        ...DEFAULT_APPEARANCE,
        ...parsed,
        colorOverrides: { ...(parsed.colorOverrides || {}) }
      };
    }
  } catch (_) {}
  applyAppearance(currentAppearance);
  return { ...currentAppearance };
}

export function saveAppearance(appearance) {
  currentAppearance = {
    ...DEFAULT_APPEARANCE,
    ...appearance,
    colorOverrides: { ...(appearance.colorOverrides || {}) }
  };
  localStorage.setItem(APPEARANCE_KEY, JSON.stringify(currentAppearance));
  applyAppearance(currentAppearance);
  return { ...currentAppearance };
}

export function getAppearance() {
  return { ...currentAppearance, colorOverrides: { ...currentAppearance.colorOverrides } };
}

export function initAppearance(opts = {}) {
  onApplyCallback = opts.onApply || null;
  return loadAppearance();
}

export function setAppearanceCallback(fn) {
  onApplyCallback = fn || null;
}

export function listThemes() {
  return Object.values(THEME_PRESETS);
}

export function hexForCssColor(value) {
  if (!value) return '#000000';
  const v = value.trim();
  if (/^#[0-9a-fA-F]{6}$/.test(v)) return v;
  const m = v.match(/rgba?\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)/);
  if (m) {
    const h = n => Math.max(0, Math.min(255, Number(n))).toString(16).padStart(2, '0');
    return `#${h(m[1])}${h(m[2])}${h(m[3])}`;
  }
  return '#888888';
}
