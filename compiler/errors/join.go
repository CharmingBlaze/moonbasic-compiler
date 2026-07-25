package errors

import (
	"fmt"
	"strings"
)

// MultiError holds multiple compiler errors (semantic CollectErrors / --check / LSP).
// Unwrap returns the underlying errors so errors.As can find *MoonError values.
type MultiError struct {
	Errs []error
}

func (m *MultiError) Error() string {
	if m == nil || len(m.Errs) == 0 {
		return ""
	}
	if len(m.Errs) == 1 {
		return m.Errs[0].Error()
	}
	var b strings.Builder
	fmt.Fprintf(&b, "[moonBASIC] %d errors:\n\n", len(m.Errs))
	for i, e := range m.Errs {
		if i > 0 {
			b.WriteString("\n---\n\n")
		}
		b.WriteString(e.Error())
	}
	return b.String()
}

// Unwrap enables errors.As / errors.Is over each contained error (Go 1.20+).
func (m *MultiError) Unwrap() []error {
	if m == nil {
		return nil
	}
	return m.Errs
}

// Join formats multiple moonBASIC errors into one report (for multi-error compiler passes).
func Join(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}
	filtered := make([]error, 0, len(errs))
	for _, e := range errs {
		if e != nil {
			filtered = append(filtered, e)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	if len(filtered) == 1 {
		return filtered[0]
	}
	return &MultiError{Errs: filtered}
}

// AsMoonErrors extracts every *MoonError from err (including MultiError members).
func AsMoonErrors(err error) []*MoonError {
	if err == nil {
		return nil
	}
	var out []*MoonError
	var walk func(error)
	walk = func(e error) {
		if e == nil {
			return
		}
		if me, ok := e.(*MoonError); ok {
			out = append(out, me)
			return
		}
		if multi, ok := e.(*MultiError); ok {
			for _, inner := range multi.Errs {
				walk(inner)
			}
			return
		}
		// Single-level unwrap for fmt.Errorf("%w", me) style wrappers.
		type unwrapper interface{ Unwrap() error }
		if u, ok := e.(unwrapper); ok {
			walk(u.Unwrap())
		}
	}
	walk(err)
	return out
}
