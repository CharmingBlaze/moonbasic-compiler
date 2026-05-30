package main

import (
	"flag"
	"fmt"
	"os"

	"moonbasic/internal/gamepack"
)

func runPack(argv []string) int {
	fs := flag.NewFlagSet("pack", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	assets := fs.String("assets", "assets", "assets folder relative to source file")
	out := fs.String("o", "", "output zip path (default <name>-pack.zip next to source)")
	noRuntime := fs.Bool("no-runtime", false, "omit moonbasic executable from the zip")
	if err := fs.Parse(argv); err != nil {
		return 2
	}
	rest := fs.Args()
	if len(rest) != 1 {
		fmt.Fprintln(os.Stderr, "usage: moonbasic pack [-assets dir] [-o out.zip] <game.mb>")
		return 2
	}
	zipPath, err := gamepack.Pack(rest[0], gamepack.PackOptions{
		AssetsDir:      *assets,
		OutZip:         *out,
		ExcludeRuntime: *noRuntime,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Printf("packed %s\n", zipPath)
	return 0
}
