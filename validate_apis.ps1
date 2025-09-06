# API Key Validation Script
Write-Host "🔍 Validating API Keys..." -ForegroundColor Green

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

Write-Host "`n📋 Current API Keys Status:" -ForegroundColor Yellow

# Test News API (should work)
$newsKey = [Environment]::GetEnvironmentVariable('NEWS_API_KEY', 'Process')
if ($newsKey -and $newsKey -ne 'your_newsapi_key_here') {
    Write-Host "✅ News API: $newsKey" -ForegroundColor Green
} else {
    Write-Host "❌ News API: Not set" -ForegroundColor Red
}

# Test OMDb API
$omdbKey = [Environment]::GetEnvironmentVariable('OMDB_API_KEY', 'Process')
if ($omdbKey -and $omdbKey -ne 'your_omdb_api_key_here') {
    Write-Host "🔍 Testing OMDb API..." -ForegroundColor Cyan
    try {
        $response = Invoke-RestMethod -Uri "http://www.omdbapi.com/?s=test&type=movie&apikey=$omdbKey" -Method Get -TimeoutSec 10
        if ($response.Response -eq "True") {
            Write-Host "✅ OMDb API: Working - Found $($response.totalResults) movies" -ForegroundColor Green
        } else {
            Write-Host "❌ OMDb API: Invalid key - $($response.Error)" -ForegroundColor Red
        }
    } catch {
        Write-Host "❌ OMDb API: Error - $($_.Exception.Message)" -ForegroundColor Red
    }
} else {
    Write-Host "❌ OMDb API: Not set" -ForegroundColor Red
}

# Test Edamam API
$edamamID = [Environment]::GetEnvironmentVariable('EDAMAM_APP_ID', 'Process')
$edamamKey = [Environment]::GetEnvironmentVariable('EDAMAM_APP_KEY', 'Process')
if ($edamamID -and $edamamID -ne 'your_edamam_app_id_here' -and $edamamKey -and $edamamKey -ne 'your_edamam_app_key_here') {
    Write-Host "🔍 Testing Edamam API..." -ForegroundColor Cyan
    try {
        $response = Invoke-RestMethod -Uri "https://api.edamam.com/search?q=chicken&app_id=$edamamID&app_key=$edamamKey&from=0&to=1" -Method Get -TimeoutSec 10
        if ($response.hits) {
            Write-Host "✅ Edamam API: Working - Found $($response.count) recipes" -ForegroundColor Green
        } else {
            Write-Host "❌ Edamam API: No results or invalid keys" -ForegroundColor Red
        }
    } catch {
        Write-Host "❌ Edamam API: Error - $($_.Exception.Message)" -ForegroundColor Red
    }
} else {
    Write-Host "❌ Edamam API: Not set" -ForegroundColor Red
}

Write-Host "`n🔧 API Key Issues Found:" -ForegroundColor Yellow
Write-Host "1. OMDb API key appears to be invalid" -ForegroundColor Red
Write-Host "2. Edamam API keys may be incorrect" -ForegroundColor Red

Write-Host "`n💡 Solutions:" -ForegroundColor Green
Write-Host "1. Get new OMDb API key: http://www.omdbapi.com/apikey.aspx" -ForegroundColor White
Write-Host "2. Get new Edamam keys: https://developer.edamam.com/" -ForegroundColor White
Write-Host "3. Or use the working News API for now" -ForegroundColor White

Write-Host "`n🚀 Starting services with current keys..." -ForegroundColor Cyan
