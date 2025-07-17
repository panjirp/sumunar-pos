package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

// Connect initializes a connection pool to PostgreSQL using DB_URL from env.
func Connect() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("❌ DB_URL is not set in environment variables")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	Pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL")
}
