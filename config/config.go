package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func GetDatabaseURL() string {
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	name := os.Getenv("DATABASE_NAME")
	sslmode := os.Getenv("DATABASE_SSLMODE")

	if user == "" || password == "" || host == "" || port == "" || name == "" || sslmode == "" {
		panic("One or more required environment variables are missing")
	}

	// Debug print of the final URL
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, name, sslmode)
	fmt.Println("Database URL:", databaseURL) // Debug output

	return databaseURL
}
