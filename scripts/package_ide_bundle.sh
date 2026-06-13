#!/usr/bin/env bash
# Package moonBASIC IDE + toolchain into a release-style archive (maintainer helper).
# Usage: ./scripts/package_ide_bundle.sh [version-tag] [linux-amd64|macos-arm64]
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TAG="${1:-dev}"
PLATFORM="${2:-linux-amd64}"
IDE_DIR="$ROOT/moonbasic ide"
STAGE="$ROOT/dist/ide-bundle"
OUT="$ROOT/moonbasic-${TAG}-ide-${PLATFORM}.tar.gz"

if [ ! -f "$IDE_DIR/build/bin/moonbasic-ide" ]; then
  echo "Build the IDE first: cd 'moonbasic ide' && npm ci && npm run langdata && wails build" >&2
  exit 1
fi
if [ ! -f "$ROOT/dist/moonbasic" ] || [ ! -f "$ROOT/dist/moonrun" ]; then
  echo "Build runtime into dist/ first (see scripts/release-windows.sh or CI release.yml)" >&2
  exit 1
fi

rm -rf "$STAGE"
mkdir -p "$STAGE"
cp "$IDE_DIR/build/bin/moonbasic-ide" "$STAGE/"
cp "$ROOT/dist/moonbasic" "$ROOT/dist/moonrun" "$STAGE/"
cp "$ROOT/packaging/README-IDE-RELEASE.txt" "$STAGE/"
cp "$ROOT/packaging/START-IDE.sh" "$STAGE/"
chmod +x "$STAGE/moonbasic-ide" "$STAGE/moonbasic" "$STAGE/moonrun" "$STAGE/START-IDE.sh"
tar czvf "$OUT" -C "$STAGE" .
echo "Wrote $OUT"
