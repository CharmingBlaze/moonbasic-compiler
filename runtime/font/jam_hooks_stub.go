//go:build !cgo && !windows

package mbfont

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerJamHooks(reg runtime.Registrar) {
	reg.Register("FONT.BUILTIN", "font", func(_ *runtime.Runtime, _ ...value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("FONT.BUILTIN: %s", hint)
	})
}
