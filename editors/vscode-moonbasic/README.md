# moonBASIC (VS Code extension)

Syntax highlighting, snippets, and **LSP** for moonBASIC (`.mb` / `.mbc`) in Visual Studio Code. The extension starts **`moonbasic --lsp`** (stdio) and talks to the same language server as other editors.

## Prerequisites

- **`moonbasic`** on your `PATH`, or configure **`moonbasic.languageServerPath`** (VS Code setting) to the full path of the executable (on Windows, e.g. `C:\path\to\moonbasic.exe`). Get a binary from [GitHub Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest) or build with `go build -o moonbasic ./cmd/moonbasic` from the repo root.

## Install from a release (easiest)

Each [GitHub Release](https://github.com/CharmingBlaze/moonbasic/releases/latest) ships **`moonbasic-<tag>-vscode.vsix`** next to the platform archives. In VS Code: **Extensions** → **⋯** → **Install from VSIX…** and pick that file. No Node.js or git clone required. If **`moonbasic`** is not on `PATH`, set **`moonbasic.languageServerPath`** to your **`moonbasic` / `moonbasic.exe`**.

User-facing steps: **[docs/GETTING_STARTED.md](../../docs/GETTING_STARTED.md#vs-code-syntax-and-lsp)**.

## Build from this repository (folder / VSIX)

1. **Node.js** required for these steps.
2. Open this folder in a terminal: `editors/vscode-moonbasic`.
3. `npm install` then `npm run compile`.
4. In VS Code: **Extensions** → `…` menu → **Install from Folder…** and choose this `vscode-moonbasic` directory.

Or package locally:

```bash
npm run package
```

Then **Extensions** → `…` → **Install from VSIX…** and select the generated `.vsix` file (release CI builds the same artifact as **`moonbasic-<tag>-vscode.vsix`**).

## Settings

| Setting | Description |
|--------|-------------|
| **`moonbasic.languageServerPath`** | Empty string uses `moonbasic` from `PATH`. Set to a full path if the compiler is not on `PATH`. |

The repo root [`.vscode/settings.json`](../../.vscode/settings.json) can set this for the workspace.

## Workspace tasks

The moonBASIC repository includes [`.vscode/tasks.json`](../../.vscode/tasks.json) for **Check**, **Compile to .mbc**, and **moonrun** on the active file. The **moonrun** task requires a **full-runtime** install (`moonrun` on `PATH`); compiler-only archives do not ship `moonrun`.
