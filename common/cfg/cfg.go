package cfg

import "os"

type Config struct {
	DBDSN string
}

func NewConfig() *Config {
	return &Config{
		DBDSN: os.Getenv("DB_DSN"),
	}
}
