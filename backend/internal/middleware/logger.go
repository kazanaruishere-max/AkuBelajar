package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestLogger logs each request with structured metadata.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Process request
		c.Next()

		// Log after response
		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		ip := c.ClientIP()
		userID, _ := c.Get("user_id")

		log.Printf("[%s] %s %s | %d | %v | ip=%s | user=%v",
			requestID[:8], method, path, status, latency, ip, userID,
		)
	}
}
