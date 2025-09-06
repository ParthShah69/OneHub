# 🚀 Complete Personalized Dashboard with Real API Integration

A comprehensive one-stop dashboard that aggregates content from multiple sources (News, YouTube, Jobs, Deals) with personalized recommendations and Web3 gamification.

## ✨ Features

- **📰 Real News Integration**: NewsAPI with fallback to static data
- **🎥 Real YouTube Videos**: YouTube API with fallback to static data  
- **💼 Real Job Postings**: LinkedIn API with fallback to static data
- **🛒 Real Deals**: Amazon API with fallback to static data
- **🤖 AI Recommendations**: Wolfram integration for personalized content
- **🎁 NFT Rewards**: Verbwire integration for Web3 gamification
- **☁️ Decentralized Deployment**: Akash Network ready
- **🔧 Clear API Status**: Shows when using static vs real data

## 🚀 Quick Start

### Option 1: Run with Static Data (Demo Mode)
```powershell
cd F:\hackathon_2
.\start_complete_dashboard.ps1
```

### Option 2: Run with Real APIs
```powershell
# 1. Copy example environment file
copy env.local.example env.local

# 2. Edit with your real API keys
notepad env.local

# 3. Start dashboard
.\start_complete_dashboard.ps1
```

## 🔑 Get API Keys (All Free)

### 📰 NewsAPI
- **URL**: https://newsapi.org/register
- **Free Tier**: 1,000 requests/day
- **Add to env.local**: `NEWS_API_KEY=your_key_here`

### 🎥 YouTube API
- **URL**: https://console.developers.google.com/
- **Free Tier**: 10,000 requests/day
- **Add to env.local**: `YOUTUBE_API_KEY=your_key_here`

### 💼 LinkedIn API
- **URL**: https://developer.linkedin.com/
- **Free Tier**: Available
- **Add to env.local**: `LINKEDIN_API_KEY=your_key_here`

### 🛒 Amazon API
- **URL**: https://webservices.amazon.com/
- **Free Tier**: Available
- **Add to env.local**: `AMAZON_API_KEY=your_key_here`

## 📁 Project Structure

```
hackathon_2/
├── frontend/                 # React Dashboard
│   ├── src/
│   │   ├── components/      # Dashboard components
│   │   │   ├── NewsSection.js
│   │   │   ├── VideosSection.js
│   │   │   ├── JobsSection.js
│   │   │   └── DealsSection.js
│   │   └── App.js
│   └── package.json
├── services/                # GoFr Microservices
│   ├── news/               # NewsAPI integration
│   ├── videos/             # YouTube API integration
│   ├── jobs/               # LinkedIn API integration
│   ├── deals/              # Amazon API integration
│   ├── recommendation/     # Wolfram integration
│   ├── user/               # User management
│   └── nft/                # Verbwire integration
├── gateway/                # API Gateway
├── shared/                 # Shared models and database
├── env.local.example       # Environment variables template
├── start_complete_dashboard.ps1  # Complete startup script
└── README_COMPLETE.md      # This file
```

## 🔧 API Status Indicators

The dashboard clearly shows when you're using static vs real data:

### ⚠️ Static Data Warning
- **Red warning banners** appear when no API keys are found
- **Clear instructions** on how to get API keys
- **Mock data** clearly labeled as "STATIC DATA"
- **Links redirect** to API registration pages

### ✅ Real API Data
- **Green indicators** show "Real API Data"
- **Actual content** from external APIs
- **Working links** to real articles/videos/jobs/deals
- **Live data** updated in real-time

## 🧪 Testing

### Test API Endpoints
```powershell
# Test news (with/without API key)
curl http://localhost:8080/api/news/trending

# Test videos (with/without API key)
curl http://localhost:8080/api/videos?category=technology

# Test jobs (with/without API key)
curl http://localhost:8080/api/jobs?category=ai

# Test deals (with/without API key)
curl http://localhost:8080/api/deals?category=electronics
```

### Expected Responses

**Without API Keys (Static Data):**
```json
{
  "error": "STATIC DATA - No API Key Found",
  "message": "This is static/mock data. Add API_KEY to env.local for real data.",
  "articles": [
    {
      "title": "⚠️ STATIC DATA: Latest News Update",
      "is_static": true,
      "url": "https://newsapi.org/register"
    }
  ]
}
```

