package handlers

import (
	"net/http"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Signup handler for user registration
func Signup(c *gin.Context) {
	var user models.User

	// Validate input
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.SendValidationError(c, "Invalid input format")
		return
	}

	// Sanitize input
	user.Username = utils.SanitizeString(user.Username)
	user.Email = utils.SanitizeString(user.Email)
	user.Phone = utils.SanitizeString(user.Phone)

	// Validate required fields
	if user.Username == "" || user.Password == "" || user.Email == "" {
		utils.SendValidationError(c, "All fields are required (username, password, email)")
		return
	}

	// Validate input data
	if !utils.ValidateUsername(user.Username) {
		utils.SendValidationError(c, "Username must be 3-30 characters long and contain only alphanumeric characters and underscores")
		return
	}

	if !utils.ValidateEmail(user.Email) {
		utils.SendValidationError(c, "Invalid email format")
		return
	}

	if !utils.ValidatePassword(user.Password) {
		utils.SendValidationError(c, "Password must be at least 8 characters long and contain uppercase, lowercase, and numeric characters")
		return
	}

	if user.Phone != "" && !utils.ValidatePhone(user.Phone) {
		utils.SendValidationError(c, "Invalid phone number format")
		return
	}

	// Set default role if not provided
	if user.Role == "" {
		user.Role = "customer"
	} else if !utils.ValidateRole(user.Role) {
		utils.SendValidationError(c, "Invalid role specified")
		return
	}

	// Check if the username or email already exists
	if err := db.DB.Where("username = ? OR email = ?", user.Username, user.Email).First(&models.User{}).Error; err == nil {
		utils.SendConflict(c, "Username or email already in use")
		return
	} else if err != gorm.ErrRecordNotFound {
		utils.AppLogger.LogError(err, "Database error during signup")
		utils.SendInternalError(c, "Internal server error")
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.AppLogger.LogError(err, "Error hashing password")
		utils.SendInternalError(c, "Failed to process password")
		return
	}
	user.Password = hashedPassword

	// Save the user to the database
	if err := db.DB.Create(&user).Error; err != nil {
		utils.AppLogger.LogError(err, "Error saving user")
		utils.SendInternalError(c, "Failed to create user")
		return
	}

	utils.Info("New user registered: %s", user.Username)

	// Respond with the created user (excluding password)
	utils.SendSuccess(c, http.StatusCreated, "User created successfully", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var user models.User

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendValidationError(c, "Username and password are required")
		return
	}

	// Sanitize input
	input.Username = utils.SanitizeString(input.Username)

	// Validate username
	if !utils.ValidateUsername(input.Username) {
		utils.SendValidationError(c, "Invalid username format")
		return
	}

	// Check if user exists
	if err := db.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendUnauthorized(c, "Invalid username or password")
		} else {
			utils.AppLogger.LogError(err, "Database error during login")
			utils.SendInternalError(c, "Internal server error")
		}
		return
	}

	// Verify password
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		utils.AppLogger.LogSecurity("Failed login attempt", c.ClientIP(), "username", input.Username)
		utils.SendUnauthorized(c, "Invalid username or password")
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user)
	if err != nil {
		utils.AppLogger.LogError(err, "Error generating token")
		utils.SendInternalError(c, "Failed to generate token")
		return
	}

	utils.Info("User logged in: %s", user.Username)

	// Respond with token
	utils.SendSuccess(c, http.StatusOK, "Login successful", gin.H{"token": token})
}

// Logout handler
func Logout(c *gin.Context) {
	// Retrieve user from the context (set by middleware)
	user, exists := c.Get("user")
	if !exists {
		utils.SendUnauthorized(c, "User not authenticated")
		return
	}

	// Ideally, implement token invalidation here (e.g., store invalidated tokens in a blacklist)
	utils.Info("User logged out: %v", user)
	utils.SendSuccess(c, http.StatusOK, "User logged out successfully", gin.H{"user": user})
}
