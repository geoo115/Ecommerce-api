package utils

import (
	"os"
	"testing"
	"time"

	"github.com/geoo115/Ecommerce/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "test_secret_key_for_testing")
	os.Exit(m.Run())
}

func TestGenerateToken(t *testing.T) {
	user := models.User{
		Role: "user",
	}
	user.ID = 1 // Set ID after creation

	token, err := GenerateToken(user)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {
	user := models.User{
		Role: "user",
	}
	user.ID = 1 // Set ID after creation

	token, err := GenerateToken(user)
	require.NoError(t, err)

	claims, err := ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Role, claims.Role)
	assert.True(t, time.Unix(claims.ExpiresAt, 0).After(time.Now()))
}

func TestValidateToken_Invalid(t *testing.T) {
	_, err := ValidateToken("invalid_token")
	assert.Error(t, err)
}

func BenchmarkGenerateToken(b *testing.B) {
	user := models.User{
		Role: "user",
	}
	user.ID = 1 // Set ID after creation
	for i := 0; i < b.N; i++ {
		GenerateToken(user)
	}
}

func BenchmarkValidateToken(b *testing.B) {
	user := models.User{
		Role: "user",
	}
	user.ID = 1 // Set ID after creation
	token, _ := GenerateToken(user)
	for i := 0; i < b.N; i++ {
		ValidateToken(token)
	}
}
