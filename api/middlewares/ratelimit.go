package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter represents a simple in-memory rate limiter
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// isAllowed checks if a request is allowed based on rate limiting rules
func (rl *RateLimiter) isAllowed(key string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Get existing requests for this key
	requests, exists := rl.requests[key]
	if !exists {
		requests = []time.Time{}
	}

	// Filter out old requests outside the window
	var validRequests []time.Time
	for _, reqTime := range requests {
		if reqTime.After(windowStart) {
			validRequests = append(validRequests, reqTime)
		}
	}

	// Check if we're under the limit
	if len(validRequests) < rl.limit {
		validRequests = append(validRequests, now)
		rl.requests[key] = validRequests
		return true
	}

	// Update the requests list (even if rejected, to maintain the window)
	rl.requests[key] = validRequests
	return false
}

// cleanup removes old entries to prevent memory leaks
func (rl *RateLimiter) cleanup() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	for key, requests := range rl.requests {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if reqTime.After(windowStart) {
				validRequests = append(validRequests, reqTime)
			}
		}
		if len(validRequests) == 0 {
			delete(rl.requests, key)
		} else {
			rl.requests[key] = validRequests
		}
	}
}

// Global rate limiter instances
var (
	// General API rate limiter: 100 requests per minute
	generalLimiter = NewRateLimiter(100, time.Minute)
	// Auth endpoints rate limiter: 10 requests per minute
	authLimiter = NewRateLimiter(10, time.Minute)
	// Admin endpoints rate limiter: 50 requests per minute
	adminLimiter = NewRateLimiter(50, time.Minute)
)

// startCleanup starts a background goroutine to clean up old rate limiter entries
func startCleanup() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			generalLimiter.cleanup()
			authLimiter.cleanup()
			adminLimiter.cleanup()
		}
	}()
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP as the key for rate limiting
		clientIP := c.ClientIP()

		if !limiter.isAllowed(clientIP) {
			c.Header("X-RateLimit-Limit", "100")
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", time.Now().Add(time.Minute).Format(time.RFC3339))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "Rate limit exceeded. Please try again later.",
				"code":    429,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GeneralRateLimit applies general rate limiting to all endpoints
func GeneralRateLimit() gin.HandlerFunc {
	return RateLimitMiddleware(generalLimiter)
}

// AuthRateLimit applies stricter rate limiting to authentication endpoints
func AuthRateLimit() gin.HandlerFunc {
	return RateLimitMiddleware(authLimiter)
}

// AdminRateLimit applies rate limiting to admin endpoints
func AdminRateLimit() gin.HandlerFunc {
	return RateLimitMiddleware(adminLimiter)
}

// Initialize rate limiting cleanup
func init() {
	startCleanup()
}
