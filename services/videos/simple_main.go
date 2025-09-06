package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port := "8003"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "service": "videos"})
	})

	http.HandleFunc("/api/videos", getVideos)
	http.HandleFunc("/api/videos/trending", getTrendingVideos)
	http.HandleFunc("/api/videos/search", searchVideos)

	log.Printf("Videos service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getVideos(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		category = "technology"
	}

	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" || apiKey == "your_youtube_api_key_here" {
		// Return error message for static data
		errorMsg := map[string]interface{}{
			"error": "STATIC DATA - No YouTube API Key Found",
			"message": "This is static/mock data. To get real YouTube videos, add YOUTUBE_API_KEY to your environment variables.",
			"instructions": "Get free API key from https://console.developers.google.com/ and add to env.local",
			"category": category,
			"count": 1,
			"videos": []map[string]interface{}{
				{
					"id":          fmt.Sprintf("static_video_%d", time.Now().UnixNano()),
					"title":       fmt.Sprintf("⚠️ STATIC DATA: Latest %s Tutorial", category),
					"description": fmt.Sprintf("This is mock data. Get real YouTube videos by adding YOUTUBE_API_KEY to env.local file.", category),
					"url":         "https://console.developers.google.com/",
					"channel":     "Mock Data - Add API Key",
					"category":    category,
					"published_at": time.Now(),
					"thumbnail":   "https://via.placeholder.com/320x180/ff6b6b/ffffff?text=STATIC+DATA",
					"duration":    "10:30",
					"views":       15000,
					"is_static":   true,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorMsg)
		return
	}

	// Map categories to YouTube search terms
	searchTerms := map[string]string{
		"technology": "tech news OR programming OR software development",
		"ai":         "artificial intelligence OR machine learning OR AI",
		"business":   "business news OR entrepreneurship OR finance",
		"entertainment": "comedy OR music OR movies",
		"education":  "tutorial OR learning OR course",
		"gaming":     "gaming OR video games OR esports",
		"lifestyle":  "lifestyle OR travel OR food",
		"science":    "science OR research OR discovery",
	}

	searchTerm := searchTerms[category]
	if searchTerm == "" {
		searchTerm = category
	}

	// Search for videos
	searchURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&q=%s&type=video&maxResults=20&order=relevance&key=%s", searchTerm, apiKey)
	
	resp, err := http.Get(searchURL)
	if err != nil {
		http.Error(w, "Failed to search videos", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "YouTube API error", http.StatusInternalServerError)
		return
	}

	var searchResp struct {
		Items []struct {
			ID struct {
				VideoID string `json:"videoId"`
			} `json:"id"`
			Snippet struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				ChannelTitle string `json:"channelTitle"`
				PublishedAt string `json:"publishedAt"`
				Thumbnails  struct {
					High struct {
						URL string `json:"url"`
					} `json:"high"`
				} `json:"thumbnails"`
			} `json:"snippet"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		http.Error(w, "Failed to decode search response", http.StatusInternalServerError)
		return
	}

	// Get video IDs for detailed information
	videoIDs := make([]string, 0, len(searchResp.Items))
	for _, item := range searchResp.Items {
		videoIDs = append(videoIDs, item.ID.VideoID)
	}

	if len(videoIDs) == 0 {
		http.Error(w, "No videos found", http.StatusNotFound)
		return
	}

	// Get detailed information for all videos
	detailsURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?id=%s&part=snippet,statistics,contentDetails&key=%s", 
		fmt.Sprintf("%s", videoIDs[0]), apiKey)
	
	detailsResp, err := http.Get(detailsURL)
	if err != nil {
		http.Error(w, "Failed to fetch video details", http.StatusInternalServerError)
		return
	}
	defer detailsResp.Body.Close()

	if detailsResp.StatusCode != http.StatusOK {
		http.Error(w, "YouTube API error", http.StatusInternalServerError)
		return
	}

	var detailsRespData struct {
		Items []struct {
			ID      string `json:"id"`
			Snippet struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				ChannelTitle string `json:"channelTitle"`
				PublishedAt string `json:"publishedAt"`
				Thumbnails  struct {
					High struct {
						URL string `json:"url"`
					} `json:"high"`
				} `json:"thumbnails"`
			} `json:"snippet"`
			Statistics struct {
				ViewCount string `json:"viewCount"`
			} `json:"statistics"`
			ContentDetails struct {
				Duration string `json:"duration"`
			} `json:"contentDetails"`
		} `json:"items"`
	}

	if err := json.NewDecoder(detailsResp.Body).Decode(&detailsRespData); err != nil {
		http.Error(w, "Failed to decode details response", http.StatusInternalServerError)
		return
	}

	// Transform to our format
	videos := make([]map[string]interface{}, 0, len(detailsRespData.Items))
	for _, video := range detailsRespData.Items {
		viewCount, _ := strconv.ParseInt(video.Statistics.ViewCount, 10, 64)
		
		videos = append(videos, map[string]interface{}{
			"id":          video.ID,
			"title":       video.Snippet.Title,
			"description": video.Snippet.Description,
			"url":         fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.ID),
			"channel":     video.Snippet.ChannelTitle,
			"category":    category,
			"published_at": video.Snippet.PublishedAt,
			"thumbnail":   video.Snippet.Thumbnails.High.URL,
			"duration":    "10:30", // Simplified for demo
			"views":       viewCount,
		})
	}

	result := map[string]interface{}{
		"category": category,
		"count":    len(videos),
		"videos":   videos,
		"source":   "REAL YOUTUBE API DATA",
		"api_key_status": "VALID",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getTrendingVideos(w http.ResponseWriter, r *http.Request) {
	videos := []map[string]interface{}{
		{
			"id":          fmt.Sprintf("trending_video_%d", time.Now().UnixNano()),
			"title":       "AI Revolution in 2024",
			"description": "Exploring the latest AI developments and their impact on society.",
			"url":         "https://youtube.com/watch?v=trending1",
			"channel":     "AI News",
			"category":    "technology",
			"published_at": time.Now(),
			"thumbnail":   "https://via.placeholder.com/320x180",
			"duration":    "15:45",
			"views":       50000,
		},
	}

	result := map[string]interface{}{
		"count":  len(videos),
		"videos": videos,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func searchVideos(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	videos := []map[string]interface{}{
		{
			"id":          fmt.Sprintf("search_video_%d", time.Now().UnixNano()),
			"title":       fmt.Sprintf("Search Results for: %s", query),
			"description": fmt.Sprintf("Videos related to %s", query),
			"url":         "https://youtube.com/watch?v=search1",
			"channel":     "Search Results",
			"category":    "search",
			"published_at": time.Now(),
			"thumbnail":   "https://via.placeholder.com/320x180",
			"duration":    "8:20",
			"views":       5000,
		},
	}

	result := map[string]interface{}{
		"query":  query,
		"count":  len(videos),
		"videos": videos,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
