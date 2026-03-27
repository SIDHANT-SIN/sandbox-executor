package main

import (
	"executor/middleware"
     	"log"
	

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
)

func main() {
	
		err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()



	r.Use(middleware.AuthMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	r.POST("/execute", execHandler)

	r.Run(":8080")
}