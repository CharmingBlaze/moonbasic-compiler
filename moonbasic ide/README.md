# moonBASIC IDE

Desktop IDE for **moonBASIC** `.mb` source files. Uses the real **moonbasic** compiler and **moonrun** runtime — not a separate transpiler.

## Download (recommended)

**Easiest setup:** download the **IDE bundle** from [github.com/CharmingBlaze/moonbasic/releases](https://github.com/CharmingBlaze/moonbasic/releases/latest):

| Platform | Archive |
|----------|---------|
| Windows | `moonbasic-<tag>-ide-windows-amd64.zip` |
| Linux x64 | `moonbasic-<tag>-ide-linux-amd64.tar.gz` |
| macOS Apple Silicon | `moonbasic-<tag>-ide-macos-arm64.tar.gz` |

Each bundle includes **moonbasic-ide**, **moonbasic**, **moonrun**, and **README-IDE-RELEASE.txt**. Documentation is built into the IDE.

1. Extract anywhere permanent.
2. Run **START-IDE.bat** (Windows), **START-IDE.command** (macOS), or **./START-IDE.sh** (Linux).
3. Open **`samples/hello.mb`** (or File → Open Samples Folder) → press **F5**.

No Go, Node.js, or VS Code required. Optional **ADD-TO-PATH** puts `moonbasic` / `moonrun` on your shell PATH; the IDE does not need it.

Engine source and release CI live in [moonbasic-compiler](https://github.com/CharmingBlaze/moonbasic-compiler).

## Features

- **Syntax highlighting** for moonBASIC (keywords, namespaces, handle methods, strings, comments)
- **Help at cursor** — documentation for `NAMESPACE.METHOD` and globals (manifest + **LSP**)
- **Built-in documentation** — full docs tree searchable inside the IDE
- **Autocomplete** — manifest index + **moonbasic --lsp**
- **Themes & appearance** — 7 presets, custom colors, font sizes (gear menu / Settings)
- **Check** — `moonbasic --check`
- **Compile** — writes `.mbc` bytecode next to source
- **Run / Stop** — launches `moonrun` (F5) and can kill it (Shift+F5 / Stop button)
- **Offline editor** — CodeMirror and fonts are local (no CDN required)
- **Cross-platform** — Windows, Linux, and macOS shortcuts (⌘ on Mac); toolchain auto-detects binaries next to the IDE

## Build from source (contributors)

| Mode | What you need |
|------|----------------|
| **Browser** (`npm run dev`) | Syntax, autocomplete, API help only |
| **Desktop** (`wails dev` / `wails build`) | Go, Wails v2, **moonbasic** + **moonrun** on PATH or in `toolchain/` |

```bash
cd "moonbasic ide"
npm install
npm run langdata
wails dev
```

Release binary:

```bash
wails build
```

Local toolchain (dev; Windows PowerShell or Unix shell):

```bash
npm run toolchain:build
```

Package IDE + runtime bundle (after building `dist/moonbasic` + `dist/moonrun`):

```powershell
.\scripts\package_ide_bundle.ps1 -Version v1.2.27
```

## Regenerate language data

When `compiler/builtinmanifest/commands.json` or `docs/` changes:

```bash
go run ./tools/ideexport
go run ./tools/docsexport
# or from moonbasic ide/:  npm run langdata
```

## VS Code alternative

For VS Code / Cursor, use the official extension in `editors/vscode-moonbasic/`.

## License

MIT — same as moonbasic-compiler.
