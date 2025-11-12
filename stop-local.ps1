# NetGuard - Stop All Local Services
Write-Host "🛑 Stopping NetGuard Services..." -ForegroundColor Red
Write-Host ""

# Kill processes on specific ports
$ports = @(3000, 8080, 8081, 8082, 8083, 8084)

foreach ($port in $ports) {
    Write-Host "Stopping service on port $port..." -ForegroundColor Yellow
    
    # Find process using the port
    $process = Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess -Unique
    
    if ($process) {
        foreach ($pid in $process) {
            try {
                Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue
                Write-Host "  ✓ Stopped process $pid on port $port" -ForegroundColor Green
            } catch {
                Write-Host "  ✗ Could not stop process $pid" -ForegroundColor Red
            }
        }
    } else {
        Write-Host "  - No process found on port $port" -ForegroundColor Gray
    }
}

Write-Host ""
Write-Host "✅ All services stopped!" -ForegroundColor Green
Write-Host ""