package academic

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers academic routes.
// All academic routes require authentication + super_admin role.
func RegisterRoutes(rg *gin.RouterGroup, h *Handler, authMW, adminMW gin.HandlerFunc) {
	acad := rg.Group("/academic")
	acad.Use(authMW, adminMW)
	{
		// Academic Years
		ay := acad.Group("/years")
		ay.GET("", h.ListAcademicYears)
		ay.POST("", h.CreateAcademicYear)
		ay.PUT("/:id", h.UpdateAcademicYear)
		ay.DELETE("/:id", h.DeleteAcademicYear)

		// Classes
		cls := acad.Group("/classes")
		cls.GET("", h.ListClasses)
		cls.POST("", h.CreateClass)
		cls.PUT("/:id", h.UpdateClass)
		cls.DELETE("/:id", h.DeleteClass)

		// Students in a class
		cls.GET("/:id/students", h.ListStudents)
		cls.POST("/:id/students", h.AssignStudents)
		cls.DELETE("/:id/students/:studentId", h.RemoveStudent)

		// Class-Subject-Teacher assignments
		cls.GET("/:id/subjects", h.ListClassSubjects)
		cls.POST("/:id/subjects", h.AssignTeacher)
		cls.DELETE("/:id/subjects/:csId", h.RemoveClassSubject)

		// Subjects
		subj := acad.Group("/subjects")
		subj.GET("", h.ListSubjects)
		subj.POST("", h.CreateSubject)
		subj.PUT("/:id", h.UpdateSubject)
		subj.DELETE("/:id", h.DeleteSubject)
	}
}
