package middleware

import (
	"net/http"
	"os"
    "log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func AuthMiddleware() gin.HandlerFunc {
	if err := godotenv.Load("go-executor/.env"); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}
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