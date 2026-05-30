// Command browser generator: embeds commands.json into a searchable static HTML page.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type cmdEntry struct {
	Key         string   `json:"key"`
	Namespace   string   `json:"namespace,omitempty"`
	Description string   `json:"description,omitempty"`
	Example     string   `json:"example,omitempty"`
	Args        []string `json:"args,omitempty"`
}

type manifest struct {
	Commands []cmdEntry `json:"commands"`
}

type row struct {
	Key         string `json:"k"`
	Namespace   string `json:"ns,omitempty"`
	Description string `json:"d,omitempty"`
	Example     string `json:"e,omitempty"`
	Args        string `json:"a,omitempty"`
}

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	inPath := filepath.Join(root, "compiler", "builtinmanifest", "commands.json")
	outPath := filepath.Join(root, "web", "command-browser.html")
	data, err := os.ReadFile(inPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	byKey := make(map[string]row)
	for _, c := range m.Commands {
		key := strings.ToUpper(strings.TrimSpace(c.Key))
		if key == "" {
			continue
		}
		args := strings.Join(c.Args, ", ")
		if prev, ok := byKey[key]; ok {
			if len(c.Description) <= len(prev.Description) {
				continue
			}
		}
		byKey[key] = row{
			Key:         key,
			Namespace:   c.Namespace,
			Description: c.Description,
			Example:     c.Example,
			Args:        args,
		}
	}
	rows := make([]row, 0, len(byKey))
	for _, r := range byKey {
		rows = append(rows, r)
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Key < rows[j].Key })
	payload, err := json.Marshal(rows)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	html := buildHTML(payload)
	if err := os.WriteFile(outPath, []byte(html), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("wrote %s (%d commands)\n", outPath, len(rows))
}

func buildHTML(payload []byte) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>moonBASIC Command Browser</title>
<style>
:root{color-scheme:dark light}
body{font-family:system-ui,sans-serif;margin:0;background:#0f1115;color:#e8eaed}
header{padding:1.25rem 1.5rem;border-bottom:1px solid #2a2f3a;background:#151922;position:sticky;top:0;z-index:1}
h1{font-size:1.15rem;margin:0 0 .75rem}
#q{width:100%;max-width:36rem;padding:.55rem .75rem;border:1px solid #3a4254;border-radius:6px;background:#0f1115;color:inherit;font-size:1rem}
main{padding:1rem 1.5rem 2rem}
table{width:100%;border-collapse:collapse;max-width:72rem}
th,td{padding:.45rem .6rem;text-align:left;border-bottom:1px solid #252a35;vertical-align:top}
th{font-size:.75rem;text-transform:uppercase;letter-spacing:.04em;color:#9aa3b5}
code{font-family:ui-monospace,Consolas,monospace;font-size:.9em;color:#8bd5ff}
.ns{color:#9aa3b5;font-size:.85rem}
.empty{padding:2rem;color:#9aa3b5}
</style>
</head>
<body>
<header>
<h1>moonBASIC Command Browser</h1>
<input id="q" type="search" placeholder="Search commands, namespaces, descriptions…" autofocus>
</header>
<main>
<table>
<thead><tr><th>Command</th><th>Namespace</th><th>Description</th><th>Example</th></tr></thead>
<tbody id="rows"></tbody>
</table>
<p class="empty" id="empty" hidden>No matches.</p>
</main>
<script>
const COMMANDS = ` + string(payload) + `;
const q = document.getElementById('q');
const tbody = document.getElementById('rows');
const empty = document.getElementById('empty');
function render() {
  const term = q.value.trim().toLowerCase();
  tbody.replaceChildren();
  let n = 0;
  for (const c of COMMANDS) {
    const hay = (c.k + ' ' + (c.ns||'') + ' ' + (c.d||'') + ' ' + (c.a||'')).toLowerCase();
    if (term && !hay.includes(term)) continue;
    const tr = document.createElement('tr');
    tr.innerHTML = '<td><code>' + esc(c.k) + '</code></td>'
      + '<td class="ns">' + esc(c.ns||'—') + '</td>'
      + '<td>' + esc(c.d||'') + (c.a ? ' <span class="ns">(' + esc(c.a) + ')</span>' : '') + '</td>'
      + '<td>' + (c.e ? '<code>' + esc(c.e) + '</code>' : '') + '</td>';
    tbody.appendChild(tr);
    n++;
  }
  empty.hidden = n > 0;
}
function esc(s){return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');}
q.addEventListener('input', render);
render();
</script>
</body>
</html>`
}
