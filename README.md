# 🚀 Personalized Dashboard - Hackathon Project

A comprehensive one-stop dashboard that aggregates content from multiple sources (News, YouTube, Jobs, Deals, Shopping, Food, OTT) and personalizes recommendations based on explicit user interests + browsing behavior.

## 🏗️ Architecture

### Backend (GoFr Microservices)
- **News Service** → NewsAPI integration for latest headlines
- **Jobs Service** → LinkedIn API integration for job listings  
- **Videos Service** → YouTube API integration for trending videos
- **Deals Service** → Real shopping APIs (Amazon, Flipkart)
- **Recommendation Service** → Wolfram API for AI-powered recommendations
- **User Service** → Profile management and behavior tracking
- **NFT Service** → Verbwire API for Web3 gamification
- **API Gateway** → Unified API endpoint

### Frontend (React)
- Modern, responsive dashboard with glassmorphism design
- Personalized content recommendations
- Real-time data updates
- User behavior tracking

### Deployment (Akash)
- Containerized microservices
- Decentralized cloud deployment
- Scalable infrastructure

## 🛠️ Tech Stack

- **Backend**: GoFr, Go
- **Frontend**: React, Styled Components
- **Database**: PostgreSQL, Redis
- **APIs**: NewsAPI, YouTube API, LinkedIn API, Wolfram API, Verbwire API
- **Deployment**: Docker, Akash Network
- **Caching**: Redis, Go Cache

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Node.js 16+
- Go 1.21+

### 1. Clone and Setup
```bash
git clone <repository>
cd personalized-dashboard
```

### 2. Environment Setup
```bash
cp env.example .env
# Add your API keys to .env file
```

### 3. Run with Docker Compose
```bash
docker-compose up --build
```

### 4. Access the Application
- **Frontend**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **Individual Services**: 
  - News: http://localhost:8001
  - Jobs: http://localhost:8002
  - Videos: http://localhost:8003
  - Deals: http://localhost:8004
  - Recommendations: http://localhost:8005
  - User: http://localhost:8006
  - NFT: http://localhost:8007

## 🔑 Required API Keys

Add these to your `.env` file:

```env
NEWS_API_KEY=your_newsapi_key
YOUTUBE_API_KEY=your_youtube_api_key
LINKEDIN_API_KEY=your_linkedin_api_key
WOLFRAM_API_KEY=your_wolfram_api_key
VERBWIRE_API_KEY=your_verbwire_api_key
AMAZON_API_KEY=your_amazon_api_key
FLIPKART_API_KEY=your_flipkart_api_key
```

## 🎯 Features

### Personalization Engine
- **Explicit Interests**: User selects categories during onboarding
- **Behavioral Tracking**: Tracks clicks, bookmarks, shares, searches
- **AI Recommendations**: Wolfram API powers intelligent content ranking
- **Dynamic Scoring**: Real-time preference adjustment

### Content Aggregation
- **News**: Latest headlines from multiple categories
- **Jobs**: Trending job opportunities from LinkedIn
- **Videos**: YouTube trending content
- **Deals**: Shopping deals from Amazon, Flipkart
- **NFT Rewards**: Web3 gamification with Verbwire

### User Experience
- **Glassmorphism Design**: Modern, beautiful UI
- **Real-time Updates**: Live content refresh
- **Responsive**: Works on all devices
- **Fast**: Optimized with caching

## 🏆 Hackathon Highlights

### GoFr Integration
- Microservices architecture with GoFr framework
- Independent service deployment
- RESTful API design

### Wolfram Integration
- AI-powered recommendation engine
- Content clustering and ranking
- Personalized scoring algorithm

### Verbwire Integration
- NFT coupon generation
- Web3 gamification
- Blockchain rewards system

### Akash Integration
- Decentralized deployment
- Container orchestration
- Scalable infrastructure

## 📊 API Endpoints

### News Service
- `GET /api/news?category=technology` - Get news by category
- `GET /api/news/trending` - Get trending news
- `GET /api/news/search?q=query` - Search news

### Jobs Service
- `GET /api/jobs?category=ai` - Get jobs by category
- `GET /api/jobs/trending` - Get trending jobs
- `GET /api/jobs/search?q=query` - Search jobs

### Videos Service
- `GET /api/videos?category=technology` - Get videos by category
- `GET /api/videos/trending` - Get trending videos
- `GET /api/videos/search?q=query` - Search videos

### Deals Service
- `GET /api/deals?category=electronics` - Get deals by category
- `GET /api/deals/trending` - Get trending deals
- `GET /api/deals/search?q=query` - Search deals

### Recommendation Service
- `GET /api/recommendations?user_id=123` - Get personalized recommendations

### User Service
- `POST /api/users` - Create user
- `GET /api/users/:id` - Get user profile
- `POST /api/users/:id/behavior` - Track user behavior

### NFT Service
- `POST /api/nft/mint` - Mint NFT coupon
- `GET /api/nft/:user_id` - Get user NFTs
- `POST /api/nft/:id/claim` - Claim NFT

## 🚀 Deployment

### Local Development
```bash
docker-compose up --build
```

### Akash Deployment
```bash
# Build and push images
docker build -t your-registry/news-service ./services/news
docker push your-registry/news-service

# Deploy to Akash
akash tx deployment create akash-deploy.yml --from your-wallet
```

## 🎨 Demo Flow

1. **User Onboarding**: Select interests and create profile
2. **Dashboard Load**: Fetch content from all services
3. **AI Recommendations**: Wolfram generates personalized content
4. **User Interaction**: Track behavior for better recommendations
5. **NFT Rewards**: Earn coupons for engagement
6. **Real-time Updates**: Live content refresh

## 🔮 Future Enhancements

- Browser extension for external behavior tracking
- Machine learning model training
- Advanced personalization algorithms
- Social features and sharing
- Mobile app development
- Advanced analytics dashboard

## 📝 License

MIT License - feel free to use this project for your hackathon!

---

**Built with ❤️ for the hackathon using GoFr, Wolfram, Verbwire, and Akash!**
