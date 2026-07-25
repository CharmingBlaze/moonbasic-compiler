#!/usr/bin/env bash
# Fail if a macOS full-runtime binary links engine dylibs that must be static.
#
# Usage:
#   bash scripts/verification/verify_macos_shared_libs.sh dist/moonrun dist/moonbasic
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
  mapfile -t lines < <(otool -L "$bin" 2>&1 || true)
  echo "== otool -L $bin =="
  printf '%s\n' "${lines[@]}"
  bad=()
  for line in "${lines[@]}"; do
    if [[ "$line" =~ $FORBIDDEN_REGEX ]]; then
      bad+=("$line")
    fi
  done
  if ((${#bad[@]} > 0)); then
    echo "ERROR: $bin links forbidden engine dylibs:" >&2
    printf '  %s\n' "${bad[@]}" >&2
    fail=1
  else
    echo "OK: $bin — no forbidden engine dylibs."
  fi
done

exit "$fail"
