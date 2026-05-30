package parser

import (
	"fmt"
	"strings"

	"moonbasic/compiler/arena"
	"moonbasic/compiler/ast"
	"moonbasic/compiler/symtable"
	"moonbasic/compiler/token"
)

func (p *Parser) parseTypeName() (string, error) {
	p.skipNewlines()
	switch p.cur().Type {
	case token.FLOAT:
		p.advance()
		return "FLOAT", nil
	case token.INT:
		p.advance()
		return "INTEGER", nil
	case token.STRING:
		p.advance()
		return "STRING", nil
	default:
		name, err := p.expectIdent()
		if err != nil {
			return "", err
		}
		return strings.ToUpper(name), nil
	}
}

func (p *Parser) parseParamList() ([]ast.Param, error) {
	var params []ast.Param
	p.skipNewlines()
	if p.cur().Type == token.RPAREN {
		return params, nil
	}
	for {
		pname, err := p.expectIdent()
		if err != nil {
			return nil, err
		}
		par := ast.Param{Name: pname}
		if p.cur().Type == token.AS {
			p.advance()
			th, err2 := p.parseTypeName()
			if err2 != nil {
				return nil, err2
			}
			par.TypeHint = th
		}
		params = append(params, par)
		p.skipNewlines()
		if p.cur().Type == token.COMMA {
			p.advance()
			p.skipNewlines()
			continue
		}
		break
	}
	return params, nil
}

func (p *Parser) parseReturnTypes() ([]string, error) {
	if p.cur().Type != token.AS {
		return nil, nil
	}
	p.advance()
	var types []string
	for {
		tn, err := p.parseTypeName()
		if err != nil {
			return nil, err
		}
		types = append(types, tn)
		if p.cur().Type == token.COMMA {
			p.advance()
			p.skipNewlines()
			continue
		}
		break
	}
	return types, nil
}

func (p *Parser) parseFunctionBody(savedSym *symtable.Table, name string, params []ast.Param) ([]ast.Stmt, error) {
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
	p.sym.DefineFunction(name)

	savedFn := p.FuncName
	p.FuncName = strings.ToUpper(name)
	defer func() { p.FuncName = savedFn }()

	var body []ast.Stmt
	for {
		p.skipNewlines()
		if p.cur().Type == token.ENDFUNCTION {
			break
		}
		if p.cur().Type == token.EOF {
			return nil, p.failf("unexpected EOF inside FUNCTION %s", name)
		}
		s, err2 := p.parseStmt()
		if err2 != nil {
			return nil, err2
		}
		if s != nil {
			body = append(body, s)
		}
		for p.cur().Type == token.COLON {
			p.advance()
			p.skipNewlines()
			s2, err3 := p.parseStmt()
			if err3 != nil {
				return nil, err3
			}
			if s2 != nil {
				body = append(body, s2)
			}
		}
	}
	p.advance() // ENDFUNCTION
	p.sym.PopScope()
	p.sym = savedSym
	p.sym.DefineFunction(name)
	return body, nil
}

// parseFunctionDef parses: FUNCTION name(params) [AS rettype] body ENDFUNCTION
func (p *Parser) parseFunctionDef() (*ast.FunctionDef, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance() // consume FUNCTION
	name, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	if err := p.expect(token.LPAREN); err != nil {
		return nil, err
	}
	params, err := p.parseParamList()
	if err != nil {
		return nil, err
	}
	if err := p.expect(token.RPAREN); err != nil {
		return nil, err
	}
	returnTypes, err := p.parseReturnTypes()
	if err != nil {
		return nil, err
	}
	if p.cur().Type != token.NEWLINE {
		return nil, p.failf("expected newline after FUNCTION header")
	}
	p.advance()

	savedSym := p.sym
	body, err := p.parseFunctionBody(savedSym, name, params)
	if err != nil {
		return nil, err
	}

	return arena.Make(p.ar, ast.FunctionDef{
		Name: name, Params: params, ReturnTypes: returnTypes, Body: body, Line: line, Col: col,
	}), nil
}

