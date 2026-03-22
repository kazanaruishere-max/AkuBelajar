package quiz

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, authMW, teacherMW, studentMW gin.HandlerFunc) {
	q := rg.Group("/quizzes")
	q.Use(authMW)

	teacher := q.Group("/teacher")
	teacher.Use(teacherMW)
	{
		teacher.GET("", h.ListTeacherQuizzes)
		teacher.POST("", h.CreateQuiz)
		teacher.POST("/:id/publish", h.PublishQuiz)
		teacher.DELETE("/:id", h.DeleteQuiz)
		teacher.POST("/:id/questions", h.AddQuestion)
		teacher.GET("/:id/questions", h.ListQuestions)
		teacher.GET("/:id/sessions", h.ListSessions)
		teacher.POST("/:id/ai-generate", h.GenerateAIQuestions)
	}

	student := q.Group("/student")
	student.Use(studentMW)
	{
		student.GET("", h.ListStudentQuizzes)
		student.POST("/:id/start", h.StartQuiz)
		student.POST("/sessions/:sessionId/answer", h.SaveAnswer)
		student.POST("/sessions/:sessionId/submit", h.SubmitQuiz)
	}
}
