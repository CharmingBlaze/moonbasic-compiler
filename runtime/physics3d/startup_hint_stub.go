//go:build (!linux && !windows) || !cgo

package mbphysics3d

import (
	"fmt"
	"os"
)

// LogJoltPhysicsBackendHint prints one line to stderr: stub builds cannot link native Jolt.
func LogJoltPhysicsBackendHint() {
	fmt.Fprintf(os.Stderr, "moonBASIC: WARN [Jolt Physics] Stub mode — native Jolt is not linked on this build. ENTITY.PHYSICS / BODY3D.COMMIT use soft stubs only. For real simulation use the official Windows or Linux fullruntime moonrun (CGO_ENABLED=1 + Jolt libs; see docs/JOLT_WINDOWS_PARITY.md).\n")
}
