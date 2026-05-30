package pkgmgr

import (
	"testing"
)

func TestDefaultRegistryLookup(t *testing.T) {
	reg := DefaultRegistryFromEnv()
	e, err := reg.Lookup("demo_extra")
	if err != nil {
		t.Fatal(err)
	}
	if e.Version != "0.1.0" {
		t.Fatalf("version %q", e.Version)
	}
	if e.URL != "moonbasic-builtin://demo_extra" {
		t.Fatalf("url %q", e.URL)
	}
}

func TestInstallBuiltinDemoExtra(t *testing.T) {
	t.Setenv("MOONBASIC_CACHE", t.TempDir())
	if err := Install("demo_extra", DefaultRegistryFromEnv()); err != nil {
		t.Fatal(err)
	}
}
