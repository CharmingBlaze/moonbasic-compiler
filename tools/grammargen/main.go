// Grammargen writes TextMate grammar for VS Code / IDE from the builtin manifest.
// Run from module root: go run ./tools/grammargen
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"moonbasic/compiler/builtinmanifest"
)

type tmGrammar struct {
	Schema    string            `json:"$schema"`
	Name      string            `json:"name"`
	ScopeName string            `json:"scopeName"`
	Patterns  []map[string]any  `json:"patterns"`
	Repo      map[string]any    `json:"repository"`
}

func regexEscape(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '.', '+', '*', '?', '^', '$', '(', ')', '[', ']', '{', '}', '|', '\\':
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

func altCI(words []string) string {
	parts := make([]string, len(words))
	for i, w := range words {
		parts[i] = "(?i)" + regexEscape(w)
	}
	sort.Slice(parts, func(i, j int) bool { return len(parts[i]) > len(parts[j]) })
	return strings.Join(parts, "|")
}

func main() {
	vscodeOut := filepath.Join("editors", "vscode-moonbasic", "syntaxes", "moonbasic.tmLanguage.json")
	if len(os.Args) > 1 {
		vscodeOut = os.Args[1]
	}

	table := builtinmanifest.Default()
	nsSet := make(map[string]struct{})
	globalSet := make(map[string]struct{})
	for _, key := range table.Keys() {
		if table.IsDeprecatedAlias(key) {
			continue
		}
		parts := strings.SplitN(key, ".", 2)
		if len(parts) == 1 {
			globalSet[strings.ToUpper(parts[0])] = struct{}{}
		} else {
			nsSet[parts[0]] = struct{}{}
		}
	}
	namespaces := make([]string, 0, len(nsSet))
	for ns := range nsSet {
		namespaces = append(namespaces, ns)
	}
	sort.Strings(namespaces)
	globals := make([]string, 0, len(globalSet))
	for g := range globalSet {
		globals = append(globals, g)
	}
	sort.Strings(globals)

	keywords := []string{
		"IF", "THEN", "ELSE", "ENDIF", "WHILE", "WEND", "ENDWHILE",
		"FOR", "TO", "DOWNTO", "STEP", "NEXT", "REPEAT", "UNTIL",
		"SELECT", "CASE", "DEFAULT", "ENDSELECT",
		"FUNCTION", "ENDFUNCTION", "RETURN", "TYPE", "FIELD", "ENDTYPE",
		"NEW", "DELETE", "EACH", "GOTO", "GOSUB",
		"DIM", "REDIM", "PRESERVE", "LOCAL", "GLOBAL", "CONST", "STATIC",
		"SWAP", "ERASE", "INCLUDE", "IMPORT", "END",
		"AND", "OR", "NOT", "XOR", "MOD", "TRUE", "FALSE", "NULL",
		"PRINT", "PRINTLN", "WRITE", "HELP", "INPUT", "CLS", "LOCATE", "TAB",
	}

	handleMethods := []string{
		"pos", "rot", "scale", "col", "alpha", "free", "move", "turn",
		"look", "fov", "visible", "name", "begin", "end", "target",
		"hide", "show", "draw", "update", "orbit", "yaw", "zoom",
		"velocity", "gravity", "bounce", "friction", "mass",
	}

	// Chunk namespaces for regex size limits (~100 per group)
	const chunk = 80
	var nsPatterns []map[string]any
	for i := 0; i < len(namespaces); i += chunk {
		end := i + chunk
		if end > len(namespaces) {
			end = len(namespaces)
		}
		nsPatterns = append(nsPatterns, map[string]any{
			"name":  "entity.name.function.moonbasic",
			"match": fmt.Sprintf(`(?i)\b(%s)\.([A-Za-z_][A-Za-z0-9_]*)`, strings.Join(namespaces[i:end], "|")),
		})
	}

	g := tmGrammar{
		Schema:    "https://raw.githubusercontent.com/martinring/tmlanguage/master/tmlanguage.json",
		Name:      "moonBASIC",
		ScopeName: "source.moonbasic",
		Patterns: []map[string]any{
			{"include": "#comments"},
			{"include": "#strings"},
			{"include": "#floats"},
			{"include": "#ints"},
			{"include": "#keyConstants"},
			{"include": "#namespaceCalls"},
			{"include": "#handleMethods"},
			{"include": "#globals"},
			{"include": "#keywords"},
			{"include": "#operators"},
		},
		Repo: map[string]any{
			"comments": map[string]any{
				"name":  "comment.line.semicolon.moonbasic",
				"begin": ";",
				"end":   "$",
			},
			"strings": map[string]any{
				"name":  "string.quoted.double.moonbasic",
				"begin": "\"",
				"end":   "\"",
				"patterns": []map[string]any{
					{"name": "constant.character.escape.moonbasic", "match": `\\.`},
				},
			},
			"floats": map[string]any{
				"name":  "constant.numeric.float.moonbasic",
				"match": `(?i)(?<![\w.])(-)?\d+\.\d+([eE][+-]?\d+)?\b`,
			},
			"ints": map[string]any{
				"name":  "constant.numeric.integer.moonbasic",
				"match": `(?i)(?<![\w.])(-)?\d+\b`,
			},
			"keyConstants": map[string]any{
				"name":  "constant.language.moonbasic",
				"match": `(?i)\bKEY_[A-Z0-9_]+\b`,
			},
			"namespaceCalls": map[string]any{
				"patterns": nsPatterns,
			},
			"handleMethods": map[string]any{
				"name":  "entity.name.tag.moonbasic",
				"match": fmt.Sprintf(`(?i)(?<=[\w\)])\s*\.(%s)\b`, strings.Join(handleMethods, "|")),
			},
			"globals": map[string]any{
				"name":  "support.function.moonbasic",
				"match": fmt.Sprintf(`(?i)\b(%s)\b`, altCI(globals)),
			},
			"keywords": map[string]any{
				"patterns": []map[string]any{
					{
						"name":  "keyword.control.moonbasic",
						"match": fmt.Sprintf(`(?i)\b(%s)\b`, strings.Join(keywords[:len(keywords)-8], "|")),
					},
					{
						"name":  "keyword.operator.word.moonbasic",
						"match": fmt.Sprintf(`(?i)\b(%s)\b`, strings.Join(keywords[len(keywords)-8:], "|")),
					},
				},
			},
			"operators": map[string]any{
				"patterns": []map[string]any{
					{"name": "keyword.operator.moonbasic", "match": `<>|<=|>=|\+=|-=|\*=|/=`},
					{"name": "keyword.operator.moonbasic", "match": `[=<>+\-*/^]`},
				},
			},
		},
	}

	raw, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.MkdirAll(filepath.Dir(vscodeOut), 0755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.WriteFile(vscodeOut, raw, 0644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Wrote %s (%d namespaces, %d globals)\n", vscodeOut, len(namespaces), len(globals))
}
