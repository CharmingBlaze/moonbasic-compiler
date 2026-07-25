# Vendored go-enet (moonBASIC)

Fork of [`github.com/codecat/go-enet`](https://github.com/codecat/go-enet) used via
`replace` in the moonBASIC `go.mod`.

## Why

Upstream go-enet uses `#cgo pkg-config: libenet` on Linux/macOS, so release
`moonrun` binaries can require `libenet.so.7` at runtime. Many distros (e.g.
Arch/CachyOS) do not make that package obvious for end users.

This tree **compiles [lsalzman/enet](https://github.com/lsalzman/enet) v1.3.18
C sources statically** into the Go package (same idea as vendored Raylib/Jolt).
Players only need the usual desktop/OpenGL stack — not a separate ENet install.

## Layout

- `*.go` — Go bindings (from codecat/go-enet)
- `enet/*.c` + `enet/include/` — ENet 1.3.18 sources/headers
- `enet/enet_amalgamation.c` — single translation unit for cgo
- `enet/enet.lib` — unused Windows leftover from upstream; safe to ignore

## Updating ENet

1. Replace `enet/*.c` and `enet/include/enet/*` from an upstream ENet tag.
2. Keep `enet_amalgamation.c` include list in sync.
3. `go test -tags fullruntime ./runtime/net/...` (CGO enabled).
