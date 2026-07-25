/**
 * Pick a project folder and write starter files (Wails, File System Access API, or in-browser)
 */

export async function pickProjectFolder() {
  if (window.go?.main?.App?.PickProjectFolder) {
    const res = await window.go.main.App.PickProjectFolder();
    if (res?.success && res.path) {
      return { kind: 'native', path: res.path, label: res.path };
    }
    return null;
  }

  if (typeof window.showDirectoryPicker === 'function') {
    try {
      const handle = await window.showDirectoryPicker({ mode: 'readwrite' });
      return { kind: 'handle', handle, path: handle.name, label: handle.name };
    } catch (_) {
      return null;
    }
  }

  return null;
}

export async function writeProjectToDisk(target, files) {
  if (!target) return false;

  if (target.kind === 'native' && window.go?.main?.App?.WriteProjectFiles) {
    const res = await window.go.main.App.WriteProjectFiles(
      target.path,
      JSON.stringify(files.map(f => ({ name: f.name, content: f.content })))
    );
    return !!res?.success;
  }

  if (target.kind === 'handle' && target.handle) {
    for (const f of files) {
      const fh = await target.handle.getFileHandle(f.name, { create: true });
      const w = await fh.createWritable();
      await w.write(f.content);
      await w.close();
    }
    return true;
  }

  return false;
}

export function folderPickerAvailable() {
  return !!(window.go?.main?.App?.PickProjectFolder || typeof window.showDirectoryPicker === 'function');
}

const MB_EXT = /\.mb$/i;

export async function readProjectFromFolder(target) {
  if (!target) return [];

  if (target.kind === 'native' && window.go?.main?.App?.ReadProjectFolder) {
    const res = await window.go.main.App.ReadProjectFolder(target.path);
    if (res?.success && res.files) return res.files;
    return [];
  }

  if (target.kind === 'handle' && target.handle) {
    const files = [];
    async function walk(dirHandle, prefix) {
      for await (const [name, handle] of dirHandle.entries()) {
        const path = prefix ? prefix + '/' + name : name;
        if (handle.kind === 'directory' && (name === 'src' || name === 'lib')) {
          await walk(handle, path);
        } else if (handle.kind === 'file' && MB_EXT.test(name)) {
          const f = await handle.getFile();
          files.push({ name: path.replace(/\\/g, '/'), content: await f.text() });
        }
      }
    }
    await walk(target.handle, '');
    return files;
  }

  return [];
}
