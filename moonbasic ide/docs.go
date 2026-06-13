package main

import (
	"os"
	"path/filepath"
	"strings"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// DocResult is bundled documentation content for the IDE.
type DocResult struct {
	Success bool   `json:"success"`
	Path    string `json:"path"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Error   string `json:"error"`
}

func bundledDocRoots() []string {
	seen := make(map[string]struct{})
	add := func(dir string) {
		dir = filepath.Clean(dir)
		if dir == "" || dir == "." {
			return
		}
		if _, ok := seen[dir]; ok {
			return
		}
		if st, err := os.Stat(dir); err != nil || !st.IsDir() {
			return
		}
		seen[dir] = struct{}{}
	}

	if cwd, err := os.Getwd(); err == nil {
		add(filepath.Join(cwd, "bundled-docs"))
		add(filepath.Join(cwd, "frontend", "bundled-docs"))
		add(filepath.Join(filepath.Dir(cwd), "moonbasic ide", "bundled-docs"))
		add(filepath.Join(filepath.Dir(cwd), "moonbasic ide", "frontend", "bundled-docs"))
	}
	if exe, err := os.Executable(); err == nil {
		cur := filepath.Dir(exe)
		for depth := 0; depth < 7; depth++ {
			add(filepath.Join(cur, "bundled-docs"))
			add(filepath.Join(cur, "frontend", "bundled-docs"))
			parent := filepath.Dir(cur)
			if parent == cur {
				break
			}
			cur = parent
		}
	}
	if root := os.Getenv("MOONBASIC_IDE_ROOT"); root != "" {
		add(filepath.Join(root, "bundled-docs"))
		add(filepath.Join(root, "frontend", "bundled-docs"))
	}
	add("bundled-docs")
	add(filepath.Join("frontend", "bundled-docs"))

	out := make([]string, 0, len(seen))
	for dir := range seen {
		out = append(out, dir)
	}
	return out
}

func bundledDocPaths(rel string) []string {
	rel = filepath.FromSlash(strings.TrimPrefix(rel, "/"))
	var paths []string
	for _, root := range bundledDocRoots() {
		paths = append(paths, filepath.Join(root, rel))
	}
	return paths
}

// ReadBundledDoc returns a markdown file shipped with the IDE.
func (a *App) ReadBundledDoc(relPath string) DocResult {
	relPath = strings.ReplaceAll(relPath, "\\", "/")
	if relPath == "" {
		return DocResult{Success: false, Error: "No document path"}
	}
	for _, p := range bundledDocPaths(relPath) {
		data, err := os.ReadFile(p)
		if err == nil {
			return DocResult{
				Success: true,
				Path:    relPath,
				Title:   filepath.Base(relPath),
				Content: string(data),
			}
		}
	}
	return DocResult{Success: false, Error: "Document not found: " + relPath}
}

// OpenFolder opens a project folder and returns all .mb files.
func (a *App) OpenFolder() ReadProjectResult {
	dir, err := wailsruntime.OpenDirectoryDialog(a.ctx, wailsruntime.OpenDialogOptions{
		Title: "Open Project Folder",
	})
	if err != nil || dir == "" {
		return ReadProjectResult{Success: false, Error: "No folder selected"}
	}
	return a.ReadProjectFolder(dir)
}
