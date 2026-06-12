package mbgame

import (
	"fmt"
	"time"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

type callbackTimer struct {
	id       int
	fn       string
	next     time.Time
	interval time.Duration
	repeat   bool
}

// SetUserInvoker wires VM.CallUserFunction for TIMER.AFTER / TIMER.EVERY callbacks.
func (m *Module) SetUserInvoker(fn func(string, []value.Value) (value.Value, error)) {
	m.invoke = fn
}

func (m *Module) registerCallbackTimers(r runtime.Registrar) {
	r.Register("TIMER.AFTER", "game", m.timerAfter)
	r.Register("TIMER.EVERY", "game", m.timerEvery)
	r.Register("TIMER.CANCEL", "game", m.timerCancel)
}

func (m *Module) timerAfter(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TIMER.AFTER expects (seconds, functionName$)")
	}
	sec, ok := argF(args[0])
	if !ok || sec < 0 {
		return value.Nil, fmt.Errorf("TIMER.AFTER: seconds must be a non-negative number")
	}
	fn, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	return m.addCallback(fn, sec, false)
}

func (m *Module) timerEvery(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TIMER.EVERY expects (seconds, functionName$)")
	}
	sec, ok := argF(args[0])
	if !ok || sec <= 0 {
		return value.Nil, fmt.Errorf("TIMER.EVERY: seconds must be a positive number")
	}
	fn, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	return m.addCallback(fn, sec, true)
}

func (m *Module) addCallback(fn string, sec float64, repeat bool) (value.Value, error) {
	m.cbNextID++
	id := m.cbNextID
	d := time.Duration(sec * float64(time.Second))
	m.callbacks = append(m.callbacks, callbackTimer{
		id:       id,
		fn:       fn,
		next:     time.Now().Add(d),
		interval: d,
		repeat:   repeat,
	})
	return value.FromInt(int64(id)), nil
}

func (m *Module) timerCancel(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TIMER.CANCEL expects (timerId)")
	}
	id64, ok := args[0].ToInt()
	if !ok {
		return value.Nil, fmt.Errorf("TIMER.CANCEL: timerId must be numeric")
	}
	id := int(id64)
	for i, t := range m.callbacks {
		if t.id == id {
			m.callbacks = append(m.callbacks[:i], m.callbacks[i+1:]...)
			return value.Nil, nil
		}
	}
	return value.Nil, nil
}

// TickCallbackTimers runs pending TIMER.AFTER / TIMER.EVERY callbacks.
func (m *Module) TickCallbackTimers() {
	if m.invoke == nil || len(m.callbacks) == 0 {
		return
	}
	now := time.Now()
	keep := make([]callbackTimer, 0, len(m.callbacks))
	for _, t := range m.callbacks {
		if now.Before(t.next) {
			keep = append(keep, t)
			continue
		}
		_, _ = m.invoke(t.fn, nil)
		if t.repeat {
			t.next = now.Add(t.interval)
			keep = append(keep, t)
		}
	}
	m.callbacks = keep
}
