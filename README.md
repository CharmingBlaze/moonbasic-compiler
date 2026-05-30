# moonBASIC

**A modern BASIC for developers who want to build without unnecessary friction** — a real compiler to bytecode, with a runtime for 2D/3D graphics, physics, and networking.

**You do not need this repository’s source code to use moonBASIC.** Install from a pre-built archive, write `.mb` files, and run them. Everything below about folders like `compiler/` or `go build` is optional and only for people hacking on the engine itself.

---

## Download and use (recommended)

The **compiled distribution** ships only from **[GitHub Releases](https://github.com/CharmingBlaze/moonbasic-compiler/releases/latest)** (Windows and Linux **x64**): official **`moonbasic`** / **`moonrun`** binaries — use these for play and compile; building the engine from source is optional.

**Optional static page** (direct archive links, no repo browser): [charmingblaze.github.io/moonbasic-compiler](https://charmingblaze.github.io/moonbasic-compiler/) — use **Releases** if that URL is not live yet. Setup: `.github/workflows/github-pages.yml`.

| Your goal | Download (replace `<tag>` with the release, e.g. `v1.2.20`) |
|-----------|-------------------------------------------------------------|
| **Run games** (window, graphics, physics, audio) | **Full runtime:** `moonbasic-<tag>-windows-amd64.zip`, `moonbasic-<tag>-linux-amd64.tar.gz`, or `moonbasic-<tag>-macos-arm64.tar.gz` |
| **Compile** `.mb` → `.mbc`, **`--check`**, **`--lsp`** only (CI, tooling, no game window) | **Compiler only:** `moonbasic-<tag>-compiler-windows-amd64.zip` or `moonbasic-<tag>-compiler-linux-amd64.tar.gz` |
| **VS Code** (syntax + LSP + debugger) | **`moonbasic-<tag>-vscode.vsix`** — [install from VSIX](docs/GETTING_STARTED.md#vs-code-syntax-and-lsp) (same release as your binaries; no repo clone) |

- **Full runtime** includes **`moonbasic`** + **`moonrun`** (+ `README-RELEASE.txt`). Use this if you want to play or develop games with a window.
- **Compiler only** is a small folder with **`moonbasic`** only (no `moonrun`). Same command names for `--check` / compile / LSP as the full build — see **[`dist/README.md`](dist/README.md)** for the full picture.

### First steps after you extract (full runtime)

**Use `moonrun` — one command to compile and play.** The full-runtime archive includes both `moonbasic` (compiler/LSP) and `moonrun` (game engine). For your first window you only need `moonrun`.

**New project:** `moonbasic new MyGame` → `cd MyGame` → `moonrun main.mb` (creates `main.mb`, `assets/`, and VS Code debug config).

1. Open a terminal in the folder that contains **`moonrun`** (and **`moonbasic`**) from the **full runtime** zip/tarball.
2. **Check it works:** `moonrun --version` (Windows: `moonrun.exe --version`).
3. **Run a game:** `moonrun path/to/game.mb` — compiles in-process if needed, then opens the engine window.
4. **Optional — lint only:** `moonbasic --check path/to/game.mb` (no window).
5. **Optional — compile to bytecode:** `moonbasic path/to/game.mb` → writes **`game.mbc`** next to the source.

**Compiler-only** installs (CI/tooling) ship **`moonbasic`** without **`moonrun`** — use those for `--check`, `--lsp`, and `.mbc` output only; install the **full runtime** to run games with a window.

**Release `moonrun` does not require Go, GCC, or Clang on your machine** — it compiles `.mb` in-process, then runs the engine. You may still need a normal GPU stack (Linux) or the [VC++ x64 redistributable](https://learn.microsoft.com/en-us/cpp/windows/latest-supported-vc-redist) on some Windows setups if a DLL is missing; see **`README-RELEASE.txt`** inside the full-runtime archive.

**Editor:** run **`moonbasic --lsp`** and attach your LSP client (stdio). **VS Code:** download **`moonbasic-<tag>-vscode.vsix`** from [Releases](https://github.com/CharmingBlaze/moonbasic-compiler/releases/latest) and **Install from VSIX…** — [quick steps](docs/GETTING_STARTED.md#vs-code-syntax-and-lsp). Repo contributors: [DEVELOPER.md — moonBASIC in VS Code](docs/DEVELOPER.md#moonbasic-in-vs-code) (workspace tasks).

More detail on what each zip contains: **[`dist/README.md`](dist/README.md)** · step-by-step install, first window, and **how to ship your game to players**: **[`docs/GETTING_STARTED.md`](docs/GETTING_STARTED.md)** (section **Ship your game**)

---

## What moonBASIC is

Many engines impose their own complexity before you can begin creating. **moonBASIC** is intentionally direct: one toolchain, one mental model, and a workflow that values clarity and speed.

- **Vertical integration** — compiler, VM, and engine in one stack.
- **2D and 3D** — same language and workflow for both.
- **Real compilation** — not an interpreter; bytecode executed by a production-oriented runtime (**Raylib**, **Jolt** 3D, **Box2D** 2D, **ENet**, … where enabled).

The compiler is stable; the standard library covers Tiled, materials, sprites, atlases, particles, audio, lighting, shaders, and more. Explore the **[documentation index](#documentation)** and **[`ARCHITECTURE.md`](ARCHITECTURE.md)** when you want to go deeper.

---

## Example

Compact 3D sample (no `#` / `$` / `?` suffixes — implicit typing only):

```moonbasic
WINDOW.OPEN(960, 540, "Spinning cube")
WINDOW.SETFPS(60)

cam = CAMERA.MAKE()
CAMERA.SETPOSITION(cam, 0, 2, 8)
CAMERA.SETTARGET(cam, 0, 0, 0)
CAMERA.SETFOV(cam, 45)

cube = MESH.MAKECUBE(2, 2, 2)
mat  = MATERIAL.MAKEDEFAULT()
cubeXform = TRANSFORM.IDENTITY()
angle = 0

WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR WINDOW.SHOULDCLOSE())
    dt = TIME.DELTA()
    angle = angle + 45 * dt
    TRANSFORM.SETROTATION(cubeXform, 0, angle, 0)
    RENDER.CLEAR(12, 14, 22)
    CAMERA.BEGIN(cam)
        MESH.DRAW(cube, mat, cubeXform)
    CAMERA.END(cam)
    RENDER.FRAME()
WEND

MESH.FREE(cube)
MATERIAL.FREE(mat)
TRANSFORM.FREE(cubeXform)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

Full sample with grid and text: **[`examples/spin_cube`](examples/spin_cube/main.mb)**

---

## Architecture in brief

1. **Compilation** produces bytecode with the **`.mbc`** extension.
2. **Execution** is handled by the virtual machine, which talks to Raylib, Jolt, Box2D, ENet, and other systems through **CGO** where enabled in full runtime builds.

---

## Documentation

| Document | What it covers |
|----------|----------------|
| [docs/GETTING_STARTED.md](docs/GETTING_STARTED.md) | Install, first window, **`moonbasic new`**, debugging |
| [docs/FIRST_HOUR.md](docs/FIRST_HOUR.md) | Friendly intro for beginners |
| [docs/reference/MIGRATION.md](docs/reference/MIGRATION.md) | BlitzBASIC porting; compile-time stub list |
| [docs/PROGRAMMING.md](docs/PROGRAMMING.md) | Game loop, modules, 2D/3D |
| [docs/LANGUAGE.md](docs/LANGUAGE.md) | Variables, control flow, functions, **`$"..."`**, **`ENUM`**, multi-return |
| [docs/ROADMAP.md](docs/ROADMAP.md) | Language roadmap — shipped vs planned |
| [docs/COMMANDS.md](docs/COMMANDS.md) | Built-in command index |
| [examples/README.md](examples/README.md) | Runnable sample programs |
| [dist/README.md](dist/README.md) | Release artifacts explained |

More: **[docs/reference/](docs/reference/)**, **[docs/reference/MULTIPLAYER.md](docs/reference/MULTIPLAYER.md)** (multiplayer hub), **[docs/JOLT_WINDOWS_PARITY.md](docs/JOLT_WINDOWS_PARITY.md)** (Windows Jolt / CGO notes for engine devs), **[ARCHITECTURE.md](ARCHITECTURE.md)**.

---

<details>
<summary><strong>For contributors: repository layout</strong></summary>

The GitHub file tree is for **engine development**. End users who only download Releases never need to open these paths.

| Path | Purpose |
|------|---------|
| [`cmd/moonbasic`](cmd/moonbasic), [`cmd/moonrun`](cmd/moonrun) | CLI entrypoints (compiler vs full runtime). |
| [`compiler/`](compiler/), [`vm/`](vm/) | Language front-end, bytecode, VM. |
| [`runtime/`](runtime/) | Engine modules (rendering, physics, audio, net, …). |
| [`docs/`](docs/) | Guides and reference. |
| [`examples/`](examples/) | Runnable projects. |
| [`dist/`](dist/) | Packaging notes — see [`dist/README.md`](dist/README.md). |
| [`scripts/`](scripts/), [`tools/`](tools/) | Release packaging and audit helpers. |

</details>

<details>
<summary><strong>Build from source</strong></summary>

Building from source requires **Go** and a **C toolchain**. Full graphical programs need the **`fullruntime`** build tag and **`moonrun`** (or **`moonbasic --run`** from a full-runtime build). See **[docs/BUILDING.md](docs/BUILDING.md)**, **[CONTRIBUTING.md](CONTRIBUTING.md)**, and **[AGENTS.md](AGENTS.md)**.

```bash
git clone https://github.com/CharmingBlaze/moonbasic-compiler
cd moonbasic
# Windows (example): set CGO_ENABLED=1 and a working gcc, then:
go build -o moonbasic.exe .

# Run a 3D sample (full runtime + CGO):
CGO_ENABLED=1 go run -tags fullruntime ./cmd/moonrun examples/spin_cube/main.mb
```

</details>

---

## Contributing

Guidelines: **[CONTRIBUTING.md](CONTRIBUTING.md)** and **[docs/DEVELOPER.md](docs/DEVELOPER.md)**. CI validates builds, tests, and representative `go run . --check` samples.

On **Windows**, a **`fullruntime`** link that pulls in Jolt may require prebuilt **`libJolt.a`** and **`libjolt_wrapper.a`** in **[third_party/jolt-go/jolt/lib/windows_amd64/](third_party/jolt-go/jolt/lib/windows_amd64/README.md)**. **`scripts/check-jolt-windows-libs.ps1`** checks that both files are present.

---

## License

**MIT** — see [LICENSE](LICENSE).
