#!/usr/bin/env bash
# Package moonBASIC IDE + toolchain into a release-style archive (maintainer helper).
# Usage: ./scripts/packaging/package_ide_bundle.sh [version-tag] [linux-amd64|macos-arm64]
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TAG="${1:-dev}"
PLATFORM="${2:-linux-amd64}"
IDE_DIR="$ROOT/ide"
STAGE="$ROOT/dist/ide-bundle"
OUT="$ROOT/moonbasic-${TAG}-ide-${PLATFORM}.tar.gz"

if [ ! -f "$IDE_DIR/build/bin/moonbasic-ide" ] && [ ! -d "$IDE_DIR/build/bin/moonbasic-ide.app" ]; then
  echo "Build the IDE first: cd 'ide' && npm ci && npm run langdata && wails build" >&2
  exit 1
fi
if [ ! -f "$ROOT/dist/moonbasic" ] || [ ! -f "$ROOT/dist/moonrun" ]; then
  echo "Build runtime into dist/ first (see scripts/release/release-windows.sh or CI release.yml)" >&2
  exit 1
fi

rm -rf "$STAGE"
mkdir -p "$STAGE"
if [ -d "$IDE_DIR/build/bin/moonbasic-ide.app" ]; then
  cp -R "$IDE_DIR/build/bin/moonbasic-ide.app" "$STAGE/"
else
  cp "$IDE_DIR/build/bin/moonbasic-ide" "$STAGE/"
  chmod +x "$STAGE/moonbasic-ide"
fi
cp "$ROOT/dist/moonbasic" "$ROOT/dist/moonrun" "$STAGE/"
chmod +x "$STAGE/moonbasic" "$STAGE/moonrun"
bash "$ROOT/scripts/packaging/stage_ide_extras.sh" "$STAGE"
tar czvf "$OUT" -C "$STAGE" .
echo "Wrote $OUT"
