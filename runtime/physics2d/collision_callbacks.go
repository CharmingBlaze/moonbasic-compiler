package mbphysics2d

import (
	"fmt"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type collRule2d struct {
	ha, hb heap.Handle
	cb     string
}

var (
	collRules2d   []collRule2d
	collPending2d []collEvent2d
)

type collEvent2d struct {
	ha, hb heap.Handle
	cb     string
}

func (m *Module) registerCollisionCallbacks(reg runtime.Registrar) {
	reg.Register("PHYSICS2D.ONCOLLISION", "physics2d", m.ph2dOnCollision)
	reg.Register("PHYSICS2D.PROCESSCOLLISIONS", "physics2d", runtime.AdaptLegacy(m.ph2dProcessCollisions))
}

func (m *Module) ph2dOnCollision(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("PHYSICS2D.ONCOLLISION expects (bodyA, bodyB, callback)")
	}
	ha := heap.Handle(args[0].IVal)
	hb := heap.Handle(args[1].IVal)
	cb, err := rt.ArgCallback(args, 2)
	if err != nil {
		return value.Nil, err
	}
	cb = strings.TrimSpace(cb)
	if cb == "" {
		return value.Nil, fmt.Errorf("PHYSICS2D.ONCOLLISION: callback required")
	}
	collRules2d = append(collRules2d, collRule2d{ha: ha, hb: hb, cb: cb})
	return value.Nil, nil
}

func (m *Module) ph2dProcessCollisions(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("PHYSICS2D.PROCESSCOLLISIONS expects 0 arguments")
	}
	q := collPending2d
	collPending2d = nil
	if m.invoke != nil {
		for _, ev := range q {
			_, _ = m.invoke(ev.cb, []value.Value{value.FromHandle(ev.ha), value.FromHandle(ev.hb)})
		}
	}
	return value.Nil, nil
}

func queueCollisionCallbacks2d(m *Module) {
	if len(collRules2d) == 0 || contact2d == nil {
		return
	}
	for _, rule := range collRules2d {
		ca, oka := contact2d[rule.ha]
		cb, okb := contact2d[rule.hb]
		if oka && ca.hit && ca.other == rule.hb {
			collPending2d = append(collPending2d, collEvent2d{ha: rule.ha, hb: rule.hb, cb: rule.cb})
		} else if okb && cb.hit && cb.other == rule.ha {
			collPending2d = append(collPending2d, collEvent2d{ha: rule.ha, hb: rule.hb, cb: rule.cb})
		}
	}
	if m != nil && m.invoke != nil && len(collPending2d) > 0 {
		q := collPending2d
		collPending2d = nil
		for _, ev := range q {
			_, _ = m.invoke(ev.cb, []value.Value{value.FromHandle(ev.ha), value.FromHandle(ev.hb)})
		}
	}
}

func clearCollisionCallbacks2d() {
	collRules2d = nil
	collPending2d = nil
}
