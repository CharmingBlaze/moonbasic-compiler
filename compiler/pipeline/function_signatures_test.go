package pipeline

import (
	"strings"
	"testing"
)

func TestFunctionSignaturesTyped(t *testing.T) {
	src := "FUNCTION Add(a AS FLOAT, b AS FLOAT) AS FLOAT\nRETURN a + b\nENDFUNCTION\n"
	sigs, err := FunctionSignatures("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	sig, ok := sigs["add"]
	if !ok {
		t.Fatal("missing add")
	}
	if len(sig.Params) != 2 || sig.Params[0].TypeHint != "FLOAT" || len(sig.ReturnTypes) != 1 {
		t.Fatalf("bad sig: %+v", sig)
	}
	got := FormatFunctionSignature(sig)
	want := "FUNCTION add(a AS FLOAT, b AS FLOAT) AS FLOAT"
	if got != want {
		t.Fatalf("format %q want %q", got, want)
	}
}

func TestFunctionSignaturesInclude(t *testing.T) {
	// No include in unit test — verify map keys are lowercase
	src := "FUNCTION Foo()\nENDFUNCTION\n"
	sigs, err := FunctionSignatures("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := sigs[strings.ToLower("Foo")]; !ok {
		t.Fatal("expected foo key")
	}
}
