package semantic

import (
	"fmt"
	"sort"
	"strings"

	"moonbasic/compiler/ast"
	"moonbasic/compiler/builtinmanifest"
	"moonbasic/compiler/entityspatial"
	"moonbasic/compiler/errors"
	"moonbasic/compiler/types"
)

// Analyzer performs constant folding and static checks after parsing.
type Analyzer struct {
	File  string
	Lines []string
	Table *builtinmanifest.Table
	Fold  bool

	// Static Analysis: caller -> set of callee names
	CallGraph   map[string]map[string]bool
	currentFunc string

	// Milestone 6: User-defined Types
	Types map[string]*ast.TypeDef

	funcNames map[string]bool // user FUNCTION names (uppercase)

	// Scopes for tracking assigned variables (Implicit Declaration/First-Assignment)
	// scopes[0] is always global.
	scopes []map[string]bool

	// typeScopes mirrors scopes with inferred static types for signature checking.
	typeScopes []map[string]types.Tag

	// handleScopes mirrors scopes with inferred heap handle tags for recv.METHOD checks.
	handleScopes []map[string]uint16

	// StrictDeprecated turns manifest migration aliases (MAKE, SETPOSITION, …) into type errors.
	StrictDeprecated bool

	// CollectErrors gathers multiple semantic errors instead of failing on the first.
	CollectErrors bool
	errors        []error

	deprecationNotices []DeprecationNotice
	deprecationSeen    map[string]bool

	enums map[string]map[string]int64 // ENUM name -> member -> value

	funcSigs   map[string]funcSig
	funcByName map[string]*ast.FunctionDef

	warnings    []SemanticWarning
	warningSeen map[string]bool
}

// DefaultAnalyzer uses the built-in command manifest and enables folding.
func DefaultAnalyzer(file string, lines []string) *Analyzer {
	return &Analyzer{
		File:  file,
		Lines: lines,
		Table: builtinmanifest.Default(),
		Fold:  true,
	}
}

// Run folds constants (if enabled) and type-checks namespace built-in calls.
func (a *Analyzer) Run(prog *ast.Program) error {
	a.CallGraph = make(map[string]map[string]bool)
	a.Types = make(map[string]*ast.TypeDef)
	a.enums = make(map[string]map[string]int64)
	a.currentFunc = "<MAIN>"

	a.collectEnums(prog.Stmts)

	if a.Fold {
		FoldConstants(prog)
	}
	a.scopes = []map[string]bool{make(map[string]bool)} // Global scope
	a.typeScopes = []map[string]types.Tag{make(map[string]types.Tag)}
	a.handleScopes = []map[string]uint16{make(map[string]uint16)}
	a.deprecationNotices = nil
	a.deprecationSeen = make(map[string]bool)
	a.warnings = nil
	a.warningSeen = make(map[string]bool)
	a.errors = nil
	err := a.checkProgram(prog)
	if a.CollectErrors && len(a.errors) > 0 {
		return errors.Join(a.errors...)
	}
	return err
}

// Warnings returns non-fatal semantic warnings (namespace shadowing, NOT/OR precedence, …).
func (a *Analyzer) Warnings() []SemanticWarning {
	return append([]SemanticWarning(nil), a.warnings...)
}

func (a *Analyzer) addWarning(line, col int, code, msg string) {
	key := fmt.Sprintf("%d:%d:%s:%s", line, col, code, msg)
	if a.warningSeen[key] {
		return
	}
	a.warningSeen[key] = true
	a.warnings = append(a.warnings, SemanticWarning{
		File: a.File, Line: line, Col: col, Code: code, Message: msg,
	})
}

// DeprecationNotices returns structured deprecation data for tooling (LSP, IDEs).
func (a *Analyzer) DeprecationNotices() []DeprecationNotice {
	return append([]DeprecationNotice(nil), a.deprecationNotices...)
}

// DeprecationWarnings returns compiler warnings collected during semantic analysis.
func (a *Analyzer) DeprecationWarnings() []string {
	notes := a.DeprecationNotices()
	out := make([]string, len(notes))
	for i, n := range notes {
		out[i] = n.String()
	}
	sort.Strings(out)
	return out
}

func (a *Analyzer) lineText(line int) string {
	if line < 1 || line > len(a.Lines) {
		return ""
	}
	return a.Lines[line-1]
}

