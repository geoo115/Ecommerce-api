package utils

import (
	"fmt"
	"sync"
	"time"
)

// MetricsCollector holds performance metrics
type MetricsCollector struct {
	mu                    sync.RWMutex
	httpRequests          map[string]int64
	httpDurations         map[string]time.Duration
	httpErrors            map[string]int64
	databaseQueries       int64
	databaseDuration      time.Duration
	cacheOperations       int64
	cacheDuration         time.Duration
	businessLogicCalls    map[string]int64
	businessLogicDuration map[string]time.Duration
	startTime             time.Time
}

var (
	metricsCollector *MetricsCollector
	metricsOnce      sync.Once
)

// GetMetricsCollector returns the singleton metrics collector
func GetMetricsCollector() *MetricsCollector {
	metricsOnce.Do(func() {
		metricsCollector = &MetricsCollector{
			httpRequests:          make(map[string]int64),
			httpDurations:         make(map[string]time.Duration),
			httpErrors:            make(map[string]int64),
			businessLogicCalls:    make(map[string]int64),
			businessLogicDuration: make(map[string]time.Duration),
			startTime:             time.Now(),
		}
	})
	return metricsCollector
}

// Reset resets all metrics (for testing)
func (c *MetricsCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.httpRequests = make(map[string]int64)
	c.httpDurations = make(map[string]time.Duration)
	c.httpErrors = make(map[string]int64)
	c.databaseQueries = 0
	c.databaseDuration = 0
	c.cacheOperations = 0
	c.cacheDuration = 0
	c.businessLogicCalls = make(map[string]int64)
	c.businessLogicDuration = make(map[string]time.Duration)
	c.startTime = time.Now()
}

// RecordHTTPRequest records HTTP request metrics
func RecordHTTPRequest(method, path string, statusCode int, duration time.Duration) {
	collector := GetMetricsCollector()
	collector.mu.Lock()
	defer collector.mu.Unlock()

	key := fmt.Sprintf("%s_%s", method, path)
	collector.httpRequests[key]++
	collector.httpDurations[key] += duration

	if statusCode >= 400 {
		collector.httpErrors[key]++
	}
}

// RecordDatabaseQuery records database query metrics
func RecordDatabaseQuery(duration time.Duration) {
	collector := GetMetricsCollector()
	collector.mu.Lock()
	defer collector.mu.Unlock()

	collector.databaseQueries++
	collector.databaseDuration += duration
}

// RecordCacheOperation records cache operation metrics
func RecordCacheOperation(duration time.Duration) {
	collector := GetMetricsCollector()
	collector.mu.Lock()
	defer collector.mu.Unlock()

	collector.cacheOperations++
	collector.cacheDuration += duration
}

// RecordBusinessLogic records business logic performance
func RecordBusinessLogic(endpoint string, duration time.Duration) {
	collector := GetMetricsCollector()
	collector.mu.Lock()
	defer collector.mu.Unlock()

	collector.businessLogicCalls[endpoint]++
	collector.businessLogicDuration[endpoint] += duration
}

