package pkgmgr

import (
	"testing"
)

func TestFallbackRegistryUsesSecondary(t *testing.T) {
	primary := &HTTTPRegistry{IndexURL: "http://127.0.0.1:1/no-such-index"}
	reg := &fallbackRegistry{primary: primary, secondary: &defaultRegistry{}}
	e, err := reg.Lookup("demo_extra")
	if err != nil {
		t.Fatal(err)
	}
	if e.Version != "0.1.0" {
		t.Fatalf("version %q", e.Version)
	}
}

func TestDefaultRegistryFromEnvUsesFallback(t *testing.T) {
	t.Setenv("MOONBASIC_REGISTRY", "")
	reg := DefaultRegistryFromEnv()
	if _, ok := reg.(*fallbackRegistry); !ok {
		t.Fatalf("want fallbackRegistry got %T", reg)
	}
}
