package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
}

func LoadConfig() Config {
	viper.AutomaticEnv()

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("DATABASE_URL", "")

	serverPort := viper.GetString("SERVER_PORT")
	databaseURL := viper.GetString("DATABASE_URL")

	if serverPort == "" || databaseURL == "" {
		log.Fatalf("Missing required configuration values")
	}

	return Config{
		ServerPort:  serverPort,
		DatabaseURL: databaseURL,
	}
}
