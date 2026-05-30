#!/usr/bin/env bash
# Build an AppImage from a moonBASIC full-runtime dist/ folder (Linux x64).
# Requires: moonbasic + moonrun in dist/, linuxdeploy + appimagetool on PATH.
#
# Usage: ./scripts/build_appimage.sh [version-tag]
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TAG="${1:-dev}"
DIST="$ROOT/dist"
APPDIR="$ROOT/build/AppDir"
OUT="$ROOT/moonbasic-${TAG}-x86_64.AppImage"

if [[ ! -f "$DIST/moonrun" && ! -f "$DIST/moonbasic" ]]; then
  echo "error: build full runtime into dist/ first (moonbasic + moonrun)" >&2
  exit 1
fi

rm -rf "$APPDIR"
mkdir -p "$APPDIR/usr/bin" "$APPDIR/usr/share/applications"

cp "$DIST/moonbasic" "$APPDIR/usr/bin/"
cp "$DIST/moonrun" "$APPDIR/usr/bin/"
chmod +x "$APPDIR/usr/bin/"*

cat > "$APPDIR/moonbasic.desktop" <<EOF
[Desktop Entry]
Name=moonBASIC
Exec=moonrun
Icon=moonbasic
Type=Application
Categories=Development;
EOF

cat > "$APPDIR/AppRun" <<'EOF'
#!/bin/sh
HERE="$(dirname "$(readlink -f "$0")")"
export PATH="$HERE/usr/bin:$PATH"
exec moonrun "$@"
EOF
chmod +x "$APPDIR/AppRun"

if command -v linuxdeploy >/dev/null 2>&1; then
  linuxdeploy --appdir "$APPDIR" --output appimage
  mv ./*.AppImage "$OUT" 2>/dev/null || true
fi

if [[ ! -f "$OUT" ]] && command -v appimagetool >/dev/null 2>&1; then
  appimagetool "$APPDIR" "$OUT"
fi

if [[ -f "$OUT" ]]; then
  echo "Wrote $OUT"
else
  echo "AppImage tools not found; staged AppDir at $APPDIR" >&2
  exit 1
fi
