package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name+tag@example.co.uk", true},
		{"invalid-email", false},
		{"", false},
		{"test@", false},
		{"@example.com", false},
		{"test..test@example.com", false},
	}

	for _, test := range tests {
		result := ValidateEmail(test.email)
		assert.Equal(t, test.expected, result, "Email: %s", test.email)
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		expected bool
	}{
		{"Password123", true},
		{"SecurePass1", true},
		{"short", false},        // too short
		{"nouppercase1", false}, // no uppercase
		{"NOLOWERCASE1", false}, // no lowercase
		{"NoNumbers", false},    // no numbers
		{"Password", false},     // no numbers
		{"12345678", false},     // no letters
		{"", false},             // empty
	}

	for _, test := range tests {
		result := ValidatePassword(test.password)
		assert.Equal(t, test.expected, result, "Password: %s", test.password)
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		username string
		expected bool
	}{
		{"testuser", true},
		{"user123", true},
		{"user_name", true},
		{"u", false}, // too short
		{"thisusernameiswaytoolongtobevalid", false}, // too long
		{"user-name", true},                          // invalid character
		{"user.name", false},                         // invalid character
		{"", false},                                  // empty
	}

	for _, test := range tests {
		result := ValidateUsername(test.username)
		assert.Equal(t, test.expected, result, "Username: %s", test.username)
	}
}

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		phone    string
		expected bool
	}{
		{"+1234567890", true},            // 10 digits
		{"+123456789012345", false},      // 15 digits (too long, max is 14)
		{"1234567890", false},            // No +
		{"+123", false},                  // Too short
		{"+12345678901234567890", false}, // Too long
		{"", false},                      // Empty
		{"+12a34567890", false},          // Contains letter
	}

	for _, test := range tests {
		result := ValidatePhone(test.phone)
		assert.Equal(t, test.expected, result, "Phone: %s", test.phone)
	}
}

func TestValidateProductName(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"Valid Product", true},
		{"Product", true},
		{"P", false},   // too short
		{"", false},    // empty
		{"   ", false}, // only spaces
	}

	for _, test := range tests {
		result := ValidateProductName(test.name)
		assert.Equal(t, test.expected, result, "Product name: %s", test.name)
	}
}

func TestValidatePrice(t *testing.T) {
	tests := []struct {
		price    float64
		expected bool
	}{
		{10.99, true},
		{0.01, true},
		{999999.99, true},
		{0, false},          // zero
		{-1.00, false},      // negative
		{1000000.00, false}, // too high
	}

	for _, test := range tests {
		result := ValidatePrice(test.price)
		assert.Equal(t, test.expected, result, "Price: %f", test.price)
	}
}

func TestValidateQuantity(t *testing.T) {
	tests := []struct {
		quantity int
		expected bool
	}{
		{1, true},
		{100, true},
		{1000, true},
		{10000, true},  // maximum allowed
		{0, false},     // zero
		{-1, false},    // negative
		{10001, false}, // too high
	}

	for _, test := range tests {
		result := ValidateQuantity(test.quantity)
		assert.Equal(t, test.expected, result, "Quantity: %d", test.quantity)
	}
}

func TestValidateCategoryName(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"Electronics", true},
		{"Books", true},
		{"A", false}, // too short
		{"This category name is way too long to be valid", false}, // too long
		{"", false},    // empty
		{"   ", false}, // only spaces
	}

	for _, test := range tests {
		result := ValidateCategoryName(test.name)
		assert.Equal(t, test.expected, result, "Category name: %s", test.name)
	}
}

func TestValidateDescription(t *testing.T) {
	tests := []struct {
		description string
		expected    bool
	}{
		{"Valid description", true},
		{"", true},                          // empty is allowed
		{string(make([]byte, 1000)), true},  // exactly 1000 chars
		{string(make([]byte, 1001)), false}, // too long
	}

	for _, test := range tests {
		result := ValidateDescription(test.description)
		assert.Equal(t, test.expected, result, "Description length: %d", len(test.description))
	}
}

