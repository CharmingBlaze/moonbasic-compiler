/**
 * moonBASIC toolchain bridge (Wails desktop app)
 */

export function isDesktopApp() {
  return !!(window.go?.main?.App);
}

export function appApi() {
  return window.go?.main?.App;
}

export async function getToolchain() {
  const api = appApi();
  if (!api?.GetToolchain) return { found: false };
  return api.GetToolchain();
}

export async function checkFile(filePath) {
  const api = appApi();
  if (!api?.CheckFile) {
    return { success: false, error: 'Check requires moonBASIC IDE desktop app with moonbasic on PATH' };
  }
  return api.CheckFile(filePath);
}

export async function compileFile(filePath) {
  const api = appApi();
  if (!api?.CompileFile) {
    return { success: false, error: 'Compile requires moonBASIC IDE desktop app with moonbasic on PATH' };
  }
  return api.CompileFile(filePath);
}

export async function runFile(filePath) {
  const api = appApi();
  if (!api?.RunFile) {
    return { success: false, error: 'Run requires moonBASIC IDE desktop app with moonrun (full runtime)' };
  }
  return api.RunFile(filePath);
}

export async function openNativeFile() {
  const api = appApi();
  if (api?.OpenFile) return api.OpenFile();
  return { success: false, error: 'Open file needs desktop app' };
}

export async function saveNativeFile(content, filename) {
  const api = appApi();
  if (api?.SaveFile) return api.SaveFile(content, filename);
  return { success: false, error: 'Save needs desktop app' };
}

export async function openProjectFolder() {
  const api = appApi();
  if (api?.OpenFolder) return api.OpenFolder();
  return { success: false, error: 'Open folder needs desktop app' };
}

export async function readBundledDoc(path) {
  const api = appApi();
  if (api?.ReadBundledDoc) return api.ReadBundledDoc(path);
  try {
    const res = await fetch('bundled-docs/' + path.replace(/\\/g, '/'));
    if (res.ok) return { success: true, path, content: await res.text() };
  } catch (_) {}
  return { success: false, error: 'Doc not found' };
}

export async function saveNativeFileAs(content, defaultName) {
  const api = appApi();
  if (api?.SaveFileAs) return api.SaveFileAs(content, defaultName);
  return { success: false, error: 'Save As needs desktop app' };
}

export async function getLSPHover(filePath, content, line, col) {
  const api = appApi();
  if (!api?.GetLSPHover) return '';
  return api.GetLSPHover(filePath, content, line, col) || '';
}

export async function getLSPCompletion(filePath, content, line, col) {
  const api = appApi();
  if (!api?.GetLSPCompletion) return [];
  return api.GetLSPCompletion(filePath, content, line, col) || [];
}

export async function getIDESettings() {
  const api = appApi();
  if (!api?.GetIDESettings) {
    try {
      return JSON.parse(localStorage.getItem('moonbasic-ide-settings') || '{}');
    } catch (_) {
      return {};
    }
  }
  return api.GetIDESettings();
}

export async function setIDESettings(settings) {
  const api = appApi();
  if (!api?.SetIDESettings) {
    try {
      localStorage.setItem('moonbasic-ide-settings', JSON.stringify(settings));
    } catch (_) {}
    return settings;
  }
  return api.SetIDESettings(settings);
}

export async function browseToolchain(which) {
  const api = appApi();
  if (!api?.BrowseToolchain) return '';
  return api.BrowseToolchain(which) || '';
}

export async function testToolchain() {
  const api = appApi();
  if (!api?.TestToolchain) return getToolchain();
  return api.TestToolchain();
}
