package db

import (
	"testing"

	"github.com/geoo115/Ecommerce/models"
	"github.com/stretchr/testify/assert"
)

func TestConnectDatabase_InMemory(t *testing.T) {
	db := SetupTestDB(t)

	user := models.User{Username: "test"}
	err := db.Create(&user).Error
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestOptimizedDB_CreateOptimizedIndexes(t *testing.T) {
	db := SetupTestDB(t)
	optimizedDB := &OptimizedDB{DB: db}

	err := optimizedDB.CreateOptimizedIndexes()
	assert.NoError(t, err)
}

func TestGetDBConfig(t *testing.T) {
	config := GetDBConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, 5432, config.Port)
	assert.Equal(t, "postgres", config.User)
	assert.Equal(t, "password", config.Password)
	assert.Equal(t, "ecommerce", config.DBName)
	assert.Equal(t, "disable", config.SSLMode)
	assert.True(t, config.MaxOpenConns > 0)
	assert.True(t, config.MaxIdleConns > 0)
	assert.True(t, config.MaxLifetime > 0)
}

func TestDatabase_Models_AutoMigration(t *testing.T) {
	db := SetupTestDB(t)

	// Verify tables were created
	assert.True(t, db.Migrator().HasTable(&models.User{}))
	assert.True(t, db.Migrator().HasTable(&models.Product{}))
	assert.True(t, db.Migrator().HasTable(&models.Category{}))
}