func TestValidateStock(t *testing.T) {
	tests := []struct {
		stock    int
		expected bool
	}{
		{0, true},
		{100, true},
		{100000, true},
		{-1, false},     // negative
		{100001, false}, // too high
	}

	for _, test := range tests {
		result := ValidateStock(test.stock)
		assert.Equal(t, test.expected, result, "Stock: %d", test.stock)
	}
}

func TestValidateRole(t *testing.T) {
	tests := []struct {
		role     string
		expected bool
	}{
		{"user", true},
		{"admin", true},
		{"customer", false}, // not a valid role
		{"", false},
		{"User", false}, // case sensitive
	}

	for _, test := range tests {
		result := ValidateRole(test.role)
		assert.Equal(t, test.expected, result, "Role: %s", test.role)
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal string", "normal string"},
		{"string with\nline break", "string with\nline break"},
		{"string with\t\ttabs", "string with\t\ttabs"},
		{"string\x00with\x01null", "stringwithnull"}, // null bytes removed
		{"  spaced string  ", "spaced string"},
	}

	for _, test := range tests {
		result := SanitizeString(test.input)
		assert.Equal(t, test.expected, result, "Input: %q", test.input)
	}
}

// Benchmark validation functions
func BenchmarkValidateEmail(b *testing.B) {
	email := "test@example.com"
	for i := 0; i < b.N; i++ {
		ValidateEmail(email)
	}
}

func BenchmarkValidatePassword(b *testing.B) {
	password := "Password123"
	for i := 0; i < b.N; i++ {
		ValidatePassword(password)
	}
}

// Additional comprehensive tests for full coverage
func TestValidateEmail_EdgeCases(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
		name     string
	}{
		{"a@b.co", true, "minimum valid email"},
		{"test+filter@example.com", true, "email with plus filter"},
		{"user.name@example.com", true, "email with dot in username"},
		{"test@sub.example.com", true, "email with subdomain"},
		{"test@example-site.com", true, "email with hyphen in domain"},
		{"test@192.168.1.1", false, "email with IP address (our regex doesn't allow this)"},
		{"test.email.with+symbol@example.com", true, "complex valid email"},
		{"plainaddress", false, "no @ symbol"},
		{"@missingdomain.com", false, "missing username"},
		{"missing-at-sign.net", false, "no @ symbol"},
		{"missing@.com", false, "missing domain name"},
		{"missing@domain", false, "missing TLD"},
		{"spaces in@email.com", false, "spaces in username"},
		{"test@spaces in.com", false, "spaces in domain"},
		{"test@.example.com", false, "dot at start of domain"},
		{"test@example..com", false, "double dot in domain"},
		{"test@", false, "missing domain entirely"},
		{".test@example.com", false, "dot at start of username"},
		{"test.@example.com", false, "dot at end of username"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ValidateEmail(test.email)
			assert.Equal(t, test.expected, result, "Email: %s", test.email)
		})
	}
}

func TestValidatePassword_EdgeCases(t *testing.T) {
	tests := []struct {
		password string
		expected bool
		name     string
	}{
		{"Aa1bcdef", true, "minimum valid password"},
		{"VeryLongPasswordWith123Numbers", true, "long valid password"},
		{"Complex123!@#", true, "password with special chars"},
		{"1234567", false, "7 characters (too short)"},
		{"12345678", false, "8 digits only"},
		{"abcdefgh", false, "8 lowercase only"},
		{"ABCDEFGH", false, "8 uppercase only"},
		{"Abcdefgh", false, "8 chars no number"},
		{"A1234567", false, "8 chars but only 1 letter"},
		{"a1234567", false, "8 chars but no uppercase"},
		{"A1bcdefghijklmnopqrstuvwxyz", true, "very long valid password"},
		{"", false, "empty password"},
		{" ", false, "single space"},
		{"        ", false, "8 spaces"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ValidatePassword(test.password)
			assert.Equal(t, test.expected, result, "Password: %s", test.password)
		})
	}
}

