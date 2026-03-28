package main

import (
	"executor/middleware"
     	

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
)

func main() {
	
    _ = godotenv.Load()

	r := gin.Default()



	r.Use(middleware.AuthMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	r.POST("/execute", execHandler)

	r.Run(":8080")
}