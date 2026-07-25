package enet

// Static ENet: compile vendored lsalzman/enet C sources into this package so
// moonrun does not require a system libenet.so / libenet.dylib (fixes Arch,
// CachyOS, and other distros where the shared library is awkward to find).

/*
#cgo CFLAGS: -I${SRCDIR}/enet/include -I${SRCDIR}/enet
#cgo !windows CFLAGS: -DHAS_FCNTL=1 -DHAS_POLL=1 -DHAS_GETADDRINFO=1 -DHAS_GETNAMEINFO=1 -DHAS_INET_PTON=1 -DHAS_INET_NTOP=1 -DHAS_MSGHDR_FLAGS=1 -DHAS_SOCKLEN_T=1 -DHAS_OFFSETOF=1
#include "enet/enet_amalgamation.c"
*/
import "C"
