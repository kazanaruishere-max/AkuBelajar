package attendance

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

// MarkAttendance marks attendance for a class on a date.
func (h *Handler) MarkAttendance(c *gin.Context) {
	var req MarkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "ATT_001", "Format request tidak valid"); return }
	if errs := h.validator.ValidateStruct(req); errs != nil { response.BadRequestWithDetails(c, "ATT_002", "Validasi gagal", errs); return }
	if err := h.repo.MarkBatch(c.Request.Context(), uid(c), req); err != nil {
		response.InternalError(c, "ATT_010", "Gagal menyimpan presensi"); return
	}
	response.OK(c, gin.H{"message": "Presensi berhasil disimpan"})
}

// ListByDate returns attendance records for a class/subject on a date.
func (h *Handler) ListByDate(c *gin.Context) {
	items, err := h.repo.ListByDate(c.Request.Context(), c.Query("class_id"), c.Query("subject_id"), c.Query("date"))
	if err != nil { response.InternalError(c, "ATT_010", "Gagal mengambil data presensi"); return }
	response.OK(c, items)
}

// Summary returns attendance summary per student.
func (h *Handler) Summary(c *gin.Context) {
	items, err := h.repo.Summary(c.Request.Context(), c.Query("class_id"), c.Query("subject_id"), c.Query("academic_year_id"))
	if err != nil { response.InternalError(c, "ATT_010", "Gagal mengambil rekap presensi"); return }
	response.OK(c, items)
}

// StudentHistory returns a student's own attendance history.
func (h *Handler) StudentHistory(c *gin.Context) {
	items, err := h.repo.StudentHistory(c.Request.Context(), uid(c), c.Query("class_id"), c.Query("subject_id"))
	if err != nil { response.InternalError(c, "ATT_010", "Gagal mengambil riwayat presensi"); return }
	response.OK(c, items)
}
