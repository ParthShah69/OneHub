package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gofr.dev/pkg/gofr"
	"github.com/patrickmn/go-cache"
)

type YouTubeResponse struct {
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
				Default struct {
					URL string `json:"url"`
				} `json:"default"`
				High struct {
					URL string `json:"url"`
				} `json:"high"`
			} `json:"thumbnails"`
		} `json:"snippet"`
		Statistics struct {
			ViewCount    string `json:"viewCount"`
			LikeCount    string `json:"likeCount"`
			CommentCount string `json:"commentCount"`
		} `json:"statistics"`
		ContentDetails struct {
			Duration string `json:"duration"`
		} `json:"contentDetails"`
	} `json:"items"`
}

type VideosService struct {
	apiKey string
	cache  *cache.Cache
}

func main() {
	app := gofr.New()

	videosService := &VideosService{
		apiKey: os.Getenv("YOUTUBE_API_KEY"),
		cache:  cache.New(5*time.Minute, 10*time.Minute),
	}

	// Health check
	app.GET("/health", func(ctx *gofr.Context) (interface{}, error) {
		return map[string]string{"status": "healthy", "service": "videos"}, nil
	})

	// Get videos by category
	app.GET("/api/videos", videosService.GetVideos)
	
	// Get trending videos
	app.GET("/api/videos/trending", videosService.GetTrendingVideos)
	
	// Search videos
	app.GET("/api/videos/search", videosService.SearchVideos)
	
	// Get video details
	app.GET("/api/videos/:id", videosService.GetVideoDetails)

	app.Run()
}

func (vs *VideosService) GetVideos(ctx *gofr.Context) (interface{}, error) {
	category := ctx.Param("category")
	if category == "" {
		category = "technology"
	}

	// Check cache first
	if cached, found := vs.cache.Get("videos_" + category); found {
		return cached, nil
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

	videos, err := vs.fetchVideosFromYouTube(searchTerm, 20)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch videos: %v", err)
	}

	result := map[string]interface{}{
		"category": category,
		"count":    len(videos),
		"videos":   videos,
	}

	// Cache the result
	vs.cache.Set("videos_"+category, result, cache.DefaultExpiration)

	return result, nil
}

func (vs *VideosService) GetTrendingVideos(ctx *gofr.Context) (interface{}, error) {
	// Check cache first
	if cached, found := vs.cache.Get("trending_videos"); found {
		return cached, nil
	}

	// Get trending videos from multiple categories
	categories := []string{"technology", "ai", "business", "entertainment"}
	allVideos := make([]map[string]interface{}, 0)

	for _, category := range categories {
		searchTerms := map[string]string{
			"technology": "tech news OR programming",
			"ai":         "artificial intelligence OR machine learning",
			"business":   "business news OR entrepreneurship",
			"entertainment": "comedy OR music",
		}
		
		searchTerm := searchTerms[category]
		videos, err := vs.fetchVideosFromYouTube(searchTerm, 5)
		if err != nil {
			log.Printf("Failed to fetch %s videos: %v", category, err)
			continue
		}
		allVideos = append(allVideos, videos...)
	}

	result := map[string]interface{}{
		"count":  len(allVideos),
		"videos": allVideos,
	}

	// Cache the result
	vs.cache.Set("trending_videos", result, cache.DefaultExpiration)

	return result, nil
}

func (vs *VideosService) SearchVideos(ctx *gofr.Context) (interface{}, error) {
	query := ctx.Param("q")
	if query == "" {
		return nil, fmt.Errorf("query parameter 'q' is required")
	}

	maxResults := ctx.Param("maxResults")
	if maxResults == "" {
		maxResults = "20"
	}

	// Check cache first
	cacheKey := "search_videos_" + query + "_" + maxResults
	if cached, found := vs.cache.Get(cacheKey); found {
		return cached, nil
	}

	videos, err := vs.fetchVideosFromYouTube(query, 20)
	if err != nil {
		return nil, fmt.Errorf("failed to search videos: %v", err)
	}

	result := map[string]interface{}{
		"query":  query,
		"count":  len(videos),
		"videos": videos,
	}

	// Cache the result
	vs.cache.Set(cacheKey, result, cache.DefaultExpiration)

	return result, nil
}

