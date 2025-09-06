package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductInput struct {
	Name        string  `json:"name" binding:"required,min=3"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	CategoryID  uint    `json:"category_id" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
}

func ValidateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("ValidateProduct middleware called")

		var input ProductInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"details": "Validation failed for product input",
			})
			c.Abort()
			return
		}

		// Store the validated input in the context
		c.Set("product_input", input)
		c.Next()
	}
}
