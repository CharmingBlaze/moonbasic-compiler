package vm

import (
	"fmt"
	"strings"

	"moonbasic/vm/callstack"
	"moonbasic/vm/heap"
	"moonbasic/vm/opcode"
	"moonbasic/vm/value"
)

type coroutineState struct {
	frames    []callstack.Frame
	done      bool
	waitUntil float64 // wall seconds; 0 = resume manually
}

// coroutineHandle is a heap object pointing at a VM coroutine id.
type coroutineHandle struct {
	id int32
}

func (c *coroutineHandle) TypeName() string { return "COROUTINE" }
func (c *coroutineHandle) TypeTag() uint16  { return heap.TagCoroutine }
func (c *coroutineHandle) Free()            {}

func (v *VM) ensureCoroutineMap() {
	if v.coroutines == nil {
		v.coroutines = make(map[int32]*coroutineState)
	}
}

func (v *VM) doYield(_ opcode.Instruction) error {
	if v.curCoID == 0 {
		return v.runtimeError("YIELD outside COROUTINE (use COROUTINE.START first)")
	}
	v.yieldPending = true
	return nil
}

func (v *VM) saveCoroutineState(id int32) {
	v.ensureCoroutineMap()
	st := v.coroutines[id]
	if st == nil {
		st = &coroutineState{}
		v.coroutines[id] = st
	}
	st.frames = copyFrames(v.CallStack.FramesCopy())
	if v.CallStack.Depth() == 0 {
		st.done = true
	}
}

func copyFrames(in []callstack.Frame) []callstack.Frame {
	out := make([]callstack.Frame, len(in))
	copy(out, in)
	return out
}

func (v *VM) runUntilYieldOrDone() error {
	v.yieldPending = false
	for !v.Halted && v.CallStack.Depth() > 0 {
		frame := v.CallStack.Top()
		if frame.IP >= len(frame.Chunk.Instructions) {
			v.CallStack.Pop()
			continue
		}
		instr := frame.Chunk.Instructions[frame.IP]
		frame.IP++
		if err := v.step(instr); err != nil {
			return err
		}
		if v.yieldPending {
			return nil
		}
	}
	return nil
}

// StartCoroutine runs fn until YIELD or completion; returns a heap handle to the coroutine.
func (v *VM) StartCoroutine(fnName string, args []value.Value) (value.Value, error) {
	if v.Program == nil {
		return value.Nil, fmt.Errorf("no program loaded")
	}
	key := strings.ToLower(strings.TrimSpace(fnName))
	chunk, ok := v.Program.Functions[key]
	if !ok {
		return value.Nil, fmt.Errorf("undefined function: %s", key)
	}
	v.ensureCoroutineMap()
	v.nextCoID++
	id := v.nextCoID
	baseDepth := v.CallStack.Depth()

	v.curCoID = id
	v.CallStack.Push(chunk, 0, 255)
	frame := v.CallStack.Top()
	for i, a := range args {
		if i < 256 {
			frame.Registers[i] = a
		}
	}
	if err := v.runUntilYieldOrDone(); err != nil {
		for v.CallStack.Depth() > baseDepth {
			v.CallStack.Pop()
		}
		v.curCoID = 0
		return value.Nil, err
	}
	v.saveCoroutineState(id)
	for v.CallStack.Depth() > baseDepth {
		v.CallStack.Pop()
	}
	v.curCoID = 0

	obj := &coroutineHandle{id: id}
	h, err := v.Heap.Alloc(obj)
	if err != nil {
		return value.Nil, err
	}
	v.activeCoroutines = append(v.activeCoroutines, h)
	return value.FromHandle(h), nil
}

// TickCoroutines resumes all auto-tracked coroutines once per frame (called from RENDER.FRAME).
func (v *VM) TickCoroutines(now float64) {
	if len(v.activeCoroutines) == 0 {
		return
	}
	pending := append([]heap.Handle(nil), v.activeCoroutines...)
	alive := v.activeCoroutines[:0]
	for _, h := range pending {
		done, err := v.CoroutineDone(h)
		if err != nil || done {
			continue
		}
		_, _ = v.ResumeCoroutine(h, now)
		if d2, err2 := v.CoroutineDone(h); err2 == nil && !d2 {
			alive = append(alive, h)
		}
	}
	v.activeCoroutines = alive
}

// ResumeCoroutine continues a suspended coroutine until the next YIELD or completion.
func (v *VM) ResumeCoroutine(hid heap.Handle, now float64) (value.Value, error) {
	obj, ok := v.Heap.Get(hid)
	if !ok {
		return value.Nil, fmt.Errorf("COROUTINE.RESUME: invalid handle")
	}
	ch, ok := obj.(*coroutineHandle)
	if !ok {
		return value.Nil, fmt.Errorf("COROUTINE.RESUME: not a coroutine handle")
	}
	st := v.coroutines[ch.id]
	if st == nil || st.done {
		return value.FromBool(false), nil
	}
	if st.waitUntil > 0 && now < st.waitUntil {
		return value.FromBool(true), nil
	}
	st.waitUntil = 0

	baseDepth := v.CallStack.Depth()
	v.curCoID = ch.id
	v.CallStack.Frames = append(v.CallStack.Frames[:0], copyFrames(st.frames)...)

	if err := v.runUntilYieldOrDone(); err != nil {
		for v.CallStack.Depth() > baseDepth {
			v.CallStack.Pop()
		}
		v.curCoID = 0
		return value.Nil, err
	}
	v.saveCoroutineState(ch.id)
	for v.CallStack.Depth() > baseDepth {
		v.CallStack.Pop()
	}
	v.curCoID = 0
	return value.FromBool(!st.done), nil
}

// CoroutineWait schedules the active coroutine to resume after seconds elapse (wall clock).
func (v *VM) CoroutineWait(seconds, now float64) error {
	if v.curCoID == 0 {
		return fmt.Errorf("COROUTINE.WAIT outside coroutine")
	}
	st := v.coroutines[v.curCoID]
	if st == nil {
		return fmt.Errorf("COROUTINE.WAIT: missing coroutine state")
	}
	st.waitUntil = now + seconds
	v.yieldPending = true
	return nil
}

// CoroutineDone reports whether a coroutine finished.
func (v *VM) CoroutineDone(hid heap.Handle) (bool, error) {
	obj, ok := v.Heap.Get(hid)
	if !ok {
		return false, fmt.Errorf("COROUTINE.DONE: invalid handle")
	}
	ch, ok := obj.(*coroutineHandle)
	if !ok {
		return false, fmt.Errorf("COROUTINE.DONE: not a coroutine handle")
	}
	st := v.coroutines[ch.id]
	return st == nil || st.done, nil
}
