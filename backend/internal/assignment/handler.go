package assignment

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/response"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/validator"
)

// Handler handles assignment HTTP requests.
type Handler struct {
	service   *Service
	repo      *Repository
	validator *validator.Validator
}

// NewHandler creates a new assignment handler.
func NewHandler(service *Service, repo *Repository, v *validator.Validator) *Handler {
	return &Handler{service: service, repo: repo, validator: v}
}

func userID(c *gin.Context) string {
	uid, _ := c.Get("user_id")
	return uid.(string)
}

// ── Teacher Endpoints ─────────────────────────────────────

// ListTeacherAssignments returns assignments created by the teacher.
func (h *Handler) ListTeacherAssignments(c *gin.Context) {
	items, err := h.repo.ListByTeacher(c.Request.Context(), userID(c))
	if err != nil {
		response.InternalError(c, "ASGN_010", "Gagal mengambil data tugas")
		return
	}
	response.OK(c, items)
}

// CreateAssignment creates a new assignment.
func (h *Handler) CreateAssignment(c *gin.Context) {
	var req CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "ASGN_001", "Format request tidak valid")
		return
	}
	if errs := h.validator.ValidateStruct(req); errs != nil {
		response.BadRequestWithDetails(c, "ASGN_002", "Validasi gagal", errs)
		return
	}
	item, err := h.service.Create(c.Request.Context(), userID(c), req)
	if err != nil {
		response.InternalError(c, "ASGN_010", "Gagal membuat tugas")
		return
	}
	response.Created(c, item)
}

// PublishAssignment publishes a draft assignment.
func (h *Handler) PublishAssignment(c *gin.Context) {
	err := h.service.Publish(c.Request.Context(), c.Param("id"), userID(c))
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, gin.H{"message": "Tugas berhasil dipublish"})
}

// CloseAssignment closes an assignment.
func (h *Handler) CloseAssignment(c *gin.Context) {
	err := h.service.Close(c.Request.Context(), c.Param("id"), userID(c))
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, gin.H{"message": "Tugas berhasil ditutup"})
}

// DeleteAssignment soft-deletes an assignment.
func (h *Handler) DeleteAssignment(c *gin.Context) {
	if err := h.repo.Delete(c.Request.Context(), c.Param("id")); err != nil {
		response.InternalError(c, "ASGN_010", "Gagal menghapus tugas")
		return
	}
	response.NoContent(c)
}

// ListSubmissions returns all submissions for an assignment (teacher view).
func (h *Handler) ListSubmissions(c *gin.Context) {
	items, err := h.repo.ListSubmissions(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, "ASGN_010", "Gagal mengambil data submission")
		return
	}
	response.OK(c, items)
}

// GradeSubmission grades a student's submission.
func (h *Handler) GradeSubmission(c *gin.Context) {
	var req GradeSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "ASGN_001", "Format request tidak valid")
		return
	}
	if errs := h.validator.ValidateStruct(req); errs != nil {
		response.BadRequestWithDetails(c, "ASGN_002", "Validasi gagal", errs)
		return
	}
	err := h.service.Grade(c.Request.Context(), c.Param("subId"), c.Param("id"), userID(c), req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, gin.H{"message": "Nilai berhasil disimpan"})
}

// ── Student Endpoints ─────────────────────────────────────

// ListStudentAssignments returns assignments for a student's class.
func (h *Handler) ListStudentAssignments(c *gin.Context) {
	classID := c.Query("class_id")
	if classID == "" {
		response.BadRequest(c, "ASGN_001", "class_id wajib diisi")
		return
	}
	items, err := h.repo.ListByClass(c.Request.Context(), classID)
	if err != nil {
		response.InternalError(c, "ASGN_010", "Gagal mengambil data tugas")
		return
	}
	response.OK(c, items)
}

// SubmitAssignment submits a student's assignment.
func (h *Handler) SubmitAssignment(c *gin.Context) {
	sub, err := h.service.Submit(c.Request.Context(), c.Param("id"), userID(c))
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, sub)
}

// GetMySubmission returns the student's submission for an assignment.
func (h *Handler) GetMySubmission(c *gin.Context) {
	sub, err := h.repo.GetSubmission(c.Request.Context(), c.Param("id"), userID(c))
	if err != nil {
		response.InternalError(c, "ASGN_010", "Gagal mengambil data submission")
		return
	}
	if sub == nil {
		response.NotFound(c, "ASGN_004", "Belum ada submission")
		return
	}
	response.OK(c, sub)
}

// ── Error Handling ────────────────────────────────────────

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		response.NotFound(c, "ASGN_004", err.Error())
	case errors.Is(err, ErrNotAuthorized):
		response.Forbidden(c, "ASGN_005", err.Error())
	case errors.Is(err, ErrDeadlinePassed):
		response.BadRequest(c, "ASGN_006", err.Error())
	case errors.Is(err, ErrNotPublished):
		response.BadRequest(c, "ASGN_007", err.Error())
	case errors.Is(err, ErrAlreadySubmitted):
		response.BadRequest(c, "ASGN_008", err.Error())
	default:
		response.InternalError(c, "ASGN_010", "Terjadi kesalahan server")
	}
}
