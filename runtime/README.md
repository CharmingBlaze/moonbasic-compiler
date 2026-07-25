# runtime/

Implementations of moonBASIC builtins (draw, window, entities, physics, audio, …).

- Prefer `hal` types and `rt.Driver` for new window/video/input code.
- Register commands on every split build path (`*_cgo.go` / `*_stub.go`) so Windows and Linux expose the same manifest keys.
- Keep headless / null-driver paths runnable without a GPU.

Domain regrouping under `graphics/`, `physics/`, etc. is **deferred** (see [`docs/development/REPOSITORY_CLEANUP_AUDIT.md`](../docs/development/REPOSITORY_CLEANUP_AUDIT.md)). Physics sync notes: [`docs/PHYSICS.md`](../docs/PHYSICS.md).
