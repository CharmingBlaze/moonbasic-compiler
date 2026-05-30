package parser

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"moonbasic/compiler/ast"
)

func TestParseReference(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	refPath := filepath.Join(filepath.Dir(file), "..", "..", "testdata", "reference.mbc")
	src, err := os.ReadFile(refPath)
	if err != nil {
		t.Fatal(err)
	}
	prog, err := ParseSource(refPath, string(src))
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Stmts) < 10 {
		t.Fatalf("expected many top-level stmts, got %d", len(prog.Stmts))
	}
	if len(prog.Functions) != 1 || prog.Functions[0].Name != "onland" {
		t.Fatalf("functions: %+v", prog.Functions)
	}
}

func TestParseSimpleAssign(t *testing.T) {
	prog, err := ParseSource("x.mbc", "x = 1\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Stmts) != 1 {
		t.Fatal(len(prog.Stmts))
	}
}

func TestSingleLineIf(t *testing.T) {
	_, err := ParseSource("x.mbc", "IF TRUE THEN x = 1\n")
	if err != nil {
		t.Fatal(err)
	}
}

func TestColonSeparatedStatements(t *testing.T) {
	prog, err := ParseSource("x.mbc", "x = 1 : y = 2 : z = 3\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Stmts) != 3 {
		t.Fatalf("want 3 stmts, got %d", len(prog.Stmts))
	}
}

func TestColonSeparatedStatementsInWhileBody(t *testing.T) {
	src := "WHILE TRUE\n" +
		"    vx = 0 : vz = 0\n" +
		"WEND\n"
	prog, err := ParseSource("x.mbc", src)
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("want 1 stmt, got %d", len(prog.Stmts))
	}
}

func TestUnknownNamespaceFallsBackToHandleCallExpr(t *testing.T) {
	prog, err := ParseSource("x.mbc", "s = SPHERE(1)\nd = s.X()\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Stmts) != 2 {
		t.Fatalf("want 2 stmts, got %d", len(prog.Stmts))
	}
	assign, ok := prog.Stmts[1].(*ast.AssignNode)
	if !ok {
		t.Fatalf("expected second stmt assign, got %T", prog.Stmts[1])
	}
	if _, ok := assign.Expr.(*ast.HandleCallExpr); !ok {
		t.Fatalf("expected s.X() to parse as HandleCallExpr, got %T", assign.Expr)
	}
}

func TestKnownNamespaceStaysNamespaceExpr(t *testing.T) {
	prog, err := ParseSource("x.mbc", "d = CAMERA.CREATE()\n")
	if err != nil {
		t.Fatal(err)
	}
	assign, ok := prog.Stmts[0].(*ast.AssignNode)
	if !ok {
		t.Fatalf("expected assign, got %T", prog.Stmts[0])
	}
	if _, ok := assign.Expr.(*ast.NamespaceCallExpr); !ok {
		t.Fatalf("expected CAMERA.CREATE() to parse as NamespaceCallExpr, got %T", assign.Expr)
	}
}

func TestIndexedHandleCallStatementParses(t *testing.T) {
	prog, err := ParseSource("x.mbc", "DIM a AS HANDLE(2)\na(1).Hide()\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Stmts) != 2 {
		t.Fatalf("want 2 stmts, got %d", len(prog.Stmts))
	}
	if _, ok := prog.Stmts[1].(*ast.HandleCallStmt); !ok {
		t.Fatalf("expected a(1).Hide() to parse as HandleCallStmt, got %T", prog.Stmts[1])
	}
}

func TestParseEnumAndForEachIn(t *testing.T) {
	src := "ENUM State\nIDLE\nWALK\nENDENUM\nDIM xs(3)\nFOR EACH x IN xs\nPRINT(x)\nNEXT\n"
	prog, err := ParseSource("x.mbc", src)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := prog.Stmts[0].(*ast.EnumDeclNode); !ok {
		t.Fatalf("stmt 0: want EnumDeclNode, got %T", prog.Stmts[0])
	}
	if _, ok := prog.Stmts[2].(*ast.ForInStmt); !ok {
		t.Fatalf("stmt 2: want ForInStmt, got %T", prog.Stmts[2])
	}
}

func TestParseInterpString(t *testing.T) {
	prog, err := ParseSource("x.mbc", "PRINT($\"hi {name}\")\n")
	if err != nil {
		t.Fatal(err)
	}
	printStmt, ok := prog.Stmts[0].(*ast.CallStmtNode)
	if !ok {
		t.Fatalf("want CallStmtNode, got %T", prog.Stmts[0])
	}
	if _, ok := printStmt.Args[0].(*ast.BinopNode); !ok {
		t.Fatalf("want desugared BinopNode, got %T", printStmt.Args[0])
	}
}

func TestParseErrorRecovery(t *testing.T) {
	// First statement is invalid; second should still parse.
	src := "PRNT 42\nPRINT(1)\n"
	prog, err := ParseSource("t.mb", src)
	if err == nil {
		t.Fatal("expected parse errors")
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("expected 1 recovered stmt, got %d", len(prog.Stmts))
	}
}

func TestParseMultiReturn(t *testing.T) {
	src := "FUNCTION GetPos()\nRETURN 1, 2, 3\nENDFUNCTION\n"
	prog, err := ParseSource("x.mbc", src)
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Functions) != 1 {
		t.Fatal("expected one function")
	}
	ret, ok := prog.Functions[0].Body[0].(*ast.ReturnNode)
	if !ok || len(ret.Exprs) != 3 {
		t.Fatalf("expected RETURN with 3 exprs, got %+v", prog.Functions[0].Body[0])
	}
}

func TestParseCoroutineBlock(t *testing.T) {
	src := "COROUTINE patrol\nYIELD\nENDCOROUTINE\n"
	prog, err := ParseSource("x.mbc", src)
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Functions) != 1 || len(prog.Stmts) != 1 {
		t.Fatalf("want 1 fn + 1 assign, got %d fn %d stmts", len(prog.Functions), len(prog.Stmts))
	}
	if !strings.HasPrefix(prog.Functions[0].Name, "__co_patrol_") {
		t.Fatalf("unexpected synthetic fn name %q", prog.Functions[0].Name)
	}
	assign, ok := prog.Stmts[0].(*ast.AssignNode)
	if !ok || assign.Name != "patrol" {
		t.Fatalf("want assign to patrol, got %+v", prog.Stmts[0])
	}
}

func TestParseTypedFunction(t *testing.T) {
	src := "FUNCTION Add(a AS FLOAT, b AS FLOAT) AS FLOAT\nRETURN 1.0\nENDFUNCTION\n"
	prog, err := ParseSource("x.mbc", src)
	if err != nil {
		t.Fatal(err)
	}
	f := prog.Functions[0]
	if len(f.Params) != 2 || f.Params[0].TypeHint != "FLOAT" || len(f.ReturnTypes) != 1 || f.ReturnTypes[0] != "FLOAT" {
		t.Fatalf("typed header not parsed: %+v", f)
	}
}
