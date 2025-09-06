# Quick Fix for 404 Errors
Write-Host "🚀 Quick Fix for Dashboard Issues..." -ForegroundColor Green

# Kill all existing Go processes
Write-Host "🧹 Stopping existing services..." -ForegroundColor Yellow
Get-Process -Name "go" -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
Get-Process -Name "node" -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
Start-Sleep 3

# Load environment variables
if (Test-Path "env.local") {
    Get-Content "env.local" | ForEach-Object {
        if ($_ -match "^([^#][^=]+)=(.*)$") {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim()
            if ($value.StartsWith('"') -and $value.EndsWith('"')) {
                $value = $value.Substring(1, $value.Length - 2)
            }
            [Environment]::SetEnvironmentVariable($key, $value, "Process")
        }
    }
}

Write-Host "✅ Environment loaded!" -ForegroundColor Green

# Start services in background
Write-Host "🔧 Starting services..." -ForegroundColor Cyan

# Start all services
$services = @(
    @{Name="User"; Path="F:\hackathon_2\services\user"; Port="8006"},
    @{Name="News"; Path="F:\hackathon_2\services\news"; Port="8001"},
    @{Name="Jobs"; Path="F:\hackathon_2\services\jobs"; Port="8002"},
    @{Name="Videos"; Path="F:\hackathon_2\services\videos"; Port="8003"},
    @{Name="Deals"; Path="F:\hackathon_2\services\deals"; Port="8004"},
    @{Name="Movies"; Path="F:\hackathon_2\services\movies"; Port="8008"},
    @{Name="Food"; Path="F:\hackathon_2\services\food"; Port="8009"},
    @{Name="Recommendation"; Path="F:\hackathon_2\services\recommendation"; Port="8005"},
    @{Name="NFT"; Path="F:\hackathon_2\services\nft"; Port="8007"},
    @{Name="Gateway"; Path="F:\hackathon_2\gateway"; Port="8080"}
)

foreach ($service in $services) {
    Write-Host "Starting $($service.Name) service..." -ForegroundColor Gray
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$($service.Path)'; go run simple_main.go" -WindowStyle Minimized
    Start-Sleep 2
}

Write-Host "⏳ Waiting for services to start..." -ForegroundColor Yellow
Start-Sleep 15

# Test critical services
Write-Host "`n🧪 Testing services..." -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get -TimeoutSec 5
    Write-Host "✅ API Gateway: Working" -ForegroundColor Green
} catch {
    Write-Host "❌ API Gateway: Not responding" -ForegroundColor Red
}

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8008/health" -Method Get -TimeoutSec 5
    Write-Host "✅ Movies Service: Working" -ForegroundColor Green
} catch {
    Write-Host "❌ Movies Service: Not responding" -ForegroundColor Red
}

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8009/health" -Method Get -TimeoutSec 5
    Write-Host "✅ Food Service: Working" -ForegroundColor Green
} catch {
    Write-Host "❌ Food Service: Not responding" -ForegroundColor Red
}

# Start frontend
Write-Host "`n📱 Starting frontend..." -ForegroundColor Cyan
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd 'F:\hackathon_2\frontend'; npm start" -WindowStyle Normal

Write-Host "`n🎉 Quick fix completed!" -ForegroundColor Green
Write-Host "Frontend: http://localhost:3000" -ForegroundColor White
Write-Host "API Gateway: http://localhost:8080" -ForegroundColor White

Write-Host "`n📋 If you still see 404 errors:" -ForegroundColor Yellow
Write-Host "1. Wait 30 seconds for all services to fully start" -ForegroundColor White
Write-Host "2. Refresh the browser page" -ForegroundColor White
Write-Host "3. Check browser console for any remaining errors" -ForegroundColor White

Write-Host "`n🔑 To get real data (not static):" -ForegroundColor Yellow
Write-Host "1. Get TMDB API key: https://www.themoviedb.org/settings/api" -ForegroundColor White
Write-Host "2. Get Recipe API key: https://spoonacular.com/food-api" -ForegroundColor White
Write-Host "3. Add keys to env.local file" -ForegroundColor White
Write-Host "4. Restart services" -ForegroundColor White
