package parser

import (
	"strings"

	"moonbasic/compiler/arena"
	"moonbasic/compiler/ast"
	"moonbasic/compiler/lexer"
	"moonbasic/compiler/token"
)

const interpSep = "\x1e"

// parseInterpString builds a + chain of string literals and STR()/FORMAT() calls.
func (p *Parser) parseInterpString(lit string, line, col int) (ast.Expr, error) {
	if lit == "" {
		return arena.Make(p.ar, ast.StringLitNode{Value: "", Line: line, Col: col}), nil
	}
	parts := strings.Split(lit, interpSep)
	if len(parts) == 1 {
		return arena.Make(p.ar, ast.StringLitNode{Value: parts[0], Line: line, Col: col}), nil
	}
	if len(parts)%2 == 0 {
		return nil, p.failf("invalid interpolated string payload")
	}

	var out ast.Expr
	for i := 0; i < len(parts); i++ {
		if i%2 == 0 {
			if parts[i] == "" && out == nil {
				continue
			}
			node := arena.Make(p.ar, ast.StringLitNode{Value: parts[i], Line: line, Col: col})
			out = p.concatInterp(out, node, line, col)
			continue
		}
		exprPart := parts[i]
		var fmtSuffix string
		if idx := strings.Index(exprPart, ":"); idx >= 0 {
			fmtSuffix = strings.TrimSpace(exprPart[idx+1:])
			exprPart = strings.TrimSpace(exprPart[:idx])
		}
		exprToks, err := lexer.Scan(p.file, exprPart)
		if err != nil {
			return nil, err
		}
		sub := &Parser{
			ar: p.ar, sym: p.sym, file: p.file, lines: p.lines,
			toks: exprToks, i: 0,
		}
		exprNode, err := sub.parseExpr()
		if err != nil {
			return nil, err
		}
		var piece ast.Expr
		if fmtSuffix != "" {
			pattern := arena.Make(p.ar, ast.StringLitNode{Value: "%" + fmtSuffix, Line: line, Col: col})
			piece = arena.Make(p.ar, ast.CallExprNode{
				Name: "FORMAT", Args: []ast.Expr{exprNode, pattern}, Line: line, Col: col,
			})
		} else {
			piece = arena.Make(p.ar, ast.CallExprNode{
				Name: "STR", Args: []ast.Expr{exprNode}, Line: line, Col: col,
			})
		}
		out = p.concatInterp(out, piece, line, col)
	}
	if out == nil {
		return arena.Make(p.ar, ast.StringLitNode{Value: "", Line: line, Col: col}), nil
	}
	return out, nil
}

func (p *Parser) concatInterp(left, right ast.Expr, line, col int) ast.Expr {
	if left == nil {
		return right
	}
	return arena.Make(p.ar, ast.BinopNode{Op: "+", Left: left, Right: right, Line: line, Col: col})
}

// parseReturnExprs parses comma-separated return values after the first expression.
func (p *Parser) parseReturnExprs(first ast.Expr, line, col int) ([]ast.Expr, error) {
	exprs := []ast.Expr{first}
	for p.cur().Type == token.COMMA {
		p.advance()
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, e)
	}
	return exprs, nil
}

func enumMemberName(enum, member string) string {
	return strings.ToUpper(enum) + "_" + strings.ToUpper(member)
}
