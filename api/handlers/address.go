package handlers

import (
	"ecommerce/db"
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddAddress(c *gin.Context) {
	var address models.Address
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, address)
}

func EditAddress(c *gin.Context) {
	var address models.Address
	id := c.Param("id")

	if err := db.DB.First(&address, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Save(&address)
	c.JSON(http.StatusOK, address)
}

func DeleteAddress(c *gin.Context) {
	var address models.Address
	id := c.Param("id")

	if err := db.DB.Delete(&address, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
}
