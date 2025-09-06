# 🔑 Complete FREE API Keys Guide - 100% Free APIs Only!

This guide provides step-by-step instructions to get **100% FREE** API keys for all services with no payment required!

## 📰 NewsAPI (Already Working!)
**Status**: ✅ **WORKING** - Your key: `46575fbd9144430bb7dce528004ec99e`

- **URL**: https://newsapi.org/register
- **Free Tier**: 1,000 requests/day
- **Setup**: Already configured in `env.local`
- **Test**: `curl http://localhost:8080/api/news/trending`

## 🎬 OMDb API (Open Movie Database) - Movies
**Status**: ⚠️ **STATIC DATA** - Need FREE API key

### Get OMDb API Key (100% FREE):
1. **Go to**: http://www.omdbapi.com/apikey.aspx
2. **Enter your email** (any email works)
3. **Get API key instantly** (no approval needed)
4. **Copy your API key**
5. **Add to env.local**: `OMDB_API_KEY=your_omdb_key_here`

### Features:
- Movie details, ratings, posters, plot summaries
- Search by title, year, genre
- IMDb ratings and Rotten Tomatoes scores
- **FREE**: 1,000 requests/day (no payment ever required!)

## 🍽️ Edamam Recipe API - Food & Recipes
**Status**: ⚠️ **STATIC DATA** - Need FREE API key

### Get Edamam API Key (100% FREE):
1. **Go to**: https://developer.edamam.com/
2. **Sign up** for free account
3. **Create new application** (Recipe Search API)
4. **Get App ID and App Key** instantly
5. **Add to env.local**: 
   ```
   EDAMAM_APP_ID=your_app_id_here
   EDAMAM_APP_KEY=your_app_key_here
   ```

### Features:
- Recipe search, nutrition info, dietary filters
- Recipe images, cooking times, health labels
- Direct links to full recipes
- **FREE**: 5 requests/minute, 10,000 requests/month (no payment required!)

## 🎥 YouTube Data API - Videos
**Status**: ⚠️ **STATIC DATA** - Need FREE API key

### Get YouTube API Key (100% FREE):
1. **Go to**: https://console.developers.google.com/
2. **Create new project** or select existing
3. **Enable YouTube Data API v3**
4. **Create credentials** → API Key
5. **Add to env.local**: `YOUTUBE_API_KEY=your_youtube_key_here`

### Features:
- Trending videos by category
- Video thumbnails, view counts, channel info
- Direct links to YouTube videos
- **FREE**: 10,000 requests/day (no payment required!)

## 💼 Adzuna Jobs API - Jobs
**Status**: ⚠️ **STATIC DATA** - Need FREE API key

### Get Adzuna API Key (100% FREE):
1. **Go to**: https://developer.adzuna.com/
2. **Sign up** for free account
3. **Get App ID and App Key** instantly
4. **Add to env.local**: 
   ```
   ADZUNA_APP_ID=your_app_id_here
   ADZUNA_APP_KEY=your_app_key_here
   ```

### Features:
- Job postings by category, location, salary
- Company information, job descriptions
- Direct links to job applications
- **FREE**: 1,000 requests/day (no payment required!)

## 🛒 RapidAPI Shopping APIs - Deals
**Status**: ⚠️ **STATIC DATA** - Need FREE API key

### Get RapidAPI Key (100% FREE):
1. **Go to**: https://rapidapi.com/
2. **Sign up** for free account
3. **Get API key** from dashboard
4. **Subscribe to free shopping APIs** (like Amazon Price API, eBay API)
5. **Add to env.local**: `RAPIDAPI_KEY=your_rapidapi_key_here`

### Features:
- Product deals, price comparisons, ratings
- Multiple shopping APIs in one platform
- Direct links to products
- **FREE**: 500 requests/month (no payment required!)

## 🤖 Hugging Face Inference API - AI Recommendations
**Status**: ⚠️ **STATIC DATA** - Need FREE API key

### Get Hugging Face API Key (100% FREE):
1. **Go to**: https://huggingface.co/settings/tokens
2. **Sign up** for free account
3. **Create new token** (read access)
4. **Add to env.local**: `HUGGINGFACE_API_KEY=your_hf_key_here`

### Features:
- AI-powered content recommendations
- Text analysis and sentiment
- Content clustering and ranking
- **FREE**: 1,000 requests/month (no payment required!)

