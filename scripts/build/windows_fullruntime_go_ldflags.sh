#!/usr/bin/env bash
# Shared by .github/workflows/release.yml and ci.yml for Windows amd64 full-runtime builds.
#
# Contract (PE link):
#   - Go uses -linkmode external with -extld=g++ so -static-libstdc++ is honored
#     (plain gcc as the external linker often still DT_NEEDs libstdc++-6.dll).
#   - -static-libgcc / -static-libstdc++ and -Bstatic -lwinpthread pull MinGW
#     runtimes into the PE — no sidecar DLLs on user machines.
#   - Raylib must NOT be linked via -lraylib (no CGO_LDFLAGS): sources compile in-tree.
#   - Jolt cgo also requests static stdc++/winpthread (cgo_windows_amd64.go).
#
# Version injection: set MOONBASIC_WINDOWS_VERSION, or rely on GITHUB_REF_NAME (tags), else "devel".
#
# shellcheck disable=SC2034
moonbasic_windows_fullruntime_go_ldflags() {
  local ver="${MOONBASIC_WINDOWS_VERSION:-${GITHUB_REF_NAME:-devel}}"
  printf '%s' \
    "-s -w -X moonbasic/internal/version.Version=${ver} " \
    "-linkmode external " \
    "-extld=g++ " \
    "-extldflags=-static-libgcc " \
    "-extldflags=-static-libstdc++ " \
    "-extldflags=-Wl,-Bstatic " \
    "-extldflags=-lstdc++ " \
    "-extldflags=-lwinpthread " \
    "-extldflags=-lgcc " \
    "-extldflags=-Wl,-Bdynamic"
}
