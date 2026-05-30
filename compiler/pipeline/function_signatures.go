package pipeline

import (
	"fmt"
	"strings"

	"moonbasic/compiler/arena"
	"moonbasic/compiler/ast"
	"moonbasic/compiler/include"
	"moonbasic/compiler/parser"
)

// ParamSignature is one formal parameter with optional type hint.
type ParamSignature struct {
	Name     string `json:"name"`
	TypeHint string `json:"typeHint,omitempty"`
}

// FunctionSignature documents a user FUNCTION header for tooling (LSP hover).
type FunctionSignature struct {
	Name        string           `json:"name"`
	Params      []ParamSignature `json:"params"`
	ReturnTypes []string         `json:"returnTypes,omitempty"`
	Line        int              `json:"line"`
	Col         int              `json:"col"`
}

// FunctionSignatures parses source (with INCLUDE expansion) and returns user function signatures keyed by lowercase name.
func FunctionSignatures(name, src string) (map[string]FunctionSignature, error) {
	SyncPackageIncludeRoots()
	ar := arena.NewArena()
	defer ar.Reset()
	prog, err := parser.ParseSourceWithArena(name, src, ar)
	if err != nil {
		return nil, err
	}
	prog, err = include.ExpandWithArena(name, prog, ar)
	if err != nil {
		return nil, err
	}
	out := make(map[string]FunctionSignature, len(prog.Functions))
	for _, f := range prog.Functions {
		key := strings.ToLower(strings.TrimSpace(f.Name))
		if key == "" {
			continue
		}
		sig := FunctionSignature{
			Name:        f.Name,
			ReturnTypes: append([]string(nil), f.ReturnTypes...),
			Line:        f.Line,
			Col:         f.Col,
		}
		for _, p := range f.Params {
			sig.Params = append(sig.Params, ParamSignature{Name: p.Name, TypeHint: p.TypeHint})
		}
		out[key] = sig
	}
	return out, nil
}

// FormatFunctionSignature renders a moonBASIC FUNCTION header for display.
func FormatFunctionSignature(sig FunctionSignature) string {
	var b strings.Builder
	b.WriteString("FUNCTION ")
	b.WriteString(sig.Name)
	b.WriteString("(")
	for i, p := range sig.Params {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(p.Name)
		if p.TypeHint != "" {
			b.WriteString(" AS ")
			b.WriteString(p.TypeHint)
		}
	}
	b.WriteString(")")
	if len(sig.ReturnTypes) > 0 {
		b.WriteString(" AS ")
		for i, t := range sig.ReturnTypes {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(t)
		}
	}
	return b.String()
}

// FormatFunctionSignatureMarkdown returns LSP hover markdown for a user function.
func FormatFunctionSignatureMarkdown(sig FunctionSignature) string {
	hdr := FormatFunctionSignature(sig)
	var b strings.Builder
	fmt.Fprintf(&b, "### `%s`\n\n", hdr)
	fmt.Fprintf(&b, "User function (line %d).\n", sig.Line)
	if len(sig.ReturnTypes) > 0 {
		b.WriteString("\n**Returns:** ")
		b.WriteString(strings.Join(sig.ReturnTypes, ", "))
		b.WriteString("\n")
	}
	return b.String()
}

// ParseProgramForTooling returns the expanded AST for LSP helpers.
func ParseProgramForTooling(name, src string) (*ast.Program, error) {
	SyncPackageIncludeRoots()
	ar := arena.NewArena()
	defer ar.Reset()
	prog, err := parser.ParseSourceWithArena(name, src, ar)
	if err != nil {
		return nil, err
	}
	return include.ExpandWithArena(name, prog, ar)
}
