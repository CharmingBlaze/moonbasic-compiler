package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerBurst(r runtime.Registrar) {
	r.Register("GAME.BURSTSPAWN", "game", m.gameBurstSpawn)
}

func (m *Module) gameBurstSpawn(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("GAME.BURSTSPAWN expects (template, count, x, y, z)")
	}
	count, ok := args[1].ToInt()
	if !ok {
		if f, okf := args[1].ToFloat(); okf {
			count = int64(f)
		} else {
			return value.Nil, fmt.Errorf("GAME.BURSTSPAWN: count must be numeric")
		}
	}
	if count < 0 {
		count = 0
	}
	if count > 100000 {
		count = 100000
	}
	x, ok1 := args[2].ToFloat()
	y, ok2 := args[3].ToFloat()
	z, ok3 := args[4].ToFloat()
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("GAME.BURSTSPAWN: x, y, z must be numeric")
	}

	// Create a transient emitter and burst at the world position.
	emitter, err := rt.Call("PARTICLE.MAKE", nil)
	if err != nil {
		return value.Nil, fmt.Errorf("GAME.BURSTSPAWN: %w", err)
	}
	if _, err := rt.Call("PARTICLE.SETPOS", []value.Value{emitter, value.FromFloat(x), value.FromFloat(y), value.FromFloat(z)}); err != nil {
		return value.Nil, err
	}
	// Optional: copy template entity tint if the handle is an entity (best-effort).
	_ = args[0]
	if _, err := rt.Call("PARTICLE.SETBURST", []value.Value{emitter, value.FromInt(count)}); err != nil {
		return value.Nil, err
	}
	if _, err := rt.Call("PARTICLE.PLAY", []value.Value{emitter}); err != nil {
		return value.Nil, err
	}
	return emitter, nil
}
