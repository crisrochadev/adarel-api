package middlewares

import (
	"sync"
	"time"

	"adarel-api/pkg/response"
	"github.com/gin-gonic/gin"
)

type visitor struct {
	count      int
	lastAccess time.Time
}

func RateLimitMiddleware(limit int, interval time.Duration) gin.HandlerFunc {
	visitors := make(map[string]*visitor)
	var mu sync.Mutex

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastAccess) > interval {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		v, exists := visitors[ip]
		if !exists || time.Since(v.lastAccess) > interval {
			visitors[ip] = &visitor{count: 1, lastAccess: time.Now()}
			mu.Unlock()
			c.Next()
			return
		}
		v.count++
		v.lastAccess = time.Now()
		count := v.count
		mu.Unlock()

		if count > limit {
			response.Error(c, 429, "too many requests")
			c.Abort()
			return
		}
		c.Next()
	}
}
