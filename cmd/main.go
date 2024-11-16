package main

import (
	"homework/config"
	"homework/pkg/db"
)

func main() {
	cfg := config.LoadConfig()
	database := db.Connect(&cfg)

	defer database.Close()

}
