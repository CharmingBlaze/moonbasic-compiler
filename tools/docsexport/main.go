// Docsexport bundles the full moonBASIC docs tree for the IDE and writes docs-index.json.
// Run from module root: go run ./tools/docsexport
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type docEntry struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Path     string `json:"path"`
	Category string `json:"category"`
}

// Preferred sidebar entries (shown first with friendly titles).
var featured = []docEntry{
	{ID: "begin", Title: "Begin Here", Path: "BEGIN_HERE.md", Category: "Start"},
	{ID: "getting-started", Title: "Getting Started", Path: "GETTING_STARTED.md", Category: "Start"},
	{ID: "first-hour", Title: "First Hour", Path: "FIRST_HOUR.md", Category: "Start"},
	{ID: "language", Title: "Language Reference", Path: "LANGUAGE.md", Category: "Language"},
	{ID: "programming", Title: "Programming Guide", Path: "PROGRAMMING.md", Category: "Language"},
	{ID: "commands", Title: "Command Index", Path: "COMMANDS.md", Category: "Reference"},
	{ID: "style", Title: "Style Guide", Path: "STYLE_GUIDE.md", Category: "Language"},
	{ID: "easy-mode", Title: "Easy Mode", Path: "EASY_MODE.md", Category: "Language"},
	{ID: "guides", Title: "Guides Index", Path: "systems/GUIDES.md", Category: "Guides"},
	{ID: "building", Title: "Building from Source", Path: "BUILDING.md", Category: "Developers"},
	{ID: "developer", Title: "Developer Guide", Path: "DEVELOPER.md", Category: "Developers"},
}

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	srcDocs := filepath.Join(root, "docs")
	dest := filepath.Join(root, "moonbasic ide", "bundled-docs")
	indexOut := filepath.Join(root, "moonbasic ide", "js", "studio", "docs-index.json")

	if err := os.RemoveAll(dest); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not clean dest: %v\n", err)
	}
	if err := copyTree(srcDocs, dest); err != nil {
		fmt.Fprintf(os.Stderr, "copy docs: %v\n", err)
		os.Exit(1)
	}

	entries := buildIndex(dest)
	raw, _ := json.MarshalIndent(entries, "", "  ")
	_ = os.MkdirAll(filepath.Dir(indexOut), 0755)
	if err := os.WriteFile(indexOut, raw, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write index: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Bundled %d docs → %s\n", len(entries), dest)
}

func copyTree(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			return nil
		}
		return copyFile(path, target)
	})
}

func buildIndex(dest string) []docEntry {
	seen := map[string]bool{}
	var entries []docEntry

	add := func(e docEntry) {
		key := strings.ReplaceAll(e.Path, "\\", "/")
		if seen[key] {
			return
		}
		seen[key] = true
		entries = append(entries, e)
	}

	for _, e := range featured {
		if fileExists(filepath.Join(dest, filepath.FromSlash(e.Path))) {
			add(e)
		}
	}

	var rest []docEntry
	_ = filepath.Walk(dest, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			return err
		}
		rel, _ := filepath.Rel(dest, path)
		rel = strings.ReplaceAll(rel, "\\", "/")
		if seen[rel] {
			return nil
		}
		rest = append(rest, docEntry{
			ID:       slugID(rel),
			Title:    titleFromPath(rel),
			Path:     rel,
			Category: categoryFromPath(rel),
		})
		return nil
	})

	sort.Slice(rest, func(i, j int) bool {
		if rest[i].Category != rest[j].Category {
			return rest[i].Category < rest[j].Category
		}
		return rest[i].Title < rest[j].Title
	})
	for _, e := range rest {
		add(e)
	}
	return entries
}

func fileExists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}

func slugID(rel string) string {
	s := strings.TrimSuffix(rel, filepath.Ext(rel))
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return strings.ToLower(s)
}

func titleFromPath(rel string) string {
	base := strings.TrimSuffix(filepath.Base(rel), ".md")
	base = strings.ReplaceAll(base, "-", " ")
	base = strings.ReplaceAll(base, "_", " ")
	if len(base) >= 2 && base[1] == '-' && base[0] >= '0' && base[0] <= '9' {
		return strings.ToUpper(base[:1]) + base[1:]
	}
	return base
}

func categoryFromPath(rel string) string {
	parts := strings.Split(rel, "/")
	if len(parts) < 2 {
		return "General"
	}
	switch parts[0] {
	case "systems":
		if len(parts) >= 3 && parts[1] == "guides" {
			return "Guides"
		}
		return "Systems"
	case "reference":
		if len(parts) >= 3 && parts[1] == "moonbasic-command-set" {
			return "Command Set"
		}
		if parts[1] == "dbpro" {
			return "DBPro Reference"
		}
		return "API Reference"
	case "tutorials":
		return "Tutorials"
	case "architecture":
		return "Architecture"
	case "audit":
		return "Audit"
	default:
		return capitalize(parts[0])
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
