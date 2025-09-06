package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"gofr.dev"
	"gofr.dev/pkg/gofr"
)

type GatewayService struct {
	newsServiceURL           string
	jobsServiceURL           string
	videosServiceURL         string
	dealsServiceURL          string
	recommendationServiceURL string
	userServiceURL           string
	nftServiceURL            string
}

func main() {
	app := gofr.New()

	gateway := &GatewayService{
		newsServiceURL:           getEnv("NEWS_SERVICE_URL", "http://localhost:8001"),
		jobsServiceURL:           getEnv("JOBS_SERVICE_URL", "http://localhost:8002"),
		videosServiceURL:         getEnv("VIDEOS_SERVICE_URL", "http://localhost:8003"),
		dealsServiceURL:          getEnv("DEALS_SERVICE_URL", "http://localhost:8004"),
		recommendationServiceURL: getEnv("RECOMMENDATION_SERVICE_URL", "http://localhost:8005"),
		userServiceURL:           getEnv("USER_SERVICE_URL", "http://localhost:8006"),
		nftServiceURL:            getEnv("NFT_SERVICE_URL", "http://localhost:8007"),
	}

	// Health check
	app.GET("/health", func(ctx *gofr.Context) (interface{}, error) {
		return map[string]string{"status": "healthy", "service": "api-gateway"}, nil
	})

	// News endpoints
	app.GET("/api/news", gateway.proxyToService(gateway.newsServiceURL+"/api/news"))
	app.GET("/api/news/trending", gateway.proxyToService(gateway.newsServiceURL+"/api/news/trending"))
	app.GET("/api/news/search", gateway.proxyToService(gateway.newsServiceURL+"/api/news/search"))

	// Jobs endpoints
	app.GET("/api/jobs", gateway.proxyToService(gateway.jobsServiceURL+"/api/jobs"))
	app.GET("/api/jobs/trending", gateway.proxyToService(gateway.jobsServiceURL+"/api/jobs/trending"))
	app.GET("/api/jobs/search", gateway.proxyToService(gateway.jobsServiceURL+"/api/jobs/search"))

	// Videos endpoints
	app.GET("/api/videos", gateway.proxyToService(gateway.videosServiceURL+"/api/videos"))
	app.GET("/api/videos/trending", gateway.proxyToService(gateway.videosServiceURL+"/api/videos/trending"))
	app.GET("/api/videos/search", gateway.proxyToService(gateway.videosServiceURL+"/api/videos/search"))

	// Deals endpoints
	app.GET("/api/deals", gateway.proxyToService(gateway.dealsServiceURL+"/api/deals"))
	app.GET("/api/deals/trending", gateway.proxyToService(gateway.dealsServiceURL+"/api/deals/trending"))
	app.GET("/api/deals/search", gateway.proxyToService(gateway.dealsServiceURL+"/api/deals/search"))

	// Recommendation endpoints
	app.GET("/api/recommendations", gateway.proxyToService(gateway.recommendationServiceURL+"/api/recommendations"))

	// User endpoints
	app.POST("/api/users", gateway.proxyToService(gateway.userServiceURL+"/api/users"))
	app.GET("/api/users/:id", gateway.proxyToService(gateway.userServiceURL+"/api/users/:id"))

	// NFT endpoints
	app.POST("/api/nft/mint", gateway.proxyToService(gateway.nftServiceURL+"/api/nft/mint"))
	app.GET("/api/nft/:user_id", gateway.proxyToService(gateway.nftServiceURL+"/api/nft/:user_id"))

	// Start server
	app.Start()
}

func (gs *GatewayService) proxyToService(targetURL string) func(ctx *gofr.Context) (interface{}, error) {
	return func(ctx *gofr.Context) (interface{}, error) {
		// Build the full URL with query parameters
		fullURL := targetURL
		if len(ctx.Request().URL.RawQuery) > 0 {
			fullURL += "?" + ctx.Request().URL.RawQuery
		}

		// Create HTTP request
		req, err := http.NewRequest(ctx.Request().Method, fullURL, ctx.Request().Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		// Copy headers
		for key, values := range ctx.Request().Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Make request
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Read response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %v", err)
		}

		// Parse JSON response
		var result interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to parse response: %v", err)
		}

		return result, nil
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
