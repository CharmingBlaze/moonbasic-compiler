package errors

import (
	"fmt"
	"strings"
)

// Join formats multiple moonBASIC errors into one report (for multi-error compiler passes).
func Join(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}
	if len(errs) == 1 {
		return errs[0]
	}
	var b strings.Builder
	fmt.Fprintf(&b, "[moonBASIC] %d errors:\n\n", len(errs))
	for i, e := range errs {
		if i > 0 {
			b.WriteString("\n---\n\n")
		}
		b.WriteString(e.Error())
	}
	return fmt.Errorf("%s", b.String())
}
