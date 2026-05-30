package pkgmgr

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//go:embed default_index.json
var defaultIndexJSON []byte

//go:embed all:bundle/demo_extra
var demoExtraBundle embed.FS

const builtinURLPrefix = "moonbasic-builtin://"

// defaultRegistry is the bundled package index used when MOONBASIC_REGISTRY is unset.
type defaultRegistry struct {
	once  sync.Once
	index Index
	err   error
}

func (r *defaultRegistry) load() {
	r.once.Do(func() {
		if err := json.Unmarshal(defaultIndexJSON, &r.index); err != nil {
			r.err = fmt.Errorf("registry: default index: %w", err)
		}
	})
}

func (r *defaultRegistry) Lookup(name string) (IndexEntry, error) {
	r.load()
	if r.err != nil {
		return IndexEntry{}, r.err
	}
	e, ok := r.index.Packages[name]
	if !ok {
		return IndexEntry{}, fmt.Errorf("registry: unknown package %q (try: moonbasic install --list)", name)
	}
	return e, nil
}

// ListDefaultPackages returns entries from the bundled registry index.
func ListDefaultPackages() ([]struct {
	Name, Version, Description string
}, error) {
	dr := &defaultRegistry{}
	dr.load()
	if dr.err != nil {
		return nil, dr.err
	}
	var out []struct {
		Name, Version, Description string
	}
	for name, e := range dr.index.Packages {
		desc := e.Description
		if desc == "" {
			desc = e.URL
		}
		out = append(out, struct {
			Name, Version, Description string
		}{name, e.Version, desc})
	}
	return out, nil
}

func installFromBuiltin(name string) error {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "demo_extra":
		return installFromEmbedFS(demoExtraBundle, "bundle/demo_extra")
	default:
		return fmt.Errorf("install: unknown builtin package %q", name)
	}
}

func installFromEmbedFS(fsys embed.FS, root string) error {
	tmp, err := os.MkdirTemp("", "moonbasic-builtin-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	if err := extractEmbedDir(fsys, root, tmp); err != nil {
		return err
	}
	return installFromDir(tmp)
}

func extractEmbedDir(fsys embed.FS, root, dest string) error {
	return fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		out := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(out, 0o755)
		}
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return err
		}
		return os.WriteFile(out, data, 0o644)
	})
}
