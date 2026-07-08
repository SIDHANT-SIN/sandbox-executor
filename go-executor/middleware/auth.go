package middleware

import (
	"github.com/gin-gonic/gin"
	
	"log"
	"net/http"
	"os"
)

func AuthMiddleware() gin.HandlerFunc {
	secret := os.Getenv("SECRET_KEY")
	
	if secret == "" {
		log.Println("WARNING: SECRET_KEY is empty! Check your .env file or docker-compose config.")
	}


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
