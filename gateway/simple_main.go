package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "service": "api-gateway"})
	})

	// News endpoints
	http.HandleFunc("/api/news", proxyToService("http://localhost:8001/api/news"))
	http.HandleFunc("/api/news/trending", proxyToService("http://localhost:8001/api/news/trending"))
	http.HandleFunc("/api/news/search", proxyToService("http://localhost:8001/api/news/search"))

	// Jobs endpoints
	http.HandleFunc("/api/jobs", proxyToService("http://localhost:8002/api/jobs"))
	http.HandleFunc("/api/jobs/trending", proxyToService("http://localhost:8002/api/jobs/trending"))
	http.HandleFunc("/api/jobs/search", proxyToService("http://localhost:8002/api/jobs/search"))

	// Videos endpoints
	http.HandleFunc("/api/videos", proxyToService("http://localhost:8003/api/videos"))
	http.HandleFunc("/api/videos/trending", proxyToService("http://localhost:8003/api/videos/trending"))
	http.HandleFunc("/api/videos/search", proxyToService("http://localhost:8003/api/videos/search"))

	// Deals endpoints
	http.HandleFunc("/api/deals", proxyToService("http://localhost:8004/api/deals"))
	http.HandleFunc("/api/deals/trending", proxyToService("http://localhost:8004/api/deals/trending"))
	http.HandleFunc("/api/deals/search", proxyToService("http://localhost:8004/api/deals/search"))

	// Movies endpoints
	http.HandleFunc("/api/movies", proxyToService("http://localhost:8008/api/movies"))
	http.HandleFunc("/api/movies/trending", proxyToService("http://localhost:8008/api/movies/trending"))
	http.HandleFunc("/api/movies/search", proxyToService("http://localhost:8008/api/movies/search"))

	// Food endpoints
	http.HandleFunc("/api/food", proxyToService("http://localhost:8009/api/food"))
	http.HandleFunc("/api/food/trending", proxyToService("http://localhost:8009/api/food/trending"))
	http.HandleFunc("/api/food/search", proxyToService("http://localhost:8009/api/food/search"))

	// Recommendation endpoints
	http.HandleFunc("/api/recommendations", proxyToService("http://localhost:8005/api/recommendations"))

	// User endpoints
	http.HandleFunc("/api/users", proxyToService("http://localhost:8006/api/users"))

	// NFT endpoints
	http.HandleFunc("/api/nft/mint", proxyToService("http://localhost:8007/api/nft/mint"))
	http.HandleFunc("/api/nft/", proxyToService("http://localhost:8007/api/nft/"))

	log.Printf("API Gateway starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func proxyToService(targetURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Build the full URL with query parameters
		fullURL := targetURL
		if len(r.URL.RawQuery) > 0 {
			fullURL += "?" + r.URL.RawQuery
		}

		// Create HTTP request
		req, err := http.NewRequest(r.Method, fullURL, r.Body)
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		// Copy headers
		for key, values := range r.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Make request
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Failed to make request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		// Copy response body
		io.Copy(w, resp.Body)
	}
}
