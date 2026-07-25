package semantic

import (
	"fmt"
	"strings"

	"moonbasic/compiler/ast"
	"moonbasic/compiler/builtinmanifest"
	"moonbasic/handlecall"
	"moonbasic/vm/heap"
)

func (a *Analyzer) setHandleTag(name string, tag uint16) {
	if tag == 0 || len(a.handleScopes) == 0 {
		return
	}
	name = strings.ToUpper(name)
	a.handleScopes[len(a.handleScopes)-1][name] = tag
}

func (a *Analyzer) lookupHandleTag(name string) (uint16, bool) {
	name = strings.ToUpper(name)
	for i := len(a.handleScopes) - 1; i >= 0; i-- {
		if t, ok := a.handleScopes[i][name]; ok && t != 0 {
			return t, true
		}
	}
	return 0, false
}

func (a *Analyzer) inferHandleTag(e ast.Expr) (uint16, bool) {
	switch n := e.(type) {
	case *ast.IdentNode:
		return a.lookupHandleTag(n.Name)
	case *ast.NamespaceCallExpr:
		m := strings.ToUpper(n.Method)
		ns := strings.ToUpper(n.NS)
		// Body builders: CREATE/MAKE return a builder; COMMIT returns a body.
		if ns == "BODY3D" && (m == "CREATE" || m == "MAKE") {
			return heap.TagPhysicsBuilder, true
		}
		if ns == "BODY2D" && (m == "CREATE" || m == "MAKE") {
			return heap.TagBody2D, true // builder shares BODY2D handle methods in dispatch
		}
		if ns == "AUDIO" {
			if strings.Contains(m, "MUSIC") {
				return heap.TagMusic, true
			}
			if strings.Contains(m, "SOUND") || strings.Contains(m, "SFX") {
				return heap.TagSound, true
			}
		}
		tag, ok := handlecall.TagFromNamespace(n.NS)
		if !ok {
			return 0, false
		}
		key := builtinmanifest.Key(n.NS, n.Method)
		if cmd, found := a.Table.FirstOverload(key); found && strings.EqualFold(strings.TrimSpace(cmd.Returns), "handle") {
			return tag, true
		}
		if strings.HasPrefix(m, "CREATE") || strings.HasPrefix(m, "LOAD") || m == "MAKE" {
			return tag, true
		}
		return 0, false
	case *ast.HandleCallExpr:
		recvTag, ok := a.inferHandleTag(n.Receiver)
		if !ok {
			return 0, false
		}
		key, _, mapped := handlecall.Dispatch(recvTag, n.Method, len(n.Args))
		if !mapped {
			return 0, false
		}
		if strings.EqualFold(n.Method, "COMMIT") {
			switch recvTag {
			case heap.TagPhysicsBuilder:
				return heap.TagPhysicsBody, true
			}
		}
		if cmd, found := a.Table.FirstOverload(key); found && strings.EqualFold(strings.TrimSpace(cmd.Returns), "handle") {
			return recvTag, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// checkHandleCall validates recv.METHOD(...) when the receiver handle type is known.
// Unknown receiver types are soft (no error) so dynamic handles still compile.
func (a *Analyzer) checkHandleCall(recv ast.Expr, method string, args []ast.Expr, line, col int) error {
	tag, known := a.inferHandleTag(recv)
	if !known {
		return nil
	}

	key, prepend, ok := handlecall.Dispatch(tag, method, len(args))
	if !ok {
		msg := fmt.Sprintf("unknown method %s on handle", strings.ToUpper(method))
		hint := "Use a method supported for this handle type (see docs/reference/UNIVERSAL_HANDLE_METHODS.md)."
		if sug := handlecall.HandleCallSuggestions(tag); len(sug) > 0 {
			// Prefer a close suggestion from the suggestion list.
			best, bestD := "", 99
			mu := strings.ToUpper(method)
			for _, s := range sug {
				d := builtinmanifest.EditDistance(mu, strings.ToUpper(s))
				if d < bestD {
					bestD = d
					best = s
				}
			}
			if bestD <= 3 && best != "" {
				hint = fmt.Sprintf("Did you mean .%s() ?", best)
			} else {
				hint = fmt.Sprintf("For this handle type use methods like %s.", strings.Join(sug, ", "))
			}
		}
		return a.typeError(line, col, msg, hint)
	}

	ns, meth, splitOK := handlecall.SplitRegistryKey(key)
	if !splitOK {
		return nil
	}
	argc := len(args)
	if prepend {
		argc++
	}
	cmd, found := a.Table.LookupArity(ns, meth, argc)
	if !found {
		if a.Table.Has(ns, meth) {
			hint := a.Table.ArityHint(ns, meth)
			scriptHint := hint
			if prepend {
				scriptHint = fmt.Sprintf("%s (handle method args exclude the receiver)", hint)
			}
			return a.typeError(line, col,
				fmt.Sprintf("%s.%s: no overload matches %d handle-method argument(s)", ns, meth, len(args)),
				scriptHint)
		}
		// Dispatch mapped to a key missing from the manifest — treat as OK (runtime may still work).
		return nil
	}
	if msg := strings.TrimSpace(cmd.Stub); msg != "" {
		return a.typeError(line, col,
			fmt.Sprintf("command %s is not yet available in this release: %s", key, msg),
			"Remove the call or use an implemented alternative (see docs/reference/MIGRATION.md).")
	}
	if repl, ok := a.Table.DeprecationReplacement(ns, meth); ok {
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
	if _, exists := a.CallGraph[a.currentFunc]; !exists {
		a.CallGraph[a.currentFunc] = make(map[string]bool)
	}
	a.CallGraph[a.currentFunc][key] = true

	// Kind-check script args against manifest args after an optional leading handle.
	manifestArgs := cmd.Args
	if prepend && len(manifestArgs) > 0 {
		manifestArgs = manifestArgs[1:]
	}
	if len(args) != len(manifestArgs) {
		return a.typeError(line, col,
			fmt.Sprintf("%s expects %d argument(s), got %d", key, len(manifestArgs), len(args)),
			fmt.Sprintf("Provide %d argument(s) matching the built-in signature.", len(manifestArgs)))
	}
	for i, want := range manifestArgs {
		got := inferKind(args[i])
		if !compatible(want, got) {
			return a.typeError(line, col,
				fmt.Sprintf("%s argument %d: expected %s, got %s", key, i+1, kindName(want), formatGotKind(args[i])),
				"Fix the argument type to match the built-in signature.")
		}
	}
	return nil
}
