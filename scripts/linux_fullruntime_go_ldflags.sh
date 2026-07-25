#!/usr/bin/env bash
# Shared by .github/workflows/release.yml for Linux amd64 full-runtime builds.
#
# Contract (ELF link):
#   - Go uses -linkmode external so g++ drives the final link (needed for Jolt C++).
#   - -static-libgcc / -static-libstdc++ embed the C++ runtime so moonrun does not
#     depend on a matching libstdc++.so.6 from the build host (Arch/containers).
#   - Desktop stack (Wayland, xkbcommon, OpenGL, glibc) remains dynamic — that is expected.
#   - Engine libs (ENet, Raylib sources, Jolt .a) must stay static; see verify_linux_shared_libs.sh.
#
# Version injection: set MOONBASIC_LINUX_VERSION, or rely on GITHUB_REF_NAME (tags), else "devel".
#
# shellcheck disable=SC2034
moonbasic_linux_fullruntime_go_ldflags() {
  local ver="${MOONBASIC_LINUX_VERSION:-${GITHUB_REF_NAME:-devel}}"
  # -extld=g++ ensures the C++ runtime (Jolt) links correctly with static-libstdc++.
  printf '%s' \
    "-s -w -X moonbasic/internal/version.Version=${ver} " \
    "-linkmode external " \
    "-extld=g++ " \
    "-extldflags=-static-libgcc " \
    "-extldflags=-static-libstdc++"
}
