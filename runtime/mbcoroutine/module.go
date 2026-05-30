package mbcoroutine

import (
	"fmt"
	"strings"
	"time"

	"moonbasic/runtime"
	"moonbasic/vm"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// Module registers COROUTINE.* commands.
type Module struct {
	machine *vm.VM
}

func NewModule() *Module { return &Module{} }

func (m *Module) BindVM(v *vm.VM) { m.machine = v }

func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("COROUTINE.START", "coroutine", m.coStart)
	reg.Register("COROUTINE.RESUME", "coroutine", m.coResume)
	reg.Register("COROUTINE.WAIT", "coroutine", m.coWait)
	reg.Register("COROUTINE.DONE", "coroutine", m.coDone)
}

func (m *Module) Reset()     {}
func (m *Module) Shutdown() {}

func (m *Module) requireVM() (*vm.VM, error) {
	if m.machine == nil {
		return nil, runtime.Errorf("COROUTINE: VM not bound")
	}
	return m.machine, nil
}

func (m *Module) coStart(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	v, err := m.requireVM()
	if err != nil {
		return value.Nil, err
	}
	if len(args) < 1 {
		return value.Nil, fmt.Errorf("COROUTINE.START expects (functionName$, ...args)")
	}
	fn, err := rt.ArgCallback(args, 0)
	if err != nil {
		return value.Nil, err
	}
	fn = strings.ToLower(strings.TrimSpace(fn))
	return v.StartCoroutine(fn, args[1:])
}

func (m *Module) coResume(_ *runtime.Runtime, args ...value.Value) (value.Value, error) {
	v, err := m.requireVM()
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("COROUTINE.RESUME expects (coroutineHandle)")
	}
	now := time.Now().UnixNano() / 1e9
	return v.ResumeCoroutine(heap.Handle(args[0].IVal), float64(now))
}

func (m *Module) coWait(_ *runtime.Runtime, args ...value.Value) (value.Value, error) {
	v, err := m.requireVM()
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("COROUTINE.WAIT expects (seconds)")
	}
	sec, ok := args[0].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("COROUTINE.WAIT: seconds must be numeric")
	}
	now := float64(time.Now().UnixNano()) / 1e9
	if err := v.CoroutineWait(sec, now); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) coDone(_ *runtime.Runtime, args ...value.Value) (value.Value, error) {
	v, err := m.requireVM()
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("COROUTINE.DONE expects (coroutineHandle)")
	}
	done, err := v.CoroutineDone(heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(done), nil
}
