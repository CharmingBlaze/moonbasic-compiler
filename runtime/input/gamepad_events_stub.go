//go:build !cgo && !windows

package input

import (
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) SetUserInvoker(fn func(string, []value.Value) (value.Value, error)) {
	m.invoke = fn
}

func (m *Module) inOnGamepad(_ *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = args
	return value.Nil, nil
}

func (m *Module) PollGamepadEvents() {}
