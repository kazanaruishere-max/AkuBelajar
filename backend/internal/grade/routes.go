package grade

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, authMW, teacherMW, studentMW gin.HandlerFunc) {
	g := rg.Group("/grades")
	g.Use(authMW)

	teacher := g.Group("/teacher")
	teacher.Use(teacherMW)
	{
		teacher.POST("", h.InsertGrade)
		teacher.GET("", h.ListGrades)
		teacher.GET("/summary", h.Summary)
	}

	student := g.Group("/student")
	student.Use(studentMW)
	{
		student.GET("", h.MyGrades)
	}
}
