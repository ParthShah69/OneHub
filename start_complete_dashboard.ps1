# Complete Dashboard Startup Script with Real API Integration
Write-Host "🚀 Starting Complete Personalized Dashboard..." -ForegroundColor Green
Write-Host "================================================" -ForegroundColor Cyan

# Load environment variables
Write-Host "📋 Loading environment variables..." -ForegroundColor Yellow
if (Test-Path "env.local") {
    Get-Content "env.local" | ForEach-Object {
        if ($_ -match "^([^#][^=]+)=(.*)$") {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim()
            
            # Remove quotes if present
            if ($value.StartsWith('"') -and $value.EndsWith('"')) {
                $value = $value.Substring(1, $value.Length - 2)
            }
            
            [Environment]::SetEnvironmentVariable($key, $value, "Process")
            Write-Host "  ✓ $key" -ForegroundColor Gray
        }
    }
    Write-Host "✅ Environment variables loaded!" -ForegroundColor Green
} else {
    Write-Host "⚠️  env.local not found. Using mock data." -ForegroundColor Yellow
    Write-Host "   Create env.local with your API keys for real data." -ForegroundColor Gray
}

# Function to start a service
function Start-Service {
    param($ServiceName, $Path, $Port)
    
    Write-Host "🔧 Starting $ServiceName on port $Port..." -ForegroundColor Cyan
    
    # Start service in new window
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$Path'; go run simple_main.go" -WindowStyle Minimized
    Start-Sleep 2
}

# Start all backend services
Write-Host "`n🔧 Starting Backend Services..." -ForegroundColor Yellow
Start-Service "News Service (Real NewsAPI)" "F:\hackathon_2\services\news" "8001"
Start-Service "Jobs Service (Real LinkedIn API)" "F:\hackathon_2\services\jobs" "8002"
Start-Service "Videos Service (Real YouTube API)" "F:\hackathon_2\services\videos" "8003"
Start-Service "Deals Service (Real Amazon API)" "F:\hackathon_2\services\deals" "8004"
Start-Service "Movies Service (Real TMDB API)" "F:\hackathon_2\services\movies" "8008"
Start-Service "Food Service (Real Spoonacular API)" "F:\hackathon_2\services\food" "8009"
Start-Service "Recommendation Service (Wolfram)" "F:\hackathon_2\services\recommendation" "8005"
Start-Service "User Service" "F:\hackathon_2\services\user" "8006"
Start-Service "NFT Service (Verbwire)" "F:\hackathon_2\services\nft" "8007"
Start-Service "API Gateway" "F:\hackathon_2\gateway" "8080"

Write-Host "✅ All backend services started!" -ForegroundColor Green

# Wait a bit for services to initialize
Write-Host "⏳ Waiting for services to initialize..." -ForegroundColor Yellow
Start-Sleep 5

# Start frontend
Write-Host "`n📱 Starting React Frontend..." -ForegroundColor Yellow
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd 'F:\hackathon_2\frontend'; npm install; npm start" -WindowStyle Normal

Write-Host "`n🎉 Dashboard is starting up!" -ForegroundColor Green
Write-Host "================================================" -ForegroundColor Cyan
Write-Host "🌐 Frontend: http://localhost:3000" -ForegroundColor White
Write-Host "🔗 API Gateway: http://localhost:8080" -ForegroundColor White
Write-Host "`n📋 API Status:" -ForegroundColor Yellow

# Check API key status
$newsKey = [Environment]::GetEnvironmentVariable('NEWS_API_KEY', 'Process')
$youtubeKey = [Environment]::GetEnvironmentVariable('YOUTUBE_API_KEY', 'Process')
$linkedinKey = [Environment]::GetEnvironmentVariable('LINKEDIN_API_KEY', 'Process')
$amazonKey = [Environment]::GetEnvironmentVariable('AMAZON_API_KEY', 'Process')
$tmdbKey = [Environment]::GetEnvironmentVariable('TMDB_API_KEY', 'Process')
$recipeKey = [Environment]::GetEnvironmentVariable('RECIPE_API_KEY', 'Process')

