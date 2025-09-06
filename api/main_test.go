package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDBForAPITests(t *testing.T) {
	t.Helper()
	testDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := testDB.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{}, &models.Inventory{}, &models.Review{}, &models.Cart{}, &models.Order{}, &models.OrderItem{}, &models.Wishlist{}, &models.Address{}, &models.Payment{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	db.DB = testDB

	// Clear existing data
	tables := []string{"users", "categories", "products", "inventories", "reviews", "carts", "orders", "order_items", "wishlists", "addresses", "payments"}
	for _, table := range tables {
		testDB.Exec(fmt.Sprintf("DELETE FROM %s", table))
	}

	// Seed required data
	seedData := []interface{}{
		&models.User{Username: "testuser", Password: "password123", Email: "testuser@example.com", Phone: "+1234567890", Role: "customer"},
		&models.Category{Name: "Electronics"},
		&models.Product{Name: "Laptop", Description: "A high-end laptop", Price: 1500.00, CategoryID: 1},
		&models.Inventory{ProductID: 1, Stock: 10},
	}
	for _, data := range seedData {
		if err := testDB.Create(data).Error; err != nil {
			t.Fatalf("failed to seed data: %v", err)
		}
	}
}

func TestCategories_ListAddDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupDBForAPITests(t)
	db.DB.Exec("DELETE FROM categories")
	db.DB.Exec("DELETE FROM users")
	utils.AppLogger = utils.NewLogger(utils.INFO)
	os.Setenv("JWT_SECRET", "test_secret_key")

	// Create admin user
	admin := models.User{Username: "admin", Email: "admin@example.com", Phone: "+15550000000", Role: "admin"}
	db.DB.Create(&admin)
	token, _ := utils.GenerateToken(admin)

	router := gin.New()
	SetupRoutes(router)

	ts := httptest.NewServer(router)
	defer ts.Close()

	// add
	payload := map[string]interface{}{"name": "c1"}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", ts.URL+"/categories", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("add category failed: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// list
	resp2, err := http.Get(ts.URL + "/categories")
	if err != nil {
		t.Fatalf("list categories failed: %v", err)
	}
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	// decode list to find id (standardized response wrapper)
	var response struct {
		Success bool                     `json:"success"`
		Message string                   `json:"message"`
		Data    []map[string]interface{} `json:"data"`
		Code    int                      `json:"code"`
	}
	body, _ := io.ReadAll(resp2.Body)
	_ = json.Unmarshal(body, &response)
	utils.AppLogger.Info("List categories response: %+v", response)
	if len(response.Data) == 0 {
		t.Fatalf("expected categories")
	}
	category := response.Data[0]
	utils.AppLogger.Info("Category from response: %+v", category)
	// gorm JSON uses ID capitalized
	id := strconv.Itoa(int(category["ID"].(float64)))

	// delete
	req3, _ := http.NewRequest("DELETE", ts.URL+"/categories/"+id, nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	resp3, err := client.Do(req3)
	if err != nil {
		t.Fatalf("delete category failed: %v", err)
	}
	defer resp3.Body.Close()
	assert.Equal(t, http.StatusOK, resp3.StatusCode)
}
