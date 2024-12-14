package main

import (
	"homework/app/internal/config"
	"homework/app/internal/storage"
)

func main() {
	cfg := config.LoadConfig()
	database := storage.Connect(&cfg)
	defer database.Close()
}
