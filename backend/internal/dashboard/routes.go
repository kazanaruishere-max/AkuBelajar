package dashboard

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc) {
	rg.GET("/dashboard/stats", authMW, h.Stats)
}
