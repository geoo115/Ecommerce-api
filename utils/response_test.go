package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSendSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := gin.H{"key": "value"}
	SendSuccess(c, http.StatusOK, "Success message", testData)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Success message", resp.Message)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestSendError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendError(c, http.StatusBadRequest, "Error message")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Error message", resp.Error)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestSendValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendValidationError(c, "Validation failed")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Validation failed", resp.Error)
}

func TestSendUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendUnauthorized(c, "Custom unauthorized message")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Custom unauthorized message", resp.Error)
}

func TestSendUnauthorizedDefaultMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendUnauthorized(c, "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Unauthorized access", resp.Error)
}

func TestSendNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendNotFound(c, "Custom not found message")

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Custom not found message", resp.Error)
}

func TestSendNotFoundDefaultMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendNotFound(c, "")

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Resource not found", resp.Error)
}

func TestSendInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendInternalError(c, "Custom internal error")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Custom internal error", resp.Error)
}

func TestSendInternalErrorDefaultMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendInternalError(c, "")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Internal server error", resp.Error)
}

func TestSendConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendConflict(c, "Resource already exists")

	assert.Equal(t, http.StatusConflict, w.Code)
	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Resource already exists", resp.Error)
}

func TestSendForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendForbidden(c, "Access denied")

	assert.Equal(t, http.StatusForbidden, w.Code)
	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Access denied", resp.Error)
	assert.Equal(t, http.StatusForbidden, resp.Code)
}

func TestSendForbidden_DefaultMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendForbidden(c, "")

	assert.Equal(t, http.StatusForbidden, w.Code)
	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Access forbidden", resp.Error)
	assert.Equal(t, http.StatusForbidden, resp.Code)
}

// Benchmark response functions for performance
func BenchmarkSendSuccess(b *testing.B) {
	gin.SetMode(gin.TestMode)
	testData := gin.H{"key": "value"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		SendSuccess(c, http.StatusOK, "Success message", testData)
	}
}

func BenchmarkSendError(b *testing.B) {
	gin.SetMode(gin.TestMode)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		SendError(c, http.StatusBadRequest, "Error message")
	}
}