Write-Host "  📰 News API: $(if ($newsKey -and $newsKey -ne 'your_newsapi_key_here') { '✅ REAL DATA' } else { '⚠️ STATIC DATA' })" -ForegroundColor $(if ($newsKey -and $newsKey -ne 'your_newsapi_key_here') { 'Green' } else { 'Yellow' })
Write-Host "  🎥 YouTube API: $(if ($youtubeKey -and $youtubeKey -ne 'your_youtube_api_key_here') { '✅ REAL DATA' } else { '⚠️ STATIC DATA' })" -ForegroundColor $(if ($youtubeKey -and $youtubeKey -ne 'your_youtube_api_key_here') { 'Green' } else { 'Yellow' })
Write-Host "  💼 LinkedIn API: $(if ($linkedinKey -and $linkedinKey -ne 'your_linkedin_api_key_here') { '✅ REAL DATA' } else { '⚠️ STATIC DATA' })" -ForegroundColor $(if ($linkedinKey -and $linkedinKey -ne 'your_linkedin_api_key_here') { 'Green' } else { 'Yellow' })
Write-Host "  🛒 Amazon API: $(if ($amazonKey -and $amazonKey -ne 'your_amazon_api_key_here') { '✅ REAL DATA' } else { '⚠️ STATIC DATA' })" -ForegroundColor $(if ($amazonKey -and $amazonKey -ne 'your_amazon_api_key_here') { 'Green' } else { 'Yellow' })
Write-Host "  🎬 TMDB API: $(if ($tmdbKey -and $tmdbKey -ne 'your_tmdb_api_key_here') { '✅ REAL DATA' } else { '⚠️ STATIC DATA' })" -ForegroundColor $(if ($tmdbKey -and $tmdbKey -ne 'your_tmdb_api_key_here') { 'Green' } else { 'Yellow' })
Write-Host "  🍽️ Recipe API: $(if ($recipeKey -and $recipeKey -ne 'your_recipe_api_key_here') { '✅ REAL DATA' } else { '⚠️ STATIC DATA' })" -ForegroundColor $(if ($recipeKey -and $recipeKey -ne 'your_recipe_api_key_here') { 'Green' } else { 'Yellow' })

Write-Host "`n🔧 To get real data:" -ForegroundColor Yellow
Write-Host "1. Get API keys from:" -ForegroundColor White
Write-Host "   📰 NewsAPI: https://newsapi.org/register (Free)" -ForegroundColor Gray
Write-Host "   🎥 YouTube: https://console.developers.google.com/ (Free)" -ForegroundColor Gray
Write-Host "   💼 LinkedIn: https://developer.linkedin.com/ (Free)" -ForegroundColor Gray
Write-Host "   🛒 Amazon: https://webservices.amazon.com/ (Free)" -ForegroundColor Gray
Write-Host "   🎬 TMDB: https://www.themoviedb.org/settings/api (Free)" -ForegroundColor Gray
Write-Host "   🍽️ Spoonacular: https://spoonacular.com/food-api (Free)" -ForegroundColor Gray
Write-Host "2. Add them to env.local file" -ForegroundColor White
Write-Host "3. Restart this script" -ForegroundColor White

Write-Host "`n🧪 Test API endpoints:" -ForegroundColor Yellow
Write-Host "curl http://localhost:8080/api/news/trending" -ForegroundColor Gray
Write-Host "curl http://localhost:8080/api/videos?category=technology" -ForegroundColor Gray
Write-Host "curl http://localhost:8080/api/jobs?category=ai" -ForegroundColor Gray
Write-Host "curl http://localhost:8080/api/deals?category=electronics" -ForegroundColor Gray
Write-Host "curl http://localhost:8080/api/movies?category=popular" -ForegroundColor Gray
Write-Host "curl http://localhost:8080/api/food?category=healthy" -ForegroundColor Gray

Write-Host "`n✨ Your personalized dashboard is ready!" -ForegroundColor Green
