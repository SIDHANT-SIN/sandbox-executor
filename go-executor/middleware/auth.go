package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	secret := os.Getenv("SECRET_KEY");

	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth != "Bearer "+secret {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "not_allowed_dude",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}