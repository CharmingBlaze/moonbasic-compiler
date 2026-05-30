//go:build cgo || (windows && !cgo)

package mbsprite

import (
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

type jamSpec struct {
	w, h int32
	r, g, b, a uint8
}

var jamSprites = map[string]jamSpec{
	"player": {32, 32, 100, 200, 255, 255},
	"enemy":  {28, 28, 255, 80, 80, 255},
	"bullet": {8, 8, 255, 255, 100, 255},
	"tile":   {32, 32, 120, 120, 130, 255},
	"coin":   {16, 16, 255, 215, 0, 255},
	"heart":  {16, 16, 255, 60, 100, 255},
	"star":   {24, 24, 255, 240, 120, 255},
	"block":  {32, 32, 80, 160, 80, 255},
}

func (m *Module) spBuiltin(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("SPRITE.BUILTIN: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SPRITE.BUILTIN expects 1 string name")
	}
	name, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	name = strings.ToLower(strings.TrimSpace(name))
	spec, ok := jamSprites[name]
	if !ok {
		keys := make([]string, 0, len(jamSprites))
		for k := range jamSprites {
			keys = append(keys, k)
		}
		return value.Nil, fmt.Errorf("SPRITE.BUILTIN: unknown %q (try: %s)", name, strings.Join(keys, ", "))
	}
	col := rl.NewColor(spec.r, spec.g, spec.b, spec.a)
	img := rl.GenImageColor(int(spec.w), int(spec.h), col)
	tex := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	obj := &spriteObj{
		tex: tex, frameW: spec.w, frameH: spec.h, numFrames: 1, curFrame: 0,
		scaleX: 1, scaleY: 1, alpha: spriteDefaultAlpha,
		tr: 255, tg: 255, tb: 255,
	}
	id, err := m.h.Alloc(obj)
	if err != nil {
		rl.UnloadTexture(tex)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func registerJamCommands(reg runtime.Registrar, m *Module) {
	reg.Register("SPRITE.BUILTIN", "sprite", m.spBuiltin)
}
