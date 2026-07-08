package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"executor/middleware"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load .env file: %v", err)
	}

	r := gin.Default()

	r.Use(middleware.AuthMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	r.POST("/execute", execHandler)

	log.Println("Server starting on :8050")
	r.Run(":8050")
}
