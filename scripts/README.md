# Scripts

Maintainer and contributor helpers, grouped by purpose. Run from the **repository root** unless noted.

| Folder | Contents |
|--------|----------|
| [`build/`](build/) | `check_builds.*`, `build_static.ps1`, `build_appimage.sh`, fullruntime `*_go_ldflags.sh` |
| [`release/`](release/) | `release-windows.*`, `release_compiler_*` |
| [`packaging/`](packaging/) | IDE/runtime zip helpers, `stage_ide_extras.sh` |
| [`verification/`](verification/) | PE/ldd/otool checks, Jolt lib preflight, Pong smoke |
| [`development/`](development/) | `dev.*` task runners, VS Code extension installers |

Examples:

```bash
bash scripts/build/check_builds.sh
powershell -File scripts/build/check_builds.ps1
powershell -File scripts/development/dev.ps1 check
bash scripts/packaging/package_ide_bundle.sh v1.2.30 linux-amd64
```

Do not put IDE-local scripts here — those live under `ide/scripts/` (Wails sync / local toolchain).
