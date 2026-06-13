package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	stdruntime "runtime"
	"strings"
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

func findToolchainAuto() ToolchainInfo {
	if info := findLocalToolchain(); info.Found {
		return info
	}

	mbName := moonbasicName()
	mrName := moonrunName()

	if p, err := exec.LookPath(strings.TrimSuffix(mbName, ".exe")); err == nil {
		info := ToolchainInfo{Moonbasic: p, Found: true}
		if mr, err := exec.LookPath(strings.TrimSuffix(mrName, ".exe")); err == nil {
			info.Moonrun = mr
		}
		return fillMissingMoonrun(info)
	}

	searchRoots := []string{}
	if exe, err := os.Executable(); err == nil {
		searchRoots = append(searchRoots, filepath.Dir(exe))
	}
	if cwd, err := os.Getwd(); err == nil {
		searchRoots = append(searchRoots, cwd)
	}
	if a := os.Getenv("MOONBASIC_ROOT"); a != "" {
		searchRoots = append(searchRoots, a)
	}

	seen := make(map[string]struct{})
	for _, root := range searchRoots {
		if root == "" {
			continue
		}
		if _, ok := seen[root]; ok {
			continue
		}
		seen[root] = struct{}{}
		cur := root
		for depth := 0; depth < 6; depth++ {
			mb := filepath.Join(cur, mbName)
			if st, err := os.Stat(mb); err == nil && !st.IsDir() {
				info := ToolchainInfo{Moonbasic: mb, Found: true}
				mr := filepath.Join(cur, mrName)
				if st2, err2 := os.Stat(mr); err2 == nil && !st2.IsDir() {
					info.Moonrun = mr
				}
				return fillMissingMoonrun(info)
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

func (a *App) resolveToolchain() ToolchainInfo {
	mb := strings.TrimSpace(a.settings.MoonbasicPath)
	mr := strings.TrimSpace(a.settings.MoonrunPath)
	if mb != "" {
		if st, err := os.Stat(mb); err == nil && !st.IsDir() {
			info := ToolchainInfo{Moonbasic: mb, Found: true}
			if mr != "" {
				if st2, err2 := os.Stat(mr); err2 == nil && !st2.IsDir() {
					info.Moonrun = mr
				}
			}
			return fillMissingMoonrun(info)
		}
	}
	info := findToolchainAuto()
	return fillMissingMoonrun(info)
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
		return ToolchainResult{Success: false, Error: "moonbasic not found on PATH or next to IDE"}
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
		return ToolchainResult{Success: false, Error: "moonbasic not found on PATH or next to IDE"}
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

func (a *App) RunFile(filePath string) ToolchainResult {
	tc := a.resolveToolchain()
	if tc.Moonrun == "" {
		return ToolchainResult{
			Success: false,
			Error:   "moonrun not found — copy moonrun.exe to toolchain/ or install the full runtime from GitHub Releases",
		}
	}
	if filePath == "" {
		return ToolchainResult{Success: false, Error: "No file path"}
	}
	cwd := filepath.Dir(filePath)
	cmd := exec.Command(tc.Moonrun, filePath)
	cmd.Dir = cwd
	if err := cmd.Start(); err != nil {
		return ToolchainResult{Success: false, Error: err.Error()}
	}
	return ToolchainResult{
		Success: true,
		Message: fmt.Sprintf("Started %s", filepath.Base(tc.Moonrun)),
	}
}
