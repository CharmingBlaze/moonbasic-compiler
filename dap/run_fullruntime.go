//go:build fullruntime

package dap

import (
	"fmt"
	"os"

	"moonbasic/compiler/pipeline"
	"moonbasic/vm"
)

func (s *session) beginDebugRun() {
	s.mu.Lock()
	if s.running || s.program == "" {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.runDone = make(chan struct{})
	s.mu.Unlock()

	go func() {
		defer close(s.runDone)
		defer func() {
			s.mu.Lock()
			s.running = false
			s.vm = nil
			s.paused = false
			s.mu.Unlock()
		}()

		prog, err := pipeline.CompileFile(s.program)
		if err != nil {
			s.sendTerminated(fmt.Errorf("compile: %w", err))
			return
		}

		err = pipeline.RunProgramDebug(prog, pipeline.Options{Out: os.Stderr}, pipeline.DebugConfig{
			BreakLines: s.breakLines(),
			OnPaused: func(v *vm.VM) {
				s.onPaused(v, v.DebugPauseReason())
			},
		})
		s.sendTerminated(err)
	}()
}
