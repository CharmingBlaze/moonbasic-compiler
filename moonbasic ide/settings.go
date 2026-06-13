package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// IDESettings stores user-configured compiler paths.
type IDESettings struct {
	MoonbasicPath string `json:"moonbasicPath"`
	MoonrunPath   string `json:"moonrunPath"`
}

func settingsFilePath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "moonbasic-ide", "settings.json")
}

func loadIDESettings() IDESettings {
	data, err := os.ReadFile(settingsFilePath())
	if err != nil {
		return IDESettings{}
	}
	var s IDESettings
	_ = json.Unmarshal(data, &s)
	return s
}

func saveIDESettings(s IDESettings) error {
	p := settingsFilePath()
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, raw, 0644)
}

func (a *App) loadSettings() {
	a.settings = loadIDESettings()
	if strings.TrimSpace(a.settings.MoonbasicPath) != "" {
		return
	}
	if info := findLocalToolchain(); info.Found {
		a.settings.MoonbasicPath = info.Moonbasic
		a.settings.MoonrunPath = info.Moonrun
	}
}

// GetIDESettings returns current toolchain path configuration.
func (a *App) GetIDESettings() IDESettings {
	return a.settings
}

// SetIDESettings saves compiler paths and restarts LSP if moonbasic changed.
func (a *App) SetIDESettings(s IDESettings) IDESettings {
	s.MoonbasicPath = strings.TrimSpace(s.MoonbasicPath)
	s.MoonrunPath = strings.TrimSpace(s.MoonrunPath)
	prev := a.settings.MoonbasicPath
	a.settings = s
	_ = saveIDESettings(s)
	if prev != s.MoonbasicPath && a.lsp != nil {
		a.lsp.stop()
		a.lsp = nil
	}
	return a.settings
}

// BrowseToolchain opens a file picker for moonbasic or moonrun.
func (a *App) BrowseToolchain(which string) string {
	title := "Choose moonbasic executable"
	if which == "moonrun" {
		title = "Choose moonrun executable"
	}
	path, err := wailsruntime.OpenFileDialog(a.ctx, wailsruntime.OpenDialogOptions{
		Title: title,
		Filters: []wailsruntime.FileFilter{
			{DisplayName: "Executables (*.exe)", Pattern: "*.exe"},
			{DisplayName: "All files", Pattern: "*.*"},
		},
	})
	if err != nil || path == "" {
		return ""
	}
	if which == "moonrun" {
		a.settings.MoonrunPath = path
	} else {
		a.settings.MoonbasicPath = path
	}
	_ = saveIDESettings(a.settings)
	return path
}

// TestToolchain verifies configured or auto-discovered binaries.
func (a *App) TestToolchain() ToolchainInfo {
	return a.resolveToolchain()
}
