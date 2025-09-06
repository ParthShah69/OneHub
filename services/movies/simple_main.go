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
	port := "8008"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "service": "movies"})
	})

	// Get movies by category
	http.HandleFunc("/api/movies", getMovies)
	
	// Get trending movies
	http.HandleFunc("/api/movies/trending", getTrendingMovies)
	
	// Get movies by query
	http.HandleFunc("/api/movies/search", searchMovies)

	log.Printf("Movies service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	userID := r.URL.Query().Get("user_id")
	
	if category == "" {
		category = "popular"
	}
	
	// If user_id is provided, get user preferences
	if userID != "" {
		preferences := getUserPreferences(userID)
		if len(preferences.MovieGenres) > 0 {
			category = preferences.MovieGenres[0] // Use first preference
		}
	}

	apiKey := os.Getenv("OMDB_API_KEY")
	if apiKey == "" || apiKey == "your_omdb_api_key_here" {
		// Return error message for static data
		errorMsg := map[string]interface{}{
			"error": "STATIC DATA - No OMDb API Key Found",
			"message": "This is static/mock data. To get real movies, add OMDB_API_KEY to your environment variables.",
			"instructions": "Get FREE API key from http://www.omdbapi.com/apikey.aspx and add to env.local",
			"category": category,
			"count": 1,
			"movies": []map[string]interface{}{
				{
					"id":          fmt.Sprintf("static_movie_%d", time.Now().UnixNano()),
					"title":       fmt.Sprintf("⚠️ STATIC DATA: Latest %s Movie", category),
					"overview":    fmt.Sprintf("This is mock data. Get real movie data by adding TMDB_API_KEY to env.local file.", category),
					"poster_path": "https://via.placeholder.com/300x450/ff6b6b/ffffff?text=STATIC+DATA",
					"backdrop_path": "https://via.placeholder.com/1920x1080/ff6b6b/ffffff?text=STATIC+DATA",
					"release_date": time.Now().Format("2006-01-02"),
					"vote_average": 7.5,
					"genre_ids":    []int{28, 12, 16},
					"is_static":    true,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorMsg)
		return
	}

	// Real OMDb API call - search for popular movies
	var url string
	switch category {
	case "popular":
		url = fmt.Sprintf("http://www.omdbapi.com/?s=action&type=movie&apikey=%s&page=1", apiKey)
	case "top_rated":
		url = fmt.Sprintf("http://www.omdbapi.com/?s=drama&type=movie&apikey=%s&page=1", apiKey)
	case "now_playing":
		url = fmt.Sprintf("http://www.omdbapi.com/?s=2024&type=movie&apikey=%s&page=1", apiKey)
	case "upcoming":
		url = fmt.Sprintf("http://www.omdbapi.com/?s=2025&type=movie&apikey=%s&page=1", apiKey)
	default:
		url = fmt.Sprintf("http://www.omdbapi.com/?s=movie&type=movie&apikey=%s&page=1", apiKey)
	}
	
	log.Printf("Fetching real movies from: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch movies: %v", err)
		http.Error(w, "Failed to fetch movies from API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("OMDb API returned status: %d", resp.StatusCode)
		http.Error(w, fmt.Sprintf("OMDb API error: %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	var omdbResp struct {
		Search []struct {
			Title  string `json:"Title"`
			Year   string `json:"Year"`
			IMDbID string `json:"imdbID"`
			Type   string `json:"Type"`
			Poster string `json:"Poster"`
		} `json:"Search"`
		TotalResults string `json:"totalResults"`
		Response     string `json:"Response"`
		Error        string `json:"Error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&omdbResp); err != nil {
		http.Error(w, "Failed to decode movies response", http.StatusInternalServerError)
		return
	}

	// Check if API key is invalid
	if omdbResp.Response == "False" && omdbResp.Error != "" {
		log.Printf("OMDb API Error: %s", omdbResp.Error)
		// Return static data with clear error message
		errorResponse := map[string]interface{}{
			"error": "INVALID API KEY",
			"message": fmt.Sprintf("OMDb API Error: %s", omdbResp.Error),
			"instructions": "Get a valid FREE API key from http://www.omdbapi.com/apikey.aspx",
			"category": category,
			"count": 1,
			"movies": []map[string]interface{}{
				{
					"id":          fmt.Sprintf("invalid_key_%d", time.Now().UnixNano()),
					"title":       "⚠️ INVALID API KEY - Get New Key",
					"overview":    fmt.Sprintf("OMDb API Error: %s. Get a valid API key to see real movies.", omdbResp.Error),
					"poster_path": "https://via.placeholder.com/300x450/ff6b6b/ffffff?text=INVALID+KEY",
					"backdrop_path": "https://via.placeholder.com/1920x1080/ff6b6b/ffffff?text=INVALID+KEY",
					"release_date": time.Now().Format("2006-01-02"),
					"vote_average": 0.0,
					"genre_ids":    []int{},
					"is_static":    true,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Transform to our format
	movies := make([]map[string]interface{}, 0, len(omdbResp.Search))
	for _, movie := range omdbResp.Search {
		posterURL := movie.Poster
		if posterURL == "N/A" {
			posterURL = "https://via.placeholder.com/300x450/cccccc/666666?text=No+Image"
		}

		movies = append(movies, map[string]interface{}{
			"id":           movie.IMDbID,
			"title":        movie.Title,
			"overview":     fmt.Sprintf("Movie from %s", movie.Year),
			"poster_path":  posterURL,
			"backdrop_path": posterURL, // Use same image for backdrop
			"release_date": movie.Year,
			"vote_average": 7.5, // Default rating since OMDb doesn't provide this in search
			"genre_ids":    []int{28, 12}, // Default genres
			"category":     category,
			"is_static":    false,
		})
	}

	result := map[string]interface{}{
		"category": category,
		"count":    len(movies),
		"movies":   movies,
		"source":   "REAL OMDb API DATA",
		"api_key_status": "VALID",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getTrendingMovies(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("OMDB_API_KEY")
	if apiKey == "" || apiKey == "your_omdb_api_key_here" {
		// Fallback to mock data
		movies := []map[string]interface{}{
			{
				"id":          fmt.Sprintf("trending_movie_%d", time.Now().UnixNano()),
				"title":       "⚠️ STATIC DATA: Trending Movie",
				"overview":    "This is mock data. Get real trending movies by adding TMDB_API_KEY to env.local file.",
				"poster_path": "https://via.placeholder.com/300x450/ff6b6b/ffffff?text=STATIC+DATA",
				"release_date": time.Now().Format("2006-01-02"),
				"vote_average": 8.5,
				"is_static":   true,
			},
		}
		result := map[string]interface{}{
			"count":  len(movies),
			"movies": movies,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	// Real TMDB trending API call
	url := fmt.Sprintf("https://api.themoviedb.org/3/trending/movie/week?api_key=%s", apiKey)
	
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch trending movies", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "TMDB API error", http.StatusInternalServerError)
		return
	}

	var tmdbResp struct {
		Results []struct {
			ID            int     `json:"id"`
			Title         string  `json:"title"`
			Overview      string  `json:"overview"`
			PosterPath    string  `json:"poster_path"`
			ReleaseDate   string  `json:"release_date"`
			VoteAverage   float64 `json:"vote_average"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tmdbResp); err != nil {
		http.Error(w, "Failed to decode trending movies response", http.StatusInternalServerError)
		return
	}

	// Transform to our format
	movies := make([]map[string]interface{}, 0, len(tmdbResp.Results))
	for _, movie := range tmdbResp.Results {
		posterURL := ""
		if movie.PosterPath != "" {
			posterURL = fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", movie.PosterPath)
		}

		movies = append(movies, map[string]interface{}{
			"id":           movie.ID,
			"title":        movie.Title,
			"overview":     movie.Overview,
			"poster_path":  posterURL,
			"release_date": movie.ReleaseDate,
			"vote_average": movie.VoteAverage,
		})
	}

	result := map[string]interface{}{
		"count":  len(movies),
		"movies": movies,
		"source": "REAL TMDB API DATA",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func searchMovies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" || apiKey == "your_tmdb_api_key_here" {
		// Fallback to mock data
		movies := []map[string]interface{}{
			{
				"id":          fmt.Sprintf("search_movie_%d", time.Now().UnixNano()),
				"title":       fmt.Sprintf("⚠️ STATIC DATA: Search Results for: %s", query),
				"overview":    fmt.Sprintf("This is mock data. Get real movie search results by adding TMDB_API_KEY to env.local file.", query),
				"poster_path": "https://via.placeholder.com/300x450/ff6b6b/ffffff?text=STATIC+DATA",
				"is_static":   true,
			},
		}
		result := map[string]interface{}{
			"query":  query,
			"count":  len(movies),
			"movies": movies,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	// Real TMDB search API call
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s&page=1", apiKey, query)
	
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to search movies", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "TMDB API error", http.StatusInternalServerError)
		return
	}

	var tmdbResp struct {
		Results []struct {
			ID            int     `json:"id"`
			Title         string  `json:"title"`
			Overview      string  `json:"overview"`
			PosterPath    string  `json:"poster_path"`
			ReleaseDate   string  `json:"release_date"`
			VoteAverage   float64 `json:"vote_average"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tmdbResp); err != nil {
		http.Error(w, "Failed to decode search response", http.StatusInternalServerError)
		return
	}

	// Transform to our format
	movies := make([]map[string]interface{}, 0, len(tmdbResp.Results))
	for _, movie := range tmdbResp.Results {
		posterURL := ""
		if movie.PosterPath != "" {
			posterURL = fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", movie.PosterPath)
		}

		movies = append(movies, map[string]interface{}{
			"id":           movie.ID,
			"title":        movie.Title,
			"overview":     movie.Overview,
			"poster_path":  posterURL,
			"release_date": movie.ReleaseDate,
			"vote_average": movie.VoteAverage,
		})
	}

	result := map[string]interface{}{
		"query":  query,
		"count":  len(movies),
		"movies": movies,
		"source": "REAL TMDB API DATA",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Helper function to get user preferences
func getUserPreferences(userID string) map[string]interface{} {
	// In a real app, this would call the user service
	// For now, return default preferences based on common interests
	return map[string]interface{}{
		"movie_genres": []string{"popular", "top_rated"},
	}
}