// parseCoroutineDef parses: COROUTINE varName body ENDCOROUTINE
// Lowers to a synthetic function plus varName = COROUTINE.START(@synthetic).
func (p *Parser) parseCoroutineDef() (*ast.FunctionDef, ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance() // COROUTINE
	varName, err := p.expectIdent()
	if err != nil {
		return nil, nil, err
	}
	if p.cur().Type != token.NEWLINE {
		return nil, nil, p.failf("expected newline after COROUTINE %s", varName)
	}
	p.advance()

	internalName := fmt.Sprintf("__co_%s_%d", strings.ToLower(varName), line)
	savedSym := p.sym

	p.sym = symtable.New()
	for k := range savedSym.Funcs() {
		p.sym.PredeclareFunction(k)
	}
	for k := range savedSym.Types() {
		p.sym.PredeclareType(k)
	}
	p.sym.PushScope()
	p.sym.DefineFunction(internalName)

	savedFn := p.FuncName
	p.FuncName = strings.ToUpper(internalName)
	defer func() { p.FuncName = savedFn }()

	var body []ast.Stmt
	for {
		p.skipNewlines()
		if p.cur().Type == token.ENDCOROUTINE {
			break
		}
		if p.cur().Type == token.EOF {
			return nil, nil, p.failf("unexpected EOF inside COROUTINE %s", varName)
		}
		s, err2 := p.parseStmt()
		if err2 != nil {
			return nil, nil, err2
		}
		if s != nil {
			body = append(body, s)
		}
		for p.cur().Type == token.COLON {
			p.advance()
			p.skipNewlines()
			s2, err3 := p.parseStmt()
			if err3 != nil {
				return nil, nil, err3
			}
			if s2 != nil {
				body = append(body, s2)
			}
		}
	}
	p.advance() // ENDCOROUTINE
	p.sym.PopScope()
	p.sym = savedSym
	p.sym.DefineFunction(internalName)

	fn := arena.Make(p.ar, ast.FunctionDef{
		Name: internalName, Body: body, Line: line, Col: col,
	})
	start := arena.Make(p.ar, ast.NamespaceCallExpr{
		NS: "COROUTINE", Method: "START",
		Args: []ast.Expr{arena.Make(p.ar, ast.FuncRefNode{Name: internalName, Line: line, Col: col})},
		Line: line, Col: col,
	})
	assign := arena.Make(p.ar, ast.AssignNode{Name: varName, Expr: start, Line: line, Col: col})
	return fn, assign, nil
}

// parseTypeDef parses: TYPE name FIELD ... ENDTYPE
// Field lines may be legacy comma-separated names, or "name AS type" / "name AS type(dim...)".
func (p *Parser) parseTypeDef() (*ast.TypeDef, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance() // consume TYPE
	tname, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	if p.cur().Type != token.NEWLINE {
		return nil, p.failf("expected newline after TYPE name")
	}
	p.advance()
	var fields []string
	var hints []string
	var arrayFlags []bool
	for {
		p.skipNewlines()
		if p.cur().Type == token.ENDTYPE {
			p.advance()
			break
		}
		if p.cur().Type == token.FIELD {
			p.advance()
		}
		if p.cur().Type != token.IDENT {
			return nil, p.failf("expected FIELD, identifier, or ENDTYPE, got %s", p.cur().Type.String())
		}
		fname, err2 := p.expectIdent()
		if err2 != nil {
			return nil, err2
		}
		if p.cur().Type == token.AS {
			p.advance()
			tn, err3 := p.expectIdent()
			if err3 != nil {
				return nil, err3
			}
			isArr := false
			if p.cur().Type == token.LPAREN {
				if _, err4 := p.parseArgList(); err4 != nil {
					return nil, err4
				}
				isArr = true
			}
			fields = append(fields, fname)
			hints = append(hints, strings.ToUpper(tn))
			arrayFlags = append(arrayFlags, isArr)
		} else {
			fields = append(fields, fname)
			hints = append(hints, "")
			arrayFlags = append(arrayFlags, false)
			for p.cur().Type == token.COMMA {
				p.advance()
				p.skipNewlines()
				fn2, err3 := p.expectIdent()
				if err3 != nil {
					return nil, err3
				}
				fields = append(fields, fn2)
				hints = append(hints, "")
				arrayFlags = append(arrayFlags, false)
			}
		}
		if p.cur().Type == token.NEWLINE {
			p.advance()
		}
	}
	p.sym.DefineType(tname)
	return arena.Make(p.ar, ast.TypeDef{Name: tname, Fields: fields, FieldTypeHints: hints, FieldIsArray: arrayFlags, Line: line, Col: col}), nil
}
