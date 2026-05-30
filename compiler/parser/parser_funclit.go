package parser

import (
	"moonbasic/compiler/arena"
	"moonbasic/compiler/ast"
	"moonbasic/compiler/symtable"
	"moonbasic/compiler/token"
)

// parseFuncLit parses FUNCTION(params) body ENDFUNCTION as an expression.
func (p *Parser) parseFuncLit(line, col int) (ast.Expr, error) {
	p.advance() // FUNCTION
	params, body, err := p.parseAnonymousFunctionBody()
	if err != nil {
		return nil, err
	}
	return arena.Make(p.ar, ast.FuncLitNode{Params: params, Body: body, Line: line, Col: col}), nil
}

func (p *Parser) parseAnonymousFunctionBody() ([]ast.Param, []ast.Stmt, error) {
	if err := p.expect(token.LPAREN); err != nil {
		return nil, nil, err
	}
	var params []ast.Param
	p.skipNewlines()
	if p.cur().Type != token.RPAREN {
		for {
			pname, err := p.expectIdent()
			if err != nil {
				return nil, nil, err
			}
			params = append(params, ast.Param{Name: pname})
			p.skipNewlines()
			if p.cur().Type == token.COMMA {
				p.advance()
				p.skipNewlines()
				continue
			}
			break
		}
	}
	if err := p.expect(token.RPAREN); err != nil {
		return nil, nil, err
	}
	if p.cur().Type != token.NEWLINE {
		return nil, nil, p.failf("expected newline after FUNCTION header")
	}
	p.advance()

	savedSym := p.sym
	p.sym = symtable.New()
	for k := range savedSym.Funcs() {
		p.sym.PredeclareFunction(k)
	}
	for k := range savedSym.Types() {
		p.sym.PredeclareType(k)
	}
	p.sym.PushScope()
	for _, par := range params {
		p.sym.DefineParam(par.Name)
	}

	var body []ast.Stmt
	for {
		p.skipNewlines()
		if p.cur().Type == token.ENDFUNCTION {
			break
		}
		if p.cur().Type == token.EOF {
			return nil, nil, p.failf("unexpected EOF inside anonymous FUNCTION")
		}
		s, err := p.parseStmt()
		if err != nil {
			return nil, nil, err
		}
		if s != nil {
			body = append(body, s)
		}
		for p.cur().Type == token.COLON {
			p.advance()
			p.skipNewlines()
			s2, err2 := p.parseStmt()
			if err2 != nil {
				return nil, nil, err2
			}
			if s2 != nil {
				body = append(body, s2)
			}
		}
	}
	p.advance() // ENDFUNCTION
	p.sym.PopScope()
	p.sym = savedSym
	return params, body, nil
}
