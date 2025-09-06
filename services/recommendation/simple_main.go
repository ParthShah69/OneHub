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
	port := "8005"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "service": "recommendation"})
	})

	http.HandleFunc("/api/recommendations", getRecommendations)

	log.Printf("Recommendation service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getRecommendations(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "default_user"
	}

	// Mock personalized recommendations
	recommendations := []map[string]interface{}{
		{
			"id":          fmt.Sprintf("rec_%d", time.Now().UnixNano()),
			"title":       "AI-Powered News Article",
			"description": "Latest developments in artificial intelligence and machine learning.",
			"url":         "https://example.com/news/ai",
			"content_type": "news",
			"recommendation_score": 0.95,
			"reason":      "Recommended based on your interest in AI",
			"image_url":   "https://via.placeholder.com/300x200",
		},
		{
			"id":          fmt.Sprintf("rec_%d", time.Now().UnixNano()+1),
			"title":       "Senior Software Engineer Position",
			"description": "Join our team as a senior software engineer working on cutting-edge projects.",
			"url":         "https://example.com/jobs/software",
			"content_type": "jobs",
			"recommendation_score": 0.88,
			"reason":      "Recommended based on your interest in technology",
			"image_url":   "https://via.placeholder.com/300x200",
		},
		{
			"id":          fmt.Sprintf("rec_%d", time.Now().UnixNano()+2),
			"title":       "Tech Tutorial Video",
			"description": "Learn the latest programming techniques and best practices.",
			"url":         "https://youtube.com/watch?v=tutorial",
			"content_type": "videos",
			"recommendation_score": 0.82,
			"reason":      "Recommended based on your interest in programming",
			"image_url":   "https://via.placeholder.com/300x200",
		},
	}

	result := map[string]interface{}{
		"user_id":        userID,
		"count":          len(recommendations),
		"recommendations": recommendations,
		"generated_at":   time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
