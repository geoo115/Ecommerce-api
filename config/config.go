package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	return nil
}

func GetDatabaseURL() (string, error) {
	// Highest precedence: full DATABASE_URL if provided
	if full := os.Getenv("DATABASE_URL"); full != "" {
		return full, nil
	}

	// Otherwise build from individual pieces, with sensible defaults
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	name := os.Getenv("DATABASE_NAME")
	sslmode := os.Getenv("DATABASE_SSLMODE")

	if port == "" {
		port = "5432"
	}
	if sslmode == "" {
		sslmode = "disable"
	}

	// Require core fields
	if user == "" || password == "" || host == "" || name == "" {
		return "", fmt.Errorf("one or more required environment variables are missing")
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, name, sslmode)
	return databaseURL, nil
}
