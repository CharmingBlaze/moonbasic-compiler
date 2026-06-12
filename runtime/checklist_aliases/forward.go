package checklistaliases

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func call(rt *runtime.Runtime, name string, args []value.Value) (value.Value, error) {
	reg := runtime.ActiveRegistry()
	if reg == nil {
		return value.Nil, fmt.Errorf("%s: registry not active", name)
	}
	return reg.Call(name, args)
}

func forward(name string) runtime.BuiltinFn {
	return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt, name, args)
	}
}

func forward0(name string) runtime.BuiltinFn {
	return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("%s expects 0 arguments", name)
		}
		return call(rt, name, nil)
	}
}
