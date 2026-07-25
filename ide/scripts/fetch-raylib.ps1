# Download raylib 5.5 DLL for purego moonrun (CGO_ENABLED=0 builds).
$ErrorActionPreference = "Stop"

$ideRoot = Split-Path $PSScriptRoot -Parent
$outDir = Join-Path $ideRoot "toolchain"
$dllPath = Join-Path $outDir "raylib.dll"

if (Test-Path $dllPath) {
    Write-Host "raylib.dll already present in toolchain/"
    exit 0
}

$url = "https://github.com/raysan5/raylib/releases/download/5.5/raylib-5.5_win64_mingw-w64.zip"
$zip = Join-Path $env:TEMP "raylib-5.5_win64_mingw-w64.zip"
$extract = Join-Path $env:TEMP "raylib55-fetch"

Write-Host "Downloading raylib 5.5 for moonrun..."
Invoke-WebRequest -Uri $url -OutFile $zip
if (Test-Path $extract) { Remove-Item $extract -Recurse -Force }
Expand-Archive -Path $zip -DestinationPath $extract -Force

$dll = Get-ChildItem -Recurse $extract -Filter "raylib.dll" | Select-Object -First 1
if (-not $dll) { throw "raylib.dll not found in $url" }

New-Item -ItemType Directory -Force -Path $outDir | Out-Null
Copy-Item $dll.FullName $dllPath -Force
Copy-Item $dll.FullName (Join-Path $ideRoot "raylib.dll") -Force
Write-Host "Installed raylib.dll to $outDir"
