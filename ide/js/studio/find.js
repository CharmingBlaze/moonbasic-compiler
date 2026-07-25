/**
 * Find in file bar
 */

export function createFindBar({ barId, inputId, closeId, getEditor }) {
  const bar = document.getElementById(barId);
  const input = document.getElementById(inputId);
  const closeBtn = document.getElementById(closeId);
  let lastIdx = 0;

  function show() {
    if (bar) bar.hidden = false;
    input?.focus();
    input?.select();
  }

  function hide() {
    if (bar) bar.hidden = true;
    clearMarks();
  }

  function clearMarks() {
    const ed = getEditor();
    if (!ed) return;
    for (let i = 0; i < ed.lineCount(); i++) ed.removeLineClass(i, 'background', 'cm-find-line');
  }

  function findNext(backward = false) {
    const ed = getEditor();
    const q = input?.value?.trim();
    if (!ed || !q) return;
    const cursor = ed.getCursor();
    const from = backward
      ? { line: cursor.line, ch: cursor.ch - 1 }
      : { line: cursor.line, ch: cursor.ch + 1 };
    const found = ed.getSearchCursor(q, from, { caseFold: true });
    if (found.find(backward)) {
      ed.setSelection(found.from(), found.to());
      ed.scrollIntoView({ from: found.from(), to: found.to() }, 80);
      ed.focus();
      clearMarks();
      ed.addLineClass(found.from().line, 'background', 'cm-find-line');
    } else {
      const wrap = ed.getSearchCursor(q, { line: 0, ch: 0 }, { caseFold: true });
      if (wrap.find(backward)) {
        ed.setSelection(wrap.from(), wrap.to());
        ed.scrollIntoView({ from: wrap.from(), to: wrap.to() }, 80);
        ed.focus();
      }
    }
  }

  closeBtn?.addEventListener('click', hide);
  input?.addEventListener('keydown', e => {
    if (e.key === 'Enter') { e.preventDefault(); findNext(e.shiftKey); }
    if (e.key === 'Escape') hide();
  });
  input?.addEventListener('input', () => findNext());

  return { show, hide, findNext };
}
