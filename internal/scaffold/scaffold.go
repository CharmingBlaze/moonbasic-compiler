// Package scaffold creates new moonBASIC game project folders from embedded templates.
package scaffold

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/new_project/main.mb
//go:embed templates/new_project/.vscode/launch.json
var templates embed.FS

const mainFallback = `WINDOW.OPEN(1280, 720, "My Game")
WINDOW.SETFPS(60)

WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR WINDOW.SHOULDCLOSE())
    RENDER.CLEAR(20, 24, 32)
    DRAW.TEXT("Edit main.mb to begin", 480, 340, 24, 220, 220, 230, 255)
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
`

// Create writes a new project directory named name in the current working directory.
func Create(name string) (absDir string, err error) {
	name = strings.TrimSpace(name)
	if name == "" || strings.ContainsAny(name, `/\`) {
		return "", fmt.Errorf("project name must be a simple folder name")
	}
	dir, err := filepath.Abs(name)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(dir); err == nil {
		return "", fmt.Errorf("directory already exists: %s", dir)
	}
	if err := os.MkdirAll(filepath.Join(dir, "assets"), 0o755); err != nil {
		return "", err
	}
	mainBody := mainFallback
	if b, err := templates.ReadFile("templates/new_project/main.mb"); err == nil {
		mainBody = string(b)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.mb"), []byte(mainBody), 0o644); err != nil {
		return "", err
	}
	if b, err := templates.ReadFile("templates/new_project/.vscode/launch.json"); err == nil {
		launchDir := filepath.Join(dir, ".vscode")
		if err := os.MkdirAll(launchDir, 0o755); err != nil {
			return "", err
		}
		if err := os.WriteFile(filepath.Join(launchDir, "launch.json"), b, 0o644); err != nil {
			return "", err
		}
	}
	readme := fmt.Sprintf("# %s\n\nRun with moonrun:\n\n```\nmoonrun main.mb\n```\n", name)
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(readme), 0o644); err != nil {
		return "", err
	}
	return dir, nil
}
