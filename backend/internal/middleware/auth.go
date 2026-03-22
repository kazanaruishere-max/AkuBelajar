package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/response"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/security"
)

// AuthMiddleware validates the Paseto token from the Authorization header
// and sets user_id, school_id, role in the Gin context.
func AuthMiddleware(tokenMaker *security.TokenMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "AUTH_003", "Token tidak ditemukan")
			c.Abort()
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(c, "AUTH_003", "Format token tidak valid")
			c.Abort()
			return
		}

		payload, err := tokenMaker.ValidateToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "AUTH_003", "Token expired atau tidak valid")
			c.Abort()
			return
		}

		// Reject refresh tokens used as access tokens
		if payload.TokenType != "access" {
			response.Unauthorized(c, "AUTH_003", "Jenis token tidak valid")
			c.Abort()
			return
		}

		// Set user info in context for downstream handlers
		c.Set("user_id", payload.UserID)
		c.Set("school_id", payload.SchoolID)
		c.Set("role", payload.Role)

		c.Next()
	}
}
