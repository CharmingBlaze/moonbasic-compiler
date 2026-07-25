//go:build windows && amd64 && cgo

/*
WINDOWS CGO LDFLAGS
The linker flags below are specifically ordered to ensure proper resolution of Jolt 
and standard library symbols on Windows. 

LIBRARIES:
- libJolt_wrapper.a: C-wrapper bridge exposing CharacterVirtual and World APIs.
- libJolt.a: The core Jolt Physics engine (compiled with SIMD).

BUILD STEPS:
Static libraries must be pre-compiled using the PowerShell script in 
third_party/jolt-go/scripts/build-libs-windows.ps1 before building the Go runtime.
*/

package jolt

/*
Force static libstdc++/libgcc/winpthread. A bare -lstdc++ after Go's -Wl,-Bdynamic
re-introduces libstdc++-6.dll / libwinpthread-1.dll load-time deps (STATUS_DLL_NOT_FOUND
on clean user PCs that do not have MinGW on PATH).
#cgo LDFLAGS: -L${SRCDIR}/lib/windows_amd64 -ljolt_wrapper -lJolt -static-libstdc++ -static-libgcc -Wl,-Bstatic -lwinpthread -Wl,-Bdynamic
*/
import "C"
