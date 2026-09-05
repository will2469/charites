# Charites Automated Installer for Windows (PowerShell)
# Usage: irm https://raw.githubusercontent.com/will2469/charites/main/scripts/install.ps1 | iex

$ErrorActionPreference = "Stop"

$Repo = "will2469/charites"
$GithubUrl = "https://github.com/$Repo"

Write-Host "==> Charites Compiler & Static Linter Windows Installer" -ForegroundColor Cyan

# 1. Architecture Check
if (-not [System.Environment]::Is64BitOperatingSystem) {
    Write-Error "Charites requires 64-bit Windows (x64/amd64)."
    exit 1
}
$Arch = "amd64"

# 2. Determine Latest Release Tag
Write-Host "==> Checking latest release..." -ForegroundColor Cyan
$Tag = ""
try {
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    $Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -UseBasicParsing
    $Tag = $Release.tag_name
} catch {
    $Tag = "v1.0.0"
}

if ([string]::IsNullOrWhiteSpace($Tag) -or $Tag -eq "latest") {
    $Tag = "v1.0.0"
}
Write-Host "==> Selected release: $Tag" -ForegroundColor Green

# 3. Prepare Download Path
$ArchiveName = "charites_${Tag}_windows_${Arch}.zip"
$DownloadUrl = "$GithubUrl/releases/download/$Tag/$ArchiveName"

$TempDir = Join-Path ([System.IO.Path]::GetTempPath()) ([System.IO.Path]::GetRandomFileName())
New-Item -ItemType Directory -Path $TempDir -Force | Out-Null
$ZipPath = Join-Path $TempDir $ArchiveName

try {
    Write-Host "==> Downloading $DownloadUrl..." -ForegroundColor Cyan
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $ZipPath -UseBasicParsing

    Write-Host "==> Extracting binary..." -ForegroundColor Cyan
    Expand-Archive -Path $ZipPath -DestinationPath $TempDir -Force

    $InstallDir = Join-Path $env:LOCALAPPDATA "charites\bin"
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null

    $BinarySrc = Join-Path $TempDir "charites.exe"
    if (-not (Test-Path $BinarySrc)) {
        Write-Error "Could not find charites.exe in downloaded archive."
        exit 1
    }

    $BinaryDest = Join-Path $InstallDir "charites.exe"
    Copy-Item -Path $BinarySrc -Destination $BinaryDest -Force

    # 4. Automatically Configure Windows User PATH
    $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($null -eq $UserPath) {
        $UserPath = ""
    }

    $CleanPath = $InstallDir.TrimEnd('\')
    if ($UserPath -split ';' -notcontains $CleanPath) {
        Write-Host "==> Adding $InstallDir to User Environment PATH..." -ForegroundColor Yellow
        $NewUserPath = if ([string]::IsNullOrWhiteSpace($UserPath)) { $CleanPath } else { "$UserPath;$CleanPath" }
        [Environment]::SetEnvironmentVariable("Path", $NewUserPath, "User")
    }

    # Inject into current PowerShell session so it's instantly usable
    if ($env:Path -split ';' -notcontains $CleanPath) {
        $env:Path = "$CleanPath;" + $env:Path
    }

    Write-Host ""
    Write-Host "==> Charites $Tag installed successfully to $BinaryDest!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Quick Start:" -ForegroundColor Cyan
    Write-Host "  charites --help"
    Write-Host "  charites scan ."
    Write-Host ""

} finally {
    if (Test-Path $TempDir) {
        Remove-Item -Path $TempDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}
