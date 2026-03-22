package academic

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/response"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/validator"
)

// Handler handles academic HTTP requests.
type Handler struct {
	repo      *Repository
	validator *validator.Validator
}

// NewHandler creates a new academic handler.
func NewHandler(repo *Repository, v *validator.Validator) *Handler {
	return &Handler{repo: repo, validator: v}
}

func (h *Handler) schoolID(c *gin.Context) string {
	sid, _ := c.Get("school_id")
	return sid.(string)
}

// ── Academic Years ────────────────────────────────────────

func (h *Handler) ListAcademicYears(c *gin.Context) {
	items, err := h.repo.ListAcademicYears(c.Request.Context(), h.schoolID(c))
	if err != nil {
		response.InternalError(c, "ACAD_010", "Gagal mengambil data tahun ajaran")
		return
	}
	response.OK(c, items)
}

func (h *Handler) CreateAcademicYear(c *gin.Context) {
	var req AcademicYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "ACAD_001", "Format request tidak valid")
		return
	}
	if errs := h.validator.ValidateStruct(req); errs != nil {
		response.BadRequestWithDetails(c, "ACAD_002", "Validasi gagal", errs)
		return
	}
	item, err := h.repo.CreateAcademicYear(c.Request.Context(), h.schoolID(c), req)
	if err != nil {
		response.InternalError(c, "ACAD_010", "Gagal membuat tahun ajaran")
		return
	}
	response.Created(c, item)
}

func (h *Handler) UpdateAcademicYear(c *gin.Context) {
	id := c.Param("id")
	var req AcademicYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "ACAD_001", "Format request tidak valid")
		return
	}
	if errs := h.validator.ValidateStruct(req); errs != nil {
		response.BadRequestWithDetails(c, "ACAD_002", "Validasi gagal", errs)
		return
	}
	item, err := h.repo.UpdateAcademicYear(c.Request.Context(), h.schoolID(c), id, req)
	if err != nil {
		response.InternalError(c, "ACAD_010", "Gagal mengupdate tahun ajaran")
		return
	}
	response.OK(c, item)
}

func (h *Handler) DeleteAcademicYear(c *gin.Context) {
	err := h.repo.DeleteAcademicYear(c.Request.Context(), h.schoolID(c), c.Param("id"))
	if err != nil {
		if errors.Is(err, ErrCannotDeleteActive) {
			response.BadRequest(c, "ACAD_003", err.Error())
			return
		}
		response.InternalError(c, "ACAD_010", "Gagal menghapus tahun ajaran")
		return
	}
	response.NoContent(c)
}

// ── Classes ───────────────────────────────────────────────

func (h *Handler) ListClasses(c *gin.Context) {
	ayID := c.Query("academic_year_id")
	if ayID == "" {
		response.BadRequest(c, "ACAD_001", "academic_year_id wajib diisi")
		return
	}
	items, err := h.repo.ListClasses(c.Request.Context(), h.schoolID(c), ayID)
	if err != nil {
		response.InternalError(c, "ACAD_010", "Gagal mengambil data kelas")
		return
	}
	response.OK(c, items)
}

func (h *Handler) CreateClass(c *gin.Context) {
	var req ClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "ACAD_001", "Format request tidak valid")
		return
	}
	if errs := h.validator.ValidateStruct(req); errs != nil {
		response.BadRequestWithDetails(c, "ACAD_002", "Validasi gagal", errs)
		return
	}
	item, err := h.repo.CreateClass(c.Request.Context(), h.schoolID(c), req)
	if err != nil {
		response.InternalError(c, "ACAD_010", "Gagal membuat kelas")
		return
	}
	response.Created(c, item)
}

func (h *Handler) UpdateClass(c *gin.Context) {
	var req ClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "ACAD_001", "Format request tidak valid")
		return
	}
	if errs := h.validator.ValidateStruct(req); errs != nil {
		response.BadRequestWithDetails(c, "ACAD_002", "Validasi gagal", errs)
		return
	}
	item, err := h.repo.UpdateClass(c.Request.Context(), h.schoolID(c), c.Param("id"), req)
	if err != nil {
		response.InternalError(c, "ACAD_010", "Gagal mengupdate kelas")
		return
	}
	response.OK(c, item)
}

