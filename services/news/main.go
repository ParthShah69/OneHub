package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gofr.dev/pkg/gofr"
	"github.com/patrickmn/go-cache"
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

type NewsService struct {
	apiKey string
	cache  *cache.Cache
}

func main() {
	app := gofr.New()

	newsService := &NewsService{
		apiKey: os.Getenv("NEWS_API_KEY"),
		cache:  cache.New(5*time.Minute, 10*time.Minute),
	}

	// Health check
	app.GET("/health", func(ctx *gofr.Context) (interface{}, error) {
		return map[string]string{"status": "healthy", "service": "news"}, nil
	})

	// Get news by category
	app.GET("/api/news", newsService.GetNews)
	
	// Get trending news
	app.GET("/api/news/trending", newsService.GetTrendingNews)
	
	// Get news by query
	app.GET("/api/news/search", newsService.SearchNews)

	app.Run()
}

func (ns *NewsService) GetNews(ctx *gofr.Context) (interface{}, error) {
	category := ctx.Param("category")
	if category == "" {
		category = "general"
	}

	// Check cache first
	if cached, found := ns.cache.Get("news_" + category); found {
		return cached, nil
	}

	url := fmt.Sprintf("https://newsapi.org/v2/top-headlines?category=%s&apiKey=%s", category, ns.apiKey)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch news: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("news API returned status: %d", resp.StatusCode)
	}

	var newsResp NewsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&newsResp); err != nil {
		return nil, fmt.Errorf("failed to decode news response: %v", err)
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
	}

	// Cache the result
	ns.cache.Set("news_"+category, result, cache.DefaultExpiration)

	return result, nil
}

func (ns *NewsService) GetTrendingNews(ctx *gofr.Context) (interface{}, error) {
	// Check cache first
	if cached, found := ns.cache.Get("trending_news"); found {
		return cached, nil
	}

	// Get trending from multiple categories
	categories := []string{"technology", "business", "entertainment", "sports"}
	allArticles := make([]map[string]interface{}, 0)

	for _, category := range categories {
		url := fmt.Sprintf("https://newsapi.org/v2/top-headlines?category=%s&pageSize=5&apiKey=%s", category, ns.apiKey)
		
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

		var newsResp NewsAPIResponse
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

	// Cache the result
	ns.cache.Set("trending_news", result, cache.DefaultExpiration)

	return result, nil
}

func (ns *NewsService) SearchNews(ctx *gofr.Context) (interface{}, error) {
	query := ctx.Param("q")
	if query == "" {
		return nil, fmt.Errorf("query parameter 'q' is required")
	}

	pageSize := ctx.Param("pageSize")
	if pageSize == "" {
		pageSize = "20"
	}

	// Check cache first
	cacheKey := "search_" + query + "_" + pageSize
	if cached, found := ns.cache.Get(cacheKey); found {
		return cached, nil
	}

	url := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&pageSize=%s&sortBy=publishedAt&apiKey=%s", 
		query, pageSize, ns.apiKey)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to search news: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("news API returned status: %d", resp.StatusCode)
	}

	var newsResp NewsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&newsResp); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %v", err)
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

	// Cache the result
	ns.cache.Set(cacheKey, result, cache.DefaultExpiration)

	return result, nil
}
