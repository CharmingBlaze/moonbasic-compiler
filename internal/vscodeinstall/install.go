// Package vscodeinstall installs the moonBASIC VS Code / Cursor extension from a local or downloaded VSIX.
package vscodeinstall

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	githubOwner = "CharmingBlaze"
	githubRepo  = "moonbasic-compiler"
)

// Options controls extension installation.
type Options struct {
	// VsixPath overrides automatic VSIX discovery (local file).
	VsixPath string
	// EditorCLI forces a specific editor command (e.g. code, cursor).
	EditorCLI string
	// MoonbasicPath sets moonbasic.languageServerPath in user settings when non-empty.
	MoonbasicPath string
	// MoonrunPath sets moonbasic.moonrunPath in user settings when non-empty.
	MoonrunPath string
	// SkipSettings skips writing VS Code user settings.json entries.
	SkipSettings bool
}

// Install finds a VSIX, installs it via the editor CLI, and optionally configures tool paths.
func Install(opts Options) error {
	moonbasic, moonrun := opts.MoonbasicPath, opts.MoonrunPath
	if moonbasic == "" || moonrun == "" {
		mb, mr := siblingBinaries()
		if moonbasic == "" {
			moonbasic = mb
		}
		if moonrun == "" {
			moonrun = mr
		}
	}

	vsix, err := resolveVsix(opts.VsixPath)
	if err != nil {
		return err
	}

	editor, err := resolveEditorCLI(opts.EditorCLI)
	if err != nil {
		return err
	}

	fmt.Printf("Installing %s into %s …\n", filepath.Base(vsix), editor.name)
	if out, err := run(editor.cmd, editor.argsPrefix, "--install-extension", vsix, "--force"); err != nil {
		if strings.TrimSpace(out) != "" {
			fmt.Fprintln(os.Stderr, out)
		}
		return fmt.Errorf("install extension: %w", err)
	}
	fmt.Println("VS Code extension installed.")

	if !opts.SkipSettings && (moonbasic != "" || moonrun != "") {
		if err := writeUserSettings(editor.settingsPath, moonbasic, moonrun); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not update editor settings: %v\n", err)
			fmt.Fprintf(os.Stderr, "Set moonbasic.languageServerPath and moonbasic.moonrunPath manually if needed.\n")
		} else {
			fmt.Println("Configured moonbasic.languageServerPath and moonbasic.moonrunPath in user settings.")
		}
	}

	fmt.Println("Done. Open VS Code or Cursor, open a .mb file, and start coding.")
	return nil
}

type editorInfo struct {
	name         string
	cmd          string
	argsPrefix   []string
	settingsPath string
}

func siblingBinaries() (moonbasic, moonrun string) {
	exe, err := os.Executable()
	if err != nil {
		return "", ""
	}
	dir := filepath.Dir(exe)
	mb := filepath.Join(dir, executableName("moonbasic"))
	mr := filepath.Join(dir, executableName("moonrun"))
	if fileExists(mb) {
		moonbasic = mb
	}
	if fileExists(mr) {
		moonrun = mr
	}
	return moonbasic, moonrun
}

func executableName(base string) string {
	if runtime.GOOS == "windows" {
		return base + ".exe"
	}
	return base
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func resolveVsix(explicit string) (string, error) {
	if explicit != "" {
		if !fileExists(explicit) {
			return "", fmt.Errorf("VSIX not found: %s", explicit)
		}
		return explicit, nil
	}
	if env := strings.TrimSpace(os.Getenv("MOONBASIC_VSIX")); env != "" {
		if fileExists(env) {
			return env, nil
		}
	}
	for _, dir := range vsixSearchDirs() {
		if p := findVsixInDir(dir); p != "" {
			return p, nil
		}
	}
	fmt.Println("No local VSIX found — downloading latest release from GitHub …")
	return downloadLatestVsix()
}

func vsixSearchDirs() []string {
	var dirs []string
	if exe, err := os.Executable(); err == nil {
		dirs = append(dirs, filepath.Dir(exe))
	}
	if cwd, err := os.Getwd(); err == nil {
		dirs = append(dirs, cwd)
	}
	if home, err := os.UserHomeDir(); err == nil {
		dirs = append(dirs, filepath.Join(home, "Downloads"))
	}
	return dirs
}

func findVsixInDir(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	var fallback string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := strings.ToLower(e.Name())
		if !strings.HasSuffix(name, ".vsix") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		if strings.Contains(name, "vscode") || strings.Contains(name, "moonbasic") {
			return path
		}
		if fallback == "" {
			fallback = path
		}
	}
	return fallback
}

