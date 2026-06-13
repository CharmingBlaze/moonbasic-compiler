//go:build fullruntime

package pipeline

import "testing"

func TestEntityCreateCubeHandleChainCompiles(t *testing.T) {
	src := `cube = ENTITY.CREATECUBE(1, 1, 1).scale(1.4, 1.4, 1.4).pos(0.0, 1.0, 0.0).col(255, 150, 70)
`
	prog, err := CompileSource("test.mb", src)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	if prog == nil {
		t.Fatal("nil program")
	}
}
