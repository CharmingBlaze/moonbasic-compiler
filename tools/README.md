# tools/

Developer utilities for auditing and maintaining the MoonBASIC compiler manifest and documentation.

## Manifest Patching

| Script | Purpose |
|--------|---------|
| `patch_manifest_missing.py` | Appends runtime-registered commands missing from `compiler/builtinmanifest/commands.json`. Entry data is split across the `manifest_entries/` subpackage (one module per namespace group). Idempotent — re-running adds only genuinely new entries. |
| `manifest_entries/` | Subpackage containing `(key, args, returns, phase)` tuples organised by namespace: `audio`, `camera`, `data`, `draw`, `draw_objects`, `entity`, `input`, `misc`, `model`, `physics`, `ray`, `rendering`, `world`. Each file has a docstring describing its scope. |
| `annotate_stubs/` | Marks compile-time stub commands in `commands.json` (`go run ./tools/annotate_stubs`). |

## Auditing

| Script | Purpose |
|--------|---------|
| `diff_manifest_runtime.py` | **Canonical (CI):** regenerates or verifies `docs/audit/manifest_keys.txt`, `docs/audit/runtime_keys.txt`, and `docs/MISSING_COMMANDS_AUDIT.md`. `runtime_keys.txt` uses a **wide** scan (`.Register("…")` literals plus `regFlat` / `regLegacy2` / `regRT0` / `reg` helpers). The markdown report uses **strict** literals for the “runtime not in manifest” section. |
| `manifest_gap_summary.py` | Prints `manifest − wide_runtime` counts **by namespace** (`python tools/manifest_gap_summary.py`); optional `--list JOLT` or `--list global` to enumerate keys in one bucket. |
| `gen_master_audit.py` | Generates a broader master-audit report. |

## Documentation

| Script | Purpose |
|--------|---------|
| `cmdaudit/` | Go tool (`go run ./tools/cmdaudit`) that audits manifest vs. doc-file coverage by namespace. |
| `apidoc/` | Go tool (`go run ./tools/apidoc`) that regenerates `docs/API_CONSISTENCY.md` from the manifest. |

## Removed (historical)

These were one-off audit artefacts or superseded scripts, now deleted from the repo:

- PowerShell manifest diff scripts (`audit_manifest.ps1`, `extract_runtime_keys.ps1`, `diff_keys.ps1`, …) — use `diff_manifest_runtime.py`
- `audit_manifest_runtime.py`, `strip_doc_typography.py`, unreferenced `check_*` PS1 one-offs
- `third_party/go-enet/` — unused vendored copy (module proxy is canonical)
- `scratch/`, `docs/audit/archives/`, `internal/bench/`, `internal/packer/`
- Orphan `testdata/examples/`, `testdata/dev_samples/`, duplicate runtime command markdown dumps
