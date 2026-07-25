package enet

/*
#cgo CFLAGS: -I${SRCDIR}/enet/include -I${SRCDIR}/enet
#cgo !windows CFLAGS: -DHAS_FCNTL=1 -DHAS_POLL=1 -DHAS_GETADDRINFO=1 -DHAS_GETNAMEINFO=1 -DHAS_INET_PTON=1 -DHAS_INET_NTOP=1 -DHAS_MSGHDR_FLAGS=1 -DHAS_SOCKLEN_T=1 -DHAS_OFFSETOF=1
#cgo windows LDFLAGS: -lws2_32 -lwinmm
#include <enet/enet.h>
*/
import "C"
import "fmt"

// Initialize enet
func Initialize() {
	C.enet_initialize()
}

// Deinitialize enet
func Deinitialize() {
	C.enet_deinitialize()
}

// LinkedVersion returns the linked version of enet currently being used.
// Returns MAJOR.MINOR.PATCH as a string.
func LinkedVersion() string {
	var version = uint32(C.enet_linked_version())
	major := uint8(version >> 16)
	minor := uint8(version >> 8)
	patch := uint8(version)
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}
