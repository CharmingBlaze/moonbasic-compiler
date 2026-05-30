//go:build fullruntime

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"moonbasic/compiler/pipeline"
	"moonbasic/internal/version"
	"moonbasic/lsp"
	"moonbasic/vm"
)

func main() {
	var (
		checkOnly        = flag.Bool("check", false, "parse and type-check only")
		strictDeprecated = flag.Bool("strict-deprecated", false, "with --check: treat deprecated MAKE/SETPOSITION aliases as errors")
		showVer          = flag.Bool("version", false, "print version and exit")
		lspMode          = flag.Bool("lsp", false, "run Language Server Protocol (stdio) for editors")
		disasm           = flag.Bool("disasm", false, "print human-readable bytecode for a .mbc file")
		runMode          = flag.Bool("run", false, "compile and run the program")
		watchMode        = flag.Bool("watch", false, "recompile and rerun when the source file changes")
		debugInfo        = flag.Bool("info", false, "enable debug diagnostics and FPS graph")
		profileMode      = flag.Bool("profile", false, "with --run: print top hot source lines (ops + wall ms)")
		profileOut       = flag.String("profile-out", "", "with --profile: write HTML profile report to path")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "moonBASIC Compiler %s\n", version.Version)
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  moonbasic [flags] <source.mb>     compile to .mbc\n")
		fmt.Fprintf(os.Stderr, "  moonbasic --check <source.mb>     parse and type-check only\n")
		fmt.Fprintf(os.Stderr, "  moonbasic --lsp                   language server on stdio\n")
		fmt.Fprintf(os.Stderr, "  moonbasic --run <source.mb>        compile and run\n")
		fmt.Fprintf(os.Stderr, "  moonbasic --run --profile <source.mb>   instruction-count profile\n")
		fmt.Fprintf(os.Stderr, "  moonbasic --watch <source.mb>       recompile and rerun on save\n")
		fmt.Fprintf(os.Stderr, "  moonbasic --disasm <file.mbc>     disassemble bytecode\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVer {
		fmt.Printf("moonBASIC Compiler %s\n", version.Version)
		return
	}

	if *lspMode {
		if err := lsp.Serve(); err != nil {
			fmt.Fprintf(os.Stderr, "lsp: %v\n", err)
			os.Exit(1)
		}
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	path := args[0]

	if *disasm {
		if !strings.EqualFold(filepath.Ext(path), ".mbc") {
			fmt.Fprintln(os.Stderr, "error: --disasm requires a .mbc file")
			os.Exit(2)
		}
		prog, err := pipeline.DecodeMOONFromFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "decode: %v\n", err)
			os.Exit(2)
		}
		pipeline.PrintProgramDisassembly(prog, os.Stdout, nil)
		return
	}

	if *checkOnly {
		if err := pipeline.CheckFileWithOptions(path, pipeline.CheckOptions{StrictDeprecated: *strictDeprecated}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		fmt.Println("Check: OK")
		return
	}

	opts := pipeline.Options{
		Debug: *debugInfo,
		Out:   os.Stderr,
	}

	if *watchMode {
		os.Exit(runWatch(path, opts))
	}

	if *runMode {
		prog, err := pipeline.CompileFile(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		var rec *vm.ProfileRecorder
		if *profileMode || *profileOut != "" {
			rec = vm.NewProfileRecorder()
			opts.ProfileRecorder = rec
		}
		if err := pipeline.RunProgram(prog, opts); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		if rec != nil {
			pipeline.PrintProfileReport(rec, path, os.Stderr, 10)
			if *profileOut != "" {
				if err := pipeline.WriteProfileHTML(rec, path, *profileOut); err != nil {
					fmt.Fprintf(os.Stderr, "profile-out: %v\n", err)
					os.Exit(2)
				}
				fmt.Fprintf(os.Stderr, "profile: wrote %s\n", *profileOut)
			}
		}
		return
	}

	if err := compileToMBC(path); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
