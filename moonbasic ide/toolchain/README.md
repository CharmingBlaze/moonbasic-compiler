# Local moonBASIC toolchain (for IDE testing)

Put compiler binaries here so the IDE finds them automatically — no Settings dialog needed.

## Quick setup

From the repo root (PowerShell):

```powershell
cd "moonbasic ide"
npm run toolchain:build
```

This builds `moonbasic.exe` into this folder. `moonrun.exe` is built when CGO/fullruntime is available.

  • Copy `raylib.dll` next to `moonrun.exe` when using a purego (CGO_ENABLED=0) build.
  • Official [GitHub Releases](https://github.com/CharmingBlaze/moonbasic-compiler/releases) ship a static `moonrun` that does not need `raylib.dll`.

## Path configuration

| File | Purpose |
|------|---------|
| `paths.json` | Default relative paths (committed) |
| `paths.local.json` | Your machine overrides (gitignored) — copy from `paths.local.example.json` |

Example `paths.local.json` pointing at a custom build:

```json
{
  "moonbasicPath": "C:\\dev\\moonbasic-compiler\\moonbasic.exe",
  "moonrunPath": "C:\\dev\\moonbasic-compiler\\moonrun.exe"
}
```

Paths in these files can be absolute or relative to this `toolchain/` folder.

## Priority

1. **File → Settings** (saved in `%APPDATA%\\moonbasic-ide\\settings.json`)
2. **`toolchain/paths.local.json`**
3. **`toolchain/paths.json`**
4. **`moonbasic.exe` / `moonrun.exe`** in this folder
5. PATH and folders near the IDE executable

`.exe` files here are gitignored — only the JSON templates are tracked.
