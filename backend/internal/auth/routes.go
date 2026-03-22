package auth

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all auth routes on the given router group.
func RegisterRoutes(rg *gin.RouterGroup, h *Handler, authMiddleware gin.HandlerFunc) {
	auth := rg.Group("/auth")
	{
		// Public routes (no auth required)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)

		// Protected routes (auth required)
		auth.POST("/logout", authMiddleware, h.Logout)
		auth.POST("/change-password", authMiddleware, h.ChangePassword)
		auth.GET("/me", authMiddleware, h.GetMe)
	}
}
