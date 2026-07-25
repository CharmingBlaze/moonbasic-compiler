/** moonBASIC IDE — resizable panel system */
let active = null;

function start(type, e, opts) {
  active = { type, startX: e.clientX, startY: e.clientY, ...opts };
  document.body.classList.add('is-resizing');
  document.body.style.cursor = type === 'panel' ? 'row-resize' : 'col-resize';
  e.preventDefault();
}

function onMove(e) {
  if (!active) return;
  const { type, startX, startY } = active;

  if (type === 'sidebar') {
    const w = Math.max(140, Math.min(360, active.startW + (e.clientX - startX)));
    active.el.style.width = w + 'px';
  }
  if (type === 'preview') {
    const row = active.row;
    const total = row.offsetWidth;
    const dx = e.clientX - startX;
    const ew = Math.max(320, Math.min(total - 280, active.startEW + dx));
    active.editor.style.flex = 'none';
    active.editor.style.width = ew + 'px';
    active.preview.style.width = (total - ew - 5) + 'px';
    if (active.onRefresh) active.onRefresh();
  }
  if (type === 'panel') {
    const stack = active.stack;
    const maxH = stack.offsetHeight - 120;
    const h = Math.max(64, Math.min(maxH, active.startH - (e.clientY - startY)));
    active.panel.style.height = h + 'px';
    active.panel.style.flexShrink = '0';
    if (active.onRefresh) active.onRefresh();
  }
}

function onUp() {
  if (!active) return;
  active = null;
  document.body.classList.remove('is-resizing');
  document.body.style.cursor = '';
}

export const PanelResize = {
  init(opts) {
    opts = opts || {};

    const sidebar = document.getElementById('sidebar');
    const sidebarHandle = document.getElementById('resize-sidebar');
    if (sidebar && sidebarHandle) {
      sidebarHandle.addEventListener('mousedown', e => {
        start('sidebar', e, { el: sidebar, startW: sidebar.offsetWidth });
      });
    }

    const editor = document.getElementById('editor-stack');
    const preview = document.getElementById('help-stack');
    const previewHandle = document.getElementById('split-h');
    const row = document.querySelector('.editor-row');
    if (editor && preview && previewHandle && row) {
      previewHandle.addEventListener('mousedown', e => {
        start('preview', e, {
          row, editor, preview,
          startEW: editor.offsetWidth,
          onRefresh: opts.onEditorRefresh
        });
      });
    }

    const panel = document.getElementById('bottom-panel');
    const panelHandle = document.getElementById('resize-panel');
    const stack = document.getElementById('editor-stack');
    if (panel && panelHandle && stack) {
      panelHandle.addEventListener('mousedown', e => {
        start('panel', e, {
          panel, stack,
          startH: panel.offsetHeight,
          onRefresh: opts.onEditorRefresh
        });
      });
    }

    document.addEventListener('mousemove', onMove);
    document.addEventListener('mouseup', onUp);
  }
};
