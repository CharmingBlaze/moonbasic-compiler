package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"moonbasic/internal/scaffold"
	"moonbasic/internal/vscodeinstall"
)

func dispatchSubcommand() (handled bool, exit int) {
	if len(os.Args) < 2 {
		return false, 0
	}
	switch strings.ToLower(strings.TrimSpace(os.Args[1])) {
	case "install", "list", "publish":
		return true, runPackageCLI(os.Args[1:])
	case "new":
		return true, runNewProject(os.Args[2:])
	case "pack":
		return true, runPack(os.Args[2:])
	case "playground":
		return true, runPlayground(os.Args[2:])
	case "test":
		return true, runTestCLI(os.Args[2:])
	case "run":
		return true, runProjectCLI("run", os.Args[2:])
	case "build":
		return true, runProjectCLI("build", os.Args[2:])
	case "package":
		return true, runPackagePlatform(os.Args[2:])
	case "install-vscode":
		return true, runInstallVscode(os.Args[2:])
	}
	return false, 0
}

func runInstallVscode(argv []string) int {
	fs := flag.NewFlagSet("install-vscode", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	vsix := fs.String("vsix", "", "path to .vsix (default: same folder as moonbasic, or download latest)")
	editor := fs.String("editor", "", "editor CLI: code, cursor, code-insiders, codium")
	skipSettings := fs.Bool("skip-settings", false, "do not write moonbasic.languageServerPath / moonbasic.moonrunPath")
	if err := fs.Parse(argv); err != nil {
		return 2
	}
	opts := vscodeinstall.Options{
		VsixPath:     *vsix,
		EditorCLI:    *editor,
		SkipSettings: *skipSettings,
	}
	if err := vscodeinstall.Install(opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func runNewProject(argv []string) int {
	fs := flag.NewFlagSet("new", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	tmpl := fs.String("template", "empty", "starter template: empty, 3d, platformer, ui, physics")
	if err := fs.Parse(argv); err != nil {
		return 2
	}
	rest := fs.Args()
	if len(rest) != 1 {
		fmt.Fprintln(os.Stderr, "usage: moonbasic new [--template name] <project-name>")
		return 2
	}
	dir, err := scaffold.Create(rest[0], *tmpl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Printf("Created project %s\n  main.mb\n  assets/\n  .vscode/ (launch, tasks, extensions)\n  README.md\n", dir)
	fmt.Println("Next: moonrun main.mb")
	return 0
}

func runTestCLI(argv []string) int {
	cmd := exec.Command("go", "test", "./compiler/...", "./vm/...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return 1
	}
	root, err := os.Getwd()
	if err != nil {
		return 1
	}
	foundation := filepath.Join(root, "examples", "foundation", "main.mb")
	if _, err := os.Stat(foundation); err == nil {
		cmd = exec.Command("go", "run", ".", "--check", foundation)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return 1
		}
	}
	return 0
}

func projectMainFile(dir string) (string, error) {
	main := filepath.Join(dir, "main.mb")
	if _, err := os.Stat(main); err == nil {
		return main, nil
	}
	return "", fmt.Errorf("no main.mb in %s (run moonbasic new <name> first)", dir)
}

func runProjectCLI(mode string, argv []string) int {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	main, err := projectMainFile(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	switch mode {
	case "run":
		cmd := exec.Command("go", "run", "-tags", "fullruntime", "./cmd/moonrun", main)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			if exit, ok := err.(*exec.ExitError); ok {
				return exit.ExitCode()
			}
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return 0
	case "build":
		if err := compileToMBC(main); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 2
		}
		out := filepath.Join(dir, "build", "main.mbc")
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		src := stringsTrimExt(main) + ".mbc"
		if err := os.Rename(src, out); err != nil {
			// If rename fails, compile directly to build/
			if err2 := compileToMBC(main); err2 != nil {
				fmt.Fprintln(os.Stderr, err2)
				return 2
			}
		}
		fmt.Printf("built %s\n", out)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown project mode %q\n", mode)
		return 2
	}
}

func runPackagePlatform(argv []string) int {
	if len(argv) == 0 {
		fmt.Fprintln(os.Stderr, "usage: moonbasic package windows|linux [game.mb]")
		return 2
	}
	platform := argv[0]
	rest := argv[1:]
	main := "main.mb"
	if len(rest) > 0 {
		main = rest[0]
	}
	switch platform {
	case "windows", "linux":
		// Platform-specific native zip is release tooling; pack ships bytecode + assets today.
		return runPack([]string{"-no-runtime", main})
	default:
		fmt.Fprintf(os.Stderr, "moonbasic package: %q not supported yet (use windows or linux)\n", platform)
		return 2
	}
}

func stringsTrimExt(path string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return path
	}
	return path[:len(path)-len(ext)]
}
