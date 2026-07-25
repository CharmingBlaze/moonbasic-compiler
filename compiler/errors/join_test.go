package errors

import (
	stderrors "errors"
	"strings"
	"testing"
)

func TestJoinPreservesMoonErrors(t *testing.T) {
	e1 := NewTypeError("a.mb", 1, 1, "first", "x", "")
	e2 := NewTypeError("a.mb", 2, 1, "second", "y", "")
	joined := Join(e1, e2)
	if joined == nil {
		t.Fatal("expected multi error")
	}
	list := AsMoonErrors(joined)
	if len(list) != 2 {
		t.Fatalf("AsMoonErrors: got %d, want 2 (%v)", len(list), joined)
	}
	if list[0].Message != "first" || list[1].Message != "second" {
		t.Fatalf("messages: %#v", list)
	}
	var me *MoonError
	if !stderrors.As(joined, &me) {
		t.Fatal("errors.As should find a *MoonError via MultiError.Unwrap")
	}
	msg := joined.Error()
	if !strings.Contains(msg, "2 errors") || !strings.Contains(msg, "first") {
		t.Fatalf("joined text: %s", msg)
	}
}

func TestJoinSingle(t *testing.T) {
	e1 := NewParseError("a.mb", 1, 1, "only", "x", "")
	joined := Join(e1)
	if joined != e1 {
		t.Fatalf("single Join should return original: %T", joined)
	}
}
