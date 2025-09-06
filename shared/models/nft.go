package models

import (
	"time"
	"github.com/google/uuid"
)

type NFTCoupon struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	TokenID     string    `json:"token_id" db:"token_id"`
	ContractAddress string `json:"contract_address" db:"contract_address"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Discount    float64   `json:"discount" db:"discount"`
	Category    string    `json:"category" db:"category"`
	Status      string    `json:"status" db:"status"` // minted, claimed, expired
	MintedAt    time.Time `json:"minted_at" db:"minted_at"`
	ClaimedAt   *time.Time `json:"claimed_at" db:"claimed_at"`
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`
}

type NFTActivity struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Action    string    `json:"action" db:"action"` // engagement, milestone, reward
	Points    int       `json:"points" db:"points"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}
