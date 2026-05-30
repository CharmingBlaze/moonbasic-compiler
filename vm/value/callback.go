package value

import "strings"

// CallbackName resolves a function callback argument (string name or @func reference).
// names is the current chunk's Names table (function keys are lowercase).
func CallbackName(v Value, names []string) (string, bool) {
	switch v.Kind {
	case KindString:
		s := v.String()
		if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
			return "", false
		}
		return s, s != ""
	case KindFunc:
		idx := int32(v.IVal)
		if idx < 0 || int(idx) >= len(names) {
			return "", false
		}
		return names[idx], true
	default:
		return "", false
	}
}
