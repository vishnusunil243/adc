package cfg

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDSN        string
	AuthUsername string
	AuthPassword string
}

var config *Config

func LoadConfig() *Config {
	if config != nil {
		return config
	}
	err := godotenv.Load(".env")
	if err != nil {
		return nil
	}
	if config == nil {
		config = &Config{
			DBDSN:        os.Getenv("DB_DSN"),
			AuthUsername: os.Getenv("AUTH_USERNAME"),
			AuthPassword: os.Getenv("AUTH_PASSWORD"),
		}
	}
	return config
}
