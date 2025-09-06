# Fix Services Startup Script
Write-Host "🔧 Fixing Services and Starting Dashboard..." -ForegroundColor Green

# Load environment variables
if (Test-Path "env.local") {
    Write-Host "📋 Loading environment variables..." -ForegroundColor Yellow
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
    Write-Host "✅ Environment variables loaded!" -ForegroundColor Green
}

# Function to start a service with error handling
function Start-Service {
    param($ServiceName, $Path, $Port)
    
    Write-Host "🔧 Starting $ServiceName on port $Port..." -ForegroundColor Cyan
    
    # Check if port is already in use
    $portInUse = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue
    if ($portInUse) {
        Write-Host "⚠️  Port $Port is already in use. Stopping existing process..." -ForegroundColor Yellow
        $process = Get-Process -Id $portInUse.OwningProcess -ErrorAction SilentlyContinue
        if ($process) {
            Stop-Process -Id $process.Id -Force -ErrorAction SilentlyContinue
            Start-Sleep 2
        }
    }
    
    # Start service in new window
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$Path'; go run simple_main.go" -WindowStyle Minimized
    Start-Sleep 3
}

# Kill any existing Go processes
Write-Host "🧹 Cleaning up existing processes..." -ForegroundColor Yellow
Get-Process -Name "go" -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
Start-Sleep 2

# Start all backend services in correct order
Write-Host "`n🔧 Starting Backend Services..." -ForegroundColor Yellow

# Start User Service first (needed for preferences)
Start-Service "User Service" "F:\hackathon_2\services\user" "8006"

# Start content services
Start-Service "News Service (Real NewsAPI)" "F:\hackathon_2\services\news" "8001"
Start-Service "Jobs Service" "F:\hackathon_2\services\jobs" "8002"
Start-Service "Videos Service" "F:\hackathon_2\services\videos" "8003"
Start-Service "Deals Service" "F:\hackathon_2\services\deals" "8004"
Start-Service "Movies Service (TMDB API)" "F:\hackathon_2\services\movies" "8008"
Start-Service "Food Service (Recipe API)" "F:\hackathon_2\services\food" "8009"

# Start supporting services
Start-Service "Recommendation Service" "F:\hackathon_2\services\recommendation" "8005"
Start-Service "NFT Service" "F:\hackathon_2\services\nft" "8007"

# Start API Gateway last
Start-Service "API Gateway" "F:\hackathon_2\gateway" "8080"

Write-Host "✅ All backend services started!" -ForegroundColor Green

# Wait for services to initialize
Write-Host "⏳ Waiting for services to initialize..." -ForegroundColor Yellow
Start-Sleep 10

# Test services
Write-Host "`n🧪 Testing Services..." -ForegroundColor Yellow

$services = @(
    @{Name="User Service"; Port="8006"; Path="/health"},
    @{Name="News Service"; Port="8001"; Path="/health"},
    @{Name="Jobs Service"; Port="8002"; Path="/health"},
    @{Name="Videos Service"; Port="8003"; Path="/health"},
    @{Name="Deals Service"; Port="8004"; Path="/health"},
    @{Name="Movies Service"; Port="8008"; Path="/health"},
    @{Name="Food Service"; Port="8009"; Path="/health"},
    @{Name="Recommendation Service"; Port="8005"; Path="/health"},
    @{Name="NFT Service"; Port="8007"; Path="/health"},
    @{Name="API Gateway"; Port="8080"; Path="/health"}
)

foreach ($service in $services) {
    try {
        $response = Invoke-RestMethod -Uri "http://localhost:$($service.Port)$($service.Path)" -Method Get -TimeoutSec 5
        Write-Host "✅ $($service.Name): Running" -ForegroundColor Green
    } catch {
        Write-Host "❌ $($service.Name): Not responding" -ForegroundColor Red
    }
}

# Start frontend
Write-Host "`n📱 Starting React Frontend..." -ForegroundColor Yellow
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd 'F:\hackathon_2\frontend'; npm start" -WindowStyle Normal

Write-Host "`n🎉 Dashboard is starting up!" -ForegroundColor Green
Write-Host "================================================" -ForegroundColor Cyan
Write-Host "🌐 Frontend: http://localhost:3000" -ForegroundColor White
Write-Host "🔗 API Gateway: http://localhost:8080" -ForegroundColor White

Write-Host "`n📋 API Status:" -ForegroundColor Yellow
$newsKey = [Environment]::GetEnvironmentVariable('NEWS_API_KEY', 'Process')
Write-Host "  📰 News API: $(if ($newsKey -and $newsKey -ne 'your_newsapi_key_here') { '✅ REAL DATA' } else { '⚠️ STATIC DATA' })" -ForegroundColor $(if ($newsKey -and $newsKey -ne 'your_newsapi_key_here') { 'Green' } else { 'Yellow' })

Write-Host "`n🔧 If you see 404 errors:" -ForegroundColor Yellow
Write-Host "1. Wait 30 seconds for all services to start" -ForegroundColor White
Write-Host "2. Refresh the frontend page" -ForegroundColor White
Write-Host "3. Check that all services are running" -ForegroundColor White
Write-Host "4. Run: .\test_complete_system.ps1" -ForegroundColor White

Write-Host "`n✨ Your personalized dashboard is ready!" -ForegroundColor Green
