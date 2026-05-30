package include

import (
	"fmt"
	"os"
	"path/filepath"
)

// Resolve resolves an INCLUDE path relative to the host file's directory, or returns a clean absolute path.
func Resolve(hostFile, includePath string) (string, error) {
	if filepath.IsAbs(includePath) {
		return filepath.Clean(includePath), nil
	}
	dir := filepath.Dir(hostFile)
	if dir == "" || dir == "." {
		dir = "."
	}
	return filepath.Join(dir, filepath.FromSlash(includePath)), nil
}

// ResolvePackage finds a package entry file under configured package roots.
// Tries pkg/main.mb, pkg/index.mb, and pkg.mb.
func ResolvePackage(pkg string) (string, error) {
	pkg = filepath.FromSlash(pkg)
	candidates := []string{
		filepath.Join(pkg, "main.mb"),
		filepath.Join(pkg, "index.mb"),
		pkg + ".mb",
	}
	for _, rel := range candidates {
		if abs, err := tryOpenFromRoots(rel); err == nil {
			return abs, nil
		}
	}
	return "", fmt.Errorf("package not found: %s (tried main.mb, index.mb, %s.mb under package roots)", pkg, pkg)
}

func tryOpenFromRoots(rel string) (string, error) {
	rel = filepath.FromSlash(rel)
	for _, root := range packageRootsSnapshot() {
		abs := filepath.Join(root, rel)
		if st, err := os.Stat(abs); err == nil && !st.IsDir() {
			return abs, nil
		}
	}
	return "", fmt.Errorf("not found")
}