func (a *Analyzer) typeError(line, col int, msg, hint string) error {
	err := errors.NewTypeError(a.File, line, col, msg, a.lineText(line), hint)
	if a.CollectErrors {
		a.errors = append(a.errors, err)
		return nil // continue walking the AST
	}
	return err
}

func (a *Analyzer) checkProgram(prog *ast.Program) error {
	a.funcNames = make(map[string]bool)
	for _, f := range prog.Functions {
		a.funcNames[strings.ToLower(f.Name)] = true
	}

	// 0. Register Types (Pass 0)
	for _, t := range prog.Types {
		if _, exists := a.Types[t.Name]; exists {
			return a.typeError(t.Line, t.Col, fmt.Sprintf("duplicate type definition %s", t.Name), "Remove or rename the duplicate type.")
		}
		a.Types[t.Name] = t
		// Verify fields are unique within type
		seen := make(map[string]bool)
		for _, f := range t.Fields {
			if seen[f] {
				return a.typeError(t.Line, t.Col, fmt.Sprintf("duplicate field %s in type %s", f, t.Name), "Ensure field names within a TYPE are unique.")
			}
			seen[f] = true
		}
		for i := range t.Fields {
			if i < len(t.FieldIsArray) && t.FieldIsArray[i] {
				return a.typeError(t.Line, t.Col, fmt.Sprintf("field %s: array-typed fields inside TYPE are not implemented yet", t.Fields[i]), "Use scalar fields only for now, or store a handle in a single field.")
			}
		}
	}

	a.buildFuncSigs(prog)
	a.funcByName = make(map[string]*ast.FunctionDef)
	for _, f := range prog.Functions {
		a.funcByName[strings.ToLower(f.Name)] = f
	}
	if err := a.validateFunctionSignatures(prog); err != nil {
		return err
	}

	a.currentFunc = "<MAIN>"
	for _, s := range prog.Stmts {
		if err := a.checkStmt(s); err != nil {
			return err
		}
	}

	// Functions
	for _, f := range prog.Functions {
		a.currentFunc = f.Name
		a.pushScope()
		// Parameters are implicitly assigned
		for _, p := range f.Params {
			a.assign(p.Name, f.Line, f.Col)
			if p.TypeHint != "" {
				if tag, _, err := normalizeTypeHint(p.TypeHint); err == nil {
					a.setVarType(p.Name, tag)
				}
			}
		}
		for _, s := range f.Body {
			if err := a.checkStmt(s); err != nil {
				return err
			}
		}
		a.popScope()
	}

	a.currentFunc = "<MAIN>"
	return nil
}

func (a *Analyzer) pushScope() {
	a.scopes = append(a.scopes, make(map[string]bool))
	a.typeScopes = append(a.typeScopes, make(map[string]types.Tag))
	a.handleScopes = append(a.handleScopes, make(map[string]uint16))
}

func (a *Analyzer) popScope() {
	if len(a.scopes) > 1 {
		a.scopes = a.scopes[:len(a.scopes)-1]
		a.typeScopes = a.typeScopes[:len(a.typeScopes)-1]
		a.handleScopes = a.handleScopes[:len(a.handleScopes)-1]
	}
}

func (a *Analyzer) setVarType(name string, tag types.Tag) {
	if tag == types.Unknown {
		return
	}
	name = strings.ToUpper(name)
	a.typeScopes[len(a.typeScopes)-1][name] = tag
}

func (a *Analyzer) lookupVarType(name string) types.Tag {
	name = strings.ToUpper(name)
	for i := len(a.typeScopes) - 1; i >= 0; i-- {
		if t, ok := a.typeScopes[i][name]; ok {
			return t
		}
	}
	return types.Unknown
}

func (a *Analyzer) assign(name string, line, col int) {
	name = strings.ToUpper(name)
	if a.Table != nil && a.Table.HasNamespace(name) {
		a.addWarning(line, col, "namespace-shadow",
			fmt.Sprintf("variable name %q shadows the %s.* namespace — built-in calls like %s.DELTA() may break; use a different name", strings.ToLower(name), name, name))
	}
	a.scopes[len(a.scopes)-1][name] = true
}

func (a *Analyzer) checkNotOrPrecedence(e ast.Expr) {
	b, ok := e.(*ast.BinopNode)
	if !ok || !strings.EqualFold(b.Op, "OR") {
		return
	}
	u, ok := b.Left.(*ast.UnaryNode)
	if !ok || !strings.EqualFold(u.Op, "NOT") {
		return
	}
	a.addWarning(b.Line, b.Col, "not-or-precedence",
		"expression NOT x OR y binds as (NOT x) OR y, not NOT (x OR y) — use NOT (x OR y) if you intend to negate the whole condition")
}

