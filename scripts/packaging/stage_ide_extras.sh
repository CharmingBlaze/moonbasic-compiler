#!/usr/bin/env bash
# Copy IDE release extras into a staging directory (arg1).
set -euo pipefail
STAGE="${1:?stage dir}"
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
PKG="$ROOT/packaging"

cp "$PKG/README-IDE-RELEASE.txt" "$STAGE/"
cp "$PKG/START-IDE.sh" "$STAGE/"
cp "$PKG/ADD-TO-PATH.sh" "$STAGE/"
chmod +x "$STAGE/START-IDE.sh" "$STAGE/ADD-TO-PATH.sh"

# macOS Finder double-click launcher (harmless on Linux).
cp "$PKG/START-IDE.command" "$STAGE/"
chmod +x "$STAGE/START-IDE.command"

# Windows launchers (useful in mixed archives / wine).
cp "$PKG/START-IDE.bat" "$STAGE/" 2>/dev/null || true
cp "$PKG/ADD-TO-PATH.bat" "$STAGE/" 2>/dev/null || true

mkdir -p "$STAGE/samples"
cp "$PKG/samples/"*.mb "$PKG/samples/README.txt" "$STAGE/samples/"

# Offline docs (same content as IDE sidebar; prefer post-docsexport bundled copy).
mkdir -p "$STAGE/docs"
if [[ -d "$ROOT/ide/bundled-docs" ]]; then
  cp -R "$ROOT/ide/bundled-docs/." "$STAGE/docs/"
else
  cp -R "$ROOT/docs/." "$STAGE/docs/"
fi
