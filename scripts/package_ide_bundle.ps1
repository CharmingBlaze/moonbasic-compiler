# Package moonBASIC IDE + toolchain for Windows (maintainer helper).
param(
  [string]$Version = "dev"
)
$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
$IdeDir = Join-Path $Root "moonbasic ide"
$Stage = Join-Path $Root "dist\ide-bundle"
$Out = Join-Path $Root "moonbasic-$Version-ide-windows-amd64.zip"

$IdeExe = Join-Path $IdeDir "build\bin\moonbasic-ide.exe"
if (-not (Test-Path $IdeExe)) {
  throw "Build the IDE first: cd 'moonbasic ide'; npm ci; npm run langdata; wails build"
}
foreach ($f in @("dist\moonbasic.exe", "dist\moonrun.exe")) {
  if (-not (Test-Path (Join-Path $Root $f))) {
    throw "Build runtime into dist/ first (see scripts/release-windows.ps1)"
  }
}

if (Test-Path $Stage) { Remove-Item -Recurse -Force $Stage }
New-Item -ItemType Directory -Path $Stage | Out-Null
Copy-Item $IdeExe (Join-Path $Stage "moonbasic-ide.exe")
Copy-Item (Join-Path $Root "dist\moonbasic.exe") $Stage
Copy-Item (Join-Path $Root "dist\moonrun.exe") $Stage
Copy-Item (Join-Path $Root "packaging\README-IDE-RELEASE.txt") $Stage
Copy-Item (Join-Path $Root "packaging\START-IDE.bat") $Stage
Copy-Item (Join-Path $Root "packaging\START-IDE.sh") $Stage
if (Test-Path $Out) { Remove-Item -Force $Out }
Compress-Archive -Path (Join-Path $Stage "*") -DestinationPath $Out -Force
Write-Host "Wrote $Out"
