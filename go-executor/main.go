// package main

// import (
// 	"executor/middleware"
//     "log"
//       "os"
// 	"github.com/joho/godotenv"
// 	"github.com/gin-gonic/gin"
// )

// func main() {

//   if err := godotenv.Load(); err != nil {
// 	log.Printf("Failed to load .env file: %v", err)
// }

// 	r := gin.Default()

// 	r.Use(middleware.AuthMiddleware())
// 	r.Use(middleware.RateLimitMiddleware())

// 	problemID := "2_sum"

// testData, err := readBlob(problemID, os.Getenv("TEST_FILE"))
// if err != nil {
// 	log.Fatalf("Failed to read test cases: %v", err)
// }

// solutionData, err := readBlob(problemID, os.Getenv("SOL_FILE"))
// if err != nil {
// 	log.Fatalf("Failed to read solution file: %v", err)
// }

// log.Println("===== TEST CASES =====")
// log.Println(testData)

// log.Println("===== SOLUTION =====")
// log.Println(solutionData)

// 	r.POST("/execute", execHandler)

// 	r.Run(":8050")
// }

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
