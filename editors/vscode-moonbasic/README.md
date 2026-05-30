# moonBASIC (VS Code extension)

Syntax highlighting, snippets, **LSP**, and a **debugger** for moonBASIC (`.mb` / `.mbc`) in Visual Studio Code.

## Prerequisites

- **`moonbasic`** on your `PATH`, or configure **`moonbasic.languageServerPath`** to the full path of the executable.
- **Debugging** requires **`moonrun`** from a **full runtime** release (not compiler-only). Set **`moonbasic.moonrunPath`** if `moonrun` is not on `PATH`.

Get binaries from [GitHub Releases](https://github.com/CharmingBlaze/moonbasic-compiler/releases/latest) or build from source (`go build ./cmd/moonbasic`, `go build -tags fullruntime ./cmd/moonrun`).

## Install from a release (easiest)

Each [GitHub Release](https://github.com/CharmingBlaze/moonbasic-compiler/releases/latest) ships **`moonbasic-<tag>-vscode.vsix`** next to the platform archives.

1. **Extensions** → **⋯** → **Install from VSIX…** and pick that file.
2. Set **`moonbasic.languageServerPath`** / **`moonbasic.moonrunPath`** if the exes are not on `PATH`.
3. Open a `.mb` file → **Run and Debug** → **Debug moonBASIC**.

User-facing steps: **[docs/GETTING_STARTED.md](../../docs/GETTING_STARTED.md#vs-code-syntax-and-lsp)**.

## Settings

| Setting | Description |
|--------|-------------|
| **`moonbasic.languageServerPath`** | Path to `moonbasic` for LSP. Empty = `moonbasic` from `PATH`. |
| **`moonbasic.moonrunPath`** | Path to `moonrun` for the debugger. Empty = `moonrun` from `PATH`. |

## Build from this repository

```bash
cd editors/vscode-moonbasic
npm install
npm run compile
npm run package   # produces .vsix
```

Release CI builds the same artifact as **`moonbasic-<tag>-vscode.vsix`**.

## Workspace tasks

The moonBASIC repository includes [`.vscode/tasks.json`](../../.vscode/tasks.json) for **Check**, **Compile to .mbc**, and **moonrun** on the active file, plus [`.vscode/launch.json`](../../.vscode/launch.json) for debugging examples.
