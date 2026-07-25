# moonBASIC distribution layout

This folder is where **local** and **CI** packaging scripts write release archives.
End users normally download pre-built zips from
[GitHub Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest) — they do
**not** need this source tree.

## What each archive is for

| Archive | Contains | Use when |
|---------|----------|----------|
| **IDE** `moonbasic-<tag>-ide-…` | `moonbasic-ide`, `moonbasic`, `moonrun`, START-IDE | Easiest: editor + run |
| **Full runtime** `moonbasic-<tag>-…` | `moonbasic`, `moonrun`, README-RELEASE | Terminal: play/run games |
| **Compiler only** `moonbasic-<tag>-compiler-…` | `moonbasic` only | CI / `--check` / `.mbc` (no window) |
| **VS Code** `moonbasic-<tag>-vscode.vsix` | Extension | Optional editor support |

## “All commands”

Every `moonbasic` binary ships the full `commands.json` catalog for `--check`, compile,
and `--lsp`. **Running** graphics/physics/net commands needs **`moonrun`** (or IDE) from a
**full runtime** / **IDE** archive.

## Unzip and run (no engine sidecars)

Official full-runtime / IDE `moonrun` links **Raylib, Jolt, ENet, and SQLite into the
binary**. Players do **not** install `libenet`, `raylib`, or `jolt` packages.

| Platform | Extra OS needs |
|----------|----------------|
| Windows | Usually none (WebView2 for IDE). Keep both `.exe` files from the **same** zip. |
| Linux | Desktop GPU stack: Mesa/OpenGL, Wayland, `libxkbcommon` (see `packaging/README-RELEASE.txt`) |
| macOS | Apple Silicon; Gatekeeper may require right-click → Open once |

## Maintainer packaging

- Windows PE policy: `scripts/verification/verify_windows_pe_imports.ps1`
- Linux ELF policy: `scripts/verification/verify_linux_shared_libs.sh` + `scripts/build/linux_fullruntime_go_ldflags.sh`
- macOS: `scripts/verification/verify_macos_shared_libs.sh`
- Release workflow: `.github/workflows/release.yml`
