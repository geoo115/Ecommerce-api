package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Uptime    time.Duration          `json:"uptime"`
	Version   string                 `json:"version"`
	Services  map[string]interface{} `json:"services"`
}

// SystemInfo represents system information
type SystemInfo struct {
	GoVersion    string `json:"go_version"`
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
	NumCPU       int    `json:"num_cpu"`
	NumGoroutine int    `json:"num_goroutine"`
}

var startTime = time.Now()

// HealthCheck provides basic health status
func HealthCheck(c *gin.Context) {
	uptime := time.Since(startTime)

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    uptime,
		Version:   "1.0.0",
		Services: map[string]interface{}{
			"api": map[string]interface{}{
				"status": "healthy",
			},
		},
	}

	utils.SendSuccess(c, http.StatusOK, "Health check passed", response)
}

// DetailedHealthCheck provides detailed health status including database
func DetailedHealthCheck(c *gin.Context) {
	uptime := time.Since(startTime)

	// Check database connectivity
	dbStatus := "healthy"
	dbError := ""

	if err := db.DB.Raw("SELECT 1").Error; err != nil {
		dbStatus = "unhealthy"
		dbError = err.Error()
		utils.AppLogger.LogError(err, "Database health check")
	}

	// Get system information
	sysInfo := SystemInfo{
		GoVersion:    runtime.Version(),
		Architecture: runtime.GOARCH,
		OS:           runtime.GOOS,
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
	}

	// Determine overall status
	overallStatus := "healthy"
	if dbStatus == "unhealthy" {
		overallStatus = "degraded"
	}

	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Uptime:    uptime,
		Version:   "1.0.0",
		Services: map[string]interface{}{
			"api": map[string]interface{}{
				"status": "healthy",
			},
			"database": map[string]interface{}{
				"status": dbStatus,
				"error":  dbError,
			},
			"system": sysInfo,
		},
	}

	statusCode := http.StatusOK
	if overallStatus == "degraded" {
		statusCode = http.StatusServiceUnavailable
	}

	utils.SendSuccess(c, statusCode, "Detailed health check completed", response)
}

// ReadinessCheck checks if the application is ready to serve traffic
func ReadinessCheck(c *gin.Context) {
	// Check database connectivity
	if err := db.DB.Raw("SELECT 1").Error; err != nil {
		utils.SendError(c, http.StatusServiceUnavailable, "Database not ready")
		return
	}

	// Check if the application has been running for at least 5 seconds
	if time.Since(startTime) < 5*time.Second {
		utils.SendError(c, http.StatusServiceUnavailable, "Application still starting up")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Application is ready", gin.H{
		"status": "ready",
		"uptime": time.Since(startTime).String(),
	})
}

// LivenessCheck checks if the application is alive
func LivenessCheck(c *gin.Context) {
	utils.SendSuccess(c, http.StatusOK, "Application is alive", gin.H{
		"status": "alive",
		"uptime": time.Since(startTime).String(),
	})
}

// Metrics provides basic application metrics
func Metrics(c *gin.Context) {
	metrics := gin.H{
		"uptime": time.Since(startTime).String(),
		"system": SystemInfo{
			GoVersion:    runtime.Version(),
			Architecture: runtime.GOARCH,
			OS:           runtime.GOOS,
			NumCPU:       runtime.NumCPU(),
			NumGoroutine: runtime.NumGoroutine(),
		},
		"memory": gin.H{
			"alloc":       runtime.MemStats{}.Alloc,
			"total_alloc": runtime.MemStats{}.TotalAlloc,
			"sys":         runtime.MemStats{}.Sys,
			"num_gc":      runtime.MemStats{}.NumGC,
		},
	}

	utils.SendSuccess(c, http.StatusOK, "Metrics retrieved", metrics)
}
