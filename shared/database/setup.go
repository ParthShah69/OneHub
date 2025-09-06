N FOR package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func SetupDatabase() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://dashboard_user:dashboard_pass@localhost:5432/dashboard_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	log.Println("Database setup completed successfully")
	return db, nil
}

func createTables(db *sql.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			interests TEXT[],
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS user_behaviors (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			action VARCHAR(50) NOT NULL,
			content_id VARCHAR(255) NOT NULL,
			category VARCHAR(100) NOT NULL,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS user_profiles (
			user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
			explicit_interests TEXT[],
			behavioral_score JSONB,
			last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS news_articles (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title TEXT NOT NULL,
			description TEXT,
			url TEXT UNIQUE NOT NULL,
			source VARCHAR(255),
			category VARCHAR(100),
			published_at TIMESTAMP,
			image_url TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS job_listings (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title TEXT NOT NULL,
			company VARCHAR(255),
			location VARCHAR(255),
			description TEXT,
			url TEXT UNIQUE NOT NULL,
			category VARCHAR(100),
			posted_at TIMESTAMP,
			salary VARCHAR(100),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS videos (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title TEXT NOT NULL,
			description TEXT,
			url TEXT UNIQUE NOT NULL,
			channel VARCHAR(255),
			category VARCHAR(100),
			published_at TIMESTAMP,
			thumbnail TEXT,
			duration VARCHAR(20),
			views BIGINT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS deals (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title TEXT NOT NULL,
			description TEXT,
			url TEXT UNIQUE NOT NULL,
			platform VARCHAR(100),
			category VARCHAR(100),
			price DECIMAL(10,2),
			original_price DECIMAL(10,2),
			discount DECIMAL(5,2),
			image_url TEXT,
			valid_until TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS recommendations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			content_type VARCHAR(50) NOT NULL,
			content_id UUID NOT NULL,
			score DECIMAL(5,4),
			reason TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS nft_coupons (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			token_id VARCHAR(255),
			contract_address VARCHAR(255),
			title VARCHAR(255),
			description TEXT,
			discount DECIMAL(5,2),
			category VARCHAR(100),
			status VARCHAR(50) DEFAULT 'minted',
			minted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			claimed_at TIMESTAMP,
			expires_at TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS nft_activities (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			action VARCHAR(100) NOT NULL,
			points INTEGER DEFAULT 0,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}

	// Create indexes for better performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_user_behaviors_user_id ON user_behaviors(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_user_behaviors_category ON user_behaviors(category)",
		"CREATE INDEX IF NOT EXISTS idx_news_category ON news_articles(category)",
		"CREATE INDEX IF NOT EXISTS idx_jobs_category ON job_listings(category)",
		"CREATE INDEX IF NOT EXISTS idx_videos_category ON videos(category)",
		"CREATE INDEX IF NOT EXISTS idx_deals_category ON deals(category)",
		"CREATE INDEX IF NOT EXISTS idx_recommendations_user_id ON recommendations(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_nft_coupons_user_id ON nft_coupons(user_id)",
	}

	for _, index := range indexes {
		if _, err := db.Exec(index); err != nil {
			log.Printf("Warning: Failed to create index: %v", err)
		}
	}

	return nil
}
