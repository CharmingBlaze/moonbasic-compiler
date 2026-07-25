#!/usr/bin/env bash
# Fail if a macOS full-runtime binary links engine dylibs that must be static.
#
# Usage:
#   bash scripts/verification/verify_macos_shared_libs.sh dist/moonrun dist/moonbasic
#
# Compatible with macOS system bash 3.2 (no mapfile).
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <path-to-mach-o> [more...]" >&2
  exit 2
fi

FORBIDDEN_REGEX='libenet\.|libraylib\.|libjolt\.|libsqlite3\.|libbox2d\.'

fail=0
for bin in "$@"; do
  if [[ ! -f "$bin" ]]; then
    echo "ERROR: not a file: $bin" >&2
    fail=1
    continue
  fi
  if ! command -v otool >/dev/null 2>&1; then
    echo "ERROR: otool not found" >&2
    exit 1
  fi
  echo "== otool -L $bin =="
  tmp="$(mktemp)"
  otool -L "$bin" >"$tmp" 2>&1 || true
  bad=""
  while IFS= read -r line || [[ -n "$line" ]]; do
    printf '%s\n' "$line"
    if [[ "$line" =~ $FORBIDDEN_REGEX ]]; then
      bad="${bad}${line}"$'\n'
    fi
  done <"$tmp"
  rm -f "$tmp"
  if [[ -n "$bad" ]]; then
    echo "ERROR: $bin links forbidden engine dylibs:" >&2
    printf '%s' "$bad" >&2
    fail=1
  else
    echo "OK: $bin — no forbidden engine dylibs."
  fi
done

exit "$fail"
