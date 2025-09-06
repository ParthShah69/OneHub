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
	port := "8007"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "service": "nft"})
	})

	http.HandleFunc("/api/nft/mint", mintCoupon)
	http.HandleFunc("/api/nft/", getUserNFTs)
	http.HandleFunc("/api/nft/claim", claimNFT)

	log.Printf("NFT service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func mintCoupon(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var mintData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&mintData); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	userID := mintData["user_id"].(string)
	category := mintData["category"].(string)
	discount := mintData["discount"].(float64)

	coupon := map[string]interface{}{
		"id":               fmt.Sprintf("nft_%d", time.Now().UnixNano()),
		"user_id":          userID,
		"token_id":         fmt.Sprintf("token_%d", time.Now().UnixNano()),
		"contract_address": "0x1234567890abcdef",
		"title":            fmt.Sprintf("%s Discount Coupon", category),
		"description":      fmt.Sprintf("Get %.0f%% off on %s items", discount, category),
		"discount":         discount,
		"category":         category,
		"status":           "minted",
		"minted_at":        time.Now(),
		"expires_at":       time.Now().Add(30 * 24 * time.Hour),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coupon)
}

func getUserNFTs(w http.ResponseWriter, r *http.Request) {
	userID := "demo_user"

	nfts := []map[string]interface{}{
		{
			"id":               fmt.Sprintf("nft_%d", time.Now().UnixNano()),
			"user_id":          userID,
			"token_id":         "12345",
			"contract_address": "0x1234567890abcdef",
			"title":            "Technology Discount Coupon",
			"description":      "Get 15% off on technology items",
			"discount":         15.0,
			"category":         "technology",
			"status":           "minted",
			"minted_at":        time.Now().Add(-5 * 24 * time.Hour),
			"expires_at":       time.Now().Add(25 * 24 * time.Hour),
		},
		{
			"id":               fmt.Sprintf("nft_%d", time.Now().UnixNano()+1),
			"user_id":          userID,
			"token_id":         "12346",
			"contract_address": "0x1234567890abcdef",
			"title":            "Fashion Discount Coupon",
			"description":      "Get 20% off on fashion items",
			"discount":         20.0,
			"category":         "fashion",
			"status":           "claimed",
			"minted_at":        time.Now().Add(-10 * 24 * time.Hour),
			"claimed_at":       time.Now().Add(-2 * 24 * time.Hour),
			"expires_at":       time.Now().Add(20 * 24 * time.Hour),
		},
	}

	result := map[string]interface{}{
		"user_id": userID,
		"count":   len(nfts),
		"nfts":    nfts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func claimNFT(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var claimData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&claimData); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	nftID := claimData["id"].(string)

	claim := map[string]interface{}{
		"id":         nftID,
		"status":     "claimed",
		"claimed_at": time.Now(),
		"message":    "NFT coupon claimed successfully!",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(claim)
}
