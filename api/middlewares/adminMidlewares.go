package middlewares

import (
	"net/http"
	"strings"

	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if the user has admin privileges
		if claims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Admins only."})
			c.Abort()
			return
		}

		// Set the user ID to context (optional)
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
