#!/usr/bin/env bash
# Fail if a Linux full-runtime binary DT_NEEDs engine libs that must be statically linked
# into moonrun (same class of bug as libenet.so.7 on Arch/CachyOS).
#
# Usage (from repo root, after building dist/moonrun):
#   bash scripts/verification/verify_linux_shared_libs.sh dist/moonrun
#   bash scripts/verification/verify_linux_shared_libs.sh dist/moonbasic   # optional; compiler may be thinner
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <path-to-elf> [more-elfs...]" >&2
  exit 2
fi

# Libraries that must NEVER appear — they mean we accidentally dynamic-linked a game engine dep
# or failed to static-link the C++ runtime used by Jolt.
FORBIDDEN_REGEX='libenet\.so|libraylib\.so|libjolt\.so|libsqlite3\.so|libbox2d\.so|libstdc\+\+\.so|libgcc_s\.so'

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
  mapfile -t lines < <(ldd "$bin" 2>&1 || true)
  echo "== ldd $bin =="
  printf '%s\n' "${lines[@]}"
  bad=()
  for line in "${lines[@]}"; do
    if [[ "$line" =~ $FORBIDDEN_REGEX ]]; then
      bad+=("$line")
    fi
    if [[ "$line" == *"not found"* ]]; then
      # Missing DT_NEEDED is always a release blocker (except linux-vdso).
      if [[ "$line" != *"linux-vdso"* ]]; then
        bad+=("$line")
      fi
    fi
  done
  if ((${#bad[@]} > 0)); then
    echo "ERROR: $bin has forbidden or missing shared libs:" >&2
    printf '  %s\n' "${bad[@]}" >&2
    echo "Expected: ENet/Raylib/Jolt/SQLite statically linked; only desktop stack (glibc, libstdc++, Wayland/X11, OpenGL, xkbcommon, …) may be dynamic." >&2
    fail=1
  else
    echo "OK: $bin — no forbidden engine shared libs; no missing DT_NEEDED."
  fi
done

exit "$fail"
