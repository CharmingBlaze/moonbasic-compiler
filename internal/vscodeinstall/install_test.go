package vscodeinstall

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindVsixInDir(t *testing.T) {
	dir := t.TempDir()
	if p := findVsixInDir(dir); p != "" {
		t.Fatalf("empty dir: got %q", p)
	}
	plain := filepath.Join(dir, "other.vsix")
	if err := os.WriteFile(plain, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if p := findVsixInDir(dir); p != plain {
		t.Fatalf("fallback vsix: got %q want %q", p, plain)
	}
	pref := filepath.Join(dir, "moonbasic-v1.0.0-vscode.vsix")
	if err := os.WriteFile(pref, []byte("y"), 0o644); err != nil {
		t.Fatal(err)
	}
	if p := findVsixInDir(dir); p != pref {
		t.Fatalf("preferred vsix: got %q want %q", p, pref)
	}
}

func TestStripJSONComments(t *testing.T) {
	in := `{
  // comment
  "a": "b"
}`
	out := stripJSONComments(in)
	if out == in {
		t.Fatal("expected comments stripped")
	}
}
