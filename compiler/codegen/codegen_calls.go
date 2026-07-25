package codegen

import (
	"strings"

	"moonbasic/compiler/ast"
	"moonbasic/vm/opcode"
)

// emitArgsStable evaluates args left-to-right and packs their final values into
// contiguous registers starting at argStart. This prevents nested expressions
// from clobbering earlier argument registers.
func (g *CodeGen) emitArgsStable(ch *opcode.Chunk, args []ast.Expr, line int) uint8 {
	argStart := g.nextReg
	if argStart >= 256 {
		g.codegenError(line, 0, "too many temporary registers (max 256)", "Simplify nested call arguments.")
		return 0
	}
	nextArg := argStart
	for _, a := range args {
		if nextArg >= 256 || g.err != nil {
			if g.err == nil {
				g.codegenError(line, 0, "too many temporary registers (max 256)", "Simplify nested call arguments.")
			}
			return 0
		}
		r := g.emitExpr(ch, a)
		if g.err != nil {
			return 0
		}
		if int(r) != nextArg {
			ch.Emit(opcode.OpMove, uint8(nextArg), r, 0, 0, line)
		}
		nextArg++
		g.nextReg = nextArg
	}
	return uint8(argStart)
}

// emitCallStmt translates a call statement (builtin or user function) into bytecode.
func (g *CodeGen) emitCallStmt(ch *opcode.Chunk, n *ast.CallStmtNode) {
	if g.err != nil {
		return
	}
	g.nextReg = g.baseReg
	argStart := g.emitArgsStable(ch, n.Args, n.Line)

	fnKey := strings.ToUpper(n.Name)
	idx := ch.AddName(fnKey)
	op := opcode.OpCallBuiltin
	if _, ok := g.Prog.Functions[fnKey]; ok {
		op = opcode.OpCallUser
	}

	dst := g.allocReg() // discard result
	operand := (int32(len(n.Args)) << 24) | idx
	ch.Emit(op, dst, 0, argStart, operand, n.Line)
	g.nextReg = g.baseReg
}
