# Using moonBASIC from a release

If you came here from the repo: **you usually do not need this folder.**  
**[Download the latest pre-built archive](https://github.com/CharmingBlaze/moonbasic/releases/latest)** — that is the normal way to get `moonbasic` and `moonrun`. The sections below describe what those downloads contain. Maintainer-only packaging notes are at the **bottom**.

---

## “All commands” — what that means

1. **Full builtin catalog (language + tooling)**  
   Every builtin name is defined in **`compiler/builtinmanifest/commands.json`** in the source tree. Any **`moonbasic`** binary from a release (compiler-only **or** full bundle) uses that catalog for **`--check`**, **`--lsp`**, and compiling to **`.mbc`**. Authors and CI get the **complete command list** at check/compile time from the compiler alone.

2. **Running games (engine at runtime)**  
   Calling **`WINDOW.*`**, **`PHYSICS3D.*`**, etc. needs the **engine**, which ships as **`moonrun`** in the **full runtime** archives — not in the compiler-only zip. Use the **full runtime** download when you need to **execute** those calls on a machine (graphics, physics, net, …).

3. **Five download families per version tag**  
   Each release publishes **IDE bundles** (Windows/Linux/macOS), **two** full-runtime archives, **two** compiler-only archives, and a VS Code extension. Together they cover beginners (IDE), tooling, and “run anywhere you install the full bundle.”

4. **moonBASIC IDE (recommended for beginners)**  
   **`moonbasic-<tag>-ide-windows-amd64.zip`**, **`moonbasic-<tag>-ide-linux-amd64.tar.gz`**, or **`moonbasic-<tag>-ide-macos-arm64.tar.gz`** — includes **`moonbasic-ide`**, **`moonbasic`**, **`moonrun`**, and **README-IDE-RELEASE.txt**. Documentation is built into the IDE. Extract and run **START-IDE** — no other tools required. Published on **[github.com/CharmingBlaze/moonbasic/releases](https://github.com/CharmingBlaze/moonbasic/releases)**.

5. **VS Code extension (optional)**  
   Each tagged release also publishes **`moonbasic-<tag>-vscode.vsix`**. Install with **Extensions → … → Install from VSIX…** in Visual Studio Code, then ensure **`moonbasic`** is on **`PATH`** or set **`moonbasic.languageServerPath`** to your **`moonbasic` / `moonbasic.exe`**. No Node.js or git clone required. Steps: **[`docs/GETTING_STARTED.md`](../docs/GETTING_STARTED.md#vs-code-syntax-and-lsp)**.

---

## Two kinds of downloads

**Release binaries do not require Go, GCC, or Clang on the user’s machine** — `moonrun` compiles `.mb` internally. (Building *this repository* from source still needs those tools; see **`docs/BUILDING.md`**.)

| Artifact | Contents | Typical use |
|----------|----------|-------------|
| **IDE bundle** (`moonbasic-<tag>-ide-…`) | `moonbasic-ide` + `moonbasic` + `moonrun` + README | Easiest setup — editor, docs, check, compile, run in one folder. |
| **Full runtime** (`moonbasic-<tag>-windows-amd64.zip` / `linux-amd64.tar.gz`) | `moonbasic` + `moonrun` + README | Play and develop games with a window; needs OS + GPU/OpenGL stack; Windows may need VC++ redist (see `README-RELEASE.txt` in the zip). |
| **Compiler only** (`moonbasic-<tag>-compiler-…`) | `moonbasic` only — **no `moonrun`** | CI, lint, compile to `.mbc`, LSP — **no** Raylib DLLs beside the compiler. |

**Authors shipping a game:** see **[`docs/GETTING_STARTED.md`](../docs/GETTING_STARTED.md)** — **§2 Ship your game**. Players need the **full runtime** archive to run `moonrun`; Windows release binaries follow the static MinGW policy in **[`docs/BUILDING.md`](../docs/BUILDING.md)**. Linux full-runtime builds use normal shared libraries (not a single fully static ELF).

---

<details>
<summary><strong>Maintainers: building bundles locally &amp; other <code>dist/</code> paths</strong></summary>

- Build the compiler bundle locally: `scripts/release_compiler_windows.ps1` / `scripts/release_compiler_linux.sh` (see **`docs/BUILDING.md`**).
- **`windows/`** — NSIS script `moonbasic.nsi` builds an installer after staging `moonbasic.exe` (full-runtime GitHub zips ship static-linked Windows binaries — no MinGW DLLs beside the exes).
- **`linux/`** — When present in the tree, packaging scripts may expect a staged tree under `dist/stage/` with `bin/moonbasic`, `share/moonbasic/{examples,assets}`. Authors usually ship the official **linux-amd64** full-runtime **`.tar.gz`** or document “install from Releases + run `moonrun` on your `.mb` / `.mbc`” — see **[`docs/GETTING_STARTED.md`](../docs/GETTING_STARTED.md)** and **[`docs/BUILDING.md`](../docs/BUILDING.md)** (**Linux full-runtime shipping**).

Release CI (`.github/workflows/release.yml`) uploads **both** full-runtime and compiler-only archives on version tags. NSIS/AppImage integration may need local paths adjusted for your Raylib/MinGW layout.

</details>