## 🎁 OpenSea API - NFT Rewards
**Status**: ⚠️ **STATIC DATA** - Need FREE API key

### Get OpenSea API Key (100% FREE):
1. **Go to**: https://docs.opensea.io/reference/api-overview
2. **Sign up** for free account
3. **Get API key** from dashboard
4. **Add to env.local**: `OPENSEA_API_KEY=your_opensea_key_here`

### Features:
- NFT metadata and collections
- Web3 gamification
- User engagement tracking
- **FREE**: 1,000 requests/day (no payment required!)

## 🔧 How to Add API Keys

### Step 1: Edit env.local file
```bash
# Open env.local file
notepad env.local
```

### Step 2: Replace placeholder values
```env
# Replace these with your actual FREE API keys
NEWS_API_KEY=46575fbd9144430bb7dce528004ec99e  # Already working!
OMDB_API_KEY=your_actual_omdb_key_here
EDAMAM_APP_ID=your_actual_edamam_app_id_here
EDAMAM_APP_KEY=your_actual_edamam_app_key_here
YOUTUBE_API_KEY=your_actual_youtube_key_here
ADZUNA_APP_ID=your_actual_adzuna_app_id_here
ADZUNA_APP_KEY=your_actual_adzuna_app_key_here
RAPIDAPI_KEY=your_actual_rapidapi_key_here
HUGGINGFACE_API_KEY=your_actual_hf_key_here
OPENSEA_API_KEY=your_actual_opensea_key_here
```

### Step 3: Restart services
```powershell
# Stop all services (Ctrl+C in each terminal)
# Then restart
.\start_complete_dashboard.ps1
```

## 🧪 Test Your API Keys

### Test Individual Services:
```powershell
# Test News (should work)
curl http://localhost:8080/api/news/trending

# Test Movies (after adding OMDb key)
curl http://localhost:8080/api/movies/trending

# Test Food (after adding Edamam key)
curl http://localhost:8080/api/food/trending

# Test Videos (after adding YouTube key)
curl http://localhost:8080/api/videos/trending
```

### Test Complete System:
```powershell
.\test_complete_system.ps1
```

## 📊 Expected Results

### With API Keys:
- ✅ **Real Data**: Actual content from external APIs
- ✅ **Green Indicators**: "Real API Data" status
- ✅ **Working Links**: Direct links to external content
- ✅ **No Warnings**: No static data warnings

### Without API Keys:
- ⚠️ **Static Data**: Mock data with clear warnings
- ⚠️ **Red Warnings**: "STATIC DATA WARNING" banners
- ⚠️ **Registration Links**: Links to get API keys

## 🚀 Quick Start Priority

### Essential APIs (Start Here - All 100% FREE):
1. **OMDb** - Movies (Instant API key, great content)
2. **Edamam** - Recipes (Instant API key, useful content)
3. **YouTube** - Videos (Popular, lots of content)
4. **Adzuna** - Jobs (Instant API key, real job data)

### Advanced APIs (Later - All 100% FREE):
5. **RapidAPI** - Shopping (Multiple APIs in one)
6. **Hugging Face** - AI (Advanced features)
7. **OpenSea** - NFTs (Web3 features)

## 💡 Pro Tips

1. **Start with OMDb and Edamam** - instant API keys, most visual impact
2. **YouTube API** - great for demo, lots of content
3. **Test one API at a time** - easier to debug
4. **All APIs are 100% FREE** - no payment ever required!
5. **Use test scripts** - verify each API works

## 🆘 Troubleshooting

### Common Issues:
1. **404 Errors**: Services not running - restart with `.\start_complete_dashboard.ps1`
2. **API Errors**: Check API key format and quotas
3. **CORS Errors**: Make sure API Gateway is running on port 8080
4. **Static Data**: API key not loaded - check env.local file

### Debug Commands:
```powershell
# Check if services are running
curl http://localhost:8080/health

# Check individual services
curl http://localhost:8001/health  # News
curl http://localhost:8008/health  # Movies
curl http://localhost:8009/health  # Food

# Check environment variables
Get-ChildItem Env: | Where-Object { $_.Name -match "API_KEY" }
```

---

**🎉 With all FREE API keys, your dashboard will show real, live data from all sources - NO PAYMENT REQUIRED!**
