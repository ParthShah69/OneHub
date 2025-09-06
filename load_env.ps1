# PowerShell script to load environment variables from env.local
param(
    [string]$EnvFile = "env.local"
)

if (Test-Path $EnvFile) {
    Write-Host "Loading environment variables from $EnvFile..." -ForegroundColor Green
    
    Get-Content $EnvFile | ForEach-Object {
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
    
    Write-Host "Environment variables loaded successfully!" -ForegroundColor Green
} else {
    Write-Host "⚠️  $EnvFile not found. Using default values." -ForegroundColor Yellow
    Write-Host "Create $EnvFile with your API keys to get real data:" -ForegroundColor Yellow
    Write-Host "  NEWS_API_KEY=your_newsapi_key" -ForegroundColor Gray
    Write-Host "  YOUTUBE_API_KEY=your_youtube_key" -ForegroundColor Gray
    Write-Host "  LINKEDIN_API_KEY=your_linkedin_key" -ForegroundColor Gray
}
