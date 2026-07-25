package pipeline

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Headless end-to-end semantic checks (parse → INCLUDE → semantic) without fullruntime.

func TestCheckSource_ControlFlowOK(t *testing.T) {
	src := `
x = 0
FOR i = 1 TO 3
	x = x + i
NEXT
IF x > 0 THEN
	x = x - 1
ELSE
	x = 0
ENDIF
SELECT x
CASE 0
	x = 1
DEFAULT
	x = 2
ENDSELECT
DO
	x = x + 1
LOOP WHILE x < 5
WHILE x > 0
	x = x - 1
WEND
`
	if err := CheckSource("ctrl.mb", src); err != nil {
		t.Fatalf("control flow should pass --check: %v", err)
	}
}

func TestCheckSource_IncludeAndSemantic(t *testing.T) {
	dir := t.TempDir()
	lib := filepath.Join(dir, "lib.mb")
	if err := os.WriteFile(lib, []byte("FUNCTION Inc(n)\nRETURN n + 1\nENDFUNCTION\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	main := filepath.Join(dir, "main.mb")
	mainSrc := "INCLUDE \"lib.mb\"\nx = Inc(1)\n"
	if err := os.WriteFile(main, []byte(mainSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := CheckFile(main); err != nil {
		t.Fatalf("INCLUDE + semantic should pass: %v", err)
	}
}

func TestCheckSource_NegativeHandleAndEasyMode(t *testing.T) {
	src := `
cam = CAMERA.CREATE()
cam.Begn()
CreateCubee()
`
	err := CheckSource("neg.mb", src)
	if err == nil {
		t.Fatal("expected semantic errors")
	}
	msg := err.Error()
	if !strings.Contains(msg, "unknown method") && !strings.Contains(strings.ToLower(msg), "begn") {
		t.Fatalf("expected handle method error, got: %v", err)
	}
	if !strings.Contains(msg, "Unknown command") {
		t.Fatalf("expected Easy Mode unknown command (multi-error), got: %v", err)
	}
}

func TestCheckSource_UnknownNamespaceRejected(t *testing.T) {
	err := CheckSource("bad.mb", "CAMERA.BEGN()\n")
	if err == nil {
		t.Fatal("expected unknown namespace method error")
	}
	if !strings.Contains(err.Error(), "Unknown command") && !strings.Contains(err.Error(), "BEGN") {
		t.Fatalf("unexpected error: %v", err)
	}
}
