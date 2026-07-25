package vm

import (
	"strings"

	"moonbasic/handlecall"
	"moonbasic/runtime"
)

func normalizeHandleMethod(mn string) string {
	return handlecall.NormalizeMethod(mn)
}

func handleCallRegistryPrefix(tag uint16) string {
	return handlecall.RegistryPrefix(tag)
}

func handleCallBuiltin(tag uint16, method string) (registryKey string, prependReceiver bool, ok bool) {
	return handlecall.Builtin(tag, method)
}

func handleCallDispatch(tag uint16, method string, argCount int) (registryKey string, prependReceiver bool, ok bool) {
	return handlecall.Dispatch(tag, method, argCount)
}

// HandleCallSuggestions lists common script-side method names for a handle type (error hints).
func HandleCallSuggestions(tag uint16) []string {
	return handlecall.HandleCallSuggestions(tag)
}

func filterRegistryKeysByPrefix(keys []string, prefix string) []string {
	pu := strings.ToUpper(prefix)
	var out []string
	for _, k := range keys {
		if strings.HasPrefix(strings.ToUpper(k), pu) {
			out = append(out, k)
		}
	}
	return out
}

// formatHandleCallError enriches a failed handle method dispatch with type-specific hints.
func (v *VM) formatHandleCallError(tag uint16, typeName, methodName, callKey string, mapped bool, err error) string {
	msg := err.Error()
	if mapped {
		return msg
	}
	prefix := handleCallRegistryPrefix(tag)
	keys := v.Registry.CommandKeys()
	prefixed := filterRegistryKeysByPrefix(keys, prefix)
	if alt, ok := runtime.BestSimilarCommand(callKey, prefixed, 3); ok {
		return msg + "\n  Did you mean " + alt + "?"
	}
	if sug := HandleCallSuggestions(tag); len(sug) > 0 {
		return msg + "\n  Hint: For this handle type use methods like " + strings.Join(sug, ", ") + "."
	}
	return msg + "\n  Hint: See docs/API_consistency.md for handle methods vs NS.COMMAND calls."
}
