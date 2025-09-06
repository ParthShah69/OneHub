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

type WolframResponse struct {
	QueryResult struct {
		Success bool `json:"success"`
		Pods    []struct {
			Title   string `json:"title"`
			SubPods []struct {
				PlainText string `json:"plaintext"`
			} `json:"subpods"`
		} `json:"pods"`
	} `json:"queryresult"`
}

type RecommendationService struct {
	wolframAPIKey string
	cache         *cache.Cache
}

func main() {
	app := gofr.New()

	recommendationService := &RecommendationService{
		wolframAPIKey: os.Getenv("WOLFRAM_API_KEY"),
		cache:         cache.New(15*time.Minute, 30*time.Minute),
	}

	// Health check
	app.GET("/health", func(ctx *gofr.Context) (interface{}, error) {
		return map[string]string{"status": "healthy", "service": "recommendation"}, nil
	})

	// Get personalized recommendations
	app.GET("/api/recommendations", recommendationService.GetRecommendations)
	
	// Get recommendations by category
	app.GET("/api/recommendations/:category", recommendationService.GetRecommendationsByCategory)

	app.Run()
}

func (rs *RecommendationService) GetRecommendations(ctx *gofr.Context) (interface{}, error) {
	userID := ctx.Param("user_id")
	if userID == "" {
		userID = "default_user"
	}

	// Check cache first
	if cached, found := rs.cache.Get("recommendations_" + userID); found {
		return cached, nil
	}

	// Get user profile and behavior data
	userProfile := rs.getUserProfile(userID)
	
	// Fetch content from all services
	content := rs.fetchAllContent()
	
	// Use Wolfram to generate personalized recommendations
	recommendations, err := rs.generateRecommendationsWithWolfram(userProfile, content)
	if err != nil {
		log.Printf("Wolfram API failed, using fallback: %v", err)
		recommendations = rs.generateFallbackRecommendations(userProfile, content)
	}

	result := map[string]interface{}{
		"user_id":        userID,
		"count":          len(recommendations),
		"recommendations": recommendations,
		"generated_at":   time.Now(),
	}

	// Cache the result
	rs.cache.Set("recommendations_"+userID, result, cache.DefaultExpiration)

	return result, nil
}

func (rs *RecommendationService) GetRecommendationsByCategory(ctx *gofr.Context) (interface{}, error) {
	userID := ctx.Param("user_id")
	if userID == "" {
		userID = "default_user"
	}
	
	category := ctx.Param("category")
	if category == "" {
		return nil, fmt.Errorf("category parameter is required")
	}

	// Check cache first
	cacheKey := fmt.Sprintf("recommendations_%s_%s", userID, category)
	if cached, found := rs.cache.Get(cacheKey); found {
		return cached, nil
	}

	// Get user profile
	userProfile := rs.getUserProfile(userID)
	
	// Fetch content for specific category
	content := rs.fetchContentByCategory(category)
	
	// Generate recommendations
	recommendations, err := rs.generateRecommendationsWithWolfram(userProfile, content)
	if err != nil {
		log.Printf("Wolfram API failed, using fallback: %v", err)
		recommendations = rs.generateFallbackRecommendations(userProfile, content)
	}

	result := map[string]interface{}{
		"user_id":        userID,
		"category":       category,
		"count":          len(recommendations),
		"recommendations": recommendations,
		"generated_at":   time.Now(),
	}

	// Cache the result
	rs.cache.Set(cacheKey, result, cache.DefaultExpiration)

	return result, nil
}

func (rs *RecommendationService) getUserProfile(userID string) map[string]interface{} {
	// In a real implementation, this would fetch from the database
	// For demo purposes, we'll return a mock profile
	return map[string]interface{}{
		"user_id":            userID,
		"explicit_interests": []string{"technology", "ai", "business"},
		"behavioral_scores": map[string]float64{
			"technology": 0.8,
			"ai":         0.9,
			"business":   0.7,
			"entertainment": 0.3,
			"fashion":    0.2,
		},
		"recent_activity": []string{"clicked_ai_article", "bookmarked_tech_job", "shared_business_news"},
	}
}

func (rs *RecommendationService) fetchAllContent() map[string][]map[string]interface{} {
	content := make(map[string][]map[string]interface{})
	
	// Fetch from news service
	newsContent := rs.fetchFromService("http://news-service:8000/api/news/trending", "news")
	content["news"] = newsContent
	
	// Fetch from jobs service
	jobsContent := rs.fetchFromService("http://jobs-service:8000/api/jobs/trending", "jobs")
	content["jobs"] = jobsContent
	
	// Fetch from videos service
	videosContent := rs.fetchFromService("http://videos-service:8000/api/videos/trending", "videos")
	content["videos"] = videosContent
	
	// Fetch from deals service
	dealsContent := rs.fetchFromService("http://deals-service:8000/api/deals/trending", "deals")
	content["deals"] = dealsContent
	
	return content
}

func (rs *RecommendationService) fetchContentByCategory(category string) map[string][]map[string]interface{} {
	content := make(map[string][]map[string]interface{})
	
	// Fetch from all services for the specific category
	services := map[string]string{
		"news":  fmt.Sprintf("http://news-service:8000/api/news?category=%s", category),
		"jobs":  fmt.Sprintf("http://jobs-service:8000/api/jobs?category=%s", category),
		"videos": fmt.Sprintf("http://videos-service:8000/api/videos?category=%s", category),
		"deals": fmt.Sprintf("http://deals-service:8000/api/deals?category=%s", category),
	}
	
	for serviceType, url := range services {
		serviceContent := rs.fetchFromService(url, serviceType)
		content[serviceType] = serviceContent
	}
	
	return content
}

