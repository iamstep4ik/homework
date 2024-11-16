package main

import (
	"homework/config"
	"homework/pkg/db"
	"log"

	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.LoadConfig()
	database := db.Connect(&cfg)
	if err := goose.Up(database, cfg.GooseDir); err != nil {
		log.Fatalf("goose up: %v", err)
	}

	defer database.Close()

}
