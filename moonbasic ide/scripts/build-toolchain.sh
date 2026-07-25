#!/usr/bin/env bash
# Build moonbasic/moonrun into moonbasic ide/toolchain/ for local IDE testing.
set -euo pipefail

IDE_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
REPO_ROOT="$(cd "$IDE_ROOT/.." && pwd)"
OUT_DIR="$IDE_ROOT/toolchain"
mkdir -p "$OUT_DIR"

cd "$REPO_ROOT"

echo "Building moonbasic (compiler-only, CGO_ENABLED=0)..."
CGO_ENABLED=0 go build -o "$OUT_DIR/moonbasic" .

echo "Building moonrun (fullruntime)..."
if CGO_ENABLED=1 go build -tags fullruntime -o "$OUT_DIR/moonrun" ./cmd/moonrun; then
  echo "Done. Binaries in: $OUT_DIR"
else
  echo "Warning: moonrun build failed (CGO may be required). moonbasic may still be usable for check/compile/LSP." >&2
fi

if [ -x "$OUT_DIR/moonbasic" ]; then
  "$OUT_DIR/moonbasic" --version || true
fi
