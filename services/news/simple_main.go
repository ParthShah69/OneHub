package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type NewsAPIResponse struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"source"`
		Author      string    `json:"author"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		URL         string    `json:"url"`
		URLToImage  string    `json:"urlToImage"`
		PublishedAt time.Time `json:"publishedAt"`
		Content     string    `json:"content"`
	} `json:"articles"`
}

func main() {
	port := "8001"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "service": "news"})
	})

	// Get news by category
	http.HandleFunc("/api/news", getNews)
	
	// Get trending news
	http.HandleFunc("/api/news/trending", getTrendingNews)
	
	// Get news by query
	http.HandleFunc("/api/news/search", searchNews)

	log.Printf("News service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getNews(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	userID := r.URL.Query().Get("user_id")
	
	if category == "" {
		category = "general"
	}
	
	// If user_id is provided, get user preferences
	if userID != "" {
		preferences := getUserPreferences(userID)
		if len(preferences.NewsCategories) > 0 {
			category = preferences.NewsCategories[0] // Use first preference
		}
	}

	// Try to load API key from multiple sources
	apiKey := os.Getenv("NEWS_API_KEY")
	log.Printf("NewsAPI Key: %s", apiKey)
	if apiKey == "" || apiKey == "your_newsapi_key_here" {
		// Return error message for static data
		errorMsg := map[string]interface{}{
			"error": "STATIC DATA - No API Key Found",
			"message": "This is static/mock data. To get real news, add NEWS_API_KEY to your environment variables.",
			"instructions": "Get free API key from https://newsapi.org/register and add to env.local",
			"category": category,
			"count": 1,
			"articles": []map[string]interface{}{
				{
					"id":          fmt.Sprintf("static_news_%d", time.Now().UnixNano()),
					"title":       fmt.Sprintf("⚠️ STATIC DATA: Latest %s News Update", category),
					"description": fmt.Sprintf("This is mock data. Get real news by adding NEWS_API_KEY to env.local file.", category),
					"url":         "https://newsapi.org/register",
					"source":      "Mock Data - Add API Key",
					"category":    category,
					"published_at": time.Now(),
					"image_url":   "https://via.placeholder.com/300x200/ff6b6b/ffffff?text=STATIC+DATA",
					"is_static":   true,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorMsg)
		return
	}

	// Real NewsAPI call
	url := fmt.Sprintf("https://newsapi.org/v2/top-headlines?category=%s&apiKey=%s&pageSize=20", category, apiKey)
	
	log.Printf("Fetching real news from: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch news: %v", err)
		http.Error(w, "Failed to fetch news from API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("News API returned status: %d", resp.StatusCode)
		http.Error(w, fmt.Sprintf("News API error: %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	var newsResp struct {
		Status       string `json:"status"`
		TotalResults int    `json:"totalResults"`
		Articles     []struct {
			Source struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"source"`
			Author      string    `json:"author"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			URL         string    `json:"url"`
			URLToImage  string    `json:"urlToImage"`
			PublishedAt time.Time `json:"publishedAt"`
			Content     string    `json:"content"`
		} `json:"articles"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&newsResp); err != nil {
		http.Error(w, "Failed to decode news response", http.StatusInternalServerError)
		return
	}

	// Transform to our format
	articles := make([]map[string]interface{}, 0, len(newsResp.Articles))
	for _, article := range newsResp.Articles {
		articles = append(articles, map[string]interface{}{
			"id":          fmt.Sprintf("news_%d", time.Now().UnixNano()),
			"title":       article.Title,
			"description": article.Description,
			"url":         article.URL,
			"source":      article.Source.Name,
			"category":    category,
			"published_at": article.PublishedAt,
			"image_url":   article.URLToImage,
		})
	}

	result := map[string]interface{}{
		"category": category,
		"count":    len(articles),
		"articles": articles,
		"source":   "REAL API DATA",
		"api_key_status": "VALID",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getTrendingNews(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		// Fallback to mock data
		articles := []map[string]interface{}{
			{
				"id":          fmt.Sprintf("trending_%d", time.Now().UnixNano()),
				"title":       "AI Breakthrough in Healthcare",
				"description": "New AI technology promises to revolutionize medical diagnosis.",
				"url":         "https://example.com/trending/1",
				"source":      "Tech News",
				"category":    "technology",
				"published_at": time.Now(),
				"image_url":   "https://via.placeholder.com/300x200",
			},
		}
		result := map[string]interface{}{
			"count":    len(articles),
			"articles": articles,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	// Get trending from multiple categories
	categories := []string{"technology", "business", "entertainment", "sports"}
	allArticles := make([]map[string]interface{}, 0)

	for _, category := range categories {
		url := fmt.Sprintf("https://newsapi.org/v2/top-headlines?category=%s&pageSize=5&apiKey=%s", category, apiKey)
		
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Failed to fetch %s news: %v", category, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("News API returned status %d for category %s", resp.StatusCode, category)
			continue
		}

		var newsResp struct {
			Articles []struct {
				Source struct {
					Name string `json:"name"`
				} `json:"source"`
				Title       string    `json:"title"`
				Description string    `json:"description"`
				URL         string    `json:"url"`
				URLToImage  string    `json:"urlToImage"`
				PublishedAt time.Time `json:"publishedAt"`
			} `json:"articles"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&newsResp); err != nil {
			log.Printf("Failed to decode %s news response: %v", category, err)
			continue
		}

		for _, article := range newsResp.Articles {
			allArticles = append(allArticles, map[string]interface{}{
				"id":          fmt.Sprintf("trending_%d", time.Now().UnixNano()),
				"title":       article.Title,
				"description": article.Description,
				"url":         article.URL,
				"source":      article.Source.Name,
				"category":    category,
				"published_at": article.PublishedAt,
				"image_url":   article.URLToImage,
			})
		}
	}

	result := map[string]interface{}{
		"count":    len(allArticles),
		"articles": allArticles,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func searchNews(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		// Fallback to mock data
		articles := []map[string]interface{}{
			{
				"id":          fmt.Sprintf("search_%d", time.Now().UnixNano()),
				"title":       fmt.Sprintf("Search Results for: %s", query),
				"description": fmt.Sprintf("Latest news and updates related to %s", query),
				"url":         "https://example.com/search/1",
				"source":      "Search News",
				"category":    "search",
				"published_at": time.Now(),
				"image_url":   "https://via.placeholder.com/300x200",
			},
		}
		result := map[string]interface{}{
			"query":    query,
			"count":    len(articles),
			"articles": articles,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	// Real NewsAPI search
	url := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&pageSize=20&sortBy=publishedAt&apiKey=%s", query, apiKey)
	
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to search news", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "News API error", http.StatusInternalServerError)
		return
	}

	var newsResp struct {
		Articles []struct {
			Source struct {
				Name string `json:"name"`
			} `json:"source"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			URL         string    `json:"url"`
			URLToImage  string    `json:"urlToImage"`
			PublishedAt time.Time `json:"publishedAt"`
		} `json:"articles"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&newsResp); err != nil {
		http.Error(w, "Failed to decode search response", http.StatusInternalServerError)
		return
	}

	// Transform to our format
	articles := make([]map[string]interface{}, 0, len(newsResp.Articles))
	for _, article := range newsResp.Articles {
		articles = append(articles, map[string]interface{}{
			"id":          fmt.Sprintf("search_%d", time.Now().UnixNano()),
			"title":       article.Title,
			"description": article.Description,
			"url":         article.URL,
			"source":      article.Source.Name,
			"category":    "search",
			"published_at": article.PublishedAt,
			"image_url":   article.URLToImage,
		})
	}

	result := map[string]interface{}{
		"query":    query,
		"count":    len(articles),
		"articles": articles,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Helper function to get user preferences
func getUserPreferences(userID string) map[string]interface{} {
	// In a real app, this would call the user service
	// For now, return default preferences based on common interests
	return map[string]interface{}{
		"news_categories": []string{"technology", "business"},
	}
}
