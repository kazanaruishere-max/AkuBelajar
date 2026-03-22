package quiz

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/response"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/validator"
)

type Handler struct {
	service   *Service
	repo      *Repository
	validator *validator.Validator
}

func NewHandler(service *Service, repo *Repository, v *validator.Validator) *Handler {
	return &Handler{service: service, repo: repo, validator: v}
}

func uid(c *gin.Context) string { v, _ := c.Get("user_id"); return v.(string) }

// ── Teacher ───────────────────────────────────────────────

func (h *Handler) ListTeacherQuizzes(c *gin.Context) {
	items, err := h.repo.ListByTeacher(c.Request.Context(), uid(c))
	if err != nil { response.InternalError(c, "QUIZ_010", "Gagal mengambil data kuis"); return }
	response.OK(c, items)
}

func (h *Handler) CreateQuiz(c *gin.Context) {
	var req CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "QUIZ_001", "Format request tidak valid"); return }
	if errs := h.validator.ValidateStruct(req); errs != nil { response.BadRequestWithDetails(c, "QUIZ_002", "Validasi gagal", errs); return }
	id, err := h.service.Create(c.Request.Context(), uid(c), req)
	if err != nil { response.InternalError(c, "QUIZ_010", "Gagal membuat kuis"); return }
	response.Created(c, gin.H{"id": id})
}

func (h *Handler) PublishQuiz(c *gin.Context) {
	if err := h.service.Publish(c.Request.Context(), c.Param("id"), uid(c)); err != nil {
		h.handleError(c, err); return
	}
	response.OK(c, gin.H{"message": "Kuis berhasil dipublish"})
}

func (h *Handler) DeleteQuiz(c *gin.Context) {
	if err := h.repo.Delete(c.Request.Context(), c.Param("id")); err != nil {
		response.InternalError(c, "QUIZ_010", "Gagal menghapus kuis"); return
	}
	response.NoContent(c)
}

func (h *Handler) AddQuestion(c *gin.Context) {
	var req QuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "QUIZ_001", "Format request tidak valid"); return }
	id, err := h.repo.AddQuestion(c.Request.Context(), c.Param("id"), req)
	if err != nil { response.InternalError(c, "QUIZ_010", "Gagal menambah soal"); return }
	response.Created(c, gin.H{"id": id})
}

func (h *Handler) ListQuestions(c *gin.Context) {
	items, err := h.repo.ListQuestions(c.Request.Context(), c.Param("id"))
	if err != nil { response.InternalError(c, "QUIZ_010", "Gagal mengambil soal"); return }
	response.OK(c, items)
}

func (h *Handler) ListSessions(c *gin.Context) {
	items, err := h.repo.ListSessions(c.Request.Context(), c.Param("id"))
	if err != nil { response.InternalError(c, "QUIZ_010", "Gagal mengambil data sesi"); return }
	response.OK(c, items)
}

// ── Student ───────────────────────────────────────────────

func (h *Handler) ListStudentQuizzes(c *gin.Context) {
	classID := c.Query("class_id")
	if classID == "" { response.BadRequest(c, "QUIZ_001", "class_id wajib diisi"); return }
	items, err := h.repo.ListByClass(c.Request.Context(), classID)
	if err != nil { response.InternalError(c, "QUIZ_010", "Gagal mengambil data kuis"); return }
	response.OK(c, items)
}

func (h *Handler) StartQuiz(c *gin.Context) {
	session, err := h.service.StartSession(c.Request.Context(), c.Param("id"), uid(c), c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil { h.handleError(c, err); return }
	response.Created(c, session)
}

func (h *Handler) SaveAnswer(c *gin.Context) {
	var req AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "QUIZ_001", "Format request tidak valid"); return }
	if err := h.service.SaveAnswer(c.Request.Context(), c.Param("sessionId"), req); err != nil {
		h.handleError(c, err); return
	}
	response.OK(c, gin.H{"message": "Jawaban tersimpan"})
}

func (h *Handler) SubmitQuiz(c *gin.Context) {
	session, err := h.service.SubmitSession(c.Request.Context(), c.Param("sessionId"))
	if err != nil { h.handleError(c, err); return }
	response.OK(c, session)
}

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):      response.NotFound(c, "QUIZ_004", err.Error())
	case errors.Is(err, ErrNotAuthorized): response.Forbidden(c, "QUIZ_005", err.Error())
	case errors.Is(err, ErrNotPublished):  response.BadRequest(c, "QUIZ_006", err.Error())
	case errors.Is(err, ErrQuizEnded):     response.BadRequest(c, "QUIZ_007", err.Error())
	case errors.Is(err, ErrMaxAttempts):   response.BadRequest(c, "QUIZ_008", err.Error())
	case errors.Is(err, ErrSessionExpired): response.BadRequest(c, "QUIZ_009", err.Error())
	default: response.InternalError(c, "QUIZ_010", "Terjadi kesalahan server")
	}
}
