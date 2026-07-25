//go:build !fullruntime

package pipeline

import (
	"errors"

	"moonbasic/compiler/builtinmanifest"
	"moonbasic/vm/opcode"
)

// RunProgram is a stub when the game runtime is not linked.
func RunProgram(prog *opcode.Program, opts Options) error {
	return errors.New("moonBASIC runtime engine is not included in this build (rebuild with -tags fullruntime)")
}

// ListBuiltins returns manifest command keys (compiler-only builds have no live registry).
func ListBuiltins() []string {
	return builtinmanifest.Default().Keys()
}
