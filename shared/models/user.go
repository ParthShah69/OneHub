package models

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Name         string    `json:"name" db:"name"`
	Interests    []string  `json:"interests" db:"interests"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserBehavior struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Action    string    `json:"action" db:"action"` // click, bookmark, share, search
	ContentID string    `json:"content_id" db:"content_id"`
	Category  string    `json:"category" db:"category"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

type UserProfile struct {
	UserID           uuid.UUID `json:"user_id" db:"user_id"`
	ExplicitInterests []string  `json:"explicit_interests" db:"explicit_interests"`
	BehavioralScore  map[string]float64 `json:"behavioral_score" db:"behavioral_score"`
	LastUpdated      time.Time `json:"last_updated" db:"last_updated"`
}
