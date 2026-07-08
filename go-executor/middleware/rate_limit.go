package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type userData struct {
	tokens     float64
	lastRefill time.Time
}

var (
	users = make(map[string]*userData)
	mu    sync.Mutex
)

func RateLimitMiddleware() gin.HandlerFunc {
	const (
		capacity   = 5.0 // max burst (5 instant requests)
		refillRate = 2.0 // tokens per second
	)

	return func(c *gin.Context) {
		user := c.ClientIP() // ✅ using IP only (as you wanted)

		mu.Lock()
		data, ok := users[user]
		if !ok {
			data = &userData{
				tokens:     capacity,
				lastRefill: time.Now(),
			}
			users[user] = data
		}

		now := time.Now()
		elapsed := now.Sub(data.lastRefill).Seconds()

		// refill tokens
		data.tokens += elapsed * refillRate
		if data.tokens > capacity {
			data.tokens = capacity
		}
		data.lastRefill = now

		// check limit
		if data.tokens < 1 {
			mu.Unlock()
			c.JSON(429, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}

		// consume token
		data.tokens -= 1
		mu.Unlock()

		c.Next()
	}
}
