# NetGuard - Start All Services Locally (No Docker)
Write-Host "🚀 Starting NetGuard Platform Locally..." -ForegroundColor Cyan
Write-Host ""

# Start API Gateway
Write-Host "1. Starting API Gateway (Go) on port 8080..." -ForegroundColor Yellow
Start-Process pwsh -ArgumentList "-NoExit", "-Command", "cd '$PSScriptRoot\services\api-gateway'; Write-Host '🚪 API Gateway Starting...' -ForegroundColor Green; go run main.go"
Start-Sleep -Seconds 2

# Start Auth Service
Write-Host "2. Starting Auth Service (Go) on port 8081..." -ForegroundColor Yellow
Start-Process pwsh -ArgumentList "-NoExit", "-Command", "cd '$PSScriptRoot\services\auth-service'; Write-Host '🔐 Auth Service Starting...' -ForegroundColor Green; go run main.go"
Start-Sleep -Seconds 2

# Start Threat Detector
Write-Host "3. Starting Threat Detector (Python) on port 8082..." -ForegroundColor Yellow
Start-Process pwsh -ArgumentList "-NoExit", "-Command", "cd '$PSScriptRoot\services\threat-detector'; Write-Host '🔥 Threat Detector Starting...' -ForegroundColor Green; python main.py"
Start-Sleep -Seconds 2

# Start Network Monitor
Write-Host "4. Starting Network Monitor (Python) on port 8083..." -ForegroundColor Yellow
Start-Process pwsh -ArgumentList "-NoExit", "-Command", "cd '$PSScriptRoot\services\network-monitor'; Write-Host '🌐 Network Monitor Starting...' -ForegroundColor Green; python main.py"
Start-Sleep -Seconds 2

# Start Firewall Service
Write-Host "5. Starting Firewall Service (Python) on port 8084..." -ForegroundColor Yellow
Start-Process pwsh -ArgumentList "-NoExit", "-Command", "cd '$PSScriptRoot\services\firewall-service'; Write-Host '🛡️ Firewall Service Starting...' -ForegroundColor Green; python main.py"
Start-Sleep -Seconds 2

# Start Frontend
Write-Host "6. Starting Frontend (Next.js) on port 3000..." -ForegroundColor Yellow
Start-Process pwsh -ArgumentList "-NoExit", "-Command", "cd '$PSScriptRoot\securecloud-nextjs'; Write-Host '🎨 Frontend Starting...' -ForegroundColor Green; npm run dev"

Write-Host ""
Write-Host "=" * 70 -ForegroundColor Cyan
Write-Host "🎉 All services are starting in separate windows!" -ForegroundColor Green
Write-Host "=" * 70 -ForegroundColor Cyan
Write-Host ""
Write-Host "⏳ Please wait 30-60 seconds for all services to fully start..." -ForegroundColor Yellow
Write-Host ""
Write-Host "📊 Access Points:" -ForegroundColor Cyan
Write-Host "  Frontend:        http://localhost:3000" -ForegroundColor White
Write-Host "  API Gateway:     http://localhost:8080" -ForegroundColor White
Write-Host "  Auth Service:    http://localhost:8081" -ForegroundColor White
Write-Host "  Threat Detector: http://localhost:8082" -ForegroundColor White
Write-Host "  Network Monitor: http://localhost:8083" -ForegroundColor White
Write-Host "  Firewall:        http://localhost:8084" -ForegroundColor White
Write-Host ""
Write-Host "🔐 Default Login:" -ForegroundColor Cyan
Write-Host "  Email:    admin@netguard.com" -ForegroundColor White
Write-Host "  Password: password" -ForegroundColor White
Write-Host ""
Write-Host "🛑 To stop: Close all PowerShell windows or press Ctrl+C in each" -ForegroundColor Yellow
Write-Host ""
Write-Host "📖 Documentation: README-COMPLETE.md" -ForegroundColor Cyan
Write-Host ""