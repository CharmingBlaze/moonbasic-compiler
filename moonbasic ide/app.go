package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	stdruntime "runtime"
	"strings"
	"time"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	projectPath string
	lastFile    string
	lsp         *LSPClient
	settings    IDESettings
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.loadSettings()
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {
	// Nothing needed here — frontend handles init
}

// beforeClose is called when the application is about to quit
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	if a.lsp != nil {
		a.lsp.stop()
	}
}

// ── Project folder ─────────────────────────────────────────────

type ProjectFolderResult struct {
	Success bool   `json:"success"`
	Path    string `json:"path"`
	Error   string `json:"error"`
}

type ProjectFileEntry struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// PickProjectFolder opens a folder picker for a new project
func (a *App) PickProjectFolder() ProjectFolderResult {
	dir, err := wailsruntime.OpenDirectoryDialog(a.ctx, wailsruntime.OpenDialogOptions{
		Title: "Choose project folder",
	})
	if err != nil || dir == "" {
		return ProjectFolderResult{Success: false, Error: "No folder selected"}
	}
	a.projectPath = dir
	return ProjectFolderResult{Success: true, Path: dir}
}

// WriteProjectFiles writes starter files into the chosen project folder
func (a *App) WriteProjectFiles(folderPath string, filesJSON string) FileResult {
	if folderPath == "" {
		return FileResult{Success: false, Error: "No project folder"}
	}
	var files []ProjectFileEntry
	if err := json.Unmarshal([]byte(filesJSON), &files); err != nil {
		return FileResult{Success: false, Error: "Invalid project files: " + err.Error()}
	}
	for _, f := range files {
		if f.Name == "" {
			continue
		}
		full := filepath.Join(folderPath, filepath.FromSlash(f.Name))
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			return FileResult{Success: false, Error: err.Error()}
		}
		if err := os.WriteFile(full, []byte(f.Content), 0644); err != nil {
			return FileResult{Success: false, Error: err.Error()}
		}
	}
	a.projectPath = folderPath
	a.lastFile = filepath.Join(folderPath, "main.mb")
	return FileResult{
		Success:  true,
		Path:     folderPath,
		Filename: filepath.Base(folderPath),
		Message:  fmt.Sprintf("Created %d file(s) in %s", len(files), folderPath),
	}
}

// GetProjectFolder returns the active project directory (desktop app)
func (a *App) GetProjectFolder() string {
	return a.projectPath
}

type ReadProjectResult struct {
	Success bool               `json:"success"`
	Path    string             `json:"path"`
	Files   []ProjectFileEntry `json:"files"`
	Error   string             `json:"error"`
}

// ReadProjectFolder reads all moonBASIC source files from a project directory
func (a *App) ReadProjectFolder(dir string) ReadProjectResult {
	if dir == "" {
		return ReadProjectResult{Success: false, Error: "No folder path"}
	}
	var files []ProjectFileEntry
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			base := filepath.Base(path)
			if base == ".git" || base == "node_modules" || base == "frontend" {
				return filepath.SkipDir
			}
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".mb" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			rel = filepath.Base(path)
		}
		files = append(files, ProjectFileEntry{
			Name:    filepath.ToSlash(rel),
			Content: string(data),
		})
		return nil
	})
	if err != nil {
		return ReadProjectResult{Success: false, Error: err.Error()}
	}
	if len(files) == 0 {
		return ReadProjectResult{Success: false, Error: "No .mb source files found in folder"}
	}
	a.projectPath = dir
	a.lastFile = filepath.Join(dir, "main.sb")
	return ReadProjectResult{Success: true, Path: dir, Files: files}
}

// ── File Operations ────────────────────────────────────────────