func (a *Analyzer) isAssigned(name string) bool {
	name = strings.ToUpper(name)
	if strings.HasPrefix(name, "KEY_") || strings.HasPrefix(name, "GAMEPAD_") {
		return true
	}
	for i := len(a.scopes) - 1; i >= 0; i-- {
		if a.scopes[i][name] {
			return true
		}
	}
	return false
}

func (a *Analyzer) checkStmt(s ast.Stmt) error {
	if ns, ok := s.(*ast.NamespaceCallStmt); ok {
		return a.checkNamespaceCall(ns.NS, ns.Method, ns.Args, ns.Line, ns.Col)
	}
	return a.walkStmtExprs(s)
}

func (a *Analyzer) walkStmtExprs(s ast.Stmt) error {
	switch n := s.(type) {
	case *ast.AssignNode:
		if err := a.checkExprCalls(n.Expr); err != nil {
			return err
		}
		a.assign(n.Name, n.Line, n.Col)
		a.setVarType(n.Name, a.inferExprTag(n.Expr))
		if ht, ok := a.inferHandleTag(n.Expr); ok {
			a.setHandleTag(n.Name, ht)
		}
		return nil
	case *ast.MultiAssignNode:
		if err := a.checkExprCalls(n.Expr); err != nil {
			return err
		}
		tag := a.inferExprTag(n.Expr)
		ht, htOK := a.inferHandleTag(n.Expr)
		for _, nm := range n.Names {
			a.assign(nm, n.Line, n.Col)
			a.setVarType(nm, tag)
			if htOK {
				a.setHandleTag(nm, ht)
			}
		}
		return nil
	case *ast.IndexAssignNode:
		if !a.isAssigned(n.Array) {
			return a.typeError(n.Line, n.Col, fmt.Sprintf("use of unassigned variable %s", n.Array), "Assign an array to the variable before subscripting.")
		}
		for _, e := range n.Index {
			if err := a.checkExprCalls(e); err != nil {
				return err
			}
		}
		return a.checkExprCalls(n.Expr)
	case *ast.IndexFieldAssignNode:
		if !a.isAssigned(n.Array) {
			return a.typeError(n.Line, n.Col, fmt.Sprintf("use of unassigned variable %s", n.Array), "Assign an array to the variable before subscripting.")
		}
		for _, e := range n.Index {
			if err := a.checkExprCalls(e); err != nil {
				return err
			}
		}
		return a.checkExprCalls(n.Expr)
	case *ast.FieldAssignNode:
		// Object.Field = expr
		if !a.isAssigned(n.Object) {
			return a.typeError(n.Line, n.Col, fmt.Sprintf("use of unassigned variable %s", n.Object), "Assign a value to the variable before accessing its fields.")
		}
		return a.checkExprCalls(n.Expr)
	case *ast.CallStmtNode:
		for _, e := range n.Args {
			if err := a.checkExprCalls(e); err != nil {
				return err
			}
		}
		if err := a.checkUserCall(n.Name, n.Args, n.Line, n.Col); err != nil {
			return err
		}
		return a.checkGlobalBuiltinCall(n.Name, n.Args, n.Line, n.Col)
	case *ast.HandleCallStmt:
		if err := a.checkExprCalls(n.Receiver); err != nil {
			return err
		}
		for _, e := range n.Args {
			if err := a.checkExprCalls(e); err != nil {
				return err
			}
		}
		if err := a.checkHandleCall(n.Receiver, n.Method, n.Args, n.Line, n.Col); err != nil {
			return err
		}
	case *ast.IfNode:
		a.checkNotOrPrecedence(n.Cond)
		if err := a.checkExprCalls(n.Cond); err != nil {
			return err
		}
		for _, t := range n.Then {
			if err := a.checkStmt(t); err != nil {
				return err
			}
		}
		for _, ei := range n.ElseIf {
			if err := a.checkExprCalls(ei.Cond); err != nil {
				return err
			}
			for _, t := range ei.Body {
				if err := a.checkStmt(t); err != nil {
					return err
				}
			}
		}
		for _, t := range n.Else {
			if err := a.checkStmt(t); err != nil {
				return err
			}
		}
	case *ast.WhileNode:
		a.checkNotOrPrecedence(n.Cond)
		if err := a.checkExprCalls(n.Cond); err != nil {
			return err
		}
		for _, t := range n.Body {
			if err := a.checkStmt(t); err != nil {
				return err
			}
		}
	case *ast.ForNode:
		a.assign(n.Var, n.Line, n.Col)
		a.setVarType(n.Var, types.Int) // Iterator is assigned
		for _, e := range []ast.Expr{n.From, n.To} {
			if err := a.checkExprCalls(e); err != nil {
				return err
			}
		}
		if n.Step != nil {
			if err := a.checkExprCalls(n.Step); err != nil {
				return err
			}
		}
		for _, t := range n.Body {
			if err := a.checkStmt(t); err != nil {
				return err
			}
		}
	case *ast.RepeatNode:
		for _, t := range n.Body {
			if err := a.checkStmt(t); err != nil {
				return err
			}
		}
		if err := a.checkExprCalls(n.Condition); err != nil {
			return err
		}
	case *ast.DoLoopNode:
		if err := a.checkExprCalls(n.Cond); err != nil {
			return err
		}
		for _, t := range n.Body {
			if err := a.checkStmt(t); err != nil {
				return err
			}
		}
	case *ast.ExitStmt, *ast.ContinueStmt:
		return nil
	case *ast.SelectNode:
		if err := a.checkExprCalls(n.Expr); err != nil {
			return err
		}
		for _, c := range n.Cases {
			if err := a.checkExprCalls(c.Value); err != nil {
				return err
			}
			for _, t := range c.Body {
				if err := a.checkStmt(t); err != nil {
					return err
				}
			}
		}
		for _, t := range n.Default {
			if err := a.checkStmt(t); err != nil {
				return err
			}
		}
	case *ast.ReturnNode:
		for _, e := range n.Exprs {
			if err := a.checkExprCalls(e); err != nil {
				return err
			}
		}
		if fn := a.funcByName[strings.ToLower(a.currentFunc)]; fn != nil {
			if err := a.checkReturnTypes(fn, n.Exprs, n.Line, n.Col); err != nil {
				return err
			}
		}
	case *ast.EnumDeclNode:
		return nil
	case *ast.ForInStmt:
		a.assign(n.Var, n.Line, n.Col)
		a.setVarType(n.Var, types.Float)
		if err := a.checkExprCalls(n.Array); err != nil {
			return err
		}
		for _, t := range n.Body {
			if err := a.checkStmt(t); err != nil {
				return err
			}
		}
	case *ast.DimNode:
		for _, e := range n.Dims {
			if err := a.checkExprCalls(e); err != nil {
				return err
			}
		}
		a.assign(n.Name, n.Line, n.Col)
		a.setVarType(n.Name, types.Array)
		return nil
	case *ast.ConstDeclNode:
		if a.currentFunc != "<MAIN>" {
			return a.typeError(n.Line, n.Col, "CONST is only allowed at module scope", "Move CONST to the top-level program, outside any FUNCTION.")
		}
		if err := a.checkExprCalls(n.Expr); err != nil {
			return err
		}
		a.assign(n.Name, n.Line, n.Col)
		a.setVarType(n.Name, a.inferExprTag(n.Expr))
		return nil
	case *ast.StaticDeclNode:
		if n.Init != nil {
			if err := a.checkExprCalls(n.Init); err != nil {
				return err
			}
			a.setVarType(n.Name, a.inferExprTag(n.Init))
		}
		a.assign(n.Name, n.Line, n.Col)
		return nil
	case *ast.SwapStmt, *ast.EraseStmt:
		return nil
	case *ast.LocalDeclNode:
		if n.Init != nil {
			if err := a.checkExprCalls(n.Init); err != nil {
				return err
			}
			a.setVarType(n.Name, a.inferExprTag(n.Init))
		}
		a.assign(n.Name, n.Line, n.Col)
		return nil
	case *ast.DeleteStmt:
		return a.checkExprCalls(n.Expr)
	case *ast.EachStmt:
		a.assign(n.Var, n.Line, n.Col)
		a.setVarType(n.Var, types.Handle)
		for _, t := range n.Body {
			if err := a.checkStmt(t); err != nil {
				return err
			}
		}
	case *ast.ExprStmt:
		return a.checkExprCalls(n.Expr)
	case *ast.NamespaceAssignNode:
		if err := a.checkEntitySpatialMacroArgs(n.NS, n.Method, n.Args, n.Line, n.Col); err != nil {
			return err
		}
		for _, arg := range n.Args {
			if err := a.checkExprCalls(arg); err != nil {
				return err
			}
		}
		return a.checkExprCalls(n.Expr)
	}
	return nil
}

