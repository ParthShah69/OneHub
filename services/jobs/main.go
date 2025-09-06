package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gofr.dev/pkg/gofr"
	"github.com/patrickmn/go-cache"
)

type LinkedInJobResponse struct {
	Elements []struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		Company   struct {
			Name string `json:"name"`
		} `json:"company"`
		Location  struct {
			City    string `json:"city"`
			Country string `json:"country"`
		} `json:"location"`
		Description struct {
			Text string `json:"text"`
		} `json:"description"`
		JobPostingURL string `json:"jobPostingUrl"`
		PostedAt      string `json:"postedAt"`
		Salary        struct {
			Min int `json:"min"`
			Max int `json:"max"`
		} `json:"salary"`
	} `json:"elements"`
	Paging struct {
		Count int `json:"count"`
		Start int `json:"start"`
	} `json:"paging"`
}

type JobsService struct {
	apiKey string
	cache  *cache.Cache
}

func main() {
	app := gofr.New()

	jobsService := &JobsService{
		apiKey: os.Getenv("LINKEDIN_API_KEY"),
		cache:  cache.New(10*time.Minute, 20*time.Minute),
	}

	// Health check
	app.GET("/health", func(ctx *gofr.Context) (interface{}, error) {
		return map[string]string{"status": "healthy", "service": "jobs"}, nil
	})

	// Get jobs by category
	app.GET("/api/jobs", jobsService.GetJobs)
	
	// Get trending jobs
	app.GET("/api/jobs/trending", jobsService.GetTrendingJobs)
	
	// Search jobs
	app.GET("/api/jobs/search", jobsService.SearchJobs)

	app.Run()
}

func (js *JobsService) GetJobs(ctx *gofr.Context) (interface{}, error) {
	category := ctx.Param("category")
	if category == "" {
		category = "technology"
	}

	// Check cache first
	if cached, found := js.cache.Get("jobs_" + category); found {
		return cached, nil
	}

	// Map categories to LinkedIn keywords
	keywords := map[string]string{
		"technology": "software engineer OR developer OR programmer",
		"ai":         "artificial intelligence OR machine learning OR AI engineer",
		"cloud":      "cloud engineer OR AWS OR Azure OR GCP",
		"startup":    "startup OR early stage OR venture",
		"remote":     "remote OR work from home",
		"finance":    "financial analyst OR investment OR banking",
		"marketing":  "marketing OR digital marketing OR growth",
		"design":     "designer OR UX OR UI OR product design",
	}

	keyword := keywords[category]
	if keyword == "" {
		keyword = category
	}

	// Use LinkedIn Jobs API (Note: This is a simplified version - real LinkedIn API requires OAuth)
	jobs := js.fetchJobsFromLinkedIn(keyword, 20)

	result := map[string]interface{}{
		"category": category,
		"count":    len(jobs),
		"jobs":     jobs,
	}

	// Cache the result
	js.cache.Set("jobs_"+category, result, cache.DefaultExpiration)

	return result, nil
}

func (js *JobsService) GetTrendingJobs(ctx *gofr.Context) (interface{}, error) {
	// Check cache first
	if cached, found := js.cache.Get("trending_jobs"); found {
		return cached, nil
	}

	// Get trending jobs from multiple categories
	categories := []string{"ai", "cloud", "remote", "startup"}
	allJobs := make([]map[string]interface{}, 0)

	for _, category := range categories {
		jobs := js.fetchJobsFromLinkedIn(category, 5)
		allJobs = append(allJobs, jobs...)
	}

	result := map[string]interface{}{
		"count": len(allJobs),
		"jobs":  allJobs,
	}

	// Cache the result
	js.cache.Set("trending_jobs", result, cache.DefaultExpiration)

	return result, nil
}

