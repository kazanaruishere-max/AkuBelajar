package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/response"
)

// RequireRole creates a middleware that restricts access to specific roles.
// Usage: router.Use(middleware.RequireRole("super_admin", "teacher"))
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	roleSet := make(map[string]bool, len(allowedRoles))
	for _, r := range allowedRoles {
		roleSet[r] = true
	}

	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.Unauthorized(c, "AUTH_003", "Tidak terautentikasi")
			c.Abort()
			return
		}

		if !roleSet[role.(string)] {
			response.Forbidden(c, "AUTH_008", "Anda tidak memiliki akses ke fitur ini")
			c.Abort()
			return
		}

		c.Next()
	}
}
