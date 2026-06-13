// Ideexport writes language data for the moonBASIC IDE from the builtin manifest.
// Run from the module root: go run ./tools/ideexport
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

type langData struct {
	Keywords       []string                    `json:"keywords"`
	Globals        []string                    `json:"globals"`
	Namespaces     []string                    `json:"namespaces"`
	NamespaceIndex map[string][]commandEntry   `json:"namespaceIndex"`
	Commands       map[string]commandEntry     `json:"commands"`
	HandleMethods  []string                    `json:"handleMethods"`
}

type commandEntry struct {
	Key         string   `json:"key"`
	Description string   `json:"description,omitempty"`
	Args        []string `json:"args,omitempty"`
	Returns     string   `json:"returns,omitempty"`
	Example     string   `json:"example,omitempty"`
	Phase       string   `json:"phase,omitempty"`
	Stub        string   `json:"stub,omitempty"`
}

var keywords = []string{
	"IF", "THEN", "ELSE", "ENDIF", "WHILE", "WEND", "ENDWHILE",
	"FOR", "TO", "DOWNTO", "STEP", "NEXT", "REPEAT", "UNTIL",
	"SELECT", "CASE", "DEFAULT", "ENDSELECT",
	"FUNCTION", "ENDFUNCTION", "RETURN", "TYPE", "FIELD", "ENDTYPE",
	"NEW", "DELETE", "EACH", "GOTO", "GOSUB",
	"DIM", "REDIM", "PRESERVE", "LOCAL", "GLOBAL", "CONST", "STATIC",
	"SWAP", "ERASE", "INCLUDE", "IMPORT", "END",
	"AND", "OR", "NOT", "XOR", "MOD", "TRUE", "FALSE", "NULL",
	"PRINT", "PRINTLN", "WRITE", "HELP", "INPUT", "CLS",
}

var handleMethods = []string{
	"pos", "rot", "scale", "col", "alpha", "free", "move", "turn",
	"look", "fov", "visible", "name", "begin", "end", "target",
	"hide", "show", "draw", "update", "orbit", "yaw", "zoom",
}

func argName(k builtinmanifest.ArgKind) string {
	switch k {
	case builtinmanifest.Int:
		return "int"
	case builtinmanifest.Float:
		return "float"
	case builtinmanifest.String:
		return "string"
	case builtinmanifest.Bool:
		return "bool"
	case builtinmanifest.Handle:
		return "handle"
	default:
		return "any"
	}
}

func main() {
	outPath := filepath.Join("moonbasic ide", "js", "studio", "lang-data.json")
	vscodeOut := filepath.Join("editors", "vscode-moonbasic", "resources", "lang-data.json")
	if len(os.Args) > 1 {
		outPath = os.Args[1]
	}

	table := builtinmanifest.Default()
	data := langData{
		Keywords:       keywords,
		NamespaceIndex: make(map[string][]commandEntry),
		Commands:       make(map[string]commandEntry),
		HandleMethods:  handleMethods,
	}

	nsSet := make(map[string]struct{})
	globalSet := make(map[string]struct{})

	for _, key := range table.Keys() {
		if table.IsDeprecatedAlias(key) {
			continue
		}
		cmd, ok := table.FirstOverload(key)
		if !ok {
			continue
		}
		entry := commandEntry{
			Key:         key,
			Description: cmd.Desc,
			Returns:     cmd.Returns,
			Phase:       cmd.Phase,
			Stub:        cmd.Stub,
		}
		for _, a := range cmd.Args {
			entry.Args = append(entry.Args, argName(a))
		}

		data.Commands[key] = entry

		parts := strings.SplitN(key, ".", 2)
		if len(parts) == 1 {
			globalSet[strings.ToUpper(parts[0])] = struct{}{}
			continue
		}
		ns := parts[0]
		nsSet[ns] = struct{}{}
		method := parts[1]
		data.NamespaceIndex[ns] = append(data.NamespaceIndex[ns], commandEntry{
			Key:         method,
			Description: cmd.Desc,
			Args:        entry.Args,
			Returns:     cmd.Returns,
			Phase:       cmd.Phase,
			Stub:        cmd.Stub,
		})
	}

	for ns := range nsSet {
		data.Namespaces = append(data.Namespaces, ns)
	}
	sort.Strings(data.Namespaces)

	for g := range globalSet {
		data.Globals = append(data.Globals, g)
	}
	sort.Strings(data.Globals)

	for ns := range data.NamespaceIndex {
		sort.Slice(data.NamespaceIndex[ns], func(i, j int) bool {
			return data.NamespaceIndex[ns][i].Key < data.NamespaceIndex[ns][j].Key
		})
	}

	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(outPath, raw, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", outPath, err)
		os.Exit(1)
	}
	if err := os.MkdirAll(filepath.Dir(vscodeOut), 0755); err == nil {
		_ = os.WriteFile(vscodeOut, raw, 0644)
	}
	fmt.Printf("Wrote %s (%d commands, %d namespaces)\n", outPath, len(data.Commands), len(data.Namespaces))
}
