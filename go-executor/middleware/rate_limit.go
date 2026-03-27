package middleware

import (
	"math/rand"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type userData struct {
	lastRequest time.Time
	failures    int
}

var (
	userLocks = make(map[string]*userData)
	mu        sync.Mutex
)

func RateLimitMiddleware() gin.HandlerFunc {
	const maxFailures = 5

	return func(c *gin.Context) {
		user := c.ClientIP() 

		mu.Lock()
		data, exists := userLocks[user]
		if !exists {
			data = &userData{}
			userLocks[user] = data
		}

		// exponential backoff up
		wait := 1 * time.Second
		if data.failures > 0 && data.failures <= maxFailures {
			wait = time.Duration(1<<data.failures) * time.Second 
		}

		// jitter
		wait += time.Duration(rand.Intn(500)) * time.Millisecond

		if time.Since(data.lastRequest) < wait {
			if data.failures < maxFailures {
				data.failures++
			}
			mu.Unlock()
			c.JSON(429, gin.H{
				"error":    "Rate limit exceeded",
				"retryIn":  wait.Seconds(),
				"failures": data.failures,
			})
			c.Abort()
			return
		}

		// Successful request
		data.lastRequest = time.Now()
		data.failures = 0
		mu.Unlock()

		c.Next()
	}
}