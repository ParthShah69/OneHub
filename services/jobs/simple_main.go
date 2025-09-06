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
	port := "8002"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "service": "jobs"})
	})

	http.HandleFunc("/api/jobs", getJobs)
	http.HandleFunc("/api/jobs/trending", getTrendingJobs)
	http.HandleFunc("/api/jobs/search", searchJobs)

	log.Printf("Jobs service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getJobs(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		category = "technology"
	}

	appID := os.Getenv("ADZUNA_APP_ID")
	appKey := os.Getenv("ADZUNA_APP_KEY")
	if appID == "" || appID == "your_adzuna_app_id_here" || appKey == "" || appKey == "your_adzuna_app_key_here" {
		// Return error message for static data
		errorMsg := map[string]interface{}{
			"error": "STATIC DATA - No Adzuna API Keys Found",
			"message": "This is static/mock data. To get real job postings, add ADZUNA_APP_ID and ADZUNA_APP_KEY to your environment variables.",
			"instructions": "Get FREE API keys from https://developer.adzuna.com/ and add to env.local",
			"category": category,
			"count": 1,
			"jobs": []map[string]interface{}{
				{
					"id":          fmt.Sprintf("static_job_%d", time.Now().UnixNano()),
					"title":       fmt.Sprintf("⚠️ STATIC DATA: Senior %s Engineer", category),
					"company":     "Mock Data - Add API Key",
					"location":    "San Francisco, CA",
					"type":        "Full-time",
					"salary":      "$120,000 - $180,000",
					"description": fmt.Sprintf("This is mock data. Get real job postings by adding LINKEDIN_API_KEY to env.local file.", category),
					"url":         "https://developer.linkedin.com/",
					"category":    category,
					"posted_at":   time.Now().Add(-24 * time.Hour),
					"is_static":   true,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorMsg)
		return
	}

	// Real Adzuna API call
	url := fmt.Sprintf("https://api.adzuna.com/v1/api/jobs/us/search/1?app_id=%s&app_key=%s&what=%s&results_per_page=20", appID, appKey, category)
	
	log.Printf("Fetching real jobs from: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch jobs: %v", err)
		http.Error(w, "Failed to fetch jobs from API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var apiResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Printf("Failed to decode API response: %v", err)
		http.Error(w, "Failed to decode API response", http.StatusInternalServerError)
		return
	}

	// Transform Adzuna response to our format
	jobs := []map[string]interface{}{}
	if results, ok := apiResponse["results"].([]interface{}); ok {
		for _, job := range results {
			if jobData, ok := job.(map[string]interface{}); ok {
				transformedJob := map[string]interface{}{
					"id":          jobData["id"],
					"title":       jobData["title"],
					"company":     jobData["company"],
					"location":    jobData["location"],
					"type":        jobData["contract_type"],
					"salary":      jobData["salary_min"],
					"description": jobData["description"],
					"url":         jobData["redirect_url"],
					"category":    category,
					"posted_at":   jobData["created"],
					"is_static":   false,
				}
				jobs = append(jobs, transformedJob)
			}
		}
	}

	result := map[string]interface{}{
		"category": category,
		"count":    len(jobs),
		"jobs":     jobs,
		"source":   "REAL ADZUNA API DATA",
		"api_key_status": "VALID",
		"message": "Real job data from Adzuna API",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getTrendingJobs(w http.ResponseWriter, r *http.Request) {
	jobs := []map[string]interface{}{
		{
			"id":          fmt.Sprintf("trending_job_%d", time.Now().UnixNano()),
			"title":       "AI/ML Engineer",
			"company":     "AI Startup Co.",
			"location":    "Remote",
			"description": "Join our AI team to build cutting-edge machine learning solutions.",
			"url":         "https://example.com/jobs/2",
			"category":    "ai",
			"posted_at":   time.Now().Add(-12 * time.Hour),
			"salary":      "$100,000 - $150,000",
		},
	}

	result := map[string]interface{}{
		"count": len(jobs),
		"jobs":  jobs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func searchJobs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	jobs := []map[string]interface{}{
		{
			"id":          fmt.Sprintf("search_job_%d", time.Now().UnixNano()),
			"title":       fmt.Sprintf("Search Results for: %s", query),
			"company":     "Various Companies",
			"location":    "Multiple Locations",
			"description": fmt.Sprintf("Job opportunities related to %s", query),
			"url":         "https://example.com/jobs/search",
			"category":    "search",
			"posted_at":   time.Now(),
			"salary":      "Competitive",
		},
	}

	result := map[string]interface{}{
		"query": query,
		"count": len(jobs),
		"jobs":  jobs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
