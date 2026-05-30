package parser

import (
	"fmt"

	"moonbasic/compiler/errors"
	"moonbasic/compiler/token"
)

func (p *Parser) recordParseError(err error) {
	if err == nil {
		return
	}
	p.errs = append(p.errs, err)
}

func (p *Parser) failfRecoverable(format string, args ...interface{}) error {
	t := p.cur()
	msg := fmt.Sprintf(format, args...)
	err := errors.NewParseError(p.file, t.Line, t.Col, msg, p.sourceLine(t.Line), "")
	p.recordParseError(err)
	return err
}

// syncStatement advances to the next likely statement boundary after a syntax error.
func (p *Parser) syncStatement() {
	for p.cur().Type != token.NEWLINE && p.cur().Type != token.EOF {
		p.advance()
	}
	for p.cur().Type != token.EOF {
		p.skipNewlines()
		if isSyncToken(p.cur().Type) || p.cur().Type == token.IDENT {
			return
		}
		p.advance()
	}
}

func isSyncToken(t token.TokenType) bool {
	switch t {
	case token.IF, token.WHILE, token.FOR, token.REPEAT, token.DO,
		token.SELECT, token.FUNCTION, token.TYPE, token.RETURN,
		token.DIM, token.LOCAL, token.GLOBAL, token.CONST, token.ENUM,
		token.INCLUDE, token.IMPORT, token.END, token.YIELD, token.COROUTINE:
		return true
	default:
		return false
	}
}
