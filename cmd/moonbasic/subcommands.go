package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"moonbasic/internal/scaffold"
)

func dispatchSubcommand() (handled bool, exit int) {
	if len(os.Args) < 2 {
		return false, 0
	}
	switch strings.ToLower(strings.TrimSpace(os.Args[1])) {
	case "new":
		return true, runNewProjectCmd(os.Args[2:])
	}
	return false, 0
}

func runNewProjectCmd(argv []string) int {
	fs := flag.NewFlagSet("new", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(argv); err != nil {
		return 2
	}
	rest := fs.Args()
	if len(rest) != 1 {
		fmt.Fprintln(os.Stderr, "usage: moonbasic new <project-name>")
		return 2
	}
	dir, err := scaffold.Create(rest[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Printf("Created project %s\n  main.mb\n  assets/\n  .vscode/launch.json\n  README.md\n", dir)
	fmt.Println("Next: moonrun main.mb")
	return 0
}
