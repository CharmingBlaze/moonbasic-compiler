package mbtime

import (
	"strings"
	"sync"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

type loopMode int

const (
	loopVariable loopMode = iota
	loopFixed
	loopSemiFixed
)

var (
	loopMu       sync.Mutex
	currentLoop  = loopVariable
	fixedStepHz  = 60.0
	semiFixedCap = 0.1 // max delta seconds in semi-fixed mode
	physicsAccum float64
)

// NoteFrameDelta records raw frame delta for fixed-step accumulation.
func NoteFrameDelta(raw float64) {
	loopMu.Lock()
	defer loopMu.Unlock()
	if currentLoop == loopFixed && fixedStepHz > 0 {
		physicsAccum += raw
	}
}

// PhysicsSteps returns how many fixed physics updates to run this frame.
func PhysicsSteps() int {
	loopMu.Lock()
	defer loopMu.Unlock()
	if currentLoop != loopFixed || fixedStepHz <= 0 {
		return 1
	}
	step := 1.0 / fixedStepHz
	n := 0
	for physicsAccum >= step && n < 8 {
		physicsAccum -= step
		n++
	}
	if n == 0 {
		return 0
	}
	return n
}

// PhysicsStepDelta returns seconds per fixed physics step.
func PhysicsStepDelta() float64 {
	loopMu.Lock()
	defer loopMu.Unlock()
	if fixedStepHz <= 0 {
		return 0
	}
	return 1.0 / fixedStepHz
}

// SetLoopMode configures how DeltaSeconds reports frame time.
// mode: "variable" (default), "fixed", or "semi-fixed".
// hz is the fixed step rate for "fixed" mode, or max-delta seconds for "semi-fixed" when > 0.
func SetLoopMode(mode string, hz float64) {
	loopMu.Lock()
	defer loopMu.Unlock()
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "fixed":
		currentLoop = loopFixed
		if hz > 0 {
			fixedStepHz = hz
		}
	case "semi-fixed", "semifixed", "semi_fixed":
		currentLoop = loopSemiFixed
		if hz > 0 {
			semiFixedCap = hz
		}
	default:
		currentLoop = loopVariable
	}
}

func applyLoopMode(dt float64) float64 {
	loopMu.Lock()
	defer loopMu.Unlock()
	switch currentLoop {
	case loopFixed:
		if fixedStepHz <= 0 {
			return dt
		}
		return 1.0 / fixedStepHz
	case loopSemiFixed:
		if dt > semiFixedCap {
			return semiFixedCap
		}
		return dt
	default:
		return dt
	}
}

func registerLoopModeCommands(reg runtime.Registrar) {
	reg.Register("WINDOW.SETLOOPMODE", "window", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) < 1 || len(args) > 2 {
			return value.Nil, runtime.Errorf("WINDOW.SETLOOPMODE expects (mode$) or (mode$, hz)")
		}
		mode, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		param := 0.0
		if len(args) == 2 {
			if f, ok := args[1].ToFloat(); ok {
				param = f
			} else if i, ok := args[1].ToInt(); ok {
				param = float64(i)
			}
		}
		SetLoopMode(mode, param)
		return value.Nil, nil
	})
	reg.Register("TIME.PHYSICSSTEPS", "time", func(_ *runtime.Runtime, _ ...value.Value) (value.Value, error) {
		return value.Int(int64(PhysicsSteps())), nil
	})
	reg.Register("TIME.PHYSICSSTEP", "time", func(_ *runtime.Runtime, _ ...value.Value) (value.Value, error) {
		return value.Float(PhysicsStepDelta()), nil
	})
}
