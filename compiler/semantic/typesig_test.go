package semantic

import (
	"strings"
	"testing"

	"moonbasic/compiler/parser"
)

func TestTypedFunctionSignature(t *testing.T) {
	src := "FUNCTION Bad() AS FLOAT\nRETURN \"no\"\nENDFUNCTION\n"
	a := DefaultAnalyzer("t.mb", strings.Split(src, "\n"))
	prog, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	if err := a.Run(prog); err == nil {
		t.Fatal("expected type error for RETURN string vs AS FLOAT")
	}
}

func TestTypedCallArity(t *testing.T) {
	src := "FUNCTION F(a AS INT, b AS INT) AS INT\nRETURN 1\nENDFUNCTION\nx = F(1)\n"
	a := DefaultAnalyzer("t.mb", strings.Split(src, "\n"))
	prog, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	if err := a.Run(prog); err == nil {
		t.Fatal("expected arity error for F(1)")
	}
}

func TestTypedCallInferredVar(t *testing.T) {
	src := `FUNCTION F(x AS INT) AS INT
RETURN x
ENDFUNCTION
n = 10
x = F(n)
`
	a := DefaultAnalyzer("t.mb", strings.Split(src, "\n"))
	prog, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	if err := a.Run(prog); err != nil {
		t.Fatalf("expected inferred INT for n=10: %v", err)
	}
}

func TestTypedCallWrongInferredVar(t *testing.T) {
	src := `FUNCTION F(x AS INT) AS INT
RETURN x
ENDFUNCTION
n = 1.5
x = F(n)
`
	a := DefaultAnalyzer("t.mb", strings.Split(src, "\n"))
	prog, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	if err := a.Run(prog); err == nil {
		t.Fatal("expected type error for F(float var)")
	}
}

func TestTypedReturnInferredExpr(t *testing.T) {
	src := `FUNCTION Sum(a AS FLOAT, b AS FLOAT) AS FLOAT
x = a + b
RETURN x
ENDFUNCTION
`
	a := DefaultAnalyzer("t.mb", strings.Split(src, "\n"))
	prog, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	if err := a.Run(prog); err != nil {
		t.Fatalf("expected inferred float return: %v", err)
	}
}
