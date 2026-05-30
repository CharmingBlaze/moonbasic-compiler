package lsp

import "testing"

func TestCallContextAt(t *testing.T) {
	fn, idx, ok := callContextAt("x = Add(1, 2)", 11)
	if !ok || fn != "Add" || idx != 1 {
		t.Fatalf("got %q %d %v", fn, idx, ok)
	}
	fn, idx, ok = callContextAt("Add(", 4)
	if !ok || fn != "Add" || idx != 0 {
		t.Fatalf("open paren got %q %d %v", fn, idx, ok)
	}
}