func (rs *RecommendationService) fetchFromService(url, contentType string) []map[string]interface{} {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch from %s: %v", url, err)
		return []map[string]interface{}{}
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		log.Printf("Service %s returned status: %d", url, resp.StatusCode)
		return []map[string]interface{}{}
	}
	
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode response from %s: %v", url, err)
		return []map[string]interface{}{}
	}
	
	// Extract the content array based on service type
	var contentArray []map[string]interface{}
	switch contentType {
	case "news":
		if articles, ok := result["articles"].([]interface{}); ok {
			for _, article := range articles {
				if articleMap, ok := article.(map[string]interface{}); ok {
					contentArray = append(contentArray, articleMap)
				}
			}
		}
	case "jobs":
		if jobs, ok := result["jobs"].([]interface{}); ok {
			for _, job := range jobs {
				if jobMap, ok := job.(map[string]interface{}); ok {
					contentArray = append(contentArray, jobMap)
				}
			}
		}
	case "videos":
		if videos, ok := result["videos"].([]interface{}); ok {
			for _, video := range videos {
				if videoMap, ok := video.(map[string]interface{}); ok {
					contentArray = append(contentArray, videoMap)
				}
			}
		}
	case "deals":
		if deals, ok := result["deals"].([]interface{}); ok {
			for _, deal := range deals {
				if dealMap, ok := deal.(map[string]interface{}); ok {
					contentArray = append(contentArray, dealMap)
				}
			}
		}
	}
	
	return contentArray
}

func (rs *RecommendationService) generateRecommendationsWithWolfram(userProfile map[string]interface{}, content map[string][]map[string]interface{}) ([]map[string]interface{}, error) {
	// Prepare data for Wolfram
	interests := userProfile["explicit_interests"].([]string)
	behavioralScores := userProfile["behavioral_scores"].(map[string]float64)
	
	// Create a query for Wolfram
	query := fmt.Sprintf("Recommend content for user with interests: %v and behavioral scores: %v", 
		interests, behavioralScores)
	
	// Call Wolfram API
	wolframURL := fmt.Sprintf("https://api.wolframalpha.com/v2/query?input=%s&appid=%s&output=json", 
		query, rs.wolframAPIKey)
	
	resp, err := http.Get(wolframURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call Wolfram API: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Wolfram API returned status: %d", resp.StatusCode)
	}
	
	var wolframResp WolframResponse
	if err := json.NewDecoder(resp.Body).Decode(&wolframResp); err != nil {
		return nil, fmt.Errorf("failed to decode Wolfram response: %v", err)
	}
	
	// Process Wolfram response and generate recommendations
	recommendations := rs.processWolframResponse(wolframResp, content, userProfile)
	
	return recommendations, nil
}

func (rs *RecommendationService) processWolframResponse(wolframResp WolframResponse, content map[string][]map[string]interface{}, userProfile map[string]interface{}) []map[string]interface{} {
	recommendations := make([]map[string]interface{}, 0)
	
	// Extract insights from Wolfram response
	insights := make([]string, 0)
	for _, pod := range wolframResp.QueryResult.Pods {
		for _, subPod := range pod.SubPods {
			if subPod.PlainText != "" {
				insights = append(insights, subPod.PlainText)
			}
		}
	}
	
	// Use insights to rank and select content
	behavioralScores := userProfile["behavioral_scores"].(map[string]float64)
	
	// Rank content based on user preferences
	for contentType, items := range content {
		userScore := behavioralScores[contentType]
		if userScore < 0.3 {
			continue // Skip low-interest categories
		}
		
		// Take top items from each category
		limit := int(userScore * 10) // Scale by interest level
		if limit > len(items) {
			limit = len(items)
		}
		if limit > 5 {
			limit = 5 // Max 5 per category
		}
		
		for i := 0; i < limit; i++ {
			item := items[i]
			item["content_type"] = contentType
			item["recommendation_score"] = userScore * (1.0 - float64(i)*0.1) // Decreasing score
			item["reason"] = fmt.Sprintf("Recommended based on your interest in %s", contentType)
			
			recommendations = append(recommendations, item)
		}
	}
	
	return recommendations
}

func (rs *RecommendationService) generateFallbackRecommendations(userProfile map[string]interface{}, content map[string][]map[string]interface{}) []map[string]interface{} {
	recommendations := make([]map[string]interface{}, 0)
	
	behavioralScores := userProfile["behavioral_scores"].(map[string]float64)
	
	// Simple fallback: rank by user interest scores
	for contentType, items := range content {
		userScore := behavioralScores[contentType]
		if userScore < 0.3 {
			continue
		}
		
		limit := int(userScore * 8)
		if limit > len(items) {
			limit = len(items)
		}
		if limit > 3 {
			limit = 3
		}
		
		for i := 0; i < limit; i++ {
			item := items[i]
			item["content_type"] = contentType
			item["recommendation_score"] = userScore
			item["reason"] = fmt.Sprintf("Recommended based on your interest in %s", contentType)
			
			recommendations = append(recommendations, item)
		}
	}
	
	return recommendations
}
