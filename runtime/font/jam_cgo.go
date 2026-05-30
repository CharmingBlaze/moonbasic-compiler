//go:build cgo || (windows && !cgo)

package mbfont

import (
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) fontBuiltinNamed(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("FONT.BUILTIN: heap not bound")
	}
	if len(args) > 1 {
		return value.Nil, fmt.Errorf("FONT.BUILTIN expects 0 or 1 arguments")
	}
	if len(args) == 1 {
		n, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		n = strings.ToLower(strings.TrimSpace(n))
		if n != "" && n != "default" {
			return value.Nil, fmt.Errorf("FONT.BUILTIN: unknown %q (only \"default\" for now)", n)
		}
	}
	o := &fontObj{f: rl.GetFontDefault(), isDefault: true}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) registerJamHooks(reg runtime.Registrar) {
	reg.Register("FONT.BUILTIN", "font", m.fontBuiltinNamed)
}