// OpenFile opens a file dialog and returns the file contents
func (a *App) OpenFile() FileResult {
	selection, err := wailsruntime.OpenFileDialog(a.ctx, wailsruntime.OpenDialogOptions{
		Title: "Open moonBASIC File",
		Filters: []wailsruntime.FileFilter{
			{DisplayName: "moonBASIC Files (*.mb)", Pattern: "*.mb"},
			{DisplayName: "Text Files (*.txt)", Pattern: "*.txt"},
			{DisplayName: "All Files (*.*)", Pattern: "*.*"},
		},
	})
	if err != nil || selection == "" {
		return FileResult{Success: false, Error: "No file selected"}
	}

	data, err := os.ReadFile(selection)
	if err != nil {
		return FileResult{Success: false, Error: err.Error()}
	}

	a.lastFile = selection
	return FileResult{
		Success:  true,
		Path:     selection,
		Filename: filepath.Base(selection),
		Content:  string(data),
	}
}

// SaveFile saves content to a file, using last path or dialog
func (a *App) SaveFile(content string, filename string) FileResult {
	savePath := a.lastFile
	if savePath == "" {
		var err error
		savePath, err = wailsruntime.SaveFileDialog(a.ctx, wailsruntime.SaveDialogOptions{
			Title:           "Save moonBASIC File",
			DefaultFilename: filename,
			Filters: []wailsruntime.FileFilter{
				{DisplayName: "moonBASIC Files (*.mb)", Pattern: "*.mb"},
				{DisplayName: "All Files (*.*)", Pattern: "*.*"},
			},
		})
		if err != nil || savePath == "" {
			return FileResult{Success: false, Error: "Save cancelled"}
		}
	}

	err := os.WriteFile(savePath, []byte(content), 0644)
	if err != nil {
		return FileResult{Success: false, Error: err.Error()}
	}

	a.lastFile = savePath
	return FileResult{
		Success:  true,
		Path:     savePath,
		Filename: filepath.Base(savePath),
	}
}

// SaveFileAs always shows a save dialog
func (a *App) SaveFileAs(content string, defaultName string) FileResult {
	savePath, err := wailsruntime.SaveFileDialog(a.ctx, wailsruntime.SaveDialogOptions{
		Title:           "Save As...",
		DefaultFilename: defaultName,
		Filters: []wailsruntime.FileFilter{
			{DisplayName: "moonBASIC Files (*.mb)", Pattern: "*.mb"},
			{DisplayName: "All Files (*.*)", Pattern: "*.*"},
		},
	})
	if err != nil || savePath == "" {
		return FileResult{Success: false, Error: "Save cancelled"}
	}

	err = os.WriteFile(savePath, []byte(content), 0644)
	if err != nil {
		return FileResult{Success: false, Error: err.Error()}
	}

	if ext := filepath.Ext(savePath); ext == ".mb" {
		a.lastFile = savePath
	}

	return FileResult{
		Success:  true,
		Path:     savePath,
		Filename: filepath.Base(savePath),
	}
}


// ── Window Controls ────────────────────────────────────────────

func (a *App) SetFullscreen(fullscreen bool) {
	if fullscreen {
		wailsruntime.WindowFullscreen(a.ctx)
	} else {
		wailsruntime.WindowUnfullscreen(a.ctx)
	}
}

func (a *App) Minimize() {
	wailsruntime.WindowMinimise(a.ctx)
}

func (a *App) SetTitle(title string) {
	wailsruntime.WindowSetTitle(a.ctx, title)
}

// ── System Info ────────────────────────────────────────────────

func (a *App) GetSystemInfo() SystemInfo {
	return SystemInfo{
		OS:      stdruntime.GOOS,
		Arch:    stdruntime.GOARCH,
		Version: "1.0.0",
		BuildAt: time.Now().Format("2006-01-02"),
	}
}

// ── Types ──────────────────────────────────────────────────────

type FileResult struct {
	Success  bool   `json:"success"`
	Path     string `json:"path"`
	Filename string `json:"filename"`
	Content  string `json:"content"`
	Error    string `json:"error"`
	Message  string `json:"message"`
}

type SystemInfo struct {
	OS      string `json:"os"`
	Arch    string `json:"arch"`
	Version string `json:"version"`
	BuildAt string `json:"buildAt"`
}

// Ensure fmt is used
var _ = fmt.Sprintf
