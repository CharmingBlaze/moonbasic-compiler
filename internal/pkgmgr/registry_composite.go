package pkgmgr

import (
	"fmt"
	"os"
	"strings"
)

// DefaultRegistryURL is the public package index used when MOONBASIC_REGISTRY is unset.
// Offline installs fall back to the bundled index (same entries, builtin packages).
const DefaultRegistryURL = "https://raw.githubusercontent.com/CharmingBlaze/moonbasic-compiler/main/internal/pkgmgr/default_index.json"

// fallbackRegistry tries primary first, then secondary on lookup failure.
type fallbackRegistry struct {
	primary   Registry
	secondary Registry
}

func (f *fallbackRegistry) Lookup(name string) (IndexEntry, error) {
	if f.primary != nil {
		if e, err := f.primary.Lookup(name); err == nil {
			return e, nil
		}
	}
	if f.secondary != nil {
		return f.secondary.Lookup(name)
	}
	return IndexEntry{}, fmt.Errorf("registry: unknown package %q", name)
}

// DefaultRegistryFromEnv returns a registry from MOONBASIC_REGISTRY, or remote+bundled fallback.
func DefaultRegistryFromEnv() Registry {
	raw := strings.TrimSpace(os.Getenv("MOONBASIC_REGISTRY"))
	if raw != "" {
		if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
			return &HTTTPRegistry{IndexURL: raw}
		}
		return &FileRegistry{Path: raw}
	}
	return &fallbackRegistry{
		primary:   &HTTTPRegistry{IndexURL: DefaultRegistryURL},
		secondary: &defaultRegistry{},
	}
}
