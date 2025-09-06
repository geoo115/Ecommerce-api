package middlewares

import (
	"net/http"
	"time"

	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		}, []string{"method", "path", "status"},
	)
	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request durations",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(httpDuration)
}

// PrometheusMetricsMiddleware collects Prometheus metrics for requests
func PrometheusMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		c.Next()

		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		httpRequests.WithLabelValues(method, path, http.StatusText(status)).Inc()
		httpDuration.WithLabelValues(method, path).Observe(duration)

		// Also record in in-app metrics collector
		utils.RecordHTTPRequest(method, path, status, time.Duration(duration*float64(time.Second)))
	}
}

// MetricsMiddleware is kept for backward compatibility and delegates to Prometheus
func MetricsMiddleware() gin.HandlerFunc {
	return PrometheusMetricsMiddleware()
}

// MetricsHandler returns the Prometheus metrics HTTP handler
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// DatabaseMetricsMiddleware tracks database query performance
func DatabaseMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		utils.RecordDatabaseQuery(duration)
	}
}

// CacheMetricsMiddleware tracks cache performance
func CacheMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		utils.RecordCacheOperation(duration)
	}
}
