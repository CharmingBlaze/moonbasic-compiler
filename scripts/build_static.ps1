# Experimental: static-linked moonrun (CGO + Zig). Intended to produce a single .exe without
# Raylib or Jolt DLLs when the toolchain and Jolt static libraries are available.
#
# Prerequisites:
#   - Go, Zig on PATH (unless you override CC/CXX to MinGW gcc/g++).
#   - CGO_ENABLED=1 (this script sets it).
#   - Raylib is compiled from source via vendored raylib-go CGO (no separate raylib.dll).
#   - Native Jolt on Windows: build static libs first — see
#     third_party/jolt-go/jolt/lib/windows_amd64/README.md
#     (third_party/jolt-go/scripts/build-libs-windows.ps1 with JPH_SRC set).
#   - Purego / CGO_ENABLED=0 builds are out of scope here (they load raylib.dll at runtime).
#
# Usage (PowerShell, from repo root):
#   .\scripts\build_static.ps1
#
# Optional overrides:
#   $env:CC = "zig cc -target x86_64-windows-gnu -static"
#   $env:CXX = "zig c++ -target x86_64-windows-gnu -static"
#   $env:MOONBASIC_SKIP_STATIC_EXTLDFLAGS = "1"   # omit -ldflags -extldflags (if link fails)
#   $env:OUTPUT = "moonrun_static.exe"

$ErrorActionPreference = "Stop"

$ZigCC = "zig cc -target x86_64-windows-gnu -static"
$ZigCXX = "zig c++ -target x86_64-windows-gnu -static"

if (-not $env:CC) {
    if (Get-Command zig -ErrorAction SilentlyContinue) {
        $env:CC = $ZigCC
        if (-not $env:CXX) { $env:CXX = $ZigCXX }
    } elseif (Get-Command gcc -ErrorAction SilentlyContinue) {
        Write-Host "Zig not found, falling back to GCC for static build." -ForegroundColor Yellow
        $env:CC = "gcc"
        if (-not $env:CXX) { $env:CXX = "g++" }
    } else {
        Write-Error "Neither zig nor gcc found on PATH. Please install a C toolchain."
        exit 1
    }
}

$env:CGO_ENABLED = "1"

$outMoonrun = "moonrun_static.exe"
$outMoonbasic = "moonbasic_static.exe"

Write-Host "--- Building Compiler: $outMoonbasic ---" -ForegroundColor Cyan
# Compiler is usually pure Go, but we use static flags for maximum portability.
& go build -ldflags="-s -w -extldflags=-static" -o $outMoonbasic .

if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to build compiler."
    exit $LASTEXITCODE
}

Write-Host "--- Building Runtime: $outMoonrun ---" -ForegroundColor Cyan
Write-Host "CC=$($env:CC)"
Write-Host "CXX=$($env:CXX)"

$goArgs = @(
    "build",
    "-tags", "fullruntime",
    "-ldflags=-linkmode external -extldflags=-static",
    "-o", $outMoonrun,
    "./cmd/moonrun"
)

& go @goArgs
if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to build runtime."
    exit $LASTEXITCODE
}

Write-Host "Done! Binaries: $outMoonbasic, $outMoonrun" -ForegroundColor Green

# Verification
$MingwBin = Split-Path (Get-Command gcc).Source
Write-Host "--- Verifying Imports ---" -ForegroundColor Yellow
powershell -File scripts/verify_windows_pe_imports.ps1 -Exe $outMoonbasic -MingwBin $MingwBin
powershell -File scripts/verify_windows_pe_imports.ps1 -Exe $outMoonrun -MingwBin $MingwBin
