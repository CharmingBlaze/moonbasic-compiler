//go:build (!linux && !windows) || !cgo

package mbphysics3d

import (
	"testing"

	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func TestStubBody3DCommitReturnsHandle(t *testing.T) {
	h := heap.New()
	m := &Module{h: h}

	bv, err := phCreateBody(m, "DYNAMIC")
	if err != nil {
		t.Fatalf("CREATE: %v", err)
	}
	if bv.Kind != value.KindHandle {
		t.Fatalf("CREATE: want handle, got %s", bv.TypeName())
	}

	if _, err := m.bdAddBox([]value.Value{bv, value.FromFloat(1), value.FromFloat(1), value.FromFloat(1)}); err != nil {
		t.Fatalf("ADDBOX: %v", err)
	}

	body, err := m.bdCommit([]value.Value{bv, value.FromFloat(0), value.FromFloat(10), value.FromFloat(0)})
	if err != nil {
		t.Fatalf("COMMIT: %v", err)
	}
	if body.Kind != value.KindHandle {
		t.Fatalf("COMMIT: want Body3D handle, got %s (nil causes SETLIFT runtime errors)", body.TypeName())
	}
	if _, err := heap.Cast[*body3dObj](h, heap.Handle(body.IVal)); err != nil {
		t.Fatalf("COMMIT handle is not Body3D: %v", err)
	}
}
