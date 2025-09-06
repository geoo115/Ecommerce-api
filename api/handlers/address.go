package handlers

import (
	"net/http"
	"strconv"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
)

func AddAddress(c *gin.Context) {
	var address models.Address
	if err := c.ShouldBindJSON(&address); err != nil {
		utils.SendValidationError(c, err.Error())
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.SendUnauthorized(c, "Unauthorized")
		return
	}
	uid, ok := userIDValue.(uint)
	if !ok {
		utils.SendUnauthorized(c, "Unauthorized")
		return
	}
	address.UserID = uid

	// Basic validation for required fields with specific messages to satisfy tests
	if address.Address == "" {
		utils.SendValidationError(c, "address is required")
		return
	}
	if address.City == "" {
		utils.SendValidationError(c, "city is required")
		return
	}
	if address.ZipCode == "" {
		utils.SendValidationError(c, "zip_code is required")
		return
	}

	if err := db.DB.Create(&address).Error; err != nil {
		utils.SendInternalError(c, "Failed to add address")
		return
	}

	// Reload the address with the associated User preloaded
	if err := db.DB.Preload("User").First(&address, address.ID).Error; err != nil {
		utils.SendInternalError(c, "Failed to load user data")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Address added successfully", address)
}

func EditAddress(c *gin.Context) {
	var address models.Address
	// Require auth and ownership
	uid, err := Base.GetUserID(c)
	if err != nil {
		return
	}

	// Validate id
	idStr := c.Param("id")
	id, convErr := strconv.Atoi(idStr)
	if convErr != nil {
		utils.SendValidationError(c, "Invalid ID parameter")
		return
	}

	if err := db.DB.Where("id = ? AND user_id = ?", id, uid).First(&address).Error; err != nil {
		utils.SendNotFound(c, "Address not found")
		return
	}

	if err := c.ShouldBindJSON(&address); err != nil {
		utils.SendValidationError(c, err.Error())
		return
	}

	if err := db.DB.Save(&address).Error; err != nil {
		utils.SendInternalError(c, "Failed to update address")
		return
	}
	// Reload the address with the associated User preloaded
	if err := db.DB.Preload("User").First(&address, address.ID).Error; err != nil {
		utils.SendInternalError(c, "Failed to load user data")
		return
	}
	utils.SendSuccess(c, http.StatusOK, "Address updated successfully", address)
}

func DeleteAddress(c *gin.Context) {
	var address models.Address
	// Require auth and ownership
	uid, err := Base.GetUserID(c)
	if err != nil {
		return
	}

	// Validate id
	idStr := c.Param("id")
	id, convErr := strconv.Atoi(idStr)
	if convErr != nil {
		utils.SendValidationError(c, "Invalid ID parameter")
		return
	}

	if err := db.DB.Where("id = ? AND user_id = ?", id, uid).First(&address).Error; err != nil {
		utils.SendNotFound(c, "Address not found")
		return
	}

	if err := db.DB.Delete(&address).Error; err != nil {
		utils.SendInternalError(c, "Failed to delete address")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Address deleted successfully", nil)
}
