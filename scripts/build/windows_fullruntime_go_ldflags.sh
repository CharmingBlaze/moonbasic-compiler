#!/usr/bin/env bash
# Shared by .github/workflows/release.yml and ci.yml for Windows amd64 full-runtime builds.
#
# Contract (PE link):
#   - Go uses -linkmode external so the last link is MinGW g++ driving ld.
#   - -static-libgcc / -static-libstdc++ and -Bstatic -lwinpthread must win over
#     any later -lstdc++ from cgo (see third_party/jolt-go/jolt/cgo_windows_amd64.go).
#   - Raylib must NOT be linked via -lraylib (no CGO_LDFLAGS): sources compile in-tree.
#
# Version injection: set MOONBASIC_WINDOWS_VERSION, or rely on GITHUB_REF_NAME (tags), else "devel".
#
# shellcheck disable=SC2034
moonbasic_windows_fullruntime_go_ldflags() {
  local ver="${MOONBASIC_WINDOWS_VERSION:-${GITHUB_REF_NAME:-devel}}"
  # Repeat static runtime flags at the end so they override cgo -lstdc++ ordering.
  printf '%s' \
    "-s -w -X moonbasic/internal/version.Version=${ver} " \
    "-linkmode external " \
    "-extldflags=-static-libgcc " \
    "-extldflags=-static-libstdc++ " \
    "-extldflags=-Wl,-Bstatic " \
    "-extldflags=-lwinpthread " \
    "-extldflags=-lstdc++ " \
    "-extldflags=-lgcc " \
    "-extldflags=-Wl,-Bdynamic " \
    "-extldflags=-static-libgcc " \
    "-extldflags=-static-libstdc++"
}
