package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	stdruntime "runtime"
	"strings"
)

// InstallLayout describes where the IDE lives and what sits beside it.
type InstallLayout struct {
	InstallDir   string   `json:"installDir"`
	Moonbasic    string   `json:"moonbasic"`
	Moonrun      string   `json:"moonrun"`
	HasMoonbasic bool     `json:"hasMoonbasic"`
	HasMoonrun   bool     `json:"hasMoonrun"`
	SamplesDir   string   `json:"samplesDir"`
	HasSamples   bool     `json:"hasSamples"`
	OS           string   `json:"os"`
	Hint         string   `json:"hint"`
	Status       string   `json:"status"` // ready | compiler_only | missing
}

func installDir() string {
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		if stdruntime.GOOS == "darwin" {
			// Foo.app/Contents/MacOS → folder containing the .app
			if strings.HasSuffix(filepath.Base(filepath.Dir(filepath.Dir(dir))), ".app") {
				return filepath.Dir(filepath.Dir(filepath.Dir(dir)))
			}
		}
		return dir
	}
	if cwd, err := os.Getwd(); err == nil {
		return cwd
	}
	return ""
}

func (a *App) GetInstallDir() string {
	return installDir()
}

// GetInstallLayout returns sidecar discovery for the status bar / first-run banner.
func (a *App) GetInstallLayout() InstallLayout {
	dir := installDir()
	tc := a.resolveToolchain()
	samples := filepath.Join(dir, "samples")
	st, err := os.Stat(samples)
	hasSamples := err == nil && st.IsDir()

	layout := InstallLayout{
		InstallDir:   dir,
		Moonbasic:    tc.Moonbasic,
		Moonrun:      tc.Moonrun,
		HasMoonbasic: tc.Found && strings.TrimSpace(tc.Moonbasic) != "",
		HasMoonrun:   strings.TrimSpace(tc.Moonrun) != "",
		SamplesDir:   samples,
		HasSamples:   hasSamples,
		OS:           stdruntime.GOOS,
	}

	switch {
	case layout.HasMoonbasic && layout.HasMoonrun:
		layout.Status = "ready"
		layout.Hint = "Ready — open a .mb file (or samples/) and press F5 to run."
	case layout.HasMoonbasic:
		layout.Status = "compiler_only"
		layout.Hint = "moonbasic found, but moonrun is missing — put " + moonrunName() +
			" in the same folder as the IDE for F5 Run (or set Settings → Toolchain)."
	default:
		layout.Status = "missing"
		layout.Hint = "Put " + moonbasicName() + " and " + moonrunName() +
			" next to the IDE (same folder as START-IDE), then restart — or File → Settings → Toolchain."
	}
	return layout
}

// SampleFileRef is a short sample listing entry.
type SampleFileRef struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// ListSampleFiles returns .mb files under installDir/samples (release bundles).
func (a *App) ListSampleFiles() []SampleFileRef {
	dir := filepath.Join(installDir(), "samples")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []SampleFileRef
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.EqualFold(filepath.Ext(name), ".mb") {
			continue
		}
		out = append(out, SampleFileRef{
			Name: name,
			Path: filepath.Join(dir, name),
		})
	}
	return out
}

// OpenInFileManager opens a file or folder in Explorer / Finder / xdg-open.
func (a *App) OpenInFileManager(path string) ToolchainResult {
	path = strings.TrimSpace(path)
	if path == "" {
		path = installDir()
	}
	if st, err := os.Stat(path); err != nil {
		return ToolchainResult{Success: false, Error: err.Error()}
	} else if !st.IsDir() {
		path = filepath.Dir(path)
	}
	var cmd *exec.Cmd
	switch stdruntime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	if err := cmd.Start(); err != nil {
		return ToolchainResult{Success: false, Error: err.Error()}
	}
	return ToolchainResult{Success: true, Message: fmt.Sprintf("Opened %s", path)}
}

// OpenSamplesFolder opens installDir/samples when present, else the install folder.
func (a *App) OpenSamplesFolder() ToolchainResult {
	samples := filepath.Join(installDir(), "samples")
	if st, err := os.Stat(samples); err == nil && st.IsDir() {
		return a.OpenInFileManager(samples)
	}
	return a.OpenInFileManager(installDir())
}
