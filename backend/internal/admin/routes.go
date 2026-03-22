package admin

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, authMW, adminMW gin.HandlerFunc) {
	a := rg.Group("/admin/users")
	a.Use(authMW, adminMW)
	{
		a.GET("", h.ListUsers)
		a.POST("", h.CreateUser)
		a.PUT("/:id", h.UpdateUser)
		a.DELETE("/:id", h.DeleteUser)
	}
}
