/**
 * moonBASIC language data — loaded from manifest-generated lang-data.json
 */
let _data = null;

export async function loadLangData() {
  if (_data) return _data;
  const res = await fetch('js/studio/lang-data.json');
  if (!res.ok) throw new Error('Failed to load lang-data.json — run: go run ./tools/ideexport');
  _data = await res.json();
  return _data;
}

export function getLangData() {
  return _data;
}

/** Build regex-safe alternation for keywords/globals/namespaces */
function alt(words) {
  return words
    .slice()
    .sort((a, b) => b.length - a.length)
    .map(w => w.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'))
    .join('|');
}

export function buildMoonBasicMode(data) {
  const kw = alt(data.keywords || []);
  const globals = alt(data.globals || []);
  const namespaces = alt(data.namespaces || []);
  const handles = alt((data.handleMethods || []).map(h => h.toLowerCase()));

  return {
    keywords: data.keywords || [],
    globals: data.globals || [],
    namespaces: data.namespaces || [],
    commands: data.commands || {},
    namespaceIndex: data.namespaceIndex || {},
    handleMethods: data.handleMethods || [],

    token(stream) {
      if (stream.match(/^;.*/)) return 'comment';
      if (stream.match(/^"([^"\\]|\\.)*"/)) return 'string';
      if (stream.match(/^-?\d+\.\d+([eE][+-]?\d+)?/)) return 'number';
      if (stream.match(/^-?\d+/)) return 'number';
      if (stream.match(/^KEY_[A-Z0-9_]+/i)) return 'atom';
      if (stream.match(new RegExp(`^(${namespaces})\\.([A-Za-z_][A-Za-z0-9_]*)`, 'i'))) return 'builtin';
      if (stream.match(new RegExp(`^\\.(${handles})\\b`, 'i'))) return 'property';
      if (stream.match(new RegExp(`^(${kw})\\b`, 'i'))) return 'keyword';
      if (stream.match(new RegExp(`^(${globals})\\b`, 'i'))) return 'builtin';
      if (stream.match(/^[A-Za-z_][A-Za-z0-9_]*/)) return 'variable';
      stream.next();
      return null;
    }
  };
}

export function formatCommandHelp(key, entry) {
  if (!entry) return '';
  const lines = [`### \`${key}\``];
  if (entry.description) lines.push('', entry.description);
  if (entry.args?.length) lines.push('', `**Arguments:** ${entry.args.join(', ')}`);
  if (entry.returns) lines.push('', `**Returns:** \`${entry.returns}\``);
  if (entry.phase) lines.push('', `**Phase:** \`${entry.phase}\``);
  if (entry.stub) lines.push('', `> ${entry.stub}`);
  return lines.join('\n');
}

export function lookupAtCursor(line, col, data) {
  const before = line.slice(0, col);

  const dotted = before.match(/([A-Za-z_][A-Za-z0-9_]*)\s*\.\s*([A-Za-z_][A-Za-z0-9_]*)$/i);
  if (dotted) {
    const ns = dotted[1].toUpperCase();
    const method = dotted[2].toUpperCase();
    const key = `${ns}.${method}`;
    if (data.commands[key]) return { key, entry: data.commands[key] };
    const methods = data.namespaceIndex[ns] || [];
    const match = methods.find(m => m.key.toUpperCase() === method)
      || methods.find(m => m.key.toUpperCase().startsWith(method));
    if (match) {
      const full = `${ns}.${match.key}`;
      return {
        key: full,
        entry: data.commands[full] || {
          description: `**${full}**`,
          args: match.args,
          returns: match.returns,
          phase: match.phase
        }
      };
    }
  }

  const handleDot = before.match(/\.([A-Za-z_][A-Za-z0-9_]*)$/i);
  if (handleDot) {
    const method = handleDot[1].toUpperCase();
    const handles = data.handleMethods || [];
    if (handles.some(h => h.toUpperCase() === method)) {
      return {
        key: `.${method}`,
        entry: {
          description: `Universal handle method **.${method}** — see reference/UNIVERSAL_HANDLE_METHODS.md`
        }
      };
    }
  }

  const word = before.match(/([A-Za-z_][A-Za-z0-9_]*)$/i);
  if (word) {
    const g = word[1].toUpperCase();
    if (data.commands[g]) return { key: g, entry: data.commands[g] };
    if (data.namespaceIndex[g]) {
      const methods = data.namespaceIndex[g];
      return {
        key: g,
        entry: {
          description: `**${g}** namespace — ${methods.length} commands. Type \`${g}.\` for methods.`
        }
      };
    }
  }
  return null;
}

export function completionItems(line, col, data) {
  const before = line.slice(0, col);
  const dot = before.lastIndexOf('.');
  if (dot >= 0) {
    let ns = before.slice(0, dot).trim();
    const sep = Math.max(ns.lastIndexOf(' '), ns.lastIndexOf('('), ns.lastIndexOf(','), ns.lastIndexOf(':'));
    if (sep >= 0) ns = ns.slice(sep + 1).trim();
    ns = ns.toUpperCase();
    const methods = data.namespaceIndex[ns] || [];
    const prefix = before.slice(dot + 1).toUpperCase();
    return methods
      .filter(m => !prefix || m.key.toUpperCase().startsWith(prefix))
      .map(m => ({
        text: m.key,
        displayText: m.key,
        hint: m.description?.slice(0, 80) || '',
        doc: formatCommandHelp(`${ns}.${m.key}`, m)
      }));
  }
  const word = (before.match(/([A-Za-z_][A-Za-z0-9_]*)$/) || ['', ''])[1].toUpperCase();
  const items = [];
  for (const ns of data.namespaces || []) {
    if (!word || ns.startsWith(word)) {
      items.push({ text: ns, displayText: ns, hint: 'namespace', doc: `**${ns}** namespace` });
    }
  }
  for (const g of data.globals || []) {
    if (!word || g.startsWith(word)) {
      const e = data.commands[g];
      items.push({ text: g, displayText: g, hint: e?.description?.slice(0, 60) || 'builtin', doc: formatCommandHelp(g, e) });
    }
  }
  for (const k of data.keywords || []) {
    if (!word || k.startsWith(word)) {
      items.push({ text: k, displayText: k, hint: 'keyword', doc: `**${k}** language keyword` });
    }
  }
  return items.slice(0, 80);
}
