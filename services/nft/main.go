package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"gofr.dev/pkg/gofr"
	"github.com/google/uuid"
)

type VerbwireResponse struct {
	Success bool `json:"success"`
	Data    struct {
		TokenID         string `json:"tokenId"`
		ContractAddress string `json:"contractAddress"`
		TransactionHash string `json:"transactionHash"`
	} `json:"data"`
}

type NFTService struct {
	verbwireAPIKey string
}

func main() {
	app := gofr.New()

	nftService := &NFTService{
		verbwireAPIKey: os.Getenv("VERBWIRE_API_KEY"),
	}

	// Health check
	app.GET("/health", func(ctx *gofr.Context) (interface{}, error) {
		return map[string]string{"status": "healthy", "service": "nft"}, nil
	})

	// NFT endpoints
	app.POST("/api/nft/mint", nftService.MintCoupon)
	app.GET("/api/nft/:user_id", nftService.GetUserNFTs)
	app.POST("/api/nft/:id/claim", nftService.ClaimNFT)

	app.Run()
}

func (ns *NFTService) MintCoupon(ctx *gofr.Context) (interface{}, error) {
	var mintData map[string]interface{}
	if err := ctx.Bind(&mintData); err != nil {
		return nil, fmt.Errorf("invalid request data: %v", err)
	}

	userID := mintData["user_id"].(string)
	category := mintData["category"].(string)
	discount := mintData["discount"].(float64)

	// Call Verbwire API to mint NFT
	nftData, err := ns.callVerbwireAPI(userID, category, discount)
	if err != nil {
		return nil, fmt.Errorf("failed to mint NFT: %v", err)
	}

	coupon := map[string]interface{}{
		"id":               uuid.New().String(),
		"user_id":          userID,
		"token_id":         nftData.Data.TokenID,
		"contract_address": nftData.Data.ContractAddress,
		"title":            fmt.Sprintf("%s Discount Coupon", category),
		"description":      fmt.Sprintf("Get %.0f%% off on %s items", discount, category),
		"discount":         discount,
		"category":         category,
		"status":           "minted",
		"minted_at":        time.Now(),
		"expires_at":       time.Now().Add(30 * 24 * time.Hour),
	}

	return coupon, nil
}

func (ns *NFTService) GetUserNFTs(ctx *gofr.Context) (interface{}, error) {
	userID := ctx.Param("user_id")

	// Mock NFT data
	nfts := []map[string]interface{}{
		{
			"id":               uuid.New().String(),
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
			"id":               uuid.New().String(),
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

	return map[string]interface{}{
		"user_id": userID,
		"count":   len(nfts),
		"nfts":    nfts,
	}, nil
}

func (ns *NFTService) ClaimNFT(ctx *gofr.Context) (interface{}, error) {
	nftID := ctx.Param("id")

	// Mock claim process
	claim := map[string]interface{}{
		"id":         nftID,
		"status":     "claimed",
		"claimed_at": time.Now(),
		"message":    "NFT coupon claimed successfully!",
	}

	return claim, nil
}

func (ns *NFTService) callVerbwireAPI(userID, category string, discount float64) (*VerbwireResponse, error) {
	// Mock Verbwire API call
	// In real implementation, this would make an HTTP request to Verbwire API
	
	mockResponse := &VerbwireResponse{
		Success: true,
		Data: struct {
			TokenID         string `json:"tokenId"`
			ContractAddress string `json:"contractAddress"`
			TransactionHash string `json:"transactionHash"`
		}{
			TokenID:         fmt.Sprintf("%d", time.Now().UnixNano()),
			ContractAddress: "0x1234567890abcdef",
			TransactionHash: fmt.Sprintf("0x%x", time.Now().UnixNano()),
		},
	}

	return mockResponse, nil
}
