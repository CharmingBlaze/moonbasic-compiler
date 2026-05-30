package pipeline

import (
	"bytes"
	"io"

	"moonbasic/runtime"
	"moonbasic/vm"
	"moonbasic/vm/heap"
	"moonbasic/vm/opcode"
)

// RunProgramHeadless executes bytecode with core builtins only (PRINT, STR, program control).
// No Raylib or game modules — suitable for the web playground and compiler-only builds.
func RunProgramHeadless(prog *opcode.Program) (stdout string, err error) {
	var buf bytes.Buffer
	return runProgramHeadlessTo(prog, &buf)
}

func runProgramHeadlessTo(prog *opcode.Program, out io.Writer) (string, error) {
	h := heap.New()
	reg := runtime.NewRegistryHeadless(h)
	reg.InitCore()
	reg.DiagOut = out

	machine := vm.New(reg, h)
	defer reg.Shutdown()

	if err := machine.Execute(prog); err != nil {
		if buf, ok := out.(*bytes.Buffer); ok {
			return buf.String(), err
		}
		return "", err
	}
	if buf, ok := out.(*bytes.Buffer); ok {
		return buf.String(), nil
	}
	return "", nil
}
