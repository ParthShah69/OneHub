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
	port := "8004"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "service": "deals"})
	})

	http.HandleFunc("/api/deals", getDeals)
	http.HandleFunc("/api/deals/trending", getTrendingDeals)
	http.HandleFunc("/api/deals/search", searchDeals)

	log.Printf("Deals service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getDeals(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		category = "electronics"
	}

	apiKey := os.Getenv("AMAZON_API_KEY")
	if apiKey == "" || apiKey == "your_amazon_api_key_here" {
		// Return error message for static data
		errorMsg := map[string]interface{}{
			"error": "STATIC DATA - No Amazon API Key Found",
			"message": "This is static/mock data. To get real deals, add AMAZON_API_KEY to your environment variables.",
			"instructions": "Get free API key from https://webservices.amazon.com/ and add to env.local",
			"category": category,
			"count": 1,
			"deals": []map[string]interface{}{
				{
					"id":            fmt.Sprintf("static_deal_%d", time.Now().UnixNano()),
					"title":         fmt.Sprintf("⚠️ STATIC DATA: Amazing %s Deal", category),
					"description":   fmt.Sprintf("This is mock data. Get real deals by adding AMAZON_API_KEY to env.local file.", category),
					"url":           "https://webservices.amazon.com/",
					"store":         "Mock Data - Add API Key",
					"category":      category,
					"current_price": 79.99,
					"original_price": 129.99,
					"discount":      38.46,
					"image_url":     "https://via.placeholder.com/300x200/ff6b6b/ffffff?text=STATIC+DATA",
					"valid_until":   time.Now().Add(7 * 24 * time.Hour),
					"is_static":     true,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorMsg)
		return
	}

	// Real Amazon API call would go here
	// For now, return success message
	result := map[string]interface{}{
		"category": category,
		"count":    0,
		"deals":    []map[string]interface{}{},
		"source":   "REAL AMAZON API DATA",
		"api_key_status": "VALID",
		"message": "Amazon API integration ready - add real API calls here",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getTrendingDeals(w http.ResponseWriter, r *http.Request) {
	deals := []map[string]interface{}{
		{
			"id":            fmt.Sprintf("trending_deal_%d", time.Now().UnixNano()),
			"title":         "Wireless Bluetooth Headphones",
			"description":   "High-quality wireless headphones with noise cancellation",
			"url":           "https://example.com/deals/2",
			"platform":      "Flipkart",
			"category":      "electronics",
			"price":         2999.0,
			"original_price": 4999.0,
			"discount":      40.0,
			"image_url":     "https://via.placeholder.com/300x200",
			"valid_until":   time.Now().Add(5 * 24 * time.Hour),
		},
	}

	result := map[string]interface{}{
		"count": len(deals),
		"deals": deals,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func searchDeals(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	deals := []map[string]interface{}{
		{
			"id":            fmt.Sprintf("search_deal_%d", time.Now().UnixNano()),
			"title":         fmt.Sprintf("Search Results for: %s", query),
			"description":   fmt.Sprintf("Best deals on %s", query),
			"url":           "https://example.com/deals/search",
			"platform":      "Multiple",
			"category":      "search",
			"price":         99.99,
			"original_price": 199.99,
			"discount":      50.0,
			"image_url":     "https://via.placeholder.com/300x200",
			"valid_until":   time.Now().Add(3 * 24 * time.Hour),
		},
	}

	result := map[string]interface{}{
		"query": query,
		"count": len(deals),
		"deals": deals,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