// GetMetrics returns current metrics in Prometheus format
func GetMetrics() string {
	collector := GetMetricsCollector()
	collector.mu.RLock()
	defer collector.mu.RUnlock()

	var output string

	// HTTP metrics
	output += "# HELP http_requests_total Total number of HTTP requests\n"
	output += "# TYPE http_requests_total counter\n"
	for key, count := range collector.httpRequests {
		output += fmt.Sprintf("http_requests_total{method=\"%s\"} %d\n", key, count)
	}

	output += "# HELP http_request_duration_seconds Duration of HTTP requests\n"
	output += "# TYPE http_request_duration_seconds histogram\n"
	for key, duration := range collector.httpDurations {
		avgDuration := duration.Seconds() / float64(collector.httpRequests[key])
		output += fmt.Sprintf("http_request_duration_seconds{method=\"%s\"} %f\n", key, avgDuration)
	}

	output += "# HELP http_requests_errors_total Total number of HTTP errors\n"
	output += "# TYPE http_requests_errors_total counter\n"
	for key, count := range collector.httpErrors {
		output += fmt.Sprintf("http_requests_errors_total{method=\"%s\"} %d\n", key, count)
	}

	// Database metrics
	output += "# HELP database_queries_total Total number of database queries\n"
	output += "# TYPE database_queries_total counter\n"
	output += fmt.Sprintf("database_queries_total %d\n", collector.databaseQueries)

	if collector.databaseQueries > 0 {
		output += "# HELP database_query_duration_seconds Average database query duration\n"
		output += "# TYPE database_query_duration_seconds gauge\n"
		avgDBDuration := collector.databaseDuration.Seconds() / float64(collector.databaseQueries)
		output += fmt.Sprintf("database_query_duration_seconds %f\n", avgDBDuration)
	}

	// Cache metrics
	output += "# HELP cache_operations_total Total number of cache operations\n"
	output += "# TYPE cache_operations_total counter\n"
	output += fmt.Sprintf("cache_operations_total %d\n", collector.cacheOperations)

	if collector.cacheOperations > 0 {
		output += "# HELP cache_operation_duration_seconds Average cache operation duration\n"
		output += "# TYPE cache_operation_duration_seconds gauge\n"
		avgCacheDuration := collector.cacheDuration.Seconds() / float64(collector.cacheOperations)
		output += fmt.Sprintf("cache_operation_duration_seconds %f\n", avgCacheDuration)
	}

	// Business logic metrics
	output += "# HELP business_logic_calls_total Total number of business logic calls\n"
	output += "# TYPE business_logic_calls_total counter\n"
	for endpoint, count := range collector.businessLogicCalls {
		output += fmt.Sprintf("business_logic_calls_total{endpoint=\"%s\"} %d\n", endpoint, count)
	}

	output += "# HELP business_logic_duration_seconds Duration of business logic calls\n"
	output += "# TYPE business_logic_duration_seconds histogram\n"
	for endpoint, duration := range collector.businessLogicDuration {
		count := collector.businessLogicCalls[endpoint]
		if count > 0 {
			avgDuration := duration.Seconds() / float64(count)
			output += fmt.Sprintf("business_logic_duration_seconds{endpoint=\"%s\"} %f\n", endpoint, avgDuration)
		}
	}

	// Uptime metric
	output += "# HELP process_uptime_seconds Process uptime in seconds\n"
	output += "# TYPE process_uptime_seconds gauge\n"
	output += fmt.Sprintf("process_uptime_seconds %f\n", time.Since(collector.startTime).Seconds())

	return output
}

// ResetMetrics resets all metrics (useful for testing)
func ResetMetrics() {
	collector := GetMetricsCollector()
	collector.mu.Lock()
	defer collector.mu.Unlock()

	collector.httpRequests = make(map[string]int64)
	collector.httpDurations = make(map[string]time.Duration)
	collector.httpErrors = make(map[string]int64)
	collector.databaseQueries = 0
	collector.databaseDuration = 0
	collector.cacheOperations = 0
	collector.cacheDuration = 0
	collector.businessLogicCalls = make(map[string]int64)
	collector.businessLogicDuration = make(map[string]time.Duration)
	collector.startTime = time.Now()
}

// GetHTTPRequestCount returns the total number of HTTP requests
func (mc *MetricsCollector) GetHTTPRequestCount() int64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	var total int64
	for _, count := range mc.httpRequests {
		total += count
	}
	return total
}

// GetDatabaseQueryCount returns the total number of database queries
func (mc *MetricsCollector) GetDatabaseQueryCount() int64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.databaseQueries
}

// GetCacheOperationCount returns the total number of cache operations
func (mc *MetricsCollector) GetCacheOperationCount() int64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.cacheOperations
}

// GetHTTPErrorCount returns the total number of HTTP errors
func (mc *MetricsCollector) GetHTTPErrorCount() int64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	var total int64
	for _, count := range mc.httpErrors {
		total += count
	}
	return total
}

// GetUptime returns the application uptime
func (mc *MetricsCollector) GetUptime() time.Duration {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return time.Since(mc.startTime)
}
