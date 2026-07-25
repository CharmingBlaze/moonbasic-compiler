/**
 * File explorer tree for open project folder
 */

export function createExplorer({ treeEl, onOpenFile }) {
  let projectPath = '';
  let files = [];

  function render() {
    if (!treeEl) return;
    treeEl.innerHTML = '';
    if (!files.length) {
      treeEl.innerHTML = '<p class="sidebar-hint">Use <strong>File → Open Folder</strong> to browse project <code>.mb</code> files, or <strong>Open File</strong> for a single script.</p>';
      return;
    }
    if (projectPath) {
      const root = document.createElement('div');
      root.className = 'tree-root';
      root.innerHTML = `<span class="tree-folder-name">${basename(projectPath)}</span>`;
      treeEl.appendChild(root);
    }
    const ul = document.createElement('div');
    ul.className = 'tree-files';
    for (const f of files.sort((a, b) => a.name.localeCompare(b.name))) {
      const btn = document.createElement('button');
      btn.type = 'button';
      btn.className = 'tree-file';
      btn.innerHTML = `<span class="tree-file-icon">📄</span><span>${f.name}</span>`;
      btn.addEventListener('click', () => onOpenFile(f));
      ul.appendChild(btn);
    }
    treeEl.appendChild(ul);
  }

  function setProject(path, fileList) {
    projectPath = path || '';
    files = fileList || [];
    render();
  }

  function addOpenFile(name) {
    if (!files.find(f => f.name === name)) {
      files.push({ name, content: '', diskPath: '' });
      render();
    }
  }

  return { setProject, addOpenFile, render };
}

function basename(p) {
  return p.replace(/\\/g, '/').split('/').pop() || p;
}
