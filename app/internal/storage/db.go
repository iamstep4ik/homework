package storage

import (
	"log"

	"homework/app/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) *sqlx.DB {
	db, err := sqlx.Open("postgres", cfg.DatabaseURL)

	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to the database successfully")

	return db
}
