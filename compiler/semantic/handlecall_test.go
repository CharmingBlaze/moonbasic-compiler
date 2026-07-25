package semantic

import (
	"strings"
	"testing"

	moonerrors "moonbasic/compiler/errors"
	"moonbasic/compiler/parser"
)

func TestHandleCallUnknownMethodRejected(t *testing.T) {
	src := "cam = CAMERA.CREATE()\ncam.Begn()\n"
	prog, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	a := DefaultAnalyzer("t.mb", parser.SplitLines(src))
	err = a.Run(prog)
	if err == nil {
		t.Fatal("expected type error for cam.Begn()")
	}
	if !strings.Contains(err.Error(), "unknown method") && !strings.Contains(strings.ToLower(err.Error()), "begn") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleCallKnownMethodOK(t *testing.T) {
	src := "cam = CAMERA.CREATE()\ncam.Begin()\ncam.End()\n"
	prog, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	a := DefaultAnalyzer("t.mb", parser.SplitLines(src))
	if err := a.Run(prog); err != nil {
		t.Fatalf("expected OK, got %v", err)
	}
}

func TestHandleCallUnknownReceiverSoft(t *testing.T) {
	// Dynamic / untyped handle: do not hard-fail unknown methods.
	src := "h = 0\nh.Begn()\n"
	prog, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	a := DefaultAnalyzer("t.mb", parser.SplitLines(src))
	if err := a.Run(prog); err != nil {
		t.Fatalf("unknown receiver should be soft, got %v", err)
	}
}

func TestHandleCallWrongArity(t *testing.T) {
	src := "cam = CAMERA.CREATE()\ncam.Move()\n"
	prog, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	a := DefaultAnalyzer("t.mb", parser.SplitLines(src))
	err = a.Run(prog)
	if err == nil {
		t.Fatal("expected arity error for cam.Move() with 0 args")
	}
}

func TestUnknownEasyModeGlobalRejected(t *testing.T) {
	src := "CreateCubee(1, 2, 3)\n"
	prog, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	a := DefaultAnalyzer("t.mb", parser.SplitLines(src))
	err = a.Run(prog)
	if err == nil {
		t.Fatal("expected unknown Easy Mode command error")
	}
	if !strings.Contains(err.Error(), "Unknown command") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCollectErrorsReportsMultiple(t *testing.T) {
	src := "CreateCubee()\nRENDER.SETFPS(\"Fast\")\n"
	prog, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	a := DefaultAnalyzer("t.mb", parser.SplitLines(src))
	a.CollectErrors = true
	err = a.Run(prog)
	if err == nil {
		t.Fatal("expected errors")
	}
	list := moonerrors.AsMoonErrors(err)
	if len(list) < 2 {
		t.Fatalf("expected >=2 structured MoonErrors, got %d: %v", len(list), err)
	}
	msg := err.Error()
	if !strings.Contains(msg, "Unknown command") {
		t.Fatalf("missing unknown command in: %v", err)
	}
}

func TestKeyPrefixPredeclared(t *testing.T) {
	src := "x = KEY_Z\n"
	prog, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	a := DefaultAnalyzer("t.mb", parser.SplitLines(src))
	if err := a.Run(prog); err != nil {
		t.Fatalf("KEY_Z should be predeclared, got %v", err)
	}
}
