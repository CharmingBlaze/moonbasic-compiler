package lsp

import (
	"encoding/json"
	"strings"

	"moonbasic/compiler/pipeline"
)

func (s *server) funcSigsFor(uri, text string) map[string]pipeline.FunctionSignature {
	s.muSig.Lock()
	defer s.muSig.Unlock()
	if s.funcSigCache == nil {
		s.funcSigCache = make(map[string]map[string]pipeline.FunctionSignature)
	}
	if cached, ok := s.funcSigCache[uri]; ok && s.funcSigSrc[uri] == text {
		return cached
	}
	path := filePathFromURI(uri)
	if path == "" {
		path = "buffer.mb"
	}
	sigs, err := pipeline.FunctionSignatures(path, text)
	if err != nil {
		s.funcSigCache[uri] = nil
	} else {
		s.funcSigCache[uri] = sigs
	}
	if s.funcSigSrc == nil {
		s.funcSigSrc = make(map[string]string)
	}
	s.funcSigSrc[uri] = text
	return sigs
}

func (s *server) invalidateFuncSig(uri string) {
	s.muSig.Lock()
	defer s.muSig.Unlock()
	delete(s.funcSigCache, uri)
	delete(s.funcSigSrc, uri)
}

func (s *server) userFunctionHover(text, uri, line string, col int) any {
	word := identifierAt(line, col)
	if word == "" {
		return nil
	}
	sigs := s.funcSigsFor(uri, text)
	sig, ok := sigs[strings.ToLower(word)]
	if !ok {
		return nil
	}
	return map[string]any{
		"contents": map[string]any{
			"kind":  "markdown",
			"value": pipeline.FormatFunctionSignatureMarkdown(sig),
		},
	}
}

func (s *server) signatureHelp(params json.RawMessage) any {
	var p struct {
		TextDocument struct {
			URI string `json:"uri"`
		} `json:"textDocument"`
		Position struct {
			Line      int `json:"line"`
			Character int `json:"character"`
		} `json:"position"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil
	}
	text := s.docs[p.TextDocument.URI]
	if text == "" {
		return nil
	}
	lines := strings.Split(text, "\n")
	if p.Position.Line < 0 || p.Position.Line >= len(lines) {
		return nil
	}
	line := lines[p.Position.Line]
	fnName, argIndex, ok := callContextAt(line, p.Position.Character)
	if !ok {
		return nil
	}
	sigs := s.funcSigsFor(p.TextDocument.URI, text)
	sig, ok := sigs[strings.ToLower(fnName)]
	if !ok {
		return nil
	}
	label := pipeline.FormatFunctionSignature(sig)
	var parameters []map[string]any
	for _, par := range sig.Params {
		lbl := par.Name
		if par.TypeHint != "" {
			lbl += " AS " + par.TypeHint
		}
		parameters = append(parameters, map[string]any{"label": lbl})
	}
	active := argIndex
	if active >= len(parameters) {
		active = len(parameters) - 1
	}
	if active < 0 {
		active = 0
	}
	return map[string]any{
		"signatures": []map[string]any{{
			"label":         label,
			"documentation": map[string]any{"kind": "markdown", "value": "User-defined function"},
			"parameters":    parameters,
		}},
		"activeSignature":  0,
		"activeParameter":  active,
	}
}

// callContextAt returns the callee name and 0-based argument index when cursor is inside a call.
func callContextAt(line string, col int) (fnName string, argIndex int, ok bool) {
	if col < 0 || col > len(line) {
		return "", 0, false
	}
	// Walk backward to find '(' then identifier before it.
	depth := 0
	inStr := false
	for i := col - 1; i >= 0; i-- {
		c := line[i]
		if c == '"' {
			inStr = !inStr
			continue
		}
		if inStr {
			continue
		}
		switch c {
		case ')':
			depth++
		case '(':
			if depth > 0 {
				depth--
				continue
			}
			// Found call open paren; identifier immediately before (allow spaces).
			j := i - 1
			for j >= 0 && (line[j] == ' ' || line[j] == '\t') {
				j--
			}
			end := j + 1
			for j >= 0 && isIdentRune(rune(line[j])) {
				j--
			}
			fnName = strings.TrimSpace(line[j+1 : end])
			if fnName == "" {
				return "", 0, false
			}
			// Count commas at depth 0 between '(' and cursor.
			open := i + 1
			if col < open {
				return fnName, 0, true
			}
			seg := line[open:col]
			depth = 0
			inStr = false
			commas := 0
			for k := 0; k < len(seg); k++ {
				ch := seg[k]
				if ch == '"' {
					inStr = !inStr
					continue
				}
				if inStr {
					continue
				}
				switch ch {
				case '(', '[':
					depth++
				case ')', ']':
					if depth > 0 {
						depth--
					}
				case ',':
					if depth == 0 {
						commas++
					}
				}
			}
			if strings.TrimSpace(seg) == "" {
				return fnName, commas, true
			}
			return fnName, commas, true
		}
	}
	return "", 0, false
}
