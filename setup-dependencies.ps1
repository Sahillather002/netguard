# NetGuard - Install All Dependencies
Write-Host "📦 Installing NetGuard Dependencies..." -ForegroundColor Cyan
Write-Host ""

# Check Python
Write-Host "Checking Python..." -ForegroundColor Yellow
try {
    $pythonVersion = python --version 2>&1
    Write-Host "  ✓ $pythonVersion" -ForegroundColor Green
} catch {
    Write-Host "  ✗ Python not found! Please install Python 3.11+" -ForegroundColor Red
    exit 1
}

# Check Go
Write-Host "Checking Go..." -ForegroundColor Yellow
try {
    $goVersion = go version
    Write-Host "  ✓ $goVersion" -ForegroundColor Green
} catch {
    Write-Host "  ✗ Go not found! Please install Go 1.21+" -ForegroundColor Red
    exit 1
}

# Check Node.js
Write-Host "Checking Node.js..." -ForegroundColor Yellow
try {
    $nodeVersion = node --version
    Write-Host "  ✓ Node.js $nodeVersion" -ForegroundColor Green
} catch {
    Write-Host "  ✗ Node.js not found! Please install Node.js 20+" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "=" * 70 -ForegroundColor Cyan
Write-Host "Installing Python Dependencies..." -ForegroundColor Cyan
Write-Host "=" * 70 -ForegroundColor Cyan

# Install Threat Detector dependencies
Write-Host ""
Write-Host "1. Threat Detector..." -ForegroundColor Yellow
Set-Location "$PSScriptRoot\services\threat-detector"
pip install -r requirements.txt
if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ Threat Detector dependencies installed" -ForegroundColor Green
} else {
    Write-Host "  ✗ Failed to install Threat Detector dependencies" -ForegroundColor Red
}

# Install Network Monitor dependencies
Write-Host ""
Write-Host "2. Network Monitor..." -ForegroundColor Yellow
Set-Location "$PSScriptRoot\services\network-monitor"
pip install -r requirements.txt
if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ Network Monitor dependencies installed" -ForegroundColor Green
} else {
    Write-Host "  ✗ Failed to install Network Monitor dependencies" -ForegroundColor Red
}

# Install Firewall Service dependencies
Write-Host ""
Write-Host "3. Firewall Service..." -ForegroundColor Yellow
Set-Location "$PSScriptRoot\services\firewall-service"
pip install -r requirements.txt
if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ Firewall Service dependencies installed" -ForegroundColor Green
} else {
    Write-Host "  ✗ Failed to install Firewall Service dependencies" -ForegroundColor Red
}

Write-Host ""
Write-Host "=" * 70 -ForegroundColor Cyan
Write-Host "Installing Go Dependencies..." -ForegroundColor Cyan
Write-Host "=" * 70 -ForegroundColor Cyan

# Install API Gateway dependencies
Write-Host ""
Write-Host "4. API Gateway..." -ForegroundColor Yellow
Set-Location "$PSScriptRoot\services\api-gateway"
go mod tidy
if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ API Gateway dependencies installed" -ForegroundColor Green
} else {
    Write-Host "  ✗ Failed to install API Gateway dependencies" -ForegroundColor Red
}

# Install Auth Service dependencies
Write-Host ""
Write-Host "5. Auth Service..." -ForegroundColor Yellow
Set-Location "$PSScriptRoot\services\auth-service"
go mod tidy
if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ Auth Service dependencies installed" -ForegroundColor Green
} else {
    Write-Host "  ✗ Failed to install Auth Service dependencies" -ForegroundColor Red
}

Write-Host ""
Write-Host "=" * 70 -ForegroundColor Cyan
Write-Host "Installing Frontend Dependencies..." -ForegroundColor Cyan
Write-Host "=" * 70 -ForegroundColor Cyan

# Install Frontend dependencies
Write-Host ""
Write-Host "6. Frontend (Next.js)..." -ForegroundColor Yellow
Set-Location "$PSScriptRoot\securecloud-nextjs"
npm install
if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ Frontend dependencies installed" -ForegroundColor Green
} else {
    Write-Host "  ✗ Failed to install Frontend dependencies" -ForegroundColor Red
}

# Return to project root
Set-Location $PSScriptRoot

Write-Host ""
Write-Host "=" * 70 -ForegroundColor Cyan
Write-Host "✅ All dependencies installed!" -ForegroundColor Green
Write-Host "=" * 70 -ForegroundColor Cyan
Write-Host ""
Write-Host "🚀 You can now start the platform with:" -ForegroundColor Cyan
Write-Host "   .\start-local.ps1" -ForegroundColor White
Write-Host ""