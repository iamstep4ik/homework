package config

import "github.com/spf13/viper"

type Config struct {
	ServerPort    string
	DatabaseURL   string
	GooseDir      string
	GooseDBString string
}

func LoadConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return Config{
		ServerPort:    viper.GetString("server.port"),
		DatabaseURL:   viper.GetString("database.url"),
		GooseDir:      viper.GetString("goose.dir"),
		GooseDBString: viper.GetString("goose.dbstring"),
	}
}
