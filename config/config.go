package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func GetDatabaseURL() string {
	return os.Getenv("DATABASE_URL")
}
