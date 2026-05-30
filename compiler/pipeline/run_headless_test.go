package pipeline

import (
	"strings"
	"testing"
)

func TestRunProgramHeadlessPrint(t *testing.T) {
	prog, err := CompileSource("t.mb", `PRINT("hi")
FOR i = 1 TO 2
    PRINT(i)
NEXT
`)
	if err != nil {
		t.Fatal(err)
	}
	out, err := RunProgramHeadless(prog)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "hi") {
		t.Fatalf("expected hi in output, got %q", out)
	}
	if !strings.Contains(out, "1") || !strings.Contains(out, "2") {
		t.Fatalf("expected loop output, got %q", out)
	}
}
