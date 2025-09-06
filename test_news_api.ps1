# Test NewsAPI Integration
Write-Host "🧪 Testing NewsAPI Integration..." -ForegroundColor Green

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

$newsApiKey = [Environment]::GetEnvironmentVariable('NEWS_API_KEY', 'Process')
Write-Host "NewsAPI Key: $newsApiKey" -ForegroundColor Yellow

if ($newsApiKey -and $newsApiKey -ne "your_newsapi_key_here") {
    Write-Host "✅ NewsAPI Key found!" -ForegroundColor Green
    
    # Test direct API call
    Write-Host "`n🔍 Testing direct NewsAPI call..." -ForegroundColor Cyan
    try {
        $response = Invoke-RestMethod -Uri "https://newsapi.org/v2/top-headlines?category=technology&apiKey=$newsApiKey&pageSize=5" -Method Get
        Write-Host "✅ Direct API call successful!" -ForegroundColor Green
        Write-Host "Found $($response.articles.Count) articles" -ForegroundColor White
        Write-Host "First article: $($response.articles[0].title)" -ForegroundColor Gray
    } catch {
        Write-Host "❌ Direct API call failed: $($_.Exception.Message)" -ForegroundColor Red
    }
    
    # Test our service
    Write-Host "`n🔍 Testing our News Service..." -ForegroundColor Cyan
    try {
        $response = Invoke-RestMethod -Uri "http://localhost:8001/api/news/trending" -Method Get
        Write-Host "✅ Our service call successful!" -ForegroundColor Green
        Write-Host "Response: $($response | ConvertTo-Json -Depth 2)" -ForegroundColor Gray
    } catch {
        Write-Host "❌ Our service call failed: $($_.Exception.Message)" -ForegroundColor Red
        Write-Host "Make sure the news service is running on port 8001" -ForegroundColor Yellow
    }
    
    # Test through gateway
    Write-Host "`n🔍 Testing through API Gateway..." -ForegroundColor Cyan
    try {
        $response = Invoke-RestMethod -Uri "http://localhost:8080/api/news/trending" -Method Get
        Write-Host "✅ Gateway call successful!" -ForegroundColor Green
        Write-Host "Response: $($response | ConvertTo-Json -Depth 2)" -ForegroundColor Gray
    } catch {
        Write-Host "❌ Gateway call failed: $($_.Exception.Message)" -ForegroundColor Red
        Write-Host "Make sure the API Gateway is running on port 8080" -ForegroundColor Yellow
    }
    
} else {
    Write-Host "❌ No valid NewsAPI key found!" -ForegroundColor Red
    Write-Host "Please add NEWS_API_KEY=46575fbd9144430bb7dce528004ec99e to env.local" -ForegroundColor Yellow
}

Write-Host "`n🎯 Test completed!" -ForegroundColor Green
