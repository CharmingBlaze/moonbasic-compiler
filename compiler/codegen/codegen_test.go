package codegen

import (
	"strings"
	"testing"

	"moonbasic/compiler/parser"
	"moonbasic/vm/opcode"
)

func TestCompilePrintLiteral(t *testing.T) {
	src := "PRINT(42)\n"
	prog, err := parser.ParseSource("t.mbc", src)
	if err != nil {
		t.Fatal(err)
	}
	g := New("t.mbc", parser.SplitLines(src))
	out, err := g.Compile(prog)
	if err != nil {
		t.Fatal(err)
	}
	if out == nil || out.Main == nil {
		t.Fatal("nil program or main chunk")
	}
	code := out.Main.Instructions
	if len(code) < 2 {
		t.Fatalf("expected instructions, got %d", len(code))
	}
	if code[0].Op != opcode.OpPushInt || code[1].Op != opcode.OpCallBuiltin {
		t.Fatalf("unexpected ops: %v", code)
	}
}

func TestCompileLogicalOrAndXor(t *testing.T) {
	src := "a = TRUE: b = FALSE\nx = a OR b\ny = a AND b\nz = a XOR b\n"
	lines := parser.SplitLines(src)
	tree, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	g := New("t.mb", lines)
	out, err := g.Compile(tree)
	if err != nil {
		t.Fatal(err)
	}
	d := out.Main.Disassemble()
	for _, op := range []string{"OR", "AND", "XOR"} {
		if !strings.Contains(d, op) {
			t.Fatalf("disassembly missing %s:\n%s", op, d)
		}
	}
}

func TestCompileAssignAdd(t *testing.T) {
	src := "a = 2: b = 3\nx = a + b\nPRINT(x)\n"
	lines := parser.SplitLines(src)
	tree, err := parser.ParseSource("t.mbc", src)
	if err != nil {
		t.Fatal(err)
	}
	g := New("t.mbc", lines)
	out, err := g.Compile(tree)
	if err != nil {
		t.Fatal(err)
	}
	d := out.Main.Disassemble()
	if !strings.Contains(d, "HALT") || !strings.Contains(d, "ADD") {
		t.Log(d)
		t.Fatal("expected runtime add in disassembly")
	}
}

func TestImplicitGlobalSelfAssignUsesLoadStoreGlobal(t *testing.T) {
	// No VAR: globals persist in VM.Globals via OpLoadGlobal / OpStoreGlobal (not frame temps).
	src := "x = 0\nx = x + 1\n"
	lines := parser.SplitLines(src)
	tree, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	g := New("t.mb", lines)
	out, err := g.Compile(tree)
	if err != nil {
		t.Fatal(err)
	}
	d := out.Main.Disassemble()
	if !strings.Contains(d, "LOAD_GLOBAL") || !strings.Contains(d, "STORE_GLOBAL") {
		t.Fatalf("expected LOAD_GLOBAL and STORE_GLOBAL for implicit global x=x+1, got:\n%s", d)
	}
	loads := strings.Count(d, "LOAD_GLOBAL")
	stores := strings.Count(d, "STORE_GLOBAL")
	if loads < 1 || stores < 2 {
		t.Fatalf("expected at least 1 LOAD_GLOBAL (x+1 rhs) and 2 STORE_GLOBAL (x=0, x=x+1), loads=%d stores=%d\n%s", loads, stores, d)
	}
}

func TestCompileRoadmapFeatures(t *testing.T) {
	src := "ENUM State\nIDLE\nWALK\nENDENUM\nDIM xs(2)\nxs(1) = 10\nFUNCTION Pair()\nRETURN 1, 2\nENDFUNCTION\na, b = Pair()\nFOR EACH v IN xs\nPRINT(v)\nNEXT\nPRINT(State.WALK)\n"
	lines := parser.SplitLines(src)
	tree, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	g := New("t.mb", lines)
	out, err := g.Compile(tree)
	if err != nil {
		t.Fatal(err)
	}
	d := out.Main.Disassemble()
	for _, want := range []string{"ARRAY_LEN", "ARRAY_GET", "ARRAY_MAKE"} {
		if !strings.Contains(d, want) {
			t.Fatalf("disassembly missing %s:\n%s", want, d)
		}
	}
}

func TestCompileFuncRefAndEachType(t *testing.T) {
	src := "TYPE Enemy\nFIELD hp\nENDTYPE\nFUNCTION OnHit(a, b)\nENDFUNCTION\nFUNCTION Main()\nPHYSICS3D.ONCOLLISION(0, 0, @OnHit)\nFOR e = EACH(Enemy)\nDELETE e\nNEXT\nENDFUNCTION\n"
	lines := parser.SplitLines(src)
	tree, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	g := New("t.mb", lines)
	out, err := g.Compile(tree)
	if err != nil {
		t.Fatal(err)
	}
	d := out.Functions["main"].Disassemble()
	for _, want := range []string{"PUSH_FUNC_REF", "TYPE_INSTANCES", "ARRAY_GET"} {
		if !strings.Contains(d, want) {
			t.Fatalf("disassembly missing %s:\n%s", want, d)
		}
	}
}

func TestCompileFuncLitAndCallRef(t *testing.T) {
	src := "FUNCTION Main()\ncb = FUNCTION(x)\nRETURN x + 1\nENDFUNCTION\ny = cb(2)\nENDFUNCTION\n"
	lines := parser.SplitLines(src)
	tree, err := parser.ParseSource("t.mb", src)
	if err != nil {
		t.Fatal(err)
	}
	g := New("t.mb", lines)
	out, err := g.Compile(tree)
	if err != nil {
		t.Fatal(err)
	}
	d := out.Functions["main"].Disassemble()
	for _, want := range []string{"PUSH_FUNC_REF", "CALL_REF"} {
		if !strings.Contains(d, want) {
			t.Fatalf("disassembly missing %s:\n%s", want, d)
		}
	}
	if _, ok := out.Functions["__anon_1"]; !ok {
		t.Fatal("expected synthetic __anon_1 function chunk")
	}
}
