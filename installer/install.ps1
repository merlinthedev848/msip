# MSIP Telephony Platform Windows Installer (PowerShell)
$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir
Set-Location $ProjectRoot

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "           MSIP Telephony Platform Installer (Windows)    " -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "Project Root: $ProjectRoot"

# Detect Docker Compose
$ComposeCmd = $null
if (Get-Command "docker" -ErrorAction SilentlyContinue) {
    try {
        $null = docker compose version 2>$null
        $ComposeCmd = "docker compose"
    } catch {}
}

if (-not $ComposeCmd -and (Get-Command "docker-compose" -ErrorAction SilentlyContinue)) {
    $ComposeCmd = "docker-compose"
}

if (-not $ComposeCmd) {
    Write-Error "Neither 'docker compose' nor 'docker-compose' was found on your system. Please install Docker Desktop."
    exit 1
}

Write-Host "Using Docker command: $ComposeCmd" -ForegroundColor Green

# Create default .env if missing
if (-not (Test-Path ".env")) {
    Write-Host "Creating default .env configuration..." -ForegroundColor Yellow
    @"
DB_USER=pbxuser
DB_PASS=pbxpass
DB_NAME=pbx
EXTERNAL_IP=127.0.0.1
"@ | Out-File -Encoding utf8 ".env"
}

# Ensure telephony-core certificates directory exists
$CertsDir = Join-Path $ProjectRoot "telephony-core\certs"
if (-not (Test-Path $CertsDir)) {
    New-Item -ItemType Directory -Path $CertsDir | Out-Null
}

$WssPem = Join-Path $CertsDir "wss.pem"
if (-not (Test-Path $WssPem) -and (Get-Command "openssl" -ErrorAction SilentlyContinue)) {
    Write-Host "Generating WebRTC WSS TLS certificates..." -ForegroundColor Yellow
    $KeyPem = Join-Path $CertsDir "key.pem"
    $CertPem = Join-Path $CertsDir "cert.pem"
    openssl req -x509 -newkey rsa:4096 -keyout $KeyPem -out $CertPem -days 365 -nodes -subj "/CN=pbx.local"
    Get-Content $CertPem, $KeyPem | Set-Content $WssPem
    Write-Host "WebRTC TLS certificates generated successfully." -ForegroundColor Green
}

# Build and start services
Write-Host "Building local services (Backend & Frontend)..." -ForegroundColor Yellow
Invoke-Expression "$ComposeCmd build"

Write-Host "Starting all telephony containers..." -ForegroundColor Yellow
Invoke-Expression "$ComposeCmd up -d --remove-orphans"

Write-Host ""
Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "   MSIP Telephony Platform Deployed Successfully!         " -ForegroundColor Green
Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host " Web Management Dashboard : http://localhost"
Write-Host " API Gateway (KrakenD)    : http://localhost:8080"
Write-Host " SIP Signaling (Kamailio) : UDP/TCP 5060"
Write-Host " FreeSWITCH Media / WSS   : TCP 5066 / WSS 7443"
Write-Host " Homer SIP Packet Capture : http://localhost:9080"
Write-Host " Turn / STUN Server       : UDP 3478"
Write-Host " Default Superadmin       : admin@pbx.local / admin123"
Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host ""
Invoke-Expression "$ComposeCmd ps"