func downloadLatestVsix() (string, error) {
	client := &http.Client{Timeout: 120 * time.Second}
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", githubOwner, githubRepo)
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "moonbasic-install-vscode")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch release metadata: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("GitHub API %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}

	var release struct {
		Assets []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("parse release JSON: %w", err)
	}

	var downloadURL string
	for _, a := range release.Assets {
		if strings.HasSuffix(strings.ToLower(a.Name), "-vscode.vsix") {
			downloadURL = a.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		return "", fmt.Errorf("no moonbasic-*-vscode.vsix asset on latest GitHub release")
	}

	cacheDir := filepath.Join(os.TempDir(), "moonbasic-vscode")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", err
	}
	dest := filepath.Join(cacheDir, "moonbasic-vscode.vsix")

	resp2, err := client.Get(downloadURL)
	if err != nil {
		return "", fmt.Errorf("download VSIX: %w", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download VSIX: HTTP %s", resp2.Status)
	}
	f, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(f, resp2.Body); err != nil {
		f.Close()
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	fmt.Printf("Downloaded %s\n", dest)
	return dest, nil
}

func resolveEditorCLI(forced string) (editorInfo, error) {
	if forced != "" {
		if p, err := exec.LookPath(forced); err == nil {
			return editorInfo{name: forced, cmd: p, settingsPath: defaultSettingsPath(forced)}, nil
		}
		return editorInfo{}, fmt.Errorf("editor command %q not found on PATH", forced)
	}

	candidates := []struct {
		cmds         []string
		name         string
		settingsPath string
	}{
		{[]string{"cursor"}, "Cursor", cursorSettingsPath()},
		{[]string{"code"}, "VS Code", vscodeSettingsPath()},
		{[]string{"code-insiders"}, "VS Code Insiders", vscodeInsidersSettingsPath()},
		{[]string{"codium"}, "VSCodium", codiumSettingsPath()},
	}

	for _, c := range candidates {
		for _, name := range c.cmds {
			if p, err := exec.LookPath(name); err == nil {
				sp := c.settingsPath
				if sp == "" {
					sp = defaultSettingsPath(name)
				}
				return editorInfo{name: c.name, cmd: p, settingsPath: sp}, nil
			}
		}
	}

	// Well-known install locations (Windows especially).
	for _, c := range wellKnownEditorPaths() {
		if fileExists(c.cmd) {
			return c, nil
		}
	}

	return editorInfo{}, fmt.Errorf(
		"no VS Code or Cursor CLI found on PATH\n\nInstall one of:\n  • Visual Studio Code (enable \"Shell Command: Install 'code' command in PATH\")\n  • Cursor\n\nThen run: moonbasic install-vscode\n\nOr pass --editor cursor or --editor code",
	)
}

func wellKnownEditorPaths() []editorInfo {
	if runtime.GOOS != "windows" {
		return nil
	}
	local := os.Getenv("LOCALAPPDATA")
	if local == "" {
		return nil
	}
	paths := []editorInfo{
		{
			name:         "VS Code",
			cmd:          filepath.Join(local, "Programs", "Microsoft VS Code", "bin", "code.cmd"),
			settingsPath: vscodeSettingsPath(),
		},
		{
			name:         "Cursor",
			cmd:          filepath.Join(local, "Programs", "cursor", "resources", "app", "bin", "cursor.cmd"),
			settingsPath: cursorSettingsPath(),
		},
	}
	return paths
}

func vscodeSettingsPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "Code", "User", "settings.json")
	}
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Code", "User", "settings.json")
	default:
		return filepath.Join(home, ".config", "Code", "User", "settings.json")
	}
}

func vscodeInsidersSettingsPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "Code - Insiders", "User", "settings.json")
	}
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Code - Insiders", "User", "settings.json")
	default:
		return filepath.Join(home, ".config", "Code - Insiders", "User", "settings.json")
	}
}

func cursorSettingsPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "Cursor", "User", "settings.json")
	}
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Cursor", "User", "settings.json")
	default:
		return filepath.Join(home, ".config", "Cursor", "User", "settings.json")
	}
}

func codiumSettingsPath() string {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "VSCodium", "User", "settings.json")
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "VSCodium", "User", "settings.json")
	default:
		return filepath.Join(home, ".config", "VSCodium", "User", "settings.json")
	}
}

func defaultSettingsPath(editor string) string {
	switch strings.ToLower(editor) {
	case "cursor":
		return cursorSettingsPath()
	case "code-insiders":
		return vscodeInsidersSettingsPath()
	case "codium":
		return codiumSettingsPath()
	default:
		return vscodeSettingsPath()
	}
}

func writeUserSettings(settingsPath, moonbasic, moonrun string) error {
	if settingsPath == "" {
		return fmt.Errorf("unknown settings path")
	}
	dir := filepath.Dir(settingsPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	merged := map[string]interface{}{}
	if b, err := os.ReadFile(settingsPath); err == nil && len(b) > 0 {
		clean := stripJSONComments(string(b))
		if err := json.Unmarshal([]byte(clean), &merged); err != nil {
			merged = map[string]interface{}{}
		}
	}

	if moonbasic != "" {
		merged["moonbasic.languageServerPath"] = moonbasic
	}
	if moonrun != "" {
		merged["moonbasic.moonrunPath"] = moonrun
	}

	out, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(settingsPath, append(out, '\n'), 0o644)
}

func stripJSONComments(src string) string {
	var lines []string
	for _, line := range strings.Split(src, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") {
			continue
		}
		if idx := strings.Index(line, "//"); idx >= 0 {
			// crude: only strip if // is outside quotes — good enough for settings files.
			if strings.Count(line[:idx], "\"")%2 == 0 {
				line = line[:idx]
			}
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func run(cmd string, prefix []string, args ...string) (string, error) {
	all := append(append([]string{}, prefix...), args...)
	c := exec.Command(cmd, all...)
	out, err := c.CombinedOutput()
	return string(out), err
}
