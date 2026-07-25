#!/usr/bin/env bash
# Fail if a Linux full-runtime binary DT_NEEDs engine libs that must be statically linked
# into moonrun (same class of bug as libenet.so.7 on Arch/CachyOS).
#
# Usage (from repo root, after building dist/moonrun):
#   bash scripts/verification/verify_linux_shared_libs.sh dist/moonrun
#
# Compatible with bash 3.2+ (no mapfile).
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <path-to-elf> [more-elfs...]" >&2
  exit 2
fi

# Engine libs must NEVER appear as DT_NEEDED (must be statically linked into moonrun).
# libstdc++/libgcc may still be dynamic on some Ubuntu link lines even with -static-libstdc++;
# that is acceptable for release (desktop glibc stack). Prefer static via ldflags when possible.
FORBIDDEN_REGEX='libenet\.so|libraylib\.so|libjolt\.so|libsqlite3\.so|libbox2d\.so'

fail=0
for bin in "$@"; do
  if [[ ! -f "$bin" ]]; then
    echo "ERROR: not a file: $bin" >&2
    fail=1
    continue
  fi
  if ! command -v ldd >/dev/null 2>&1; then
    echo "ERROR: ldd not found" >&2
    exit 1
  fi
  echo "== ldd $bin =="
  tmp="$(mktemp)"
  ldd "$bin" >"$tmp" 2>&1 || true
  bad=""
  while IFS= read -r line || [[ -n "$line" ]]; do
    printf '%s\n' "$line"
    if [[ "$line" =~ $FORBIDDEN_REGEX ]]; then
      bad="${bad}${line}"$'\n'
    fi
    if [[ "$line" == *"not found"* ]]; then
      if [[ "$line" != *"linux-vdso"* ]]; then
        bad="${bad}${line}"$'\n'
      fi
    fi
  done <"$tmp"
  rm -f "$tmp"
  if [[ -n "$bad" ]]; then
    echo "ERROR: $bin has forbidden or missing shared libs:" >&2
    printf '%s' "$bad" >&2
    echo "Expected: ENet/Raylib/Jolt/SQLite statically linked; C++ runtime static (-static-libstdc++); only desktop stack (glibc, Wayland/X11, OpenGL, xkbcommon, …) may be dynamic." >&2
    fail=1
  else
    echo "OK: $bin — no forbidden engine shared libs; no missing DT_NEEDED."
  fi
done

exit "$fail"
