package handlecall

import (
	"testing"

	"moonbasic/vm/heap"
)

func TestTagFromNamespace(t *testing.T) {
	tag, ok := TagFromNamespace("CAMERA")
	if !ok || tag != heap.TagCamera {
		t.Fatalf("CAMERA: got %d ok=%v", tag, ok)
	}
	tag, ok = TagFromNamespace("entity")
	if !ok || tag != heap.TagEntityRef {
		t.Fatalf("entity: got %d ok=%v", tag, ok)
	}
	if _, ok := TagFromNamespace("NOTANAMESSPACE"); ok {
		t.Fatal("expected unknown namespace")
	}
}

func TestDispatchBegin(t *testing.T) {
	key, prep, ok := Dispatch(heap.TagCamera, "Begin", 0)
	if !ok || key != "CAMERA.BEGIN" || !prep {
		t.Fatalf("Begin: key=%q prep=%v ok=%v", key, prep, ok)
	}
	_, _, ok = Dispatch(heap.TagCamera, "Begn", 0)
	if ok {
		t.Fatal("Begn should not dispatch")
	}
}
