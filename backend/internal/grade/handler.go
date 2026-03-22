package grade

import (
	"github.com/gin-gonic/gin"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/response"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/validator"
)

type Handler struct {
	repo      *Repository
	validator *validator.Validator
}

func NewHandler(repo *Repository, v *validator.Validator) *Handler {
	return &Handler{repo: repo, validator: v}
}

func uid(c *gin.Context) string { v, _ := c.Get("user_id"); return v.(string) }

// InsertGrade adds a new grade entry (teacher).
func (h *Handler) InsertGrade(c *gin.Context) {
	var req InsertGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "GRADE_001", "Format request tidak valid"); return }
	if errs := h.validator.ValidateStruct(req); errs != nil { response.BadRequestWithDetails(c, "GRADE_002", "Validasi gagal", errs); return }
	if err := h.repo.Insert(c.Request.Context(), req); err != nil {
		response.InternalError(c, "GRADE_010", "Gagal menyimpan nilai"); return
	}
	response.Created(c, gin.H{"message": "Nilai berhasil disimpan"})
}

// ListGrades returns grades for a student (teacher view).
func (h *Handler) ListGrades(c *gin.Context) {
	items, err := h.repo.ListByStudent(c.Request.Context(), c.Query("student_id"), c.Query("class_id"), c.Query("subject_id"))
	if err != nil { response.InternalError(c, "GRADE_010", "Gagal mengambil data nilai"); return }
	response.OK(c, items)
}

// Summary returns grade summary for a class.
func (h *Handler) Summary(c *gin.Context) {
	items, err := h.repo.Summary(c.Request.Context(), c.Query("class_id"), c.Query("academic_year_id"))
	if err != nil { response.InternalError(c, "GRADE_010", "Gagal mengambil rekap nilai"); return }
	response.OK(c, items)
}

// MyGrades returns a student's own grades.
func (h *Handler) MyGrades(c *gin.Context) {
	items, err := h.repo.ListByStudent(c.Request.Context(), uid(c), c.Query("class_id"), c.Query("subject_id"))
	if err != nil { response.InternalError(c, "GRADE_010", "Gagal mengambil data nilai"); return }
	response.OK(c, items)
}