**With API Keys (Real Data):**
```json
{
  "source": "REAL API DATA",
  "api_key_status": "VALID",
  "articles": [
    {
      "title": "Real News Article Title",
      "url": "https://real-news-site.com/article",
      "source": "BBC News"
    }
  ]
}
```

## 🌐 Frontend Features

### 📰 News Section
- Real news from NewsAPI
- Category filtering (technology, business, sports, etc.)
- Trending news aggregation
- Search functionality
- Clear static data warnings

### 🎥 Videos Section
- Real YouTube videos
- Category-based search
- Video thumbnails and metadata
- Channel information
- View counts and duration

### 💼 Jobs Section
- Real job postings from LinkedIn
- Location and salary information
- Job type filtering
- Company details
- Application links

### 🛒 Deals Section
- Real deals from Amazon
- Price comparisons
- Discount percentages
- Store information
- Expiration dates

## 🔧 Backend Services

### News Service (Port 8001)
- **Real API**: NewsAPI integration
- **Endpoints**: `/api/news`, `/api/news/trending`, `/api/news/search`
- **Fallback**: Static data with clear warnings

### Videos Service (Port 8003)
- **Real API**: YouTube API integration
- **Endpoints**: `/api/videos`, `/api/videos/trending`
- **Fallback**: Static data with clear warnings

### Jobs Service (Port 8002)
- **Real API**: LinkedIn API integration
- **Endpoints**: `/api/jobs`, `/api/jobs/trending`, `/api/jobs/search`
- **Fallback**: Static data with clear warnings

### Deals Service (Port 8004)
- **Real API**: Amazon API integration
- **Endpoints**: `/api/deals`, `/api/deals/trending`, `/api/deals/search`
- **Fallback**: Static data with clear warnings

### API Gateway (Port 8080)
- **Routes**: All API requests
- **Health**: `/health`
- **CORS**: Enabled for frontend

## 🚀 Deployment

### Local Development
```powershell
.\start_complete_dashboard.ps1
```

### Docker Deployment
```powershell
docker-compose up --build -d
```

### Akash Network Deployment
```powershell
# Deploy to decentralized cloud
akash tx deployment create akash-deploy.yml --from wallet
```

## 🔧 Troubleshooting

### Common Issues

1. **Services not starting**
   - Check if ports are available
   - Ensure Go is installed
   - Check environment variables

2. **API errors**
   - Verify API keys in env.local
   - Check API quotas
   - Review API documentation

3. **Frontend not loading**
   - Ensure all backend services are running
   - Check API Gateway health
   - Verify CORS settings

### Debug Commands
```powershell
# Check service health
curl http://localhost:8080/health

# Check individual services
curl http://localhost:8001/health  # News
curl http://localhost:8002/health  # Jobs
curl http://localhost:8003/health  # Videos
curl http://localhost:8004/health  # Deals

# Check environment variables
Get-ChildItem Env: | Where-Object { $_.Name -match "API_KEY" }
```

## 📊 Performance

- **Response Time**: < 200ms for cached data
- **API Limits**: Respects rate limits
- **Fallback**: Graceful degradation to static data
- **Caching**: In-memory caching for performance

## 🔒 Security

- **API Keys**: Stored in environment variables
- **CORS**: Properly configured
- **Input Validation**: All inputs sanitized
- **Error Handling**: No sensitive data exposed

## 🎯 Hackathon Demo

### Demo Flow
1. **Start without API keys** - Show static data warnings
2. **Add API keys** - Show real data integration
3. **Test all sections** - News, Videos, Jobs, Deals
4. **Show personalization** - User preferences
5. **Demonstrate NFT rewards** - Web3 integration
6. **Deploy to Akash** - Decentralized deployment

### Key Features to Highlight
- ✅ **Real API Integration** with clear status indicators
- ✅ **Static Data Warnings** for transparency
- ✅ **Working Redirects** for all content types
- ✅ **Personalized Recommendations** using Wolfram
- ✅ **Web3 Gamification** with Verbwire NFTs
- ✅ **Decentralized Deployment** on Akash Network

## 📞 Support

For issues or questions:
1. Check the troubleshooting section
2. Review API documentation
3. Verify environment setup
4. Test individual services

---

**🎉 Your complete personalized dashboard is ready!**

Start with static data to see the warnings, then add real API keys for live data integration.
