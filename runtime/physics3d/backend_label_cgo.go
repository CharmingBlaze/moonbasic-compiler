//go:build (linux || windows) && cgo

package mbphysics3d

// JoltBackendLabel is printed by moonrun -version so users can tell stub vs native.
func JoltBackendLabel() string {
	return "native (Windows/Linux CGO + Jolt)"
}