func (js *JobsService) SearchJobs(ctx *gofr.Context) (interface{}, error) {
	query := ctx.Param("q")
	if query == "" {
		return nil, fmt.Errorf("query parameter 'q' is required")
	}

	location := ctx.Param("location")
	limit := ctx.Param("limit")
	if limit == "" {
		limit = "20"
	}

	// Check cache first
	cacheKey := "search_jobs_" + query + "_" + location + "_" + limit
	if cached, found := js.cache.Get(cacheKey); found {
		return cached, nil
	}

	jobs := js.fetchJobsFromLinkedIn(query+" "+location, 20)

	result := map[string]interface{}{
		"query":  query,
		"count":  len(jobs),
		"jobs":   jobs,
	}

	// Cache the result
	js.cache.Set(cacheKey, result, cache.DefaultExpiration)

	return result, nil
}

func (js *JobsService) fetchJobsFromLinkedIn(keyword string, limit int) []map[string]interface{} {
	// Note: This is a mock implementation since LinkedIn API requires OAuth
	// In a real implementation, you would use the LinkedIn Jobs API with proper authentication
	
	// For demo purposes, we'll return mock data that simulates real job listings
	mockJobs := []map[string]interface{}{
		{
			"id":          fmt.Sprintf("job_%d", time.Now().UnixNano()),
			"title":       "Senior Software Engineer - " + keyword,
			"company":     "TechCorp Inc.",
			"location":    "San Francisco, CA",
			"description": "We are looking for a talented software engineer to join our team...",
			"url":         "https://linkedin.com/jobs/view/123456",
			"category":    "technology",
			"posted_at":   time.Now().Add(-24 * time.Hour),
			"salary":      "$120,000 - $180,000",
		},
		{
			"id":          fmt.Sprintf("job_%d", time.Now().UnixNano()+1),
			"title":       "AI/ML Engineer - " + keyword,
			"company":     "AI Startup Co.",
			"location":    "Remote",
			"description": "Join our AI team to build cutting-edge machine learning solutions...",
			"url":         "https://linkedin.com/jobs/view/123457",
			"category":    "ai",
			"posted_at":   time.Now().Add(-12 * time.Hour),
			"salary":      "$100,000 - $150,000",
		},
		{
			"id":          fmt.Sprintf("job_%d", time.Now().UnixNano()+2),
			"title":       "Cloud Solutions Architect - " + keyword,
			"company":     "CloudTech Solutions",
			"location":    "New York, NY",
			"description": "Design and implement cloud infrastructure solutions...",
			"url":         "https://linkedin.com/jobs/view/123458",
			"category":    "cloud",
			"posted_at":   time.Now().Add(-6 * time.Hour),
			"salary":      "$130,000 - $200,000",
		},
	}

	// Filter and limit results
	if len(mockJobs) > limit {
		mockJobs = mockJobs[:limit]
	}

	return mockJobs
}

// Real LinkedIn API integration would look like this:
func (js *JobsService) fetchRealLinkedInJobs(keyword string, limit int) ([]map[string]interface{}, error) {
	// This would be the actual LinkedIn API call
	// Note: LinkedIn API requires OAuth 2.0 authentication
	
	url := fmt.Sprintf("https://api.linkedin.com/v2/jobSearch?keywords=%s&count=%d", keyword, limit)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+js.apiKey)
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LinkedIn API returned status: %d", resp.StatusCode)
	}
	
	var linkedInResp LinkedInJobResponse
	if err := json.NewDecoder(resp.Body).Decode(&linkedInResp); err != nil {
		return nil, err
	}
	
	// Transform to our format
	jobs := make([]map[string]interface{}, 0, len(linkedInResp.Elements))
	for _, job := range linkedInResp.Elements {
		salary := ""
		if job.Salary.Min > 0 && job.Salary.Max > 0 {
			salary = fmt.Sprintf("$%d - $%d", job.Salary.Min, job.Salary.Max)
		}
		
		jobs = append(jobs, map[string]interface{}{
			"id":          job.ID,
			"title":       job.Title,
			"company":     job.Company.Name,
			"location":    fmt.Sprintf("%s, %s", job.Location.City, job.Location.Country),
			"description": job.Description.Text,
			"url":         job.JobPostingURL,
			"category":    "technology", // You'd determine this based on job content
			"posted_at":   job.PostedAt,
			"salary":      salary,
		})
	}
	
	return jobs, nil
}
