package vm

import (
	"moonbasic/vm/callstack"
)

func (v *VM) SetBreakLines(lines []int) {
	if v.BreakLines == nil {
		v.BreakLines = make(map[int]bool)
	}
	for k := range v.BreakLines {
		delete(v.BreakLines, k)
	}
	for _, ln := range lines {
		if ln > 0 {
			v.BreakLines[ln] = true
		}
	}
}

func (v *VM) EnableDebugMode(on bool) {
	v.DebugMode = on
	if on && v.debugContinue == nil {
		v.debugContinue = make(chan struct{}, 1)
	}
}

func (v *VM) DebugContinue() {
	if v.debugContinue != nil {
		select {
		case v.debugContinue <- struct{}{}:
		default:
		}
	}
}

func (v *VM) DebugStepOnce() {
	v.debugStepOnce = true
	v.DebugContinue()
}

func (v *VM) shouldBreak(frame *callstack.Frame) bool {
	if !v.DebugMode || frame == nil || v.BreakLines == nil {
		return false
	}
	if v.debugStepOnce {
		v.debugStepOnce = false
		v.debugPauseReason = "step"
		return true
	}
	ip := frame.IP - 1
	if ip < 0 || ip >= len(frame.Chunk.SourceLines) {
		return false
	}
	line := int(frame.Chunk.SourceLines[ip])
	if v.BreakLines[line] {
		v.debugPauseReason = "breakpoint"
		return true
	}
	return false
}

func (v *VM) waitDebugContinue() error {
	if v.OnPaused != nil {
		v.OnPaused(v)
	}
	if v.debugContinue == nil {
		v.debugContinue = make(chan struct{}, 1)
	}
	<-v.debugContinue
	return nil
}

// DebugGlobals returns a snapshot of global variables for debug adapters.
func (v *VM) DebugGlobals() map[string]string {
	out := make(map[string]string, len(v.Globals))
	for k, val := range v.Globals {
		out[k] = val.String()
	}
	return out
}

// CurrentLine returns the 1-based source line at the current IP, or 0.
func (v *VM) CurrentLine() int {
	frame := v.CallStack.Top()
	if frame == nil || frame.IP == 0 {
		return 0
	}
	ip := frame.IP - 1
	if ip < 0 || ip >= len(frame.Chunk.SourceLines) {
		return 0
	}
	return int(frame.Chunk.SourceLines[ip])
}

// CurrentSourcePath returns the program source path when known.
func (v *VM) CurrentSourcePath() string {
	if v.Program != nil && v.Program.SourcePath != "" {
		return v.Program.SourcePath
	}
	return ""
}

// DebugPauseReason returns the DAP stop reason for the current pause ("step" or "breakpoint").
func (v *VM) DebugPauseReason() string {
	if v.debugPauseReason == "" {
		return "breakpoint"
	}
	return v.debugPauseReason
}
