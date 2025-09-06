package middlewares

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Set JWT secret for tests
	os.Setenv("JWT_SECRET", "test_secret_key")
	utils.AppLogger = utils.NewLogger(utils.INFO)

	// Run tests
	code := m.Run()

	// Clean up
	os.Exit(code)
}

func TestAdminMiddleware_NoAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "GET",
		Header: http.Header{},
	}

	AdminMiddleware()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header is required")
}

func TestAdminMiddleware_EmptyToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "GET",
		Header: http.Header{
			"Authorization": []string{"Bearer "},
		},
	}

	AdminMiddleware()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Token is required")
}

func TestAdminMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "GET",
		Header: http.Header{
			"Authorization": []string{"Bearer invalid_token"},
		},
	}

	AdminMiddleware()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")
}

func TestAdminMiddleware_NonAdminUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set JWT secret for token generation
	os.Setenv("JWT_SECRET", "test_secret_key")

	// Create a regular user token
	user := models.User{
		Username: "regularuser",
		Email:    "regular@example.com",
		Role:     "customer",
	}
	token, _ := utils.GenerateToken(user)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "GET",
		Header: http.Header{
			"Authorization": []string{"Bearer " + token},
		},
	}

	AdminMiddleware()(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Access denied. Admins only")
}

func TestAdminMiddleware_AdminUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set JWT secret for token generation
	os.Setenv("JWT_SECRET", "test_secret_key")

	// Create an admin user token
	adminUser := models.User{
		Username: "adminuser",
		Email:    "admin@example.com",
		Role:     "admin",
	}
	token, _ := utils.GenerateToken(adminUser)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "GET",
		Header: http.Header{
			"Authorization": []string{"Bearer " + token},
		},
	}

	// Create a test handler that will be called if middleware passes
	testHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Admin access granted"})
	}

	// Set up router with middleware and test handler
	r := gin.New()
	r.Use(AdminMiddleware())
	r.GET("/admin-test", testHandler)

	// Create a test server
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Make request to test endpoint
	req, _ := http.NewRequest("GET", ts.URL+"/admin-test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read response body for debugging
	body := make([]byte, 1024)
	n, _ := resp.Body.Read(body)
	responseBody := string(body[:n])

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, responseBody, "Admin access granted")
}