func (a *Analyzer) checkExprCalls(e ast.Expr) error {
	switch n := e.(type) {
	case *ast.NamespaceCallExpr:
		if err := a.checkNamespaceCall(n.NS, n.Method, n.Args, n.Line, n.Col); err != nil {
			return err
		}
	case *ast.BinopNode:
		if err := a.checkExprCalls(n.Left); err != nil {
			return err
		}
		return a.checkExprCalls(n.Right)
	case *ast.UnaryNode:
		return a.checkExprCalls(n.Expr)
	case *ast.GroupedExpr:
		return a.checkExprCalls(n.Inner)
	case *ast.CallExprNode:
		if td, ok := a.Types[n.Name]; ok && !a.funcNames[n.Name] {
			if len(n.Args) != len(td.Fields) {
				return a.typeError(n.Line, n.Col,
					fmt.Sprintf("type %s constructor expects %d arguments, got %d", n.Name, len(td.Fields), len(n.Args)),
					"Pass one value per field in declaration order.")
			}
			for _, arg := range n.Args {
				if err := a.checkExprCalls(arg); err != nil {
					return err
				}
			}
			return nil
		}
		for _, arg := range n.Args {
			if err := a.checkExprCalls(arg); err != nil {
				return err
			}
		}
		if err := a.checkUserCall(n.Name, n.Args, n.Line, n.Col); err != nil {
			return err
		}
		if err := a.checkGlobalBuiltinCall(n.Name, n.Args, n.Line, n.Col); err != nil {
			return err
		}
	case *ast.IndexFieldExpr:
		for _, arg := range n.Index {
			if err := a.checkExprCalls(arg); err != nil {
				return err
			}
		}
	case *ast.HandleCallExpr:
		if err := a.checkExprCalls(n.Receiver); err != nil {
			return err
		}
		for _, arg := range n.Args {
			if err := a.checkExprCalls(arg); err != nil {
				return err
			}
		}
		if err := a.checkHandleCall(n.Receiver, n.Method, n.Args, n.Line, n.Col); err != nil {
			return err
		}
	case *ast.IndexExpr:
		if err := a.checkExprCalls(n.Base); err != nil {
			return err
		}
		for _, x := range n.Index {
			if err := a.checkExprCalls(x); err != nil {
				return err
			}
		}
	case *ast.IdentNode:
		if !a.isAssigned(n.Name) {
			return a.typeError(n.Line, n.Col, fmt.Sprintf("use of unassigned variable %s", n.Name), "Assign a value to the variable before using it in an expression.")
		}
	case *ast.NewNode:
		if _, exists := a.Types[n.TypeName]; !exists {
			return a.typeError(n.Line, n.Col, fmt.Sprintf("unknown type %s", n.TypeName), "Ensure the type is defined with TYPE ... END TYPE before use.")
		}
	}
	return nil
}

