package semantic

import "fmt"

// SemanticWarning is a non-fatal issue found during static analysis.
type SemanticWarning struct {
	File    string
	Line    int
	Col     int
	Code    string
	Message string
}

// String matches the CLI stderr warning format.
func (w SemanticWarning) String() string {
	if w.File != "" && w.Line > 0 {
		return fmt.Sprintf("%s:%d:%d: %s", w.File, w.Line, w.Col, w.Message)
	}
	return w.Message
}
