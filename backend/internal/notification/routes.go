package notification

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, authMW, teacherMW gin.HandlerFunc) {
	n := rg.Group("/notifications")
	n.Use(authMW)
	{
		n.GET("", h.ListNotifications)
		n.GET("/unread-count", h.UnreadCount)
		n.POST("/:id/read", h.MarkRead)
		n.POST("/read-all", h.MarkAllRead)
	}

	// Teacher/Admin can send notifications
	admin := n.Group("")
	admin.Use(teacherMW)
	{
		admin.POST("/send", h.Send)
		admin.POST("/broadcast", h.Broadcast)
	}
}
