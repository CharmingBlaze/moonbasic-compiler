package main

import (
	"flag"
	"fmt"
	"os"

	"moonbasic/internal/vscodeinstall"
)

func runInstallVscodeCmd(argv []string) int {
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
