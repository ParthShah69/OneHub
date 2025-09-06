# PowerShell script to start all services
Write-Host "🚀 Starting Personalized Dashboard Services..." -ForegroundColor Green

# Function to start a service in background
function Start-Service {
    param($ServiceName, $Path, $Port)
    
    Write-Host "Starting $ServiceName on port $Port..." -ForegroundColor Yellow
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$Path'; go run simple_main.go" -WindowStyle Minimized
    Start-Sleep 2
}

# Start all services
Start-Service "News Service" "F:\hackathon_2\services\news" "8001"
Start-Service "Jobs Service" "F:\hackathon_2\services\jobs" "8002"
Start-Service "Videos Service" "F:\hackathon_2\services\videos" "8003"
Start-Service "Deals Service" "F:\hackathon_2\services\deals" "8004"
Start-Service "Recommendation Service" "F:\hackathon_2\services\recommendation" "8005"
Start-Service "User Service" "F:\hackathon_2\services\user" "8006"
Start-Service "NFT Service" "F:\hackathon_2\services\nft" "8007"
Start-Service "API Gateway" "F:\hackathon_2\gateway" "8080"

Write-Host "✅ All services started!" -ForegroundColor Green
Write-Host "📱 Frontend will start next..." -ForegroundColor Cyan

# Start frontend
Write-Host "Starting React Frontend..." -ForegroundColor Yellow
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd 'F:\hackathon_2\frontend'; npm install; npm start" -WindowStyle Normal

Write-Host "🎉 Dashboard is starting up!" -ForegroundColor Green
Write-Host "Frontend: http://localhost:3000" -ForegroundColor Cyan
Write-Host "API Gateway: http://localhost:8080" -ForegroundColor Cyan
