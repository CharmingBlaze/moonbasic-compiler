package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	stdruntime "runtime"
	"strings"
	"sync"
)

// ToolchainInfo describes discovered moonbasic/moonrun binaries.
type ToolchainInfo struct {
	Moonbasic string `json:"moonbasic"`
	Moonrun   string `json:"moonrun"`
	Found     bool   `json:"found"`
}

// ToolchainResult is the outcome of check, compile, or run.
type ToolchainResult struct {
	Success  bool   `json:"success"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	Message  string `json:"message"`
	Error    string `json:"error"`
	ExitCode int    `json:"exitCode"`
}

func moonbasicName() string {
	if stdruntime.GOOS == "windows" {
		return "moonbasic.exe"
	}
	return "moonbasic"
}

func moonrunName() string {
	if stdruntime.GOOS == "windows" {
		return "moonrun.exe"
	}
	return "moonrun"
}

func lookPathBin(name string) (string, bool) {
	// Try with and without .exe so Windows/Unix both work.
	candidates := []string{name}
	if stdruntime.GOOS == "windows" {
		if !strings.HasSuffix(strings.ToLower(name), ".exe") {
			candidates = append(candidates, name+".exe")
		}
	} else {
		candidates = append(candidates, strings.TrimSuffix(name, ".exe"))
	}
	for _, c := range candidates {
		if p, err := exec.LookPath(c); err == nil {
			return p, true
		}
	}
	return "", false
}

func toolchainAtDir(dir string) ToolchainInfo {
	mbName := moonbasicName()
	mrName := moonrunName()
	mb := filepath.Join(dir, mbName)
	if st, err := os.Stat(mb); err != nil || st.IsDir() {
		// macOS .app: binaries may sit next to the .app bundle, not inside Contents/MacOS
		return ToolchainInfo{}
	}
	info := ToolchainInfo{Moonbasic: mb, Found: true}
	mr := filepath.Join(dir, mrName)
	if st, err := os.Stat(mr); err == nil && !st.IsDir() {
		info.Moonrun = mr
	}
	return fillMissingMoonrun(info)
}

func findSiblingToolchain() ToolchainInfo {
	roots := []string{}
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		roots = append(roots, exeDir)
		// If running from Foo.app/Contents/MacOS, also check the folder containing the .app
		if stdruntime.GOOS == "darwin" {
			if strings.HasSuffix(filepath.Base(filepath.Dir(filepath.Dir(exeDir))), ".app") {
				roots = append(roots, filepath.Dir(filepath.Dir(filepath.Dir(exeDir))))
			}
		}
	}
	if cwd, err := os.Getwd(); err == nil {
		roots = append(roots, cwd)
	}
	seen := map[string]struct{}{}
	for _, root := range roots {
		if root == "" {
			continue
		}
		if _, ok := seen[root]; ok {
			continue
		}
		seen[root] = struct{}{}
		if info := toolchainAtDir(root); info.Found {
			return info
		}
		// Walk a few parents (release extract / monorepo layouts).
		cur := root
		for depth := 0; depth < 6; depth++ {
			if info := toolchainAtDir(cur); info.Found {
				return info
			}
			parent := filepath.Dir(cur)
			if parent == cur {
				break
			}
			cur = parent
		}
	}
	return ToolchainInfo{}
}

func findToolchainAuto() ToolchainInfo {
	// 1) Binaries next to the IDE (release zip/tar) — highest priority.
	if info := findSiblingToolchain(); info.Found {
		return info
	}
	// 2) Dev toolchain/ folder under the repo.
	if info := findLocalToolchain(); info.Found {
		return info
	}
	// 3) PATH
	if p, ok := lookPathBin(moonbasicName()); ok {
		info := ToolchainInfo{Moonbasic: p, Found: true}
		if mr, ok2 := lookPathBin(moonrunName()); ok2 {
			info.Moonrun = mr
		}
		return fillMissingMoonrun(info)
	}
	if a := os.Getenv("MOONBASIC_ROOT"); a != "" {
		if info := toolchainAtDir(a); info.Found {
			return info
		}
	}
	return ToolchainInfo{}
}

func settingsToolchainValid(mb, mr string) ToolchainInfo {
	mb = strings.TrimSpace(mb)
	if mb == "" {
		return ToolchainInfo{}
	}
	if st, err := os.Stat(mb); err != nil || st.IsDir() {
		return ToolchainInfo{}
	}
	info := ToolchainInfo{Moonbasic: mb, Found: true}
	mr = strings.TrimSpace(mr)
	if mr != "" {
		if st, err := os.Stat(mr); err == nil && !st.IsDir() {
			info.Moonrun = mr
		}
	}
	return fillMissingMoonrun(info)
}

func (a *App) resolveToolchain() ToolchainInfo {
	// Prefer binaries beside the IDE so moving the extract folder always works.
	if info := findSiblingToolchain(); info.Found {
		return info
	}
	if info := settingsToolchainValid(a.settings.MoonbasicPath, a.settings.MoonrunPath); info.Found {
		return info
	}
	return findToolchainAuto()
}

func runTool(exe string, args []string, cwd string) ToolchainResult {
	if exe == "" {
		return ToolchainResult{
			Success: false,
			Error:   "moonbasic toolchain not found — install from GitHub Releases or build from source",
		}
	}
	cmd := exec.Command(exe, args...)
	if cwd != "" {
		cmd.Dir = cwd
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	code := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code = exitErr.ExitCode()
		} else {
			return ToolchainResult{
				Success: false,
				Error:   err.Error(),
				Stdout:  stdout.String(),
				Stderr:  stderr.String(),
			}
		}
	}
	return ToolchainResult{
		Success:  code == 0,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: code,
	}
}

func (a *App) GetToolchain() ToolchainInfo {
	return a.resolveToolchain()
}

func (a *App) CheckFile(filePath string) ToolchainResult {
	tc := a.resolveToolchain()
	if !tc.Found {
		return ToolchainResult{Success: false, Error: "moonbasic not found — keep it in the same folder as the IDE (START-IDE), or set Settings → Toolchain"}
	}
	if filePath == "" {
		return ToolchainResult{Success: false, Error: "No file path"}
	}
	cwd := filepath.Dir(filePath)
	res := runTool(tc.Moonbasic, []string{"--check", filePath}, cwd)
	if res.Success {
		res.Message = "Check OK"
	} else if res.Stderr == "" && res.Error != "" {
		res.Stderr = res.Error
	}
	return res
}

func (a *App) CompileFile(filePath string) ToolchainResult {
	tc := a.resolveToolchain()
	if !tc.Found {
		return ToolchainResult{Success: false, Error: "moonbasic not found — keep it in the same folder as the IDE (START-IDE), or set Settings → Toolchain"}
	}
	if filePath == "" {
		return ToolchainResult{Success: false, Error: "No file path"}
	}
	cwd := filepath.Dir(filePath)
	res := runTool(tc.Moonbasic, []string{filePath}, cwd)
	if res.Success {
		mbc := strings.TrimSuffix(filePath, filepath.Ext(filePath)) + ".mbc"
		res.Message = fmt.Sprintf("Wrote %s", filepath.Base(mbc))
	} else if res.Stderr == "" && res.Error != "" {
		res.Stderr = res.Error
	}
	return res
}

var (
	runMu   sync.Mutex
	runCmd  *exec.Cmd
)

func (a *App) RunFile(filePath string) ToolchainResult {
	tc := a.resolveToolchain()
	if tc.Moonrun == "" {
		return ToolchainResult{
			Success: false,
			Error:   "moonrun not found — put " + moonrunName() + " in the same folder as the IDE (re-extract the IDE zip), or Settings → Toolchain",
		}
	}
	if filePath == "" {
		return ToolchainResult{Success: false, Error: "No file path"}
	}
	cwd := filepath.Dir(filePath)
	cmd := exec.Command(tc.Moonrun, filePath)
	cmd.Dir = cwd
	runMu.Lock()
	if runCmd != nil && runCmd.Process != nil {
		_ = runCmd.Process.Kill()
		runCmd = nil
	}
	if err := cmd.Start(); err != nil {
		runMu.Unlock()
		return ToolchainResult{Success: false, Error: err.Error()}
	}
	runCmd = cmd
	runMu.Unlock()
	go func() {
		_ = cmd.Wait()
		runMu.Lock()
		if runCmd == cmd {
			runCmd = nil
		}
		runMu.Unlock()
	}()
	return ToolchainResult{
		Success: true,
		Message: fmt.Sprintf("Started %s", filepath.Base(tc.Moonrun)),
	}
}

// StopRun kills the last moonrun process started from the IDE (if still running).
func (a *App) StopRun() ToolchainResult {
	runMu.Lock()
	defer runMu.Unlock()
	if runCmd == nil || runCmd.Process == nil {
		return ToolchainResult{Success: false, Error: "No running game from this IDE"}
	}
	if err := runCmd.Process.Kill(); err != nil {
		return ToolchainResult{Success: false, Error: err.Error()}
	}
	runCmd = nil
	return ToolchainResult{Success: true, Message: "Stopped game"}
}
