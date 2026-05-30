package semantic

import (
	"fmt"
	"strings"

	"moonbasic/compiler/ast"
	"moonbasic/compiler/types"
)

type funcSig struct {
	Params  []string // type hint per param; "" = any
	Returns []string // expected RETURN value types; empty = any
}

func (a *Analyzer) buildFuncSigs(prog *ast.Program) {
	a.funcSigs = make(map[string]funcSig)
	for _, f := range prog.Functions {
		key := strings.ToLower(f.Name)
		sig := funcSig{Returns: append([]string(nil), f.ReturnTypes...)}
		for _, p := range f.Params {
			sig.Params = append(sig.Params, p.TypeHint)
		}
		a.funcSigs[key] = sig
	}
}

func normalizeTypeHint(h string) (types.Tag, string, error) {
	h = strings.ToUpper(strings.TrimSpace(h))
	if h == "" {
		return types.Unknown, "", nil
	}
	switch h {
	case "INTEGER", "INT", "LONG":
		return types.Int, "", nil
	case "FLOAT", "SINGLE", "DOUBLE", "NUMBER":
		return types.Float, "", nil
	case "STRING", "STR":
		return types.String, "", nil
	case "BOOL", "BOOLEAN":
		return types.Bool, "", nil
	default:
		return types.UserType, h, nil
	}
}

func (a *Analyzer) checkFuncTypeHint(hint string, line, col int, what string) error {
	if hint == "" {
		return nil
	}
	tag, user, err := normalizeTypeHint(hint)
	if err != nil {
		return err
	}
	if tag == types.UserType {
		if _, ok := a.Types[user]; !ok {
			return a.typeError(line, col, fmt.Sprintf("unknown type %s in %s", user, what),
				"Use INTEGER, FLOAT, STRING, BOOL, or a defined TYPE name.")
		}
	}
	return nil
}

func (a *Analyzer) inferExprTag(e ast.Expr) types.Tag {
	switch n := e.(type) {
	case *ast.IntLitNode:
		return types.Int
	case *ast.FloatLitNode:
		return types.Float
	case *ast.StringLitNode:
		return types.String
	case *ast.BoolLitNode:
		return types.Bool
	case *ast.FuncRefNode, *ast.FuncLitNode:
		return types.FuncRef
	case *ast.IdentNode:
		if t := a.lookupVarType(n.Name); t != types.Unknown {
			return t
		}
		return types.Unknown
	case *ast.UnaryNode:
		if n.Op == "-" {
			if t := a.inferExprTag(n.Expr); t == types.Int || t == types.Float {
				return t
			}
		}
		if strings.EqualFold(n.Op, "NOT") {
			return types.Bool
		}
		return a.inferExprTag(n.Expr)
	case *ast.BinopNode:
		lt, rt := a.inferExprTag(n.Left), a.inferExprTag(n.Right)
		if strings.EqualFold(n.Op, "AND") || strings.EqualFold(n.Op, "OR") {
			return types.Bool
		}
		return types.Promote(lt, rt)
	case *ast.GroupedExpr:
		return a.inferExprTag(n.Inner)
	case *ast.CallExprNode:
		key := strings.ToLower(n.Name)
		if sig, ok := a.funcSigs[key]; ok && len(sig.Returns) == 1 {
			tag, _, err := normalizeTypeHint(sig.Returns[0])
			if err == nil {
				return tag
			}
		}
	}
	return types.Unknown
}

func tagsCompatible(got, want types.Tag) bool {
	if want == types.Unknown || got == types.Unknown {
		return true
	}
	if got == want {
		return true
	}
	if (got == types.Int && want == types.Float) || (got == types.Float && want == types.Int) {
		return want == types.Float
	}
	return false
}

func (a *Analyzer) checkUserCall(name string, args []ast.Expr, line, col int) error {
	key := strings.ToLower(name)
	sig, ok := a.funcSigs[key]
	if !ok {
		return nil
	}
	if len(sig.Params) > 0 && len(args) != len(sig.Params) {
		return a.typeError(line, col,
			fmt.Sprintf("function %s expects %d argument(s), got %d", name, len(sig.Params), len(args)),
			"Match the parameter list in the FUNCTION declaration.")
	}
	hasHints := false
	for _, hint := range sig.Params {
		if hint != "" {
			hasHints = true
			break
		}
	}
	if !hasHints {
		return nil
	}
	for i, hint := range sig.Params {
		if hint == "" {
			continue
		}
		wantTag, _, err := normalizeTypeHint(hint)
		if err != nil {
			return err
		}
		got := a.inferExprTag(args[i])
		if !tagsCompatible(got, wantTag) {
			return a.typeError(line, col,
				fmt.Sprintf("argument %d to %s: expected %s, got %s expression", i+1, name, hint, got.String()),
				"Pass a value matching the declared parameter type.")
		}
	}
	return nil
}

func (a *Analyzer) checkReturnTypes(fn *ast.FunctionDef, exprs []ast.Expr, line, col int) error {
	if len(fn.ReturnTypes) == 0 {
		return nil
	}
	if len(exprs) != len(fn.ReturnTypes) {
		return a.typeError(line, col,
			fmt.Sprintf("RETURN expects %d value(s) for %s, got %d", len(fn.ReturnTypes), fn.Name, len(exprs)),
			"Match the AS types on the FUNCTION header.")
	}
	for i, hint := range fn.ReturnTypes {
		wantTag, _, err := normalizeTypeHint(hint)
		if err != nil {
			return err
		}
		got := a.inferExprTag(exprs[i])
		if !tagsCompatible(got, wantTag) {
			return a.typeError(line, col,
				fmt.Sprintf("RETURN value %d: expected %s, got %s expression", i+1, hint, got.String()),
				"Return values must match the declared AS types.")
		}
	}
	return nil
}

func (a *Analyzer) validateFunctionSignatures(prog *ast.Program) error {
	for _, f := range prog.Functions {
		for i, p := range f.Params {
			if err := a.checkFuncTypeHint(p.TypeHint, f.Line, f.Col, fmt.Sprintf("parameter %d of %s", i+1, f.Name)); err != nil {
				return err
			}
		}
		for i, hint := range f.ReturnTypes {
			if err := a.checkFuncTypeHint(hint, f.Line, f.Col, fmt.Sprintf("return type %d of %s", i+1, f.Name)); err != nil {
				return err
			}
		}
	}
	return nil
}
