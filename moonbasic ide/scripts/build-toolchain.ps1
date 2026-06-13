# Build moonbasic/moonrun into moonbasic ide/toolchain/ for local IDE testing.
$ErrorActionPreference = "Stop"

$ideRoot = Split-Path $PSScriptRoot -Parent
$repoRoot = Split-Path $ideRoot -Parent
$outDir = Join-Path $ideRoot "toolchain"

New-Item -ItemType Directory -Force -Path $outDir | Out-Null

Push-Location $repoRoot
try {
    Write-Host "Building moonbasic.exe (compiler-only, CGO_ENABLED=0)..."
    $env:CGO_ENABLED = "0"
    go build -o (Join-Path $outDir "moonbasic.exe") .

    Write-Host "Building moonrun.exe (fullruntime)..."
    $env:CGO_ENABLED = "1"
    go build -tags fullruntime -o (Join-Path $outDir "moonrun.exe") ./cmd/moonrun
    Write-Host "Done. Binaries in: $outDir"
}
catch {
    Write-Warning "moonrun build failed (CGO may be required). moonbasic.exe may still be usable for check/compile/LSP."
    Write-Warning $_.Exception.Message
    if (Test-Path (Join-Path $outDir "moonrun.exe")) {
        Write-Host "Fetching raylib.dll for existing moonrun..."
        & (Join-Path $PSScriptRoot "fetch-raylib.ps1")
    }
}
finally {
    Pop-Location
}

if (Test-Path (Join-Path $outDir "moonrun.exe")) {
    $dll = Join-Path $outDir "raylib.dll"
    if (-not (Test-Path $dll)) {
        & (Join-Path $PSScriptRoot "fetch-raylib.ps1")
    }
}

if (Test-Path (Join-Path $outDir "moonbasic.exe")) {
    & (Join-Path $outDir "moonbasic.exe") --version
}
