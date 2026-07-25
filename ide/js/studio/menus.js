/**
 * Application menu bar with dropdown menus
 */

export function buildMenubar(containerId, menus) {
  const bar = document.getElementById(containerId);
  if (!bar) return;
  bar.innerHTML = '';

  menus.forEach(menu => {
    const root = document.createElement('div');
    root.className = 'menu-root';

    const trigger = document.createElement('button');
    trigger.type = 'button';
    trigger.className = 'menu-trigger';
    trigger.textContent = menu.label;

    const drop = document.createElement('div');
    drop.className = 'menu-dropdown';

    for (const item of menu.items) {
      if (item.sep) {
        const hr = document.createElement('hr');
        hr.className = 'menu-sep';
        drop.appendChild(hr);
        continue;
      }
      if (item.submenu) {
        const sub = document.createElement('div');
        sub.className = 'menu-submenu';
        const subBtn = document.createElement('button');
        subBtn.type = 'button';
        subBtn.className = 'menu-item menu-has-sub';
        subBtn.innerHTML = `${item.label}<span class="menu-sub-arrow">›</span>`;
        const subDrop = document.createElement('div');
        subDrop.className = 'menu-sub-dropdown';
        for (const subItem of item.submenu) {
          if (subItem.sep) {
            subDrop.appendChild(document.createElement('hr'));
            continue;
          }
          subDrop.appendChild(makeMenuButton(subItem, root));
        }
        sub.append(subBtn, subDrop);
        drop.appendChild(sub);
        continue;
      }
      drop.appendChild(makeMenuButton(item, root));
    }

    trigger.addEventListener('click', e => {
      e.stopPropagation();
      document.querySelectorAll('.menu-root.open').forEach(r => r !== root && r.classList.remove('open'));
      root.classList.toggle('open');
    });

    root.append(trigger, drop);
    bar.appendChild(root);
  });

  document.addEventListener('click', () => {
    document.querySelectorAll('.menu-root.open').forEach(r => r.classList.remove('open'));
  });
}

function makeMenuButton(item, root) {
  const btn = document.createElement('button');
  btn.type = 'button';
  btn.className = 'menu-item';
  if (item.disabled) btn.disabled = true;
  btn.innerHTML = item.label + (item.shortcut ? `<span class="shortcut">${item.shortcut}</span>` : '');
  btn.addEventListener('click', e => {
    e.stopPropagation();
    root.classList.remove('open');
    if (!item.disabled && item.action) item.action();
  });
  return btn;
}
