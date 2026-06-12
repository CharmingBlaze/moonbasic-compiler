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
//go:embed templates/new_project/.vscode/tasks.json
//go:embed templates/new_project/.vscode/extensions.json
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
// template selects embedded starter content: empty, 3d, platformer, ui (default empty).
func Create(name string, template string) (absDir string, err error) {
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
	tmpl := strings.ToLower(strings.TrimSpace(template))
	mainBody := templateMain(template)
	if mainBody == "" {
		mainBody = mainFallback
	}
	if tmpl == "" || tmpl == "empty" {
		if b, err := templates.ReadFile("templates/new_project/main.mb"); err == nil {
			mainBody = string(b)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "main.mb"), []byte(mainBody), 0o644); err != nil {
		return "", err
	}
	vscodeDir := filepath.Join(dir, ".vscode")
	if err := os.MkdirAll(vscodeDir, 0o755); err != nil {
		return "", err
	}
	for _, name := range []string{"launch.json", "tasks.json", "extensions.json"} {
		embedPath := "templates/new_project/.vscode/" + name
		if b, err := templates.ReadFile(embedPath); err == nil {
			if err := os.WriteFile(filepath.Join(vscodeDir, name), b, 0o644); err != nil {
				return "", err
			}
		}
	}
	readme := fmt.Sprintf("# %s\n\nRun with moonrun:\n\n```\nmoonrun main.mb\n```\n", name)
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(readme), 0o644); err != nil {
		return "", err
	}
	return dir, nil
}

func templateMain(template string) string {
	switch strings.ToLower(strings.TrimSpace(template)) {
	case "3d", "first-person", "third-person":
		return `APP.OPEN(1280, 720, "My 3D Game")
APP.SETFPS(60)

cam = CAMERA.CREATE("Main Camera")
CAMERA.SETACTIVE(cam)
ENTITY.SETPOSITION(cam, 0, 2, -8)
CAMERA.LOOKAT(cam, 0, 0, 0)

cube = ENTITY.CREATECUBE("Cube", 2)
ENTITY.SETPOSITION(cube, 0, 1, 5)

WHILE NOT APP.SHOULDCLOSE()
    ENTITY.TURN(cube, 0, 45 * APP.DELTA(), 0)
    RENDER.CLEAR(20, 24, 32)
    RENDER.BEGIN(cam)
    SCENE.DRAW()
    RENDER.END()
    RENDER.FRAME()
WEND

APP.CLOSE()
`
	case "platformer", "2d", "top-down":
		return `APP.OPEN(960, 540, "My Platformer")
APP.SETFPS(60)

player = ENTITY.CREATE("Player")
ENTITY.SETPOSITION(player, 100, 200, 0)

WHILE NOT APP.SHOULDCLOSE()
    IF INPUT.KEYDOWN(KEY_LEFT) THEN ENTITY.MOVE(player, -200 * APP.DELTA(), 0, 0)
    IF INPUT.KEYDOWN(KEY_RIGHT) THEN ENTITY.MOVE(player, 200 * APP.DELTA(), 0, 0)
    RENDER.CLEAR(30, 40, 55)
    DRAW.TEXT("Platformer template — edit main.mb", 20, 20, 20, 220, 220, 230, 255)
    RENDER.FRAME()
WEND

APP.CLOSE()
`
	case "ui", "ui-menu", "menu":
		return `APP.OPEN(800, 600, "My Menu")
APP.SETFPS(60)

WHILE NOT APP.SHOULDCLOSE()
    RENDER.CLEAR(18, 20, 28)
    DRAW.TEXT("UI Menu template", 280, 40, 28, 255, 255, 255, 255)
    IF GUI.BUTTON(300, 200, 200, 40, "Start") THEN PRINT "Start clicked"
    RENDER.FRAME()
WEND

APP.CLOSE()
`
	case "physics":
		return `APP.OPEN(960, 540, "Physics Playground")
APP.SETFPS(60)
PHYSICS.START()
PHYSICS.SETGRAVITY(0, -9.8, 0)

floor = ENTITY.CREATECUBE("Floor", 10)
ENTITY.SETPOSITION(floor, 0, -1, 0)
BODY.ADDSTATICBOX(floor, 10, 1, 10)

crate = ENTITY.CREATECUBE("Crate", 1)
ENTITY.SETPOSITION(crate, 0, 5, 0)
BODY.ADDDYNAMICBOX(crate, 1, 1, 1)

WHILE NOT APP.SHOULDCLOSE()
    PHYSICS.STEP()
    RENDER.CLEAR(25, 28, 35)
    RENDER.BEGIN()
    SCENE.DRAW()
    RENDER.END()
    RENDER.FRAME()
WEND

APP.CLOSE()
`
	default:
		return ""
	}
}
