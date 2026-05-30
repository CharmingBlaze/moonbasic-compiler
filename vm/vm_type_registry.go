package vm

import (
	"strings"

	"moonbasic/vm/heap"
	"moonbasic/vm/opcode"
	"moonbasic/vm/value"
)

func (v *VM) registerTypeInstance(typeName string, h heap.Handle) {
	if v.typeInstances == nil {
		v.typeInstances = make(map[string][]heap.Handle)
	}
	key := strings.ToUpper(strings.TrimSpace(typeName))
	v.typeInstances[key] = append(v.typeInstances[key], h)
}

func (v *VM) unregisterTypeInstance(h heap.Handle) {
	if v.typeInstances == nil {
		return
	}
	obj, ok := v.Heap.Get(h)
	if !ok {
		return
	}
	inst, ok := obj.(*heap.Instance)
	if !ok {
		return
	}
	key := strings.ToUpper(inst.Type)
	slice := v.typeInstances[key]
	for i, eh := range slice {
		if eh == h {
			v.typeInstances[key] = append(slice[:i], slice[i+1:]...)
			return
		}
	}
}

func (v *VM) doTypeInstances(i opcode.Instruction) error {
	frame := v.CallStack.Top()
	typeName := frame.Chunk.Names[i.Operand]
	key := strings.ToUpper(typeName)
	handles := v.typeInstances[key]
	n := len(handles)
	arr, err := heap.NewArrayOfKind([]int64{int64(n)}, heap.ArrayKindHandle, 0)
	if err != nil {
		return v.runtimeError(err.Error())
	}
	for j, h := range handles {
		if err := arr.SetHandle([]int64{int64(j + 1)}, int32(h)); err != nil {
			return v.runtimeError(err.Error())
		}
	}
	ah, err := v.Heap.Alloc(arr)
	if err != nil {
		return v.runtimeError(err.Error())
	}
	v.setReg(i.Dst, value.FromHandle(ah))
	return nil
}
