package playground

import (
	"strings"
	"testing"
)

func TestRunSourceHello(t *testing.T) {
	res := runSource(`PRINT("ok")`)
	if !res.OK {
		t.Fatalf("compile/run failed: %v", res.Errors)
	}
	if !strings.Contains(res.Output, "ok") {
		t.Fatalf("expected ok in output, got %q", res.Output)
	}
}

func TestRunSourceCompileError(t *testing.T) {
	res := runSource(`FUNCTION`)
	if res.OK {
		t.Fatal("expected compile error")
	}
}
