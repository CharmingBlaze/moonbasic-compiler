package main

import (
	"flag"
	"fmt"
	"os"

	"moonbasic/internal/playground"
)

func runPlayground(argv []string) int {
	fs := flag.NewFlagSet("playground", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	addr := fs.String("addr", "127.0.0.1:8765", "listen address")
	web := fs.String("web", "", "path to web/playground (default auto)")
	if err := fs.Parse(argv); err != nil {
		return 2
	}
	s := &playground.Server{Addr: *addr, WebRoot: *web}
	if err := s.ListenAndServe(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
