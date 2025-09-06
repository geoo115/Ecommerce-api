package services

import (
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	// This function is used by tests that don't have access to testing.T
	// In a real scenario, this should be replaced with proper test setup
	return nil
}
