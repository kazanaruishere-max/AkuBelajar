package academic

import "time"

// ── Request DTOs ──────────────────────────────────────────

// AcademicYearRequest is the payload for creating/updating an academic year.
type AcademicYearRequest struct {
	Name      string `json:"name" validate:"required,max=20"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
	IsActive  *bool  `json:"is_active"`
}

// ClassRequest is the payload for creating/updating a class.
type ClassRequest struct {
	Name              string  `json:"name" validate:"required,max=20"`
	GradeLevel        int     `json:"grade_level" validate:"required,min=1,max=12"`
	AcademicYearID    string  `json:"academic_year_id" validate:"required,uuid"`
	HomeroomTeacherID *string `json:"homeroom_teacher_id" validate:"omitempty,uuid"`
}

// SubjectRequest is the payload for creating/updating a subject.
type SubjectRequest struct {
	Name        string  `json:"name" validate:"required,max=100"`
	Code        *string `json:"code" validate:"omitempty,max=10"`
	Description *string `json:"description"`
}

// AssignStudentRequest adds students to a class.
type AssignStudentRequest struct {
	StudentIDs []string `json:"student_ids" validate:"required,min=1,dive,uuid"`
}

// ClassSubjectRequest assigns a teacher to teach a subject in a class.
type ClassSubjectRequest struct {
	SubjectID string `json:"subject_id" validate:"required,uuid"`
	TeacherID string `json:"teacher_id" validate:"required,uuid"`
}

// ── Response DTOs ─────────────────────────────────────────

// AcademicYearResponse is the API response for an academic year.
type AcademicYearResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// ClassResponse is the API response for a class.
type ClassResponse struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	GradeLevel        int     `json:"grade_level"`
	AcademicYearID    string  `json:"academic_year_id"`
	AcademicYearName  string  `json:"academic_year_name,omitempty"`
	HomeroomTeacherID *string `json:"homeroom_teacher_id,omitempty"`
	HomeroomTeacher   *string `json:"homeroom_teacher,omitempty"`
	StudentCount      int     `json:"student_count"`
	CreatedAt         string  `json:"created_at"`
}

// SubjectResponse is the API response for a subject.
type SubjectResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Code        *string `json:"code,omitempty"`
	Description *string `json:"description,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

// StudentListItem is a student in a class list.
type StudentListItem struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// ClassSubjectResponse shows a subject-teacher assignment in a class.
type ClassSubjectResponse struct {
	ID           string `json:"id"`
	SubjectID    string `json:"subject_id"`
	SubjectName  string `json:"subject_name"`
	TeacherID    string `json:"teacher_id"`
	TeacherEmail string `json:"teacher_email"`
}
