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
	port := "8009"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "service": "food"})
	})

	// Get recipes by category
	http.HandleFunc("/api/food", getRecipes)
	
	// Get trending recipes
	http.HandleFunc("/api/food/trending", getTrendingRecipes)
	
	// Get recipes by query
	http.HandleFunc("/api/food/search", searchRecipes)

	log.Printf("Food service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getRecipes(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	userID := r.URL.Query().Get("user_id")
	
	if category == "" {
		category = "popular"
	}
	
	// If user_id is provided, get user preferences
	if userID != "" {
		preferences := getUserPreferences(userID)
		if len(preferences.FoodCategories) > 0 {
			category = preferences.FoodCategories[0] // Use first preference
		}
	}

	appID := os.Getenv("EDAMAM_APP_ID")
	appKey := os.Getenv("EDAMAM_APP_KEY")
	if appID == "" || appID == "your_edamam_app_id_here" || appKey == "" || appKey == "your_edamam_app_key_here" {
		// Return error message for static data
		errorMsg := map[string]interface{}{
			"error": "STATIC DATA - No Edamam API Keys Found",
			"message": "This is static/mock data. To get real recipes, add EDAMAM_APP_ID and EDAMAM_APP_KEY to your environment variables.",
			"instructions": "Get FREE API keys from https://developer.edamam.com/ and add to env.local",
			"category": category,
			"count": 1,
			"recipes": []map[string]interface{}{
				{
					"id":          fmt.Sprintf("static_recipe_%d", time.Now().UnixNano()),
					"title":       fmt.Sprintf("⚠️ STATIC DATA: Delicious %s Recipe", category),
					"summary":     fmt.Sprintf("This is mock data. Get real recipe data by adding RECIPE_API_KEY to env.local file.", category),
					"image":       "https://via.placeholder.com/300x200/ff6b6b/ffffff?text=STATIC+DATA",
					"readyInMinutes": 30,
					"servings":    4,
					"healthScore": 85,
					"is_static":   true,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorMsg)
		return
	}

	// Real Edamam API call
	var url string
	switch category {
	case "popular":
		url = fmt.Sprintf("https://api.edamam.com/search?q=chicken&app_id=%s&app_key=%s&from=0&to=20", appID, appKey)
	case "healthy":
		url = fmt.Sprintf("https://api.edamam.com/search?q=healthy&app_id=%s&app_key=%s&from=0&to=20", appID, appKey)
	case "vegetarian":
		url = fmt.Sprintf("https://api.edamam.com/search?q=vegetarian&app_id=%s&app_key=%s&from=0&to=20", appID, appKey)
	case "quick":
		url = fmt.Sprintf("https://api.edamam.com/search?q=quick&app_id=%s&app_key=%s&from=0&to=20", appID, appKey)
	default:
		url = fmt.Sprintf("https://api.edamam.com/search?q=recipe&app_id=%s&app_key=%s&from=0&to=20", appID, appKey)
	}
	
	log.Printf("Fetching real recipes from: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch recipes: %v", err)
		http.Error(w, "Failed to fetch recipes from API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Recipe API returned status: %d", resp.StatusCode)
		http.Error(w, fmt.Sprintf("Recipe API error: %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	var recipeResp struct {
		Results []struct {
			ID           int    `json:"id"`
			Title        string `json:"title"`
			Image        string `json:"image"`
			ReadyInMinutes int  `json:"readyInMinutes"`
			Servings     int    `json:"servings"`
			HealthScore  int    `json:"healthScore"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&recipeResp); err != nil {
		http.Error(w, "Failed to decode recipes response", http.StatusInternalServerError)
		return
	}

	// Get detailed information for each recipe
	recipes := make([]map[string]interface{}, 0, len(recipeResp.Results))
	for _, recipe := range recipeResp.Results {
		// Get recipe summary
		summaryURL := fmt.Sprintf("https://api.spoonacular.com/recipes/%d/summary?apiKey=%s", recipe.ID, apiKey)
		summaryResp, err := http.Get(summaryURL)
		var summary string
		if err == nil && summaryResp.StatusCode == http.StatusOK {
			var summaryData struct {
				Summary string `json:"summary"`
			}
			if json.NewDecoder(summaryResp.Body).Decode(&summaryData) == nil {
				summary = summaryData.Summary
			}
			summaryResp.Body.Close()
		}

		recipes = append(recipes, map[string]interface{}{
			"id":              recipe.ID,
			"title":           recipe.Title,
			"summary":         summary,
			"image":           recipe.Image,
			"readyInMinutes":  recipe.ReadyInMinutes,
			"servings":        recipe.Servings,
			"healthScore":     recipe.HealthScore,
			"category":        category,
		})
	}

	result := map[string]interface{}{
		"category": category,
		"count":    len(recipes),
		"recipes":  recipes,
		"source":   "REAL SPOONACULAR API DATA",
		"api_key_status": "VALID",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getTrendingRecipes(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("RECIPE_API_KEY")
	if apiKey == "" || apiKey == "your_recipe_api_key_here" {
		// Fallback to mock data
		recipes := []map[string]interface{}{
			{
				"id":          fmt.Sprintf("trending_recipe_%d", time.Now().UnixNano()),
				"title":       "⚠️ STATIC DATA: Trending Recipe",
				"summary":     "This is mock data. Get real trending recipes by adding RECIPE_API_KEY to env.local file.",
				"image":       "https://via.placeholder.com/300x200/ff6b6b/ffffff?text=STATIC+DATA",
				"readyInMinutes": 25,
				"servings":    6,
				"healthScore": 90,
				"is_static":   true,
			},
		}
		result := map[string]interface{}{
			"count":   len(recipes),
			"recipes": recipes,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	// Real Spoonacular trending API call
	url := fmt.Sprintf("https://api.spoonacular.com/recipes/complexSearch?apiKey=%s&number=20&sort=popularity", apiKey)
	
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch trending recipes", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Recipe API error", http.StatusInternalServerError)
		return
	}

	var recipeResp struct {
		Results []struct {
			ID           int    `json:"id"`
			Title        string `json:"title"`
			Image        string `json:"image"`
			ReadyInMinutes int  `json:"readyInMinutes"`
			Servings     int    `json:"servings"`
			HealthScore  int    `json:"healthScore"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&recipeResp); err != nil {
		http.Error(w, "Failed to decode trending recipes response", http.StatusInternalServerError)
		return
	}

	// Transform to our format
	recipes := make([]map[string]interface{}, 0, len(recipeResp.Results))
	for _, recipe := range recipeResp.Results {
		recipes = append(recipes, map[string]interface{}{
			"id":              recipe.ID,
			"title":           recipe.Title,
			"image":           recipe.Image,
			"readyInMinutes":  recipe.ReadyInMinutes,
			"servings":        recipe.Servings,
			"healthScore":     recipe.HealthScore,
		})
	}

	result := map[string]interface{}{
		"count":   len(recipes),
		"recipes": recipes,
		"source":  "REAL SPOONACULAR API DATA",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func searchRecipes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("RECIPE_API_KEY")
	if apiKey == "" || apiKey == "your_recipe_api_key_here" {
		// Fallback to mock data
		recipes := []map[string]interface{}{
			{
				"id":          fmt.Sprintf("search_recipe_%d", time.Now().UnixNano()),
				"title":       fmt.Sprintf("⚠️ STATIC DATA: Search Results for: %s", query),
				"summary":     fmt.Sprintf("This is mock data. Get real recipe search results by adding RECIPE_API_KEY to env.local file.", query),
				"image":       "https://via.placeholder.com/300x200/ff6b6b/ffffff?text=STATIC+DATA",
				"is_static":   true,
			},
		}
		result := map[string]interface{}{
			"query":   query,
			"count":   len(recipes),
			"recipes": recipes,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	// Real Spoonacular search API call
	url := fmt.Sprintf("https://api.spoonacular.com/recipes/complexSearch?apiKey=%s&query=%s&number=20", apiKey, query)
	
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to search recipes", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Recipe API error", http.StatusInternalServerError)
		return
	}

	var recipeResp struct {
		Results []struct {
			ID           int    `json:"id"`
			Title        string `json:"title"`
			Image        string `json:"image"`
			ReadyInMinutes int  `json:"readyInMinutes"`
			Servings     int    `json:"servings"`
			HealthScore  int    `json:"healthScore"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&recipeResp); err != nil {
		http.Error(w, "Failed to decode search response", http.StatusInternalServerError)
		return
	}

	// Transform to our format
	recipes := make([]map[string]interface{}, 0, len(recipeResp.Results))
	for _, recipe := range recipeResp.Results {
		recipes = append(recipes, map[string]interface{}{
			"id":              recipe.ID,
			"title":           recipe.Title,
			"image":           recipe.Image,
			"readyInMinutes":  recipe.ReadyInMinutes,
			"servings":        recipe.Servings,
			"healthScore":     recipe.HealthScore,
		})
	}

	result := map[string]interface{}{
		"query":   query,
		"count":   len(recipes),
		"recipes": recipes,
		"source":  "REAL SPOONACULAR API DATA",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Helper function to get user preferences
func getUserPreferences(userID string) map[string]interface{} {
	// In a real app, this would call the user service
	// For now, return default preferences based on common interests
	return map[string]interface{}{
		"food_categories": []string{"popular", "healthy"},
	}
}
