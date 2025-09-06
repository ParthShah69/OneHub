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

type UserService struct {
	// Add database connection here
}

func main() {
	app := gofr.New()

	userService := &UserService{}

	// Health check
	app.GET("/health", func(ctx *gofr.Context) (interface{}, error) {
		return map[string]string{"status": "healthy", "service": "user"}, nil
	})

	// User endpoints
	app.POST("/api/users", userService.CreateUser)
	app.GET("/api/users/:id", userService.GetUser)
	app.PUT("/api/users/:id", userService.UpdateUser)
	app.POST("/api/users/:id/behavior", userService.TrackBehavior)

	app.Run()
}

func (us *UserService) CreateUser(ctx *gofr.Context) (interface{}, error) {
	var userData map[string]interface{}
	if err := ctx.Bind(&userData); err != nil {
		return nil, fmt.Errorf("invalid request data: %v", err)
	}

	user := map[string]interface{}{
		"id":        uuid.New().String(),
		"email":     userData["email"],
		"name":      userData["name"],
		"interests": userData["interests"],
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	return user, nil
}

func (us *UserService) GetUser(ctx *gofr.Context) (interface{}, error) {
	userID := ctx.Param("id")
	
	// Mock user data
	user := map[string]interface{}{
		"id":        userID,
		"email":     "user@example.com",
		"name":      "Demo User",
		"interests": []string{"technology", "ai", "business"},
		"created_at": time.Now().Add(-30 * 24 * time.Hour),
		"updated_at": time.Now(),
	}

	return user, nil
}

func (us *UserService) UpdateUser(ctx *gofr.Context) (interface{}, error) {
	userID := ctx.Param("id")
	
	var updateData map[string]interface{}
	if err := ctx.Bind(&updateData); err != nil {
		return nil, fmt.Errorf("invalid request data: %v", err)
	}

	// Mock update
	user := map[string]interface{}{
		"id":        userID,
		"email":     updateData["email"],
		"name":      updateData["name"],
		"interests": updateData["interests"],
		"updated_at": time.Now(),
	}

	return user, nil
}

func (us *UserService) TrackBehavior(ctx *gofr.Context) (interface{}, error) {
	userID := ctx.Param("id")
	
	var behaviorData map[string]interface{}
	if err := ctx.Bind(&behaviorData); err != nil {
		return nil, fmt.Errorf("invalid request data: %v", err)
	}

	behavior := map[string]interface{}{
		"id":         uuid.New().String(),
		"user_id":    userID,
		"action":     behaviorData["action"],
		"content_id": behaviorData["content_id"],
		"category":   behaviorData["category"],
		"timestamp":  time.Now(),
	}

	return behavior, nil
}
