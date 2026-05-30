//go:build cgo || (windows && !cgo)

package input

import (
	"fmt"
	"strings"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

type gamepadListener struct {
	pad int // -1 = any pad
	fn  string
}

var (
	gpListenMu sync.RWMutex
	gpListen   []gamepadListener
	gpWasConn  [8]bool
)

func (m *Module) SetUserInvoker(fn func(string, []value.Value) (value.Value, error)) {
	m.invoke = fn
}

func (m *Module) inOnGamepad(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("INPUT.ONGAMEPAD expects (gamepadIndex, callback)")
	}
	pad := int64(-1)
	if i, ok := args[0].ToInt(); ok {
		pad = i
	} else if f, ok := args[0].ToFloat(); ok {
		pad = int64(f)
	} else {
		return value.Nil, fmt.Errorf("INPUT.ONGAMEPAD: gamepad index must be numeric (-1 = any)")
	}
	fn, err := rt.ArgCallback(args, 1)
	if err != nil {
		return value.Nil, err
	}
	fn = strings.ToLower(strings.TrimSpace(fn))
	if fn == "" {
		return value.Nil, fmt.Errorf("INPUT.ONGAMEPAD: callback required")
	}
	gpListenMu.Lock()
	gpListen = append(gpListen, gamepadListener{pad: int(pad), fn: fn})
	gpListenMu.Unlock()
	return value.Nil, nil
}

// PollGamepadEvents detects connect/disconnect and invokes registered callbacks.
func (m *Module) PollGamepadEvents() {
	if m.invoke == nil {
		return
	}
	gpListenMu.RLock()
	if len(gpListen) == 0 {
		gpListenMu.RUnlock()
		return
	}
	listeners := append([]gamepadListener(nil), gpListen...)
	gpListenMu.RUnlock()

	for pad := 0; pad < len(gpWasConn); pad++ {
		now := rl.IsGamepadAvailable(int32(pad))
		if now == gpWasConn[pad] {
			continue
		}
		gpWasConn[pad] = now
		for _, l := range listeners {
			if l.pad >= 0 && l.pad != pad {
				continue
			}
			connected := value.FromBool(now)
			idx := value.FromInt(int64(pad))
			_, _ = m.invoke(l.fn, []value.Value{idx, connected})
		}
	}
}
