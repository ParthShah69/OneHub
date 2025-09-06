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

type DealsService struct {
	amazonAPIKey  string
	flipkartAPIKey string
	cache         *cache.Cache
}

func main() {
	app := gofr.New()

	dealsService := &DealsService{
		amazonAPIKey:  os.Getenv("AMAZON_API_KEY"),
		flipkartAPIKey: os.Getenv("FLIPKART_API_KEY"),
		cache:         cache.New(10*time.Minute, 20*time.Minute),
	}

	// Health check
	app.GET("/health", func(ctx *gofr.Context) (interface{}, error) {
		return map[string]string{"status": "healthy", "service": "deals"}, nil
	})

	// Get deals by category
	app.GET("/api/deals", dealsService.GetDeals)
	
	// Get trending deals
	app.GET("/api/deals/trending", dealsService.GetTrendingDeals)
	
	// Search deals
	app.GET("/api/deals/search", dealsService.SearchDeals)

	app.Run()
}

func (ds *DealsService) GetDeals(ctx *gofr.Context) (interface{}, error) {
	category := ctx.Param("category")
	if category == "" {
		category = "electronics"
	}

	// Check cache first
	if cached, found := ds.cache.Get("deals_" + category); found {
		return cached, nil
	}

	// Fetch deals from multiple platforms
	deals := ds.fetchDealsFromMultipleSources(category, 20)

	result := map[string]interface{}{
		"category": category,
		"count":    len(deals),
		"deals":    deals,
	}

	// Cache the result
	ds.cache.Set("deals_"+category, result, cache.DefaultExpiration)

	return result, nil
}

func (ds *DealsService) GetTrendingDeals(ctx *gofr.Context) (interface{}, error) {
	// Check cache first
	if cached, found := ds.cache.Get("trending_deals"); found {
		return cached, nil
	}

	// Get trending deals from multiple categories
	categories := []string{"electronics", "fashion", "home", "books"}
	allDeals := make([]map[string]interface{}, 0)

	for _, category := range categories {
		deals := ds.fetchDealsFromMultipleSources(category, 5)
		allDeals = append(allDeals, deals...)
	}

	result := map[string]interface{}{
		"count": len(allDeals),
		"deals": allDeals,
	}

	// Cache the result
	ds.cache.Set("trending_deals", result, cache.DefaultExpiration)

	return result, nil
}

func (ds *DealsService) SearchDeals(ctx *gofr.Context) (interface{}, error) {
	query := ctx.Param("q")
	if query == "" {
		return nil, fmt.Errorf("query parameter 'q' is required")
	}

	// Check cache first
	cacheKey := "search_deals_" + query
	if cached, found := ds.cache.Get(cacheKey); found {
		return cached, nil
	}

	deals := ds.searchDealsFromMultipleSources(query, 20)

	result := map[string]interface{}{
		"query": query,
		"count": len(deals),
		"deals": deals,
	}

	// Cache the result
	ds.cache.Set(cacheKey, result, cache.DefaultExpiration)

	return result, nil
}

func (ds *DealsService) fetchDealsFromMultipleSources(category string, limit int) []map[string]interface{} {
	deals := make([]map[string]interface{}, 0)
	
	// Amazon deals
	amazonDeals := ds.fetchAmazonDeals(category, limit/2)
	deals = append(deals, amazonDeals...)
	
	// Flipkart deals
	flipkartDeals := ds.fetchFlipkartDeals(category, limit/2)
	deals = append(deals, flipkartDeals...)
	
	// Mock additional deals for demo
	mockDeals := ds.generateMockDeals(category, limit-len(deals))
	deals = append(deals, mockDeals...)
	
	return deals
}

func (ds *DealsService) fetchAmazonDeals(category string, limit int) []map[string]interface{} {
	// Real Amazon Product Advertising API integration
	// Note: This requires Amazon Associate credentials
	deals := make([]map[string]interface{}, 0)
	
	// Mock Amazon deals for demo
	mockDeals := []map[string]interface{}{
		{
			"id":            fmt.Sprintf("amazon_%d", time.Now().UnixNano()),
			"title":         "Wireless Bluetooth Headphones - " + category,
			"description":   "High-quality wireless headphones with noise cancellation",
			"url":           "https://amazon.com/dp/B08XYZ123",
			"platform":      "Amazon",
			"category":      category,
			"price":         79.99,
			"original_price": 129.99,
			"discount":      38.46,
			"image_url":     "https://images.amazon.com/headphones.jpg",
			"valid_until":   time.Now().Add(7 * 24 * time.Hour),
		},
		{
			"id":            fmt.Sprintf("amazon_%d", time.Now().UnixNano()+1),
			"title":         "Smart Watch Series 8 - " + category,
			"description":   "Latest smartwatch with health monitoring features",
			"url":           "https://amazon.com/dp/B08XYZ124",
			"platform":      "Amazon",
			"category":      category,
			"price":         299.99,
			"original_price": 399.99,
			"discount":      25.0,
			"image_url":     "https://images.amazon.com/smartwatch.jpg",
			"valid_until":   time.Now().Add(5 * 24 * time.Hour),
		},
	}
	
	if len(mockDeals) > limit {
		mockDeals = mockDeals[:limit]
	}
	
	return mockDeals
}