func (h *Handler) DeleteClass(c *gin.Context) {
	if err := h.repo.DeleteClass(c.Request.Context(), h.schoolID(c), c.Param("id")); err != nil {
		response.InternalError(c, "ACAD_010", "Gagal menghapus kelas")
		return
	}
	response.NoContent(c)
}

func (h *Handler) ListStudents(c *gin.Context) {
	items, err := h.repo.ListStudentsInClass(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, "ACAD_010", "Gagal mengambil data siswa")
		return
	}
	response.OK(c, items)
}

func (h *Handler) AssignStudents(c *gin.Context) {
	var req AssignStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "ACAD_001", "Format request tidak valid")
		return
	}
	// Need academic_year_id from query or class lookup
	ayID := c.Query("academic_year_id")
	if ayID == "" {
		response.BadRequest(c, "ACAD_001", "academic_year_id wajib diisi")
		return
	}
	if err := h.repo.AssignStudents(c.Request.Context(), c.Param("id"), ayID, req.StudentIDs); err != nil {
		response.InternalError(c, "ACAD_010", "Gagal menambahkan siswa")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Siswa berhasil ditambahkan"}})
}

func (h *Handler) RemoveStudent(c *gin.Context) {
	if err := h.repo.RemoveStudent(c.Request.Context(), c.Param("id"), c.Param("studentId")); err != nil {
		response.InternalError(c, "ACAD_010", "Gagal menghapus siswa dari kelas")
		return
	}
	response.NoContent(c)
}

// ── Subjects ──────────────────────────────────────────────

func (h *Handler) ListSubjects(c *gin.Context) {
	items, err := h.repo.ListSubjects(c.Request.Context(), h.schoolID(c))
	if err != nil {
		response.InternalError(c, "ACAD_010", "Gagal mengambil data mata pelajaran")
		return
	}
	response.OK(c, items)
}

func (h *Handler) CreateSubject(c *gin.Context) {
	var req SubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "ACAD_001", "Format request tidak valid")
		return
	}
	if errs := h.validator.ValidateStruct(req); errs != nil {
		response.BadRequestWithDetails(c, "ACAD_002", "Validasi gagal", errs)
		return
	}
	item, err := h.repo.CreateSubject(c.Request.Context(), h.schoolID(c), req)
	if err != nil {
		response.InternalError(c, "ACAD_010", "Gagal membuat mata pelajaran")
		return
	}
	response.Created(c, item)
}

func (h *Handler) UpdateSubject(c *gin.Context) {
	var req SubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "ACAD_001", "Format request tidak valid")
		return
	}
	item, err := h.repo.UpdateSubject(c.Request.Context(), h.schoolID(c), c.Param("id"), req)
	if err != nil {
		response.InternalError(c, "ACAD_010", "Gagal mengupdate mata pelajaran")
		return
	}
	response.OK(c, item)
}

func (h *Handler) DeleteSubject(c *gin.Context) {
	if err := h.repo.DeleteSubject(c.Request.Context(), h.schoolID(c), c.Param("id")); err != nil {
		response.InternalError(c, "ACAD_010", "Gagal menghapus mata pelajaran")
		return
	}
	response.NoContent(c)
}

// ── Class-Subject Assignments ─────────────────────────────

func (h *Handler) ListClassSubjects(c *gin.Context) {
	items, err := h.repo.ListClassSubjects(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, "ACAD_010", "Gagal mengambil data pengajar")
		return
	}
	response.OK(c, items)
}

func (h *Handler) AssignTeacher(c *gin.Context) {
	var req ClassSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "ACAD_001", "Format request tidak valid")
		return
	}
	item, err := h.repo.AssignTeacherToSubject(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		response.InternalError(c, "ACAD_010", "Gagal assign guru ke mapel")
		return
	}
	response.Created(c, item)
}

func (h *Handler) RemoveClassSubject(c *gin.Context) {
	if err := h.repo.RemoveClassSubject(c.Request.Context(), c.Param("csId")); err != nil {
		response.InternalError(c, "ACAD_010", "Gagal menghapus pengajar dari mapel")
		return
	}
	response.NoContent(c)
}
