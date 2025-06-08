package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("Authorization")
		expectedAPIKey := os.Getenv("API_KEY")
		if expectedAPIKey == "" {
			expectedAPIKey = "BONGA_SERVER"
		}

		if apiKey == "" || apiKey != expectedAPIKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing API key"})
			c.Abort()
			return
		}

		c.Next()
	}
} 