func (ds *DealsService) fetchFlipkartDeals(category string, limit int) []map[string]interface{} {
	// Real Flipkart API integration
	deals := make([]map[string]interface{}, 0)
	
	// Mock Flipkart deals for demo
	mockDeals := []map[string]interface{}{
		{
			"id":            fmt.Sprintf("flipkart_%d", time.Now().UnixNano()),
			"title":         "Laptop Gaming Pro - " + category,
			"description":   "High-performance gaming laptop with RTX graphics",
			"url":           "https://flipkart.com/laptop-gaming-pro",
			"platform":      "Flipkart",
			"category":      category,
			"price":         89999.0,
			"original_price": 119999.0,
			"discount":      25.0,
			"image_url":     "https://images.flipkart.com/laptop.jpg",
			"valid_until":   time.Now().Add(3 * 24 * time.Hour),
		},
		{
			"id":            fmt.Sprintf("flipkart_%d", time.Now().UnixNano()+1),
			"title":         "Smartphone Galaxy S23 - " + category,
			"description":   "Latest flagship smartphone with advanced camera",
			"url":           "https://flipkart.com/galaxy-s23",
			"platform":      "Flipkart",
			"category":      category,
			"price":         69999.0,
			"original_price": 89999.0,
			"discount":      22.22,
			"image_url":     "https://images.flipkart.com/galaxy-s23.jpg",
			"valid_until":   time.Now().Add(2 * 24 * time.Hour),
		},
	}
	
	if len(mockDeals) > limit {
		mockDeals = mockDeals[:limit]
	}
	
	return mockDeals
}

func (ds *DealsService) generateMockDeals(category string, count int) []map[string]interface{} {
	deals := make([]map[string]interface{}, 0, count)
	
	platforms := []string{"Amazon", "Flipkart", "Myntra", "Nykaa", "BigBasket"}
	products := map[string][]string{
		"electronics": {"Smartphone", "Laptop", "Headphones", "Tablet", "Smart Watch"},
		"fashion":     {"T-Shirt", "Jeans", "Shoes", "Dress", "Jacket"},
		"home":        {"Coffee Maker", "Vacuum Cleaner", "Air Purifier", "Blender", "Microwave"},
		"books":       {"Programming Book", "Novel", "Biography", "Cookbook", "Self-Help"},
	}
	
	productList := products[category]
	if len(productList) == 0 {
		productList = []string{"Product A", "Product B", "Product C"}
	}
	
	for i := 0; i < count; i++ {
		platform := platforms[i%len(platforms)]
		product := productList[i%len(productList)]
		
		originalPrice := float64(1000 + (i*500))
		discount := float64(10 + (i*5))
		price := originalPrice * (1 - discount/100)
		
		deals = append(deals, map[string]interface{}{
			"id":            fmt.Sprintf("mock_%d", time.Now().UnixNano()+int64(i)),
			"title":         fmt.Sprintf("%s - %s", product, category),
			"description":   fmt.Sprintf("High-quality %s for %s category", product, category),
			"url":           fmt.Sprintf("https://%s.com/product-%d", platform, i),
			"platform":      platform,
			"category":      category,
			"price":         price,
			"original_price": originalPrice,
			"discount":      discount,
			"image_url":     fmt.Sprintf("https://images.%s.com/product-%d.jpg", platform, i),
			"valid_until":   time.Now().Add(time.Duration(1+i) * 24 * time.Hour),
		})
	}
	
	return deals
}

func (ds *DealsService) searchDealsFromMultipleSources(query string, limit int) []map[string]interface{} {
	deals := make([]map[string]interface{}, 0)
	
	// Search across all categories
	categories := []string{"electronics", "fashion", "home", "books"}
	
	for _, category := range categories {
		categoryDeals := ds.fetchDealsFromMultipleSources(category, limit/len(categories))
		// Filter deals that match the query
		for _, deal := range categoryDeals {
			title := deal["title"].(string)
			if containsIgnoreCase(title, query) {
				deals = append(deals, deal)
			}
		}
	}
	
	if len(deals) > limit {
		deals = deals[:limit]
	}
	
	return deals
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    len(s) > len(substr) && 
		    (s[:len(substr)] == substr || 
		     s[len(s)-len(substr):] == substr ||
		     containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
