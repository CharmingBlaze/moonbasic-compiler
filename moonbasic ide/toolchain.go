package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// toolchainPathsFile is optional JSON in moonbasic ide/toolchain/ for local testing.
type toolchainPathsFile struct {
	MoonbasicPath string `json:"moonbasicPath"`
	MoonrunPath   string `json:"moonrunPath"`
}

func localToolchainDirs() []string {
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
		add(filepath.Join(cwd, "toolchain"))
		add(filepath.Join(filepath.Dir(cwd), "moonbasic ide", "toolchain"))
	}
	if exe, err := os.Executable(); err == nil {
		cur := filepath.Dir(exe)
		for depth := 0; depth < 5; depth++ {
			add(filepath.Join(cur, "toolchain"))
			parent := filepath.Dir(cur)
			if parent == cur {
				break
			}
			cur = parent
		}
	}
	if root := os.Getenv("MOONBASIC_IDE_ROOT"); root != "" {
		add(filepath.Join(root, "toolchain"))
	}

	out := make([]string, 0, len(seen))
	for dir := range seen {
		out = append(out, dir)
	}
	return out
}

func resolveToolchainPath(baseDir, p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return ""
	}
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(baseDir, p)
}

func statExe(p string) bool {
	if p == "" {
		return false
	}
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}

func loadToolchainFromDir(dir string) ToolchainInfo {
	for _, name := range []string{"paths.local.json", "paths.json"} {
		raw, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			continue
		}
		var cfg toolchainPathsFile
		if json.Unmarshal(raw, &cfg) != nil {
			continue
		}
		mb := resolveToolchainPath(dir, cfg.MoonbasicPath)
		mr := resolveToolchainPath(dir, cfg.MoonrunPath)
		if st, err := os.Stat(mb); err == nil && !st.IsDir() {
			info := ToolchainInfo{Moonbasic: mb, Found: true}
			if statExe(mr) {
				info.Moonrun = mr
			}
			return fillMissingMoonrun(info)
		}
	}

	mbName := moonbasicName()
	mrName := moonrunName()
	mb := filepath.Join(dir, mbName)
	if !statExe(mb) {
		return ToolchainInfo{}
	}
	info := ToolchainInfo{Moonbasic: mb, Found: true}
	mr := filepath.Join(dir, mrName)
	if statExe(mr) {
		info.Moonrun = mr
	}
	return fillMissingMoonrun(info)
}

func findLocalToolchain() ToolchainInfo {
	for _, dir := range localToolchainDirs() {
		if info := loadToolchainFromDir(dir); info.Found {
			return fillMissingMoonrun(info)
		}
	}
	return ToolchainInfo{}
}

func fillMissingMoonrun(info ToolchainInfo) ToolchainInfo {
	if !info.Found || info.Moonrun != "" {
		return info
	}
	mrName := moonrunName()

	if mb := strings.TrimSpace(info.Moonbasic); mb != "" {
		for _, dir := range []string{filepath.Dir(mb), filepath.Dir(filepath.Dir(mb))} {
			candidate := filepath.Join(dir, mrName)
			if statExe(candidate) {
				info.Moonrun = candidate
				return info
			}
		}
	}

	if p, err := exec.LookPath(strings.TrimSuffix(mrName, ".exe")); err == nil && statExe(p) {
		info.Moonrun = p
		return info
	}

	searchRoots := []string{}
	if exe, err := os.Executable(); err == nil {
		searchRoots = append(searchRoots, filepath.Dir(exe))
	}
	if cwd, err := os.Getwd(); err == nil {
		searchRoots = append(searchRoots, cwd)
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
			candidate := filepath.Join(cur, mrName)
			if statExe(candidate) {
				info.Moonrun = candidate
				return info
			}
			parent := filepath.Dir(cur)
			if parent == cur {
				break
			}
			cur = parent
		}
	}
	return info
}

// GetLocalToolchainDir returns the first toolchain folder that exists (for UI hints).
func (a *App) GetLocalToolchainDir() string {
	for _, dir := range localToolchainDirs() {
		return dir
	}
	return ""
}