func (a *Analyzer) checkGlobalBuiltinCall(name string, args []ast.Expr, line, col int) error {
	if a.funcNames[strings.ToLower(name)] {
		return nil
	}
	key := builtinmanifest.NormalizeCommand(name)
	cmd, ok := a.Table.LookupGlobalArity(name, len(args))
	if !ok {
		ovs := a.Table.Overloads(key)
		if len(ovs) == 0 {
			msg := fmt.Sprintf("Unknown command '%s'", key)
			hint := "Use Namespace.Method form, or a registered Easy Mode global (see docs/EASY_MODE.md)."
			if alt, ok := a.Table.BestSimilarKey(key, 3); ok {
				hint = fmt.Sprintf("Did you mean %s ?", alt)
			}
			return a.typeError(line, col, msg, hint)
		}
		var parts []string
		for _, c := range ovs {
			parts = append(parts, fmt.Sprintf("%d", len(c.Args)))
		}
		return a.typeError(line, col,
			fmt.Sprintf("%s: no overload matches %d argument(s)", key, len(args)),
			fmt.Sprintf("Overloads expect argument count(s): %s.", strings.Join(parts, ", ")))
	}
	if msg := strings.TrimSpace(cmd.Stub); msg != "" {
		return a.typeError(line, col,
			fmt.Sprintf("command %s is not yet available in this release: %s", key, msg),
			"Remove the call or use an implemented alternative (see docs/reference/MIGRATION.md).")
	}
	return nil
}

