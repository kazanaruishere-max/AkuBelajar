package assignment

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers assignment routes.
func RegisterRoutes(rg *gin.RouterGroup, h *Handler, authMW, teacherMW, studentMW gin.HandlerFunc) {
	asgn := rg.Group("/assignments")
	asgn.Use(authMW)

	// Teacher routes
	teacher := asgn.Group("/teacher")
	teacher.Use(teacherMW)
	{
		teacher.GET("", h.ListTeacherAssignments)
		teacher.POST("", h.CreateAssignment)
		teacher.POST("/:id/publish", h.PublishAssignment)
		teacher.POST("/:id/close", h.CloseAssignment)
		teacher.DELETE("/:id", h.DeleteAssignment)
		teacher.GET("/:id/submissions", h.ListSubmissions)
		teacher.POST("/:id/submissions/:subId/grade", h.GradeSubmission)
	}

	// Student routes
	student := asgn.Group("/student")
	student.Use(studentMW)
	{
		student.GET("", h.ListStudentAssignments)
		student.POST("/:id/submit", h.SubmitAssignment)
		student.GET("/:id/my-submission", h.GetMySubmission)
	}
}
