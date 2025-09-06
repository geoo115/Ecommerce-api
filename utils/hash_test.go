package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword"

	hash, err := HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword"

	hash, err := HashPassword(password)
	require.NoError(t, err)

	assert.True(t, CheckPasswordHash(password, hash))
	assert.False(t, CheckPasswordHash("wrongpassword", hash))
}

func BenchmarkHashPassword(b *testing.B) {
	password := "testpassword"
	for i := 0; i < b.N; i++ {
		HashPassword(password)
	}
}

func BenchmarkCheckPasswordHash(b *testing.B) {
	password := "testpassword"
	hash, _ := HashPassword(password)
	for i := 0; i < b.N; i++ {
		CheckPasswordHash(password, hash)
	}
}
