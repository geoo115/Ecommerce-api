package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGzipMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(GzipMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(200, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	// Check if gzip encoding is applied
	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
}
