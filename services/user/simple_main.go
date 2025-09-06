package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Interests []string  `json:"interests"`
	CreatedAt time.Time `json:"created_at"`
}

type UserPreferences struct {
	NewsCategories    []string `json:"news_categories"`
	VideoCategories   []string `json:"video_categories"`
	JobCategories     []string `json:"job_categories"`
	DealCategories    []string `json:"deal_categories"`
	MovieGenres       []string `json:"movie_genres"`
	FoodCategories    []string `json:"food_categories"`
	PreferredSources  []string `json:"preferred_sources"`
}

// In-memory storage for demo
var users = make(map[string]User)
var userPreferences = make(map[string]UserPreferences)

func main() {
	port := "8006"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "service": "user"})
	})

	// User endpoints
	http.HandleFunc("/api/users", createUser)
	http.HandleFunc("/api/users/", getUser)
	http.HandleFunc("/api/users/preferences/", getUserPreferences)
	http.HandleFunc("/api/users/preferences/update/", updateUserPreferences)

	log.Printf("User service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userData struct {
		Name      string   `json:"name"`
		Email     string   `json:"email"`
		Interests []string `json:"interests"`
	}

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userID := fmt.Sprintf("user_%d", time.Now().UnixNano())
	user := User{
		ID:        userID,
		Name:      userData.Name,
		Email:     userData.Email,
		Interests: userData.Interests,
		CreatedAt: time.Now(),
	}

	users[userID] = user

	// Set default preferences based on interests
	preferences := UserPreferences{
		NewsCategories:   getDefaultNewsCategories(userData.Interests),
		VideoCategories:  getDefaultVideoCategories(userData.Interests),
		JobCategories:    getDefaultJobCategories(userData.Interests),
		DealCategories:   getDefaultDealCategories(userData.Interests),
		MovieGenres:      getDefaultMovieGenres(userData.Interests),
		FoodCategories:   getDefaultFoodCategories(userData.Interests),
		PreferredSources: getDefaultSources(userData.Interests),
	}

	userPreferences[userID] = preferences

	response := map[string]interface{}{
		"user":        user,
		"preferences": preferences,
		"message":     "User created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Path[len("/api/users/"):]
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	user, exists := users[userID]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func getUserPreferences(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Path[len("/api/users/preferences/"):]
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	preferences, exists := userPreferences[userID]
	if !exists {
		http.Error(w, "User preferences not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(preferences)
}

func updateUserPreferences(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Path[len("/api/users/preferences/update/"):]
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	var preferences UserPreferences
	if err := json.NewDecoder(r.Body).Decode(&preferences); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userPreferences[userID] = preferences

	response := map[string]interface{}{
		"preferences": preferences,
		"message":     "Preferences updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper functions to set default preferences based on interests
func getDefaultNewsCategories(interests []string) []string {
	categories := []string{"general"}
	
	for _, interest := range interests {
		switch interest {
		case "technology":
			categories = append(categories, "technology", "science")
		case "business":
			categories = append(categories, "business", "finance")
		case "entertainment":
			categories = append(categories, "entertainment")
		case "sports":
			categories = append(categories, "sports")
		case "health":
			categories = append(categories, "health")
		}
	}
	
	return removeDuplicates(categories)
}

func getDefaultVideoCategories(interests []string) []string {
	categories := []string{"technology"}
	
	for _, interest := range interests {
		switch interest {
		case "technology":
			categories = append(categories, "technology", "programming", "ai")
		case "business":
			categories = append(categories, "business", "entrepreneurship")
		case "entertainment":
			categories = append(categories, "entertainment", "comedy", "music")
		case "education":
			categories = append(categories, "education", "tutorial")
		case "gaming":
			categories = append(categories, "gaming")
		}
	}
	
	return removeDuplicates(categories)
}

func getDefaultJobCategories(interests []string) []string {
	categories := []string{"technology"}
	
	for _, interest := range interests {
		switch interest {
		case "technology":
			categories = append(categories, "software", "ai", "data-science")
		case "business":
			categories = append(categories, "business", "marketing", "finance")
		case "design":
			categories = append(categories, "design", "ui-ux")
		case "health":
			categories = append(categories, "healthcare", "medical")
		}
	}
	
	return removeDuplicates(categories)
}

func getDefaultDealCategories(interests []string) []string {
	categories := []string{"electronics"}
	
	for _, interest := range interests {
		switch interest {
		case "technology":
			categories = append(categories, "electronics", "computers", "gadgets")
		case "fashion":
			categories = append(categories, "fashion", "clothing")
		case "home":
			categories = append(categories, "home", "furniture")
		case "fitness":
			categories = append(categories, "fitness", "sports")
		}
	}
	
	return removeDuplicates(categories)
}

func getDefaultMovieGenres(interests []string) []string {
	genres := []string{"popular"}
	
	for _, interest := range interests {
		switch interest {
		case "entertainment":
			genres = append(genres, "popular", "top_rated", "now_playing")
		case "action":
			genres = append(genres, "action", "adventure")
		case "comedy":
			genres = append(genres, "comedy")
		case "drama":
			genres = append(genres, "drama")
		case "horror":
			genres = append(genres, "horror")
		case "romance":
			genres = append(genres, "romance")
		}
	}
	
	return removeDuplicates(genres)
}

func getDefaultFoodCategories(interests []string) []string {
	categories := []string{"popular"}
	
	for _, interest := range interests {
		switch interest {
		case "cooking":
			categories = append(categories, "popular", "healthy", "quick")
		case "healthy":
			categories = append(categories, "healthy", "vegetarian")
		case "baking":
			categories = append(categories, "dessert", "baking")
		case "international":
			categories = append(categories, "italian", "asian", "mexican")
		}
	}
	
	return removeDuplicates(categories)
}

func getDefaultSources(interests []string) []string {
	sources := []string{"general"}
	
	for _, interest := range interests {
		switch interest {
		case "technology":
			sources = append(sources, "techcrunch", "wired", "the-verge")
		case "business":
			sources = append(sources, "bloomberg", "reuters", "cnbc")
		case "entertainment":
			sources = append(sources, "entertainment-weekly", "variety")
		}
	}
	
	return removeDuplicates(sources)
}

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	result := []string{}
	
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	
	return result
}