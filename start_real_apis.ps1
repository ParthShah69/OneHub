# PowerShell script to start all services with real API integration
Write-Host "🚀 Starting Personalized Dashboard with Real API Integration..." -ForegroundColor Green

# Load environment variables
Write-Host "Loading environment variables..." -ForegroundColor Cyan
& ".\load_env.ps1"

# Function to start a service with environment variables
function Start-Service {
    param($ServiceName, $Path, $Port)
    
    Write-Host "Starting $ServiceName on port $Port..." -ForegroundColor Cyan
    
    # Start service in new window with environment variables
    $envVars = @()
    Get-ChildItem Env: | Where-Object { $_.Name -match "API_KEY|SERVICE_URL" } | ForEach-Object {
        $envVars += "$($_.Name)=$($_.Value)"
    }
    
    $envString = $envVars -join "; "
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$Path'; $envString; go run simple_main.go" -WindowStyle Minimized
    Start-Sleep 2
}

# Start all services
Start-Service "News Service (Real NewsAPI)" "F:\hackathon_2\services\news" "8001"
Start-Service "Jobs Service" "F:\hackathon_2\services\jobs" "8002"
Start-Service "Videos Service (Real YouTube API)" "F:\hackathon_2\services\videos" "8003"
Start-Service "Deals Service" "F:\hackathon_2\services\deals" "8004"
Start-Service "Recommendation Service" "F:\hackathon_2\services\recommendation" "8005"
Start-Service "User Service" "F:\hackathon_2\services\user" "8006"
Start-Service "NFT Service" "F:\hackathon_2\services\nft" "8007"
Start-Service "API Gateway" "F:\hackathon_2\gateway" "8080"

Write-Host "✅ All services started!" -ForegroundColor Green
Write-Host "📱 Starting React Frontend..." -ForegroundColor Cyan

# Start frontend
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd 'F:\hackathon_2\frontend'; npm install; npm start" -WindowStyle Normal

Write-Host "🎉 Dashboard is starting up!" -ForegroundColor Green
Write-Host "Frontend: http://localhost:3000" -ForegroundColor Cyan
Write-Host "API Gateway: http://localhost:8080" -ForegroundColor Cyan
Write-Host ""
Write-Host "📋 API Status:" -ForegroundColor Yellow
Write-Host "  News API: $([Environment]::GetEnvironmentVariable('NEWS_API_KEY', 'Process'))" -ForegroundColor White
Write-Host "  YouTube API: $([Environment]::GetEnvironmentVariable('YOUTUBE_API_KEY', 'Process'))" -ForegroundColor White
Write-Host ""
Write-Host "🔧 To get real data:" -ForegroundColor Yellow
Write-Host "1. Get API keys from:" -ForegroundColor White
Write-Host "   - NewsAPI: https://newsapi.org/register" -ForegroundColor Gray
Write-Host "   - YouTube: https://console.developers.google.com/" -ForegroundColor Gray
Write-Host "2. Add them to env.local file" -ForegroundColor White
Write-Host "3. Restart this script" -ForegroundColor White
