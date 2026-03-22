package attendance

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, authMW, teacherMW, studentMW gin.HandlerFunc) {
	att := rg.Group("/attendance")
	att.Use(authMW)

	teacher := att.Group("/teacher")
	teacher.Use(teacherMW)
	{
		teacher.POST("", h.MarkAttendance)
		teacher.GET("", h.ListByDate)
		teacher.GET("/summary", h.Summary)
	}

	student := att.Group("/student")
	student.Use(studentMW)
	{
		student.GET("/history", h.StudentHistory)
	}
}
