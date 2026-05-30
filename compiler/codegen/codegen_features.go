package codegen

import (
	"fmt"
	"strings"

	"moonbasic/compiler/ast"
	"moonbasic/vm/opcode"
)

func (g *CodeGen) collectAllEnums(tree *ast.Program) {
	g.walkEnumStmts(tree.Stmts)
	for _, fn := range tree.Functions {
		g.walkEnumStmts(fn.Body)
	}
}

func (g *CodeGen) walkEnumStmts(stmts []ast.Stmt) {
	for _, s := range stmts {
		en, ok := s.(*ast.EnumDeclNode)
		if !ok {
			continue
		}
		enumName := strings.ToUpper(en.Name)
		for i, m := range en.Members {
			member := strings.ToUpper(m)
			constName := enumName + "_" + member
			val := int64(i)
			g.enumValues[constName] = val
			g.enumValues[enumName+"."+member] = val
		}
	}
}

func (g *CodeGen) emitEnumDecl(ch *opcode.Chunk, n *ast.EnumDeclNode) {
	g.nextReg = g.baseReg
	enumName := strings.ToUpper(n.Name)
	for i, m := range n.Members {
		member := strings.ToUpper(m)
		constName := enumName + "_" + member
		valReg := g.allocReg()
		ch.Emit(opcode.OpPushInt, valReg, 0, 0, int32(i), n.Line)
		idx := ch.AddName(constName)
		ch.Emit(opcode.OpStoreGlobal, 0, valReg, 0, idx, n.Line)
		val := int64(i)
		g.enumValues[constName] = val
		g.enumValues[enumName+"."+member] = val
	}
	g.nextReg = g.baseReg
}

func (g *CodeGen) emitReturn(ch *opcode.Chunk, n *ast.ReturnNode) {
	g.nextReg = g.baseReg
	switch len(n.Exprs) {
	case 0:
		ch.Emit(opcode.OpReturnVoid, 0, 0, 0, 0, n.Line)
	case 1:
		r := g.emitExpr(ch, n.Exprs[0])
		ch.Emit(opcode.OpReturn, 0, r, 0, 1, n.Line)
	default:
		count := len(n.Exprs)
		dimReg := g.allocReg()
		ch.Emit(opcode.OpPushInt, dimReg, 0, 0, int32(count), n.Line)
		arrReg := g.allocReg()
		ch.Emit(opcode.OpArrayMake, arrReg, dimReg, 0, 1, n.Line)
		for i, e := range n.Exprs {
			valReg := g.emitExpr(ch, e)
			idxReg := g.allocReg()
			ch.Emit(opcode.OpPushInt, idxReg, 0, 0, int32(i+1), n.Line)
			ch.Emit(opcode.OpArraySet, valReg, arrReg, idxReg, 1, n.Line)
		}
		ch.Emit(opcode.OpReturn, 0, arrReg, 0, 1, n.Line)
	}
}

func (g *CodeGen) emitForIn(ch *opcode.Chunk, n *ast.ForInStmt) {
	g.nextReg = g.baseReg
	arrReg := g.emitExpr(ch, n.Array)

	lenReg := g.allocReg()
	ch.Emit(opcode.OpArrayLen, lenReg, arrReg, 0, 0, n.Line)

	iReg := g.allocReg()
	ch.Emit(opcode.OpPushInt, iReg, 0, 0, 1, n.Line)

	startIdx := len(ch.Instructions)

	cmpReg := g.allocReg()
	ch.Emit(opcode.OpLte, cmpReg, iReg, lenReg, 0, n.Line)
	exitJump := ch.Emit(opcode.OpJumpIfFalse, 0, cmpReg, 0, 0, n.Line)

	idxReg := g.allocReg()
	ch.Emit(opcode.OpMove, idxReg, iReg, 0, 0, n.Line)
	valReg := g.allocReg()
	ch.Emit(opcode.OpArrayGet, valReg, arrReg, idxReg, 1, n.Line)
	g.emitStoreNamed(ch, n.Var, n.Line, valReg)

	g.beginLoop("for", int32(startIdx))
	for _, st := range n.Body {
		g.emitStmt(ch, st)
	}

	oneReg := g.allocReg()
	ch.Emit(opcode.OpPushInt, oneReg, 0, 0, 1, n.Line)
	ch.Emit(opcode.OpAdd, iReg, iReg, oneReg, 0, n.Line)
	ch.Emit(opcode.OpJump, 0, 0, 0, int32(startIdx), n.Line)

	leave := len(ch.Instructions)
	g.endLoop(ch, leave)
	ch.Instructions[exitJump].Operand = int32(leave)
	g.nextReg = g.baseReg
}

func (g *CodeGen) emitEach(ch *opcode.Chunk, n *ast.EachStmt) {
	g.nextReg = g.baseReg
	typeName := strings.ToUpper(strings.TrimSpace(n.TypeName))
	if _, ok := g.Prog.Types[typeName]; !ok {
		g.codegenError(n.Line, n.Col, fmt.Sprintf("unknown TYPE %s for EACH", n.TypeName), "Declare TYPE before FOR ... = EACH(Type).")
		return
	}
	tidx := ch.AddName(typeName)
	arrReg := g.allocReg()
	ch.Emit(opcode.OpTypeInstances, arrReg, 0, 0, tidx, n.Line)

	lenReg := g.allocReg()
	ch.Emit(opcode.OpArrayLen, lenReg, arrReg, 0, 0, n.Line)

	iReg := g.allocReg()
	ch.Emit(opcode.OpPushInt, iReg, 0, 0, 1, n.Line)

	startIdx := len(ch.Instructions)

	cmpReg := g.allocReg()
	ch.Emit(opcode.OpLte, cmpReg, iReg, lenReg, 0, n.Line)
	exitJump := ch.Emit(opcode.OpJumpIfFalse, 0, cmpReg, 0, 0, n.Line)

	idxReg := g.allocReg()
	ch.Emit(opcode.OpMove, idxReg, iReg, 0, 0, n.Line)
	valReg := g.allocReg()
	ch.Emit(opcode.OpArrayGet, valReg, arrReg, idxReg, 1, n.Line)
	g.emitStoreNamed(ch, n.Var, n.Line, valReg)

	g.beginLoop("for", int32(startIdx))
	for _, st := range n.Body {
		g.emitStmt(ch, st)
	}

	oneReg := g.allocReg()
	ch.Emit(opcode.OpPushInt, oneReg, 0, 0, 1, n.Line)
	ch.Emit(opcode.OpAdd, iReg, iReg, oneReg, 0, n.Line)
	ch.Emit(opcode.OpJump, 0, 0, 0, int32(startIdx), n.Line)

	leave := len(ch.Instructions)
	g.endLoop(ch, leave)
	ch.Instructions[exitJump].Operand = int32(leave)
	g.nextReg = g.baseReg
}

func (g *CodeGen) tryEmitEnumMember(ch *opcode.Chunk, ns, method string, argCount, line int) (uint8, bool) {
	if argCount != 0 {
		return 0, false
	}
	key := strings.ToUpper(ns) + "." + strings.ToUpper(method)
	val, ok := g.enumValues[key]
	if !ok {
		return 0, false
	}
	dst := g.allocReg()
	ch.Emit(opcode.OpPushInt, dst, 0, 0, int32(val), line)
	return dst, true
}
