//go:build (!linux && !windows) || !cgo

package mbphysics3d

// JoltBackendLabel is printed by moonrun -version so users can tell stub vs native.
func JoltBackendLabel() string {
	return "stub (no native Jolt — soft BODY3D only; need Windows/Linux + CGO)"
}
