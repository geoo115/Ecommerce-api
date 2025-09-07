package db

import (
	"context"
	"time"
)

// IsHealthy checks if the database connection is healthy
func IsHealthy() bool {
	if DB == nil {
		return false
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return false
	}

	// Quick ping with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx) == nil
}

// GetConnectionStatus returns a detailed status of the database connection
func GetConnectionStatus() map[string]interface{} {
	status := map[string]interface{}{
		"connected": false,
		"error":     nil,
	}

	if DB == nil {
		status["error"] = "Database connection not initialized"
		return status
	}

	sqlDB, err := DB.DB()
	if err != nil {
		status["error"] = err.Error()
		return status
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		status["error"] = err.Error()
		return status
	}

	// Get connection stats if available
	stats := sqlDB.Stats()
	status["connected"] = true
	status["open_connections"] = stats.OpenConnections
	status["idle_connections"] = stats.Idle
	status["in_use_connections"] = stats.InUse

	return status
}
