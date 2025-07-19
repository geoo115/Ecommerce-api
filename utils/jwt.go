package utils

import (
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/geoo115/Ecommerce/models"
)

// Claims structure
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// JWT key management with lazy initialization
var (
	jwtKey     []byte
	jwtKeyOnce sync.Once
)

// getJWTKey retrieves JWT key from environment variable with lazy initialization
func getJWTKey() []byte {
	jwtKeyOnce.Do(func() {
		key := os.Getenv("JWT_SECRET")
		if key == "" {
			log.Fatal("JWT_SECRET environment variable is required")
		}
		jwtKey = []byte(key)
	})
	return jwtKey
}

// GenerateToken generates a JWT token for a given user
func GenerateToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getJWTKey())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates the provided JWT token string
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return getJWTKey(), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
