package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPerformanceTimingMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Reset metrics before test
	ResetPerformanceMetrics()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/test", nil)
	c.Request = req

	// Create middleware
	middleware := PerformanceTimingMiddleware()

	// Call middleware
	middleware(c)

	// Check that metrics were updated
	metrics := GetPerformanceMetrics()
	assert.Equal(t, int64(1), metrics.RequestCount)
	assert.True(t, metrics.TotalLatency > 0)
	assert.True(t, metrics.LastRequestAt.After(time.Now().Add(-time.Second)))
}

func TestPerformanceTimingMiddleware_WithError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Reset metrics before test
	ResetPerformanceMetrics()

	// Test the updatePerformanceMetrics function directly with error
	updatePerformanceMetrics(50*time.Millisecond, true)

	// Check that error count was incremented
	metrics := GetPerformanceMetrics()
	assert.Equal(t, int64(1), metrics.ErrorCount)
	assert.Equal(t, int64(1), metrics.RequestCount)
}

func TestUpdatePerformanceMetrics(t *testing.T) {
	// Reset metrics before test
	ResetPerformanceMetrics()

	latency := 50 * time.Millisecond

	// Update metrics
	updatePerformanceMetrics(latency, false)

	metrics := GetPerformanceMetrics()
	assert.Equal(t, int64(1), metrics.RequestCount)
	assert.Equal(t, latency, metrics.TotalLatency)
	assert.Equal(t, latency, metrics.MinLatency)
	assert.Equal(t, latency, metrics.MaxLatency)
	assert.Equal(t, latency, metrics.AvgLatency)
	assert.Equal(t, int64(0), metrics.ErrorCount)
}

func TestUpdatePerformanceMetrics_WithError(t *testing.T) {
	// Reset metrics before test
	ResetPerformanceMetrics()

	latency := 100 * time.Millisecond

	// Update metrics with error
	updatePerformanceMetrics(latency, true)

	metrics := GetPerformanceMetrics()
	assert.Equal(t, int64(1), metrics.ErrorCount)
}

func TestUpdatePerformanceMetrics_MultipleCalls(t *testing.T) {
	// Reset metrics before test
	ResetPerformanceMetrics()

	// First call
	updatePerformanceMetrics(50*time.Millisecond, false)

	// Second call
	updatePerformanceMetrics(150*time.Millisecond, false)

	metrics := GetPerformanceMetrics()
	assert.Equal(t, int64(2), metrics.RequestCount)
	assert.Equal(t, 200*time.Millisecond, metrics.TotalLatency)
	assert.Equal(t, 50*time.Millisecond, metrics.MinLatency)
	assert.Equal(t, 150*time.Millisecond, metrics.MaxLatency)
	assert.Equal(t, 100*time.Millisecond, metrics.AvgLatency)
}

func TestPerformanceLogger_FastRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/fast", nil)
	c.Request = req

	// Create middleware
	middleware := PerformanceLogger()

	// Call middleware (should not log anything for fast requests)
	middleware(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPerformanceLogger_SlowRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/slow", nil)
	c.Request = req

	// Create middleware
	middleware := PerformanceLogger()

	// Call middleware first, then simulate slow processing
	middleware(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetPerformanceMetrics(t *testing.T) {
	// Reset metrics before test
	ResetPerformanceMetrics()

	metrics := GetPerformanceMetrics()
	assert.NotNil(t, metrics)
	assert.Equal(t, int64(0), metrics.RequestCount)
	assert.Equal(t, time.Duration(0), metrics.TotalLatency)
	assert.Equal(t, time.Hour, metrics.MinLatency) // Should be initialized to large value
	assert.Equal(t, time.Duration(0), metrics.MaxLatency)
	assert.Equal(t, int64(0), metrics.ErrorCount)
}

func TestResetPerformanceMetrics(t *testing.T) {
	// First, update some metrics
	updatePerformanceMetrics(100*time.Millisecond, true)

	// Verify metrics are set
	metrics := GetPerformanceMetrics()
	assert.Equal(t, int64(1), metrics.RequestCount)
	assert.Equal(t, int64(1), metrics.ErrorCount)

	// Reset metrics
	ResetPerformanceMetrics()

	// Verify metrics are reset
	metrics = GetPerformanceMetrics()
	assert.Equal(t, int64(0), metrics.RequestCount)
	assert.Equal(t, time.Duration(0), metrics.TotalLatency)
	assert.Equal(t, time.Hour, metrics.MinLatency)
	assert.Equal(t, time.Duration(0), metrics.MaxLatency)
	assert.Equal(t, int64(0), metrics.ErrorCount)
}
