package codegen

import (
	"fmt"
	"strings"

	"moonbasic/compiler/ast"
	"moonbasic/vm/opcode"
)

func (g *CodeGen) emitAnonFunction(params []ast.Param, body []ast.Stmt, line, col int) string {
	g.anonFnSeq++
	name := fmt.Sprintf("__anon_%d", g.anonFnSeq)
	if _, exists := g.Prog.Functions[name]; exists {
		return name
	}
	ch := opcode.NewChunk(name)
	g.Prog.Functions[name] = ch

	savedDepth := g.fnDepth
	savedName := g.funcName
	savedBase := g.baseReg
	savedNext := g.nextReg
	savedLoops := g.loopStack

	g.fnDepth = 1
	g.funcName = strings.ToUpper(name)
	g.loopStack = nil
	g.Symbols.PushScope()
	for _, p := range params {
		g.Symbols.DefineParam(p.Name)
	}
	g.predeclareFunctionBody(body)
	g.baseReg = uint8(g.Symbols.NextLocal())
	g.nextReg = g.baseReg
	for _, st := range body {
		g.emitStmt(ch, st)
		if g.err != nil {
			break
		}
	}
	ch.Emit(opcode.OpReturnVoid, 0, 0, 0, 0, line)

	g.Symbols.PopScope()
	g.fnDepth = savedDepth
	g.funcName = savedName
	g.baseReg = savedBase
	g.nextReg = savedNext
	g.loopStack = savedLoops
	return name
}

func (g *CodeGen) emitFuncLit(ch *opcode.Chunk, n *ast.FuncLitNode) uint8 {
	name := g.emitAnonFunction(n.Params, n.Body, n.Line, n.Col)
	idx := ch.AddName(name)
	dst := g.allocReg()
	ch.Emit(opcode.OpPushFuncRef, dst, 0, 0, idx, n.Line)
	g.nextReg = dst + 1
	return dst
}

func (g *CodeGen) emitCallRef(ch *opcode.Chunk, n *ast.CallRefExpr) uint8 {
	refReg := g.emitExpr(ch, n.Receiver)
	argStart := g.emitArgsStable(ch, n.Args, n.Line)
	dst := g.allocReg()
	ch.Emit(opcode.OpCallRef, dst, refReg, argStart, int32(len(n.Args)), n.Line)
	g.nextReg = dst + 1
	return dst
}

func (g *CodeGen) emitCallRefStmt(ch *opcode.Chunk, recv ast.Expr, args []ast.Expr, line, col int) {
	g.nextReg = g.baseReg
	n := &ast.CallRefExpr{Receiver: recv, Args: args, Line: line, Col: col}
	g.emitCallRef(ch, n)
	g.nextReg = g.baseReg
}
