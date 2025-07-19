package middlewares

import (
	"time"

	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs HTTP requests with timing and details
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Log the request details using the logger
		utils.AppLogger.LogRequest(
			param.Method,
			param.Path,
			param.ClientIP,
			param.StatusCode,
			param.Latency,
		)

		// Return empty string since we're handling logging ourselves
		return ""
	})
}

// RequestTimingMiddleware adds timing information to requests
func RequestTimingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log request timing
		utils.Debug("Request completed: %s %s - Duration: %v",
			c.Request.Method, c.Request.URL.Path, duration)
	}
}

// ErrorLoggingMiddleware logs errors that occur during request processing
func ErrorLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were any errors
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				utils.AppLogger.LogError(err.Err, "Request processing")
			}
		}
	}
}
