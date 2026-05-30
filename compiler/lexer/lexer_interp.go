package lexer

import (
	"strings"

	"moonbasic/compiler/errors"
	"moonbasic/compiler/token"
)

// interpSep separates literal and expression segments in INTERP_STRING token Lit.
const interpSep = "\x1e"

// lexInterpString scans $"text {expr} more" after the leading $ has been consumed.
func (l *Lexer) lexInterpString() (token.Token, error) {
	startLine, startCol := l.line, l.col
	if l.peek() != '"' {
		return token.Token{}, errors.NewLexerError(l.file, l.line, l.col,
			"expected '\"' after $ for interpolated string",
			l.currentLineText(),
			"Use the form: $\"text {expression} more text\"")
	}
	l.advance() // opening "

	var parts []string
	var lit strings.Builder

	flushLit := func() {
		parts = append(parts, lit.String())
		lit.Reset()
	}

	for !l.eof() {
		c := l.peek()
		if c == '"' {
			l.advance()
			flushLit()
			return token.Token{
				Type: token.INTERP_STRING,
				Lit:  strings.Join(parts, interpSep),
				Line: startLine,
				Col:  startCol,
			}, nil
		}
		if c == '{' {
			l.advance()
			flushLit()
			exprStart := l.pos
			depth := 1
			for !l.eof() && depth > 0 {
				ch := l.peek()
				if ch == '{' {
					depth++
				} else if ch == '}' {
					depth--
					if depth == 0 {
						break
					}
				}
				l.advance()
			}
			if depth != 0 {
				return token.Token{}, errors.NewLexerError(l.file, startLine, startCol,
					"unclosed { in interpolated string",
					l.currentLineText(),
					"Close each {expression} with }")
			}
			exprSrc := strings.TrimSpace(l.input[exprStart:l.pos])
			l.advance() // closing }
			if exprSrc == "" {
				return token.Token{}, errors.NewLexerError(l.file, startLine, startCol,
					"empty {} in interpolated string",
					l.currentLineText(),
					"Put an expression inside { ... }")
			}
			parts = append(parts, exprSrc)
			continue
		}
		if c == '\\' {
			l.advance()
			if l.eof() {
				break
			}
			switch l.peek() {
			case '"':
				l.advance()
				lit.WriteByte('"')
			case '\\':
				l.advance()
				lit.WriteByte('\\')
			case '{', '}':
				b := l.advance()
				lit.WriteByte(b)
			default:
				lit.WriteByte('\\')
			}
			continue
		}
		if c == '\n' {
			return token.Token{}, errors.NewLexerError(l.file, startLine, startCol,
				"unterminated interpolated string",
				l.currentLineText(),
				"Close the string with \" before end of line.")
		}
		lit.WriteByte(l.advance())
	}
	return token.Token{}, errors.NewLexerError(l.file, startLine, startCol,
		"unterminated interpolated string",
		l.currentLineText(),
		"Close the string with \" before end of file.")
}