func (a *Analyzer) checkNamespaceCall(ns, method string, args []ast.Expr, line, col int) error {
	if err := a.checkEntitySpatialMacroArgs(ns, method, args, line, col); err != nil {
		return err
	}
	for _, arg := range args {
		if err := a.checkExprCalls(arg); err != nil {
			return err
		}
	}
	if members, ok := a.enums[strings.ToUpper(ns)]; ok {
		if _, ok2 := members[strings.ToUpper(method)]; ok2 && len(args) == 0 {
			return nil
		}
	}

	cmd, ok := a.Table.LookupArity(ns, method, len(args))
	if !ok {
		if a.Table.Has(ns, method) {
			hint := a.Table.ArityHint(ns, method)
			return a.typeError(line, col,
				fmt.Sprintf("%s.%s: no overload matches %d argument(s)", ns, method, len(args)),
				hint)
		}
		msg, hint := unknownCommandMessageAndHint(a.Table, ns, method)
		return a.typeError(line, col, msg, hint)
	}
	key := builtinmanifest.Key(ns, method)
	if msg := strings.TrimSpace(cmd.Stub); msg != "" {
		return a.typeError(line, col,
			fmt.Sprintf("command %s is not yet available in this release: %s", key, msg),
			"Remove the call or use an implemented alternative (see docs/reference/MIGRATION.md).")
	}
	if repl, ok := a.Table.DeprecationReplacement(ns, method); ok {
		if a.StrictDeprecated {
			return a.typeError(line, col,
				fmt.Sprintf("deprecated command %s (strict: use %s)", key, repl),
				"Remove --strict-deprecated or migrate to the canonical name.")
		}
		dk := fmt.Sprintf("%d:%d:%s:%s", line, col, key, repl)
		if !a.deprecationSeen[dk] {
			a.deprecationSeen[dk] = true
			a.deprecationNotices = append(a.deprecationNotices, DeprecationNotice{
				File:           a.File,
				Line:           line,
				Col:            col,
				DeprecatedKey:  key,
				ReplacementKey: repl,
			})
		}
	}

	// Record CallGraph edge
	if _, exists := a.CallGraph[a.currentFunc]; !exists {
		a.CallGraph[a.currentFunc] = make(map[string]bool)
	}
	a.CallGraph[a.currentFunc][key] = true

	if len(args) != len(cmd.Args) {
		return a.typeError(line, col,
			fmt.Sprintf("%s.%s expects %d argument(s), got %d", ns, method, len(cmd.Args), len(args)),
			fmt.Sprintf("Provide %d argument(s) matching the built-in signature.", len(cmd.Args)))
	}
	for i, want := range cmd.Args {
		got := inferKind(args[i])
		if !compatible(want, got) {
			return a.typeError(line, col,
				fmt.Sprintf("%s.%s argument %d: expected %s, got %s", ns, method, i+1, kindName(want), formatGotKind(args[i])),
				"Fix the argument type to match the built-in signature.")
		}
	}
	return nil
}

func (a *Analyzer) collectEnums(stmts []ast.Stmt) {
	for _, s := range stmts {
		en, ok := s.(*ast.EnumDeclNode)
		if !ok {
			continue
		}
		enumName := strings.ToUpper(en.Name)
		if a.enums[enumName] == nil {
			a.enums[enumName] = make(map[string]int64)
		}
		for i, m := range en.Members {
			a.enums[enumName][strings.ToUpper(m)] = int64(i)
		}
	}
}
func (a *Analyzer) checkEntitySpatialMacroArgs(ns, method string, args []ast.Expr, line, col int) error {
	if !strings.EqualFold(ns, "ENTITY") {
		return nil
	}
	if _, ok := entityspatial.SpatialPropID(method); !ok {
		return nil
	}
	if len(args) < 1 {
		return nil
	}
	id, ok := entityspatial.ConstEntitySlotID(args[0])
	if !ok {
		return nil
	}
	if err := entityspatial.ValidateLiteralSlot(id); err != nil {
		return a.typeError(line, col, err.Error(), entityspatial.LiteralSlotHint())
	}
	return nil
}
