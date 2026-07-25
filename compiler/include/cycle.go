package include

import (
	"path/filepath"
	stdruntime "runtime"
	"strings"
)

// pathKey normalizes absolute paths for INCLUDE cycle / duplicate detection.
// On Windows, paths are compared case-insensitively so the same file is not
// expanded twice under different casing.
func pathKey(abs string) string {
	abs = filepath.Clean(abs)
	if stdruntime.GOOS == "windows" {
		return strings.ToLower(abs)
	}
	return abs
}

func stackContains(stack []string, abs string) bool {
	key := pathKey(abs)
	for _, p := range stack {
		if pathKey(p) == key {
			return true
		}
	}
	return false
}

func formatCircularChain(stack []string, closesTo string) string {
	parts := make([]string, 0, len(stack)+1)
	for _, p := range stack {
		parts = append(parts, filepath.Base(p))
	}
	parts = append(parts, filepath.Base(closesTo))
	return strings.Join(parts, " → ")
}
