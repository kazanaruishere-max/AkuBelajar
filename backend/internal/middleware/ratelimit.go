package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiter creates a Redis-backed sliding window rate limiter.
// maxRequests: max allowed in the window
// window: time window duration
// prefix: Redis key prefix (e.g., "login", "api")
func RateLimiter(rdb *redis.Client, maxRequests int, window time.Duration, prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rdb == nil {
			// Skip rate limiting if Redis is unavailable (degraded mode)
			c.Next()
			return
		}

		identifier := c.ClientIP()
		// Use user_id if authenticated, for more granular limiting
		if uid, exists := c.Get("user_id"); exists {
			identifier = uid.(string)
		}

		key := fmt.Sprintf("rl:%s:%s", prefix, identifier)
		ctx := c.Request.Context()

		// Increment counter
		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			// If Redis errors, allow the request (fail-open)
			c.Next()
			return
		}

		// Set TTL on first request
		if count == 1 {
			rdb.Expire(ctx, key, window)
		}

		// Check limit
		if count > int64(maxRequests) {
			ttl, _ := rdb.TTL(ctx, key).Result()
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("Retry-After", fmt.Sprintf("%d", int(ttl.Seconds())))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"code":    "RATE_001",
					"message": "Terlalu banyak permintaan. Coba lagi nanti.",
				},
			})
			c.Abort()
			return
		}

		// Set rate limit headers
		remaining := int64(maxRequests) - count
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		c.Next()
	}
}
