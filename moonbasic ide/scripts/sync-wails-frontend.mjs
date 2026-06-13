/**
 * Copy web IDE assets into frontend/ for Wails //go:embed
 * Overwrites in place (no rmdir) to avoid EBUSY on Windows when Wails holds the folder.
 */
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const root = path.resolve(__dirname, '..');
const dest = path.join(root, 'frontend');

const SKIP_DIRS = new Set(['node_modules', 'frontend', 'legacy', 'templates', 'scripts', '.git']);

function copyTree(src, dst) {
  fs.mkdirSync(dst, { recursive: true });
  for (const ent of fs.readdirSync(src, { withFileTypes: true })) {
    if (SKIP_DIRS.has(ent.name)) continue;
    const s = path.join(src, ent.name);
    const d = path.join(dst, ent.name);
    if (ent.isDirectory()) copyTree(s, d);
    else fs.copyFileSync(s, d);
  }
}

fs.mkdirSync(dest, { recursive: true });
fs.copyFileSync(path.join(root, 'index.html'), path.join(dest, 'index.html'));
copyTree(path.join(root, 'css'), path.join(dest, 'css'));
copyTree(path.join(root, 'js'), path.join(dest, 'js'));
if (fs.existsSync(path.join(root, 'bundled-docs'))) {
  copyTree(path.join(root, 'bundled-docs'), path.join(dest, 'bundled-docs'));
}

console.log('Synced frontend/ for Wails embed (' + dest + ')');
