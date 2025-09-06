# PowerShell script to fix GoFr dependencies and run all services
Write-Host "🔧 Fixing GoFr Dependencies..." -ForegroundColor Yellow

# Function to fix go.mod and run service
function Fix-And-Run-Service {
    param($ServiceName, $Path, $Port)
    
    Write-Host "Fixing $ServiceName..." -ForegroundColor Cyan
    
    # Navigate to service directory
    Set-Location $Path
    
    # Initialize go.mod if it doesn't exist
    if (-not (Test-Path "go.mod")) {
        go mod init $ServiceName
    }
    
    # Download GoFr dependency
    Write-Host "Downloading GoFr for $ServiceName..." -ForegroundColor Yellow
    go mod download gofr.dev
    
    # Tidy up dependencies
    go mod tidy
    
    # Run the service
    Write-Host "Starting $ServiceName on port $Port..." -ForegroundColor Green
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$Path'; go run main.go" -WindowStyle Minimized
    Start-Sleep 3
}

# Fix and start all services
Write-Host "🚀 Starting all services with GoFr..." -ForegroundColor Green

Fix-And-Run-Service "news-service" "F:\hackathon_2\services\news" "8001"
Fix-And-Run-Service "jobs-service" "F:\hackathon_2\services\jobs" "8002"
Fix-And-Run-Service "videos-service" "F:\hackathon_2\services\videos" "8003"
Fix-And-Run-Service "deals-service" "F:\hackathon_2\services\deals" "8004"
Fix-And-Run-Service "recommendation-service" "F:\hackathon_2\services\recommendation" "8005"
Fix-And-Run-Service "user-service" "F:\hackathon_2\services\user" "8006"
Fix-And-Run-Service "nft-service" "F:\hackathon_2\services\nft" "8007"
Fix-And-Run-Service "api-gateway" "F:\hackathon_2\gateway" "8080"

Write-Host "✅ All services started with GoFr!" -ForegroundColor Green
Write-Host "📱 Starting React Frontend..." -ForegroundColor Cyan

# Start frontend
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd 'F:\hackathon_2\frontend'; npm install; npm start" -WindowStyle Normal

Write-Host "🎉 Dashboard is starting up!" -ForegroundColor Green
Write-Host "Frontend: http://localhost:3000" -ForegroundColor Cyan
Write-Host "API Gateway: http://localhost:8080" -ForegroundColor Cyan
