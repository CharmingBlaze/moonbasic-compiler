//go:build !cgo && !windows

package mbaudio

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerJamHooks(reg runtime.Registrar) {
	reg.Register("SOUND.BUILTIN", "audio", func(_ *runtime.Runtime, _ ...value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("SOUND.BUILTIN: %s", hint)
	})
}