func TestValidatePhone_EdgeCases(t *testing.T) {
	tests := []struct {
		phone    string
		expected bool
		name     string
	}{
		{"+12345678901", true, "11 digit international"},
		{"+123456789012345", false, "15 digit international (too long, max is 14)"},
		{"+1234567890123456", false, "16 digits (too long)"},
		{"+123456789", false, "9 digits (too short)"},
		{"+", false, "just plus sign"},
		{"12345678901", false, "no plus sign"},
		{"+abc1234567890", false, "letters in number"},
		{"+123-456-7890", false, "dashes not allowed"},
		{"+123 456 7890", false, "spaces not allowed"},
		{"+1(234)567-8901", false, "parentheses not allowed"},
		{"", false, "empty phone"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ValidatePhone(test.phone)
			assert.Equal(t, test.expected, result, "Phone: %s", test.phone)
		})
	}
}

func TestValidatePrice_EdgeCases(t *testing.T) {
	tests := []struct {
		price    float64
		expected bool
		name     string
	}{
		{0.01, true, "minimum valid price"},
		{999999.99, true, "maximum valid price"},
		{1000000.00, false, "above maximum"},
		{0.0, false, "zero price"},
		{-0.01, false, "negative price"},
		{-100.0, false, "negative price"},
		{500.5, true, "decimal price"},
		{999999.999, false, "too many decimal places"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ValidatePrice(test.price)
			assert.Equal(t, test.expected, result, "Price: %f", test.price)
		})
	}
}

func TestValidateQuantity_EdgeCases(t *testing.T) {
	tests := []struct {
		quantity int
		expected bool
		name     string
	}{
		{1, true, "minimum valid quantity"},
		{10000, true, "maximum valid quantity"},
		{10001, false, "above maximum"},
		{0, false, "zero quantity"},
		{-1, false, "negative quantity"},
		{5000, true, "mid-range quantity"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ValidateQuantity(test.quantity)
			assert.Equal(t, test.expected, result, "Quantity: %d", test.quantity)
		})
	}
}

func TestValidateStock_EdgeCases(t *testing.T) {
	tests := []struct {
		stock    int
		expected bool
		name     string
	}{
		{0, true, "zero stock (valid)"},
		{1, true, "minimum positive stock"},
		{100000, true, "maximum valid stock"},
		{100001, false, "above maximum"},
		{-1, false, "negative stock"},
		{50000, true, "mid-range stock"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ValidateStock(test.stock)
			assert.Equal(t, test.expected, result, "Stock: %d", test.stock)
		})
	}
}

func TestValidateRole_EdgeCases(t *testing.T) {
	tests := []struct {
		role     string
		expected bool
		name     string
	}{
		{"user", true, "valid user role"},
		{"admin", true, "valid admin role"},
		{"User", false, "case sensitive - User"},
		{"Admin", false, "case sensitive - Admin"},
		{"USER", false, "case sensitive - USER"},
		{"ADMIN", false, "case sensitive - ADMIN"},
		{"customer", false, "customer not supported"},
		{"moderator", false, "invalid role"},
		{"", false, "empty role"},
		{" user", false, "role with leading space"},
		{"user ", false, "role with trailing space"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ValidateRole(test.role)
			assert.Equal(t, test.expected, result, "Role: %s", test.role)
		})
	}
}

func TestSanitizeString_EdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		name     string
	}{
		{"normal string", "normal string", "normal string unchanged"},
		{"string with\x00null", "string withnull", "null byte removed"},
		{"string\x01with\x02control", "stringwithcontrol", "control chars removed"},
		{"\x1fstring\x7f", "string", "more control chars removed"},
		{"", "", "empty string unchanged"},
		{"only normal chars!", "only normal chars!", "punctuation preserved"},
		{"tabs\tand\nnewlines", "tabs\tand\nnewlines", "tabs and newlines preserved"},
		{"\x00\x01\x02\x03", "", "all control chars removed"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := SanitizeString(test.input)
			assert.Equal(t, test.expected, result, "Input: %q", test.input)
		})
	}
}

func BenchmarkValidateUsername(b *testing.B) {
	username := "testuser"
	for i := 0; i < b.N; i++ {
		ValidateUsername(username)
	}
}

func BenchmarkSanitizeString(b *testing.B) {
	input := "test string with\x00null bytes"
	for i := 0; i < b.N; i++ {
		SanitizeString(input)
	}
}
