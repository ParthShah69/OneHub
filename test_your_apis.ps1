# Test Your API Keys
Write-Host "🧪 Testing Your API Keys..." -ForegroundColor Green

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

Write-Host "`n📋 Your Current API Keys:" -ForegroundColor Yellow

# Test News API (should work)
$newsKey = [Environment]::GetEnvironmentVariable('NEWS_API_KEY', 'Process')
Write-Host "📰 News API: $newsKey" -ForegroundColor $(if ($newsKey -and $newsKey -ne 'your_newsapi_key_here') { 'Green' } else { 'Red' })

# Test OMDb API
$omdbKey = [Environment]::GetEnvironmentVariable('OMDB_API_KEY', 'Process')
Write-Host "🎬 OMDb API: $omdbKey" -ForegroundColor $(if ($omdbKey -and $omdbKey -ne 'your_omdb_api_key_here') { 'Green' } else { 'Red' })

# Test Edamam API
$edamamID = [Environment]::GetEnvironmentVariable('EDAMAM_APP_ID', 'Process')
$edamamKey = [Environment]::GetEnvironmentVariable('EDAMAM_APP_KEY', 'Process')
Write-Host "🍽️ Edamam App ID: $edamamID" -ForegroundColor $(if ($edamamID -and $edamamID -ne 'your_edamam_app_id_here') { 'Green' } else { 'Red' })
Write-Host "🍽️ Edamam App Key: $edamamKey" -ForegroundColor $(if ($edamamKey -and $edamamKey -ne 'your_edamam_app_key_here') { 'Green' } else { 'Red' })

Write-Host "`n🔍 Testing APIs..." -ForegroundColor Cyan

# Test OMDb API
if ($omdbKey -and $omdbKey -ne 'your_omdb_api_key_here') {
    Write-Host "Testing OMDb API..." -ForegroundColor Gray
    try {
        $response = Invoke-RestMethod -Uri "http://www.omdbapi.com/?s=action&type=movie&apikey=$omdbKey" -Method Get -TimeoutSec 10
        if ($response.Response -eq "True") {
            Write-Host "✅ OMDb API: Working! Found $($response.totalResults) movies" -ForegroundColor Green
        } else {
            Write-Host "❌ OMDb API: Error - $($response.Error)" -ForegroundColor Red
        }
    } catch {
        Write-Host "❌ OMDb API: Connection Error - $($_.Exception.Message)" -ForegroundColor Red
    }
} else {
    Write-Host "❌ OMDb API: No key provided" -ForegroundColor Red
}

# Test Edamam API
if ($edamamID -and $edamamID -ne 'your_edamam_app_id_here' -and $edamamKey -and $edamamKey -ne 'your_edamam_app_key_here') {
    Write-Host "Testing Edamam API..." -ForegroundColor Gray
    try {
        $response = Invoke-RestMethod -Uri "https://api.edamam.com/search?q=chicken&app_id=$edamamID&app_key=$edamamKey&from=0&to=1" -Method Get -TimeoutSec 10
        if ($response.hits) {
            Write-Host "✅ Edamam API: Working! Found $($response.count) recipes" -ForegroundColor Green
        } else {
            Write-Host "❌ Edamam API: No results or invalid keys" -ForegroundColor Red
        }
    } catch {
        Write-Host "❌ Edamam API: Connection Error - $($_.Exception.Message)" -ForegroundColor Red
    }
} else {
    Write-Host "❌ Edamam API: Keys not provided" -ForegroundColor Red
}

Write-Host "`n🚀 Starting services..." -ForegroundColor Cyan
Write-Host "Run: .\quick_fix.ps1" -ForegroundColor White
Write-Host "Then open: http://localhost:3000" -ForegroundColor White
