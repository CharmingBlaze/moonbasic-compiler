package vm

import "strings"

// NormalizeName folds identifiers (functions, globals, command keys) to uppercase.
func NormalizeName(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}
