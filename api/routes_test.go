package api

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRoutes_RegistersExpectedRoutes(t *testing.T) {
	// Use a fresh gin engine
	r := gin.New()
	SetupRoutes(r)

	routes := r.Routes()
	seen := map[string]bool{}
	for _, rt := range routes {
		seen[rt.Method+" "+rt.Path] = true
	}

	// Check a handful of important routes are registered
	assert.True(t, seen["GET /health"], "expected GET /health to be registered")
	assert.True(t, seen["GET /metrics"], "expected GET /metrics to be registered")
	assert.True(t, seen["GET /products"], "expected GET /products to be registered")
	assert.True(t, seen["POST /signup"], "expected POST /signup to be registered")
	assert.True(t, seen["POST /checkout"], "expected POST /checkout to be registered")
}
