//go:build windows && amd64 && cgo

/*
WINDOWS CGO LDFLAGS

Libraries:
- libJolt_wrapper.a: C-wrapper bridge exposing CharacterVirtual and World APIs.
- libJolt.a: The core Jolt Physics engine (compiled with SIMD).

Static libraries must be pre-compiled using
third_party/jolt-go/scripts/build-libs-windows.ps1 before building the Go runtime.

Important: use -static-libstdc++ / -static-libgcc / -Bstatic -lwinpthread here.
A bare -lstdc++ after Go's -Wl,-Bdynamic re-introduces MinGW DLL load-time deps
(libstdc++-6.dll, libwinpthread-1.dll) and moonrun exits with STATUS_DLL_NOT_FOUND
on clean Windows PCs.
*/

package jolt

/*
#cgo LDFLAGS: -L${SRCDIR}/lib/windows_amd64 -ljolt_wrapper -lJolt -static-libstdc++ -static-libgcc -Wl,-Bstatic -lwinpthread -Wl,-Bdynamic
*/
import "C"
