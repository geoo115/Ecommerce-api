package utils

import (
	"regexp"
	"strings"
	"unicode"
)

// ValidateEmail checks if the email format is valid
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePassword checks if the password meets security requirements
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasDigit bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

// ValidateUsername checks if the username is valid
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 30 {
		return false
	}

	// Only allow alphanumeric characters and underscores
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

// ValidatePhone checks if the phone number is valid
func ValidatePhone(phone string) bool {
	if phone == "" {
		return false
	}

	// Remove all non-digit characters
	digits := regexp.MustCompile(`[^0-9]`).ReplaceAllString(phone, "")

	// Check if it's between 10-15 digits
	return len(digits) >= 10 && len(digits) <= 15
}

// ValidateProductName checks if the product name is valid
func ValidateProductName(name string) bool {
	name = strings.TrimSpace(name)
	return len(name) >= 2 && len(name) <= 200
}

// ValidatePrice checks if the price is valid
func ValidatePrice(price float64) bool {
	return price > 0 && price <= 999999.99
}

// ValidateQuantity checks if the quantity is valid
func ValidateQuantity(quantity int) bool {
	return quantity > 0 && quantity <= 1000
}

// ValidateCategoryName checks if the category name is valid
func ValidateCategoryName(name string) bool {
	name = strings.TrimSpace(name)
	return len(name) >= 2 && len(name) <= 50
}

// ValidateDescription checks if the description is valid
func ValidateDescription(description string) bool {
	description = strings.TrimSpace(description)
	return len(description) <= 1000
}

// ValidateStock checks if the stock quantity is valid
func ValidateStock(stock int) bool {
	return stock >= 0 && stock <= 100000
}

// ValidateRole checks if the role is valid
func ValidateRole(role string) bool {
	validRoles := []string{"customer", "admin"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// SanitizeString removes potentially dangerous characters
func SanitizeString(input string) string {
	// Remove null bytes and other control characters
	input = strings.Map(func(r rune) rune {
		if r < 32 && r != 9 && r != 10 && r != 13 {
			return -1
		}
		return r
	}, input)

	return strings.TrimSpace(input)
}
