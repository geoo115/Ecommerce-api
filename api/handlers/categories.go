package handlers

import (
	"net/http"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
)

func ListCategories(c *gin.Context) {
	var categories []models.Category
	if err := db.DB.Find(&categories).Error; err != nil {
		utils.SendInternalError(c, "Failed to fetch categories")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Categories retrieved successfully", categories)
}

func AddCategory(c *gin.Context) {
	var category models.Category
	// Only admin can add categories
	role, _ := c.Get("userRole")
	if role != "admin" {
		utils.SendForbidden(c, "")
		return
	}

	if err := c.ShouldBindJSON(&category); err != nil {
		// tests expect the raw JSON parsing error message
		utils.SendValidationError(c, err.Error())
		return
	}

	// Validate name
	if !utils.ValidateCategoryName(category.Name) {
		utils.SendValidationError(c, "Invalid category name")
		return
	}

	if err := db.DB.Create(&category).Error; err != nil {
		utils.SendInternalError(c, "Failed to create category")
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "Category added successfully", category)
}

func DeleteCategory(c *gin.Context) {
	// Only admin can delete categories
	role, _ := c.Get("userRole")
	if role != "admin" {
		utils.SendForbidden(c, "")
		return
	}

	// Validate id param
	idUint, err := Base.ValidateIDParam(c, "id")
	if err != nil {
		// response already sent with appropriate message
		return
	}

	tx := db.DB.Delete(&models.Category{}, idUint)
	if tx.Error != nil {
		utils.SendInternalError(c, "Failed to delete category")
		return
	}
	if tx.RowsAffected == 0 {
		utils.SendNotFound(c, "Category not found")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Category deleted successfully", nil)
}
