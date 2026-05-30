package vm

import (
	"moonbasic/vm/opcode"
	"moonbasic/vm/value"
)

func (v *VM) doPushFuncRef(i opcode.Instruction) error {
	frame := v.CallStack.Top()
	if int(i.Operand) < 0 || int(i.Operand) >= len(frame.Chunk.Names) {
		return v.runtimeError("PUSH_FUNC_REF: bad name index")
	}
	v.setReg(i.Dst, value.FuncRef(int32(i.Operand)))
	return nil
}

func (v *VM) doCallRef(i opcode.Instruction) error {
	frame := v.CallStack.Top()
	ref := v.reg(i.SrcA)
	if ref.Kind != value.KindFunc {
		return v.runtimeError("CALL_REF: expected function reference")
	}
	nameIdx := int32(ref.IVal)
	if int(nameIdx) < 0 || int(nameIdx) >= len(frame.Chunk.Names) {
		return v.runtimeError("CALL_REF: invalid function reference")
	}
	name := frame.Chunk.Names[nameIdx]
	targetChunk, ok := v.Program.Functions[name]
	if !ok {
		return v.runtimeError("undefined function: " + name)
	}

	argCount := int(i.Operand)
	args := make([]value.Value, argCount)
	argStart := i.SrcB
	for idx := 0; idx < argCount; idx++ {
		args[idx] = v.reg(argStart + uint8(idx))
	}

	v.CallStack.Push(targetChunk, 0, i.Dst)
	newFrame := v.CallStack.Top()
	for idx := 0; idx < argCount; idx++ {
		newFrame.Registers[idx] = args[idx]
	}
	return nil
}
