//go:build fullruntime

package pipeline

import (
	"fmt"
	goruntime "runtime"
	"os"

	"moonbasic/runtime"
	"moonbasic/vm"
	"moonbasic/vm/heap"
	"moonbasic/vm/opcode"
)

// DebugConfig wires breakpoint and pause notifications for DAP / interactive debuggers.
type DebugConfig struct {
	BreakLines []int
	OnPaused   func(*vm.VM)
}

// RunProgramDebug executes prog with VM breakpoints and optional pause callbacks.
func RunProgramDebug(prog *opcode.Program, opts Options, dbg DebugConfig) error {
	if opts.Out == nil {
		opts.Out = os.Stderr
	}
	if opts.Debug {
		fmt.Fprintln(opts.Out, prog.Main.Disassemble())
	}

	goruntime.LockOSThread()

	h := heap.New()
	d := DefaultDriver()
	reg := runtime.NewRegistry(h, d)
	setupRegistry(reg, h, opts)

	machine := vm.New(reg, h)
	wireRegistryCallbacks(reg, machine)

	machine.Trace = opts.Trace
	machine.TraceOut = opts.Out
	machine.Profiler = opts.ProfileRecorder
	machine.EnableDebugMode(true)
	machine.SetBreakLines(dbg.BreakLines)
	machine.OnPaused = dbg.OnPaused

	defer reg.Shutdown()

	return machine.Execute(prog)
}
