package middlewares

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidateProduct_ValidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ValidateProduct())

	router.POST("/test", func(c *gin.Context) {
		input, exists := c.Get("product_input")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "input not found"})
			return
		}
		c.JSON(http.StatusOK, input)
	})

	validProduct := ProductInput{
		Name:        "Test Product",
		Price:       29.99,
		CategoryID:  1,
		Description: "Test description",
		Stock:       10,
	}

	jsonData, _ := json.Marshal(validProduct)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response ProductInput
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, validProduct.Name, response.Name)
	assert.Equal(t, validProduct.Price, response.Price)
}

func TestValidateProduct_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ValidateProduct())

	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Test missing required fields
	invalidProduct := map[string]interface{}{
		"name": "Test",
		// missing price, category_id, description, stock
	}

	jsonData, _ := json.Marshal(invalidProduct)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestValidateProduct_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ValidateProduct())

	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestValidateProduct_NegativePrice(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ValidateProduct())

	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	invalidProduct := ProductInput{
		Name:        "Test Product",
		Price:       -10.00, // negative price
		CategoryID:  1,
		Description: "Test description",
		Stock:       10,
	}

	jsonData, _ := json.Marshal(invalidProduct)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestValidateProduct_NegativeStock(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ValidateProduct())

	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	invalidProduct := ProductInput{
		Name:        "Test Product",
		Price:       29.99,
		CategoryID:  1,
		Description: "Test description",
		Stock:       -5, // negative stock
	}

	jsonData, _ := json.Marshal(invalidProduct)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

// Benchmark product validation middleware
func BenchmarkValidateProduct(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ValidateProduct())

	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	validProduct := ProductInput{
		Name:        "Test Product",
		Price:       29.99,
		CategoryID:  1,
		Description: "Test description",
		Stock:       10,
	}
	jsonData, _ := json.Marshal(validProduct)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
	}
}