func (vs *VideosService) GetVideoDetails(ctx *gofr.Context) (interface{}, error) {
	videoID := ctx.Param("id")
	if videoID == "" {
		return nil, fmt.Errorf("video ID is required")
	}

	// Check cache first
	if cached, found := vs.cache.Get("video_details_" + videoID); found {
		return cached, nil
	}

	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?id=%s&part=snippet,statistics,contentDetails&key=%s", 
		videoID, vs.apiKey)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch video details: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("YouTube API returned status: %d", resp.StatusCode)
	}

	var youtubeResp YouTubeResponse
	if err := json.NewDecoder(resp.Body).Decode(&youtubeResp); err != nil {
		return nil, fmt.Errorf("failed to decode video response: %v", err)
	}

	if len(youtubeResp.Items) == 0 {
		return nil, fmt.Errorf("video not found")
	}

	video := youtubeResp.Items[0]
	
	// Parse duration (ISO 8601 format)
	duration := vs.parseDuration(video.ContentDetails.Duration)
	
	// Parse view count
	viewCount, _ := strconv.ParseInt(video.Statistics.ViewCount, 10, 64)
	
	result := map[string]interface{}{
		"id":          video.ID.VideoID,
		"title":       video.Snippet.Title,
		"description": video.Snippet.Description,
		"url":         fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.ID.VideoID),
		"channel":     video.Snippet.ChannelTitle,
		"category":    "video",
		"published_at": video.Snippet.PublishedAt,
		"thumbnail":   video.Snippet.Thumbnails.High.URL,
		"duration":    duration,
		"views":       viewCount,
		"likes":       video.Statistics.LikeCount,
		"comments":    video.Statistics.CommentCount,
	}

	// Cache the result
	vs.cache.Set("video_details_"+videoID, result, cache.DefaultExpiration)

	return result, nil
}

func (vs *VideosService) fetchVideosFromYouTube(query string, maxResults int) ([]map[string]interface{}, error) {
	// First, search for videos
	searchURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&q=%s&type=video&maxResults=%d&order=relevance&key=%s", 
		query, maxResults, vs.apiKey)
	
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search videos: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("YouTube API returned status: %d", resp.StatusCode)
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
					Default struct {
						URL string `json:"url"`
					} `json:"default"`
					High struct {
						URL string `json:"url"`
					} `json:"high"`
				} `json:"thumbnails"`
			} `json:"snippet"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %v", err)
	}

	// Get video IDs for detailed information
	videoIDs := make([]string, 0, len(searchResp.Items))
	for _, item := range searchResp.Items {
		videoIDs = append(videoIDs, item.ID.VideoID)
	}

	// Get detailed information for all videos
	detailsURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?id=%s&part=snippet,statistics,contentDetails&key=%s", 
		strings.Join(videoIDs, ","), vs.apiKey)
	
	detailsResp, err := http.Get(detailsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch video details: %v", err)
	}
	defer detailsResp.Body.Close()

	if detailsResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("YouTube API returned status: %d", detailsResp.StatusCode)
	}

	var detailsRespData YouTubeResponse
	if err := json.NewDecoder(detailsResp.Body).Decode(&detailsRespData); err != nil {
		return nil, fmt.Errorf("failed to decode details response: %v", err)
	}

	// Transform to our format
	videos := make([]map[string]interface{}, 0, len(detailsRespData.Items))
	for _, video := range detailsRespData.Items {
		duration := vs.parseDuration(video.ContentDetails.Duration)
		viewCount, _ := strconv.ParseInt(video.Statistics.ViewCount, 10, 64)
		
		videos = append(videos, map[string]interface{}{
			"id":          video.ID.VideoID,
			"title":       video.Snippet.Title,
			"description": video.Snippet.Description,
			"url":         fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.ID.VideoID),
			"channel":     video.Snippet.ChannelTitle,
			"category":    "video",
			"published_at": video.Snippet.PublishedAt,
			"thumbnail":   video.Snippet.Thumbnails.High.URL,
			"duration":    duration,
			"views":       viewCount,
		})
	}

	return videos, nil
}

func (vs *VideosService) parseDuration(duration string) string {
	// Parse ISO 8601 duration format (PT4M13S -> 4:13)
	if !strings.HasPrefix(duration, "PT") {
		return duration
	}
	
	duration = strings.TrimPrefix(duration, "PT")
	
	var hours, minutes, seconds int
	
	if strings.Contains(duration, "H") {
		parts := strings.Split(duration, "H")
		fmt.Sscanf(parts[0], "%d", &hours)
		duration = parts[1]
	}
	
	if strings.Contains(duration, "M") {
		parts := strings.Split(duration, "M")
		fmt.Sscanf(parts[0], "%d", &minutes)
		duration = parts[1]
	}
	
	if strings.Contains(duration, "S") {
		parts := strings.Split(duration, "S")
		fmt.Sscanf(parts[0], "%d", &seconds)
	}
	
	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
