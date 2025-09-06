package utils

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetMetricsCollector(t *testing.T) {
	collector := GetMetricsCollector()
	assert.NotNil(t, collector)
	assert.NotNil(t, collector.httpRequests)
	assert.NotNil(t, collector.httpDurations)
	assert.NotNil(t, collector.httpErrors)
	assert.NotNil(t, collector.businessLogicCalls)
	assert.NotNil(t, collector.businessLogicDuration)
}

func TestRecordHTTPRequest(t *testing.T) {
	collector := GetMetricsCollector()
	collector.Reset()

	// Record a request
	RecordHTTPRequest("GET", "/api/test", 200, 100*time.Millisecond)

	// Check metrics
	key := "GET_/api/test"
	assert.Equal(t, int64(1), collector.httpRequests[key])
	assert.Equal(t, 100*time.Millisecond, collector.httpDurations[key])
	assert.Equal(t, int64(0), collector.httpErrors[key])

	// Record an error
	RecordHTTPRequest("POST", "/api/test", 500, 50*time.Millisecond)
	key2 := "POST_/api/test"
	assert.Equal(t, int64(1), collector.httpErrors[key2])
}

func TestRecordDatabaseQuery(t *testing.T) {
	collector := GetMetricsCollector()

	initialQueries := collector.databaseQueries
	initialDuration := collector.databaseDuration

	RecordDatabaseQuery(10 * time.Millisecond)

	assert.Equal(t, initialQueries+1, collector.databaseQueries)
	assert.Equal(t, initialDuration+10*time.Millisecond, collector.databaseDuration)
}

func TestRecordCacheOperation(t *testing.T) {
	collector := GetMetricsCollector()

	initialOps := collector.cacheOperations
	initialDuration := collector.cacheDuration

	RecordCacheOperation(5 * time.Millisecond)

	assert.Equal(t, initialOps+1, collector.cacheOperations)
	assert.Equal(t, initialDuration+5*time.Millisecond, collector.cacheDuration)
}

func TestRecordBusinessLogic(t *testing.T) {
	collector := GetMetricsCollector()

	RecordBusinessLogic("/api/users", 20*time.Millisecond)

	assert.Equal(t, int64(1), collector.businessLogicCalls["/api/users"])
	assert.Equal(t, 20*time.Millisecond, collector.businessLogicDuration["/api/users"])
}

func TestGetMetrics(t *testing.T) {
	// Record some metrics
	RecordHTTPRequest("GET", "/health", 200, 1*time.Millisecond)
	RecordDatabaseQuery(2 * time.Millisecond)

	metrics := GetMetrics()
	assert.NotEmpty(t, metrics)
	assert.True(t, strings.Contains(metrics, "http_requests_total"))
	assert.True(t, strings.Contains(metrics, "database_queries_total"))
}

func TestGetHTTPRequestCount(t *testing.T) {
	collector := GetMetricsCollector()
	collector.Reset()

	RecordHTTPRequest("GET", "/test", 200, 1*time.Millisecond)
	RecordHTTPRequest("GET", "/test", 200, 1*time.Millisecond)

	count := collector.GetHTTPRequestCount()
	assert.Equal(t, int64(2), count)
}

func TestGetDatabaseQueryCount(t *testing.T) {
	collector := GetMetricsCollector()
	collector.Reset()

	RecordDatabaseQuery(1 * time.Millisecond)
	RecordDatabaseQuery(1 * time.Millisecond)

	count := collector.GetDatabaseQueryCount()
	assert.Equal(t, int64(2), count)
}

func TestGetCacheOperationCount(t *testing.T) {
	collector := GetMetricsCollector()
	collector.Reset()

	RecordCacheOperation(1 * time.Millisecond)
	RecordCacheOperation(1 * time.Millisecond)

	count := collector.GetCacheOperationCount()
	assert.Equal(t, int64(2), count)
}

func TestGetHTTPErrorCount(t *testing.T) {
	ResetMetrics()
	collector := GetMetricsCollector()

	RecordHTTPRequest("GET", "/error", 404, 1*time.Millisecond)
	RecordHTTPRequest("GET", "/error", 500, 1*time.Millisecond)

	count := collector.GetHTTPErrorCount()
	assert.Equal(t, int64(2), count)
}

func TestGetUptime(t *testing.T) {
	collector := GetMetricsCollector()
	uptime := collector.GetUptime()
	assert.True(t, uptime > 0)
}

func TestResetMetrics(t *testing.T) {
	collector := GetMetricsCollector()

	RecordHTTPRequest("GET", "/reset", 200, 1*time.Millisecond)
	assert.Equal(t, int64(1), collector.httpRequests["GET_/reset"])

	ResetMetrics()

	assert.Equal(t, int64(0), collector.httpRequests["GET_/reset"])
	assert.Equal(t, int64(0), collector.databaseQueries)
}
