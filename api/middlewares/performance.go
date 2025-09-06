package middlewares

import (
	"sync"
	"time"

	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
)

// PerformanceMetrics holds performance data
type PerformanceMetrics struct {
	RequestCount  int64
	TotalLatency  time.Duration
	AvgLatency    time.Duration
	MinLatency    time.Duration
	MaxLatency    time.Duration
	ErrorCount    int64
	LastRequestAt time.Time
}

// Global performance metrics
var globalMetrics = &PerformanceMetrics{
	MinLatency: time.Hour, // Initialize to a large value
}
var metricsMu sync.RWMutex

// PerformanceTimingMiddleware measures request processing time with detailed tracking
func PerformanceTimingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Update metrics
		updatePerformanceMetrics(latency, c.Writer.Status() >= 400)
	}
}

// updatePerformanceMetrics updates global performance metrics
func updatePerformanceMetrics(latency time.Duration, isError bool) {
	metricsMu.Lock()
	defer metricsMu.Unlock()

	globalMetrics.RequestCount++
	globalMetrics.TotalLatency += latency
	globalMetrics.LastRequestAt = time.Now()

	if latency < globalMetrics.MinLatency {
		globalMetrics.MinLatency = latency
	}
	if latency > globalMetrics.MaxLatency {
		globalMetrics.MaxLatency = latency
	}

	if globalMetrics.RequestCount > 0 {
		globalMetrics.AvgLatency = globalMetrics.TotalLatency / time.Duration(globalMetrics.RequestCount)
	}

	if isError {
		globalMetrics.ErrorCount++
	}
}

// PerformanceLogger logs slow requests with detailed information
func PerformanceLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)

		// Log slow requests (>100ms)
		if latency > 100*time.Millisecond {
			utils.AppLogger.Warn("Slow request detected: %s %s - latency: %s, status: %d",
				c.Request.Method,
				c.Request.URL.Path,
				latency.String(),
				c.Writer.Status(),
			)
		}

		// Log very slow requests (>500ms)
		if latency > 500*time.Millisecond {
			utils.AppLogger.Error("Very slow request detected: %s %s - latency: %s, status: %d",
				c.Request.Method,
				c.Request.URL.Path,
				latency.String(),
				c.Writer.Status(),
			)
		}
	}
}

// GetPerformanceMetrics returns current performance metrics
func GetPerformanceMetrics() *PerformanceMetrics {
	metricsMu.RLock()
	defer metricsMu.RUnlock()

	// Return a copy to avoid races if the caller mutates fields
	copy := *globalMetrics
	return &copy
}

// ResetPerformanceMetrics resets all metrics (useful for testing)
func ResetPerformanceMetrics() {
	metricsMu.Lock()
	defer metricsMu.Unlock()

	globalMetrics = &PerformanceMetrics{
		RequestCount:  0,
		TotalLatency:  0,
		AvgLatency:    0,
		MinLatency:    time.Hour,
		MaxLatency:    0,
		ErrorCount:    0,
		LastRequestAt: time.Time{},
	}
}
