package models

import (
	"time"
	"github.com/google/uuid"
)

type NewsArticle struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	URL         string    `json:"url" db:"url"`
	Source      string    `json:"source" db:"source"`
	Category    string    `json:"category" db:"category"`
	PublishedAt time.Time `json:"published_at" db:"published_at"`
	ImageURL    string    `json:"image_url" db:"image_url"`
}

type JobListing struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Company     string    `json:"company" db:"company"`
	Location    string    `json:"location" db:"location"`
	Description string    `json:"description" db:"description"`
	URL         string    `json:"url" db:"url"`
	Category    string    `json:"category" db:"category"`
	PostedAt    time.Time `json:"posted_at" db:"posted_at"`
	Salary      string    `json:"salary" db:"salary"`
}

type Video struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	URL         string    `json:"url" db:"url"`
	Channel     string    `json:"channel" db:"channel"`
	Category    string    `json:"category" db:"category"`
	PublishedAt time.Time `json:"published_at" db:"published_at"`
	Thumbnail   string    `json:"thumbnail" db:"thumbnail"`
	Duration    string    `json:"duration" db:"duration"`
	Views       int64     `json:"views" db:"views"`
}

type Deal struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	URL         string    `json:"url" db:"url"`
	Platform    string    `json:"platform" db:"platform"`
	Category    string    `json:"category" db:"category"`
	Price       float64   `json:"price" db:"price"`
	OriginalPrice float64 `json:"original_price" db:"original_price"`
	Discount    float64   `json:"discount" db:"discount"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	ValidUntil  time.Time `json:"valid_until" db:"valid_until"`
}

type Recommendation struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	ContentType string    `json:"content_type" db:"content_type"` // news, job, video, deal
	ContentID   uuid.UUID `json:"content_id" db:"content_id"`
	Score       float64   `json:"score" db:"score"`
	Reason      string    `json:"reason" db:"reason"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
