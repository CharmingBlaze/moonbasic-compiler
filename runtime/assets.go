package runtime

import (
	"path/filepath"
	"strings"
)

// ResolveAssetPath maps a relative asset path to an absolute filesystem path.
// Relative paths resolve against SourceFile's directory plus optional AssetRelBase
// (set by ASSET.PATH). Absolute paths are returned unchanged.
func (r *Registry) ResolveAssetPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" || filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	if r == nil || r.SourceFile == "" {
		return filepath.FromSlash(path)
	}
	root := filepath.Dir(r.SourceFile)
	if r.AssetRelBase != "" {
		root = filepath.Join(root, filepath.FromSlash(strings.Trim(r.AssetRelBase, "/\\")))
	}
	return filepath.Clean(filepath.Join(root, filepath.FromSlash(path)))
}

// SetSourceFile records the primary .mb path for asset resolution and error messages.
func (r *Registry) SetSourceFile(absPath string) {
	if r == nil {
		return
	}
	r.SourceFile = filepath.Clean(absPath)
}

// SetAssetBase sets a subdirectory (relative to the source file) for asset loads.
func (r *Registry) SetAssetBase(rel string) {
	if r == nil {
		return
	}
	r.AssetRelBase = strings.TrimSpace(rel)
}
