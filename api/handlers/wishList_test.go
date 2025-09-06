package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/geoo115/Ecommerce/api/middlewares"
	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupWishlistTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	testDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := testDB.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{}, &models.Inventory{}, &models.Wishlist{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	db.DB = testDB
	return testDB
}

func TestAddToWishlist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupWishlistTestDB(t)

	os.Setenv("JWT_SECRET", "test_secret_key")
	utils.AppLogger = utils.NewLogger(utils.INFO)

	// create user
	user := models.User{Username: "wishlistuser", Email: "wishlist@example.com", Phone: "+15550000010"}
	db.DB.Create(&user)

	// create product
	cat := models.Category{Name: "wcat"}
	db.DB.Create(&cat)
	prod := models.Product{Name: "wprod", Price: 25.0, CategoryID: cat.ID, Description: "wishlist product"}
	db.DB.Create(&prod)

	// generate token
	token, _ := utils.GenerateToken(user)

	r := gin.New()
	r.POST("/wishlist", middlewares.AuthMiddleware(), AddToWishlist)

	ts := httptest.NewServer(r)
	defer ts.Close()

	// add to wishlist
	wishlistReq := map[string]interface{}{"product_id": prod.ID}
	b, _ := json.Marshal(wishlistReq)
	req, _ := http.NewRequest("POST", ts.URL+"/wishlist", io.NopCloser(bytes.NewReader(b)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("add to wishlist request failed: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestListWishlist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupWishlistTestDB(t)

	os.Setenv("JWT_SECRET", "test_secret_key")
	utils.AppLogger = utils.NewLogger(utils.INFO)

	// create user
	user := models.User{Username: "listwishlistuser", Email: "listwishlist@example.com", Phone: "+15550000011"}
	db.DB.Create(&user)

	// create product and add to wishlist
	cat := models.Category{Name: "lwcat"}
	db.DB.Create(&cat)
	prod := models.Product{Name: "lwprod", Price: 30.0, CategoryID: cat.ID, Description: "list wishlist product"}
	db.DB.Create(&prod)

	wishlist := models.Wishlist{UserID: user.ID, ProductID: prod.ID}
	db.DB.Create(&wishlist)

	// generate token
	token, _ := utils.GenerateToken(user)

	r := gin.New()
	r.GET("/wishlist", middlewares.AuthMiddleware(), ListWishlist)

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/wishlist", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("list wishlist request failed: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRemoveFromWishlist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupWishlistTestDB(t)

	os.Setenv("JWT_SECRET", "test_secret_key")
	utils.AppLogger = utils.NewLogger(utils.INFO)

	// create user
	user := models.User{Username: "removewishlistuser", Email: "removewishlist@example.com", Phone: "+15550000012"}
	db.DB.Create(&user)

	// create product and add to wishlist
	cat := models.Category{Name: "rwcat"}
	db.DB.Create(&cat)
	prod := models.Product{Name: "rwprod", Price: 35.0, CategoryID: cat.ID, Description: "remove wishlist product"}
	db.DB.Create(&prod)

	wishlist := models.Wishlist{UserID: user.ID, ProductID: prod.ID}
	db.DB.Create(&wishlist)

	// generate token
	token, _ := utils.GenerateToken(user)

	r := gin.New()
	r.DELETE("/wishlist/:id", middlewares.AuthMiddleware(), RemoveFromWishlist)

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, _ := http.NewRequest("DELETE", ts.URL+"/wishlist/"+strconv.Itoa(int(wishlist.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("remove from wishlist request failed: %v", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// verify item was removed
	var count int64
	db.DB.Model(&models.Wishlist{}).Where("id = ?", wishlist.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}
