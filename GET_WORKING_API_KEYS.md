# 🔑 Get Working API Keys - Step by Step

## 🎬 **OMDb API (Movies) - FREE**

### Step 1: Get API Key
1. **Go to**: http://www.omdbapi.com/apikey.aspx
2. **Enter your email** (any valid email)
3. **Click "Get API Key"**
4. **Copy the API key** (it will be a long string like: `a1b2c3d4-e5f6-7890-abcd-ef1234567890`)

### Step 2: Test Your Key
```powershell
# Replace YOUR_KEY_HERE with your actual key
Invoke-RestMethod -Uri "http://www.omdbapi.com/?s=action&type=movie&apikey=YOUR_KEY_HERE" -Method Get
```

### Step 3: Add to env.local
```env
OMDB_API_KEY=your_actual_omdb_key_here
```

## 🍽️ **Edamam Recipe API (Food) - FREE**

### Step 1: Get API Keys
1. **Go to**: https://developer.edamam.com/
2. **Sign up** for free account
3. **Click "Create New Application"**
4. **Select "Recipe Search API"**
5. **Fill in application details**:
   - Application Name: "Dashboard App"
   - Description: "Personalized dashboard"
6. **Click "Create Application"**
7. **Copy your App ID and App Key**

### Step 2: Test Your Keys
```powershell
# Replace YOUR_APP_ID and YOUR_APP_KEY with your actual keys
Invoke-RestMethod -Uri "https://api.edamam.com/search?q=chicken&app_id=YOUR_APP_ID&app_key=YOUR_APP_KEY&from=0&to=1" -Method Get
```

### Step 3: Add to env.local
```env
EDAMAM_APP_ID=your_actual_app_id_here
EDAMAM_APP_KEY=your_actual_app_key_here
```

## 🎥 **YouTube Data API (Videos) - FREE**

### Step 1: Get API Key
1. **Go to**: https://console.developers.google.com/
2. **Create new project** or select existing
3. **Enable YouTube Data API v3**:
   - Go to "APIs & Services" → "Library"
   - Search for "YouTube Data API v3"
   - Click "Enable"
4. **Create credentials**:
   - Go to "APIs & Services" → "Credentials"
   - Click "Create Credentials" → "API Key"
   - Copy the API key

### Step 2: Test Your Key
```powershell
# Replace YOUR_KEY_HERE with your actual key
Invoke-RestMethod -Uri "https://www.googleapis.com/youtube/v3/search?part=snippet&q=programming&key=YOUR_KEY_HERE&maxResults=1" -Method Get
```

### Step 3: Add to env.local
```env
YOUTUBE_API_KEY=your_actual_youtube_key_here
```

## 🧪 **Test All APIs**

### Run Validation Script
```powershell
.\validate_apis.ps1
```

### Test Individual APIs
```powershell
# Test Movies
curl http://localhost:8080/api/movies/trending

# Test Food
curl http://localhost:8080/api/food/trending

# Test Videos
curl http://localhost:8080/api/videos/trending

# Test News (should work)
curl http://localhost:8080/api/news/trending
```

## 🚀 **Quick Start (Recommended Order)**

### 1. Start with OMDb (Easiest)
- Takes 2 minutes
- Just need email
- Instant API key

### 2. Then Edamam (Good for food)
- Takes 5 minutes
- Need to create account
- Get App ID + App Key

### 3. Then YouTube (Popular)
- Takes 10 minutes
- Need Google account
- Enable API + create key

## 🎯 **Expected Results**

### With Working API Keys:
- ✅ **Real movie data** from OMDb
- ✅ **Real recipe data** from Edamam
- ✅ **Real video data** from YouTube
- ✅ **No static data warnings**

### Without API Keys:
- ⚠️ **Static data** with clear warnings
- ⚠️ **Red warning banners**
- ⚠️ **Links to get API keys**

## 🔧 **After Getting Keys**

### 1. Update env.local
```env
# Replace with your actual keys
OMDB_API_KEY=your_actual_omdb_key_here
EDAMAM_APP_ID=your_actual_app_id_here
EDAMAM_APP_KEY=your_actual_app_key_here
YOUTUBE_API_KEY=your_actual_youtube_key_here
```

### 2. Restart Services
```powershell
.\quick_fix.ps1
```

### 3. Test Frontend
- Open: http://localhost:3000
- Check Movies section for real data
- Check Food section for real data
- Check Videos section for real data

---

**🎉 Once you have working API keys, your dashboard will show real, live data!**
