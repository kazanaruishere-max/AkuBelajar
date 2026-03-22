package assignment

import "time"

// ── Request DTOs ──────────────────────────────────────────

// CreateAssignmentRequest is the payload for creating an assignment.
type CreateAssignmentRequest struct {
	ClassID           string   `json:"class_id" validate:"required,uuid"`
	SubjectID         string   `json:"subject_id" validate:"required,uuid"`
	Title             string   `json:"title" validate:"required,max=200"`
	Description       *string  `json:"description"`
	DeadlineAt        string   `json:"deadline_at" validate:"required"`
	AllowLate         bool     `json:"allow_late"`
	LatePenaltyPct    int      `json:"late_penalty_pct" validate:"min=0,max=100"`
	MaxLateDays       int      `json:"max_late_days" validate:"min=0,max=30"`
	MaxFileCount      int      `json:"max_file_count" validate:"min=1,max=10"`
	MaxFileSizeMB     int      `json:"max_file_size_mb" validate:"min=1,max=50"`
	AllowedExtensions []string `json:"allowed_extensions"`
	WeightPct         int      `json:"weight_pct" validate:"min=1,max=100"`
}

// UpdateAssignmentRequest is the payload for updating an assignment.
type UpdateAssignmentRequest struct {
	Title             *string  `json:"title" validate:"omitempty,max=200"`
	Description       *string  `json:"description"`
	DeadlineAt        *string  `json:"deadline_at"`
	AllowLate         *bool    `json:"allow_late"`
	LatePenaltyPct    *int     `json:"late_penalty_pct" validate:"omitempty,min=0,max=100"`
	MaxLateDays       *int     `json:"max_late_days" validate:"omitempty,min=0,max=30"`
	WeightPct         *int     `json:"weight_pct" validate:"omitempty,min=1,max=100"`
	Status            *string  `json:"status" validate:"omitempty,oneof=draft published closed"`
}

// GradeSubmissionRequest is the payload for grading a submission.
type GradeSubmissionRequest struct {
	Grade    int     `json:"grade" validate:"required,min=0,max=100"`
	Feedback *string `json:"feedback"`
}

// ── Response DTOs ─────────────────────────────────────────

// AssignmentResponse is the API response for an assignment.
type AssignmentResponse struct {
	ID                string     `json:"id"`
	ClassID           string     `json:"class_id"`
	ClassName         string     `json:"class_name,omitempty"`
	SubjectID         string     `json:"subject_id"`
	SubjectName       string     `json:"subject_name,omitempty"`
	TeacherID         string     `json:"teacher_id"`
	TeacherEmail      string     `json:"teacher_email,omitempty"`
	Title             string     `json:"title"`
	Description       *string    `json:"description,omitempty"`
	DeadlineAt        time.Time  `json:"deadline_at"`
	AllowLate         bool       `json:"allow_late"`
	LatePenaltyPct    int        `json:"late_penalty_pct"`
	MaxLateDays       int        `json:"max_late_days"`
	MaxFileCount      int        `json:"max_file_count"`
	MaxFileSizeMB     int        `json:"max_file_size_mb"`
	WeightPct         int        `json:"weight_pct"`
	Status            string     `json:"status"`
	SubmissionCount   int        `json:"submission_count"`
	GradedCount       int        `json:"graded_count"`
	CreatedAt         time.Time  `json:"created_at"`
}

// SubmissionResponse is the API response for a student submission.
type SubmissionResponse struct {
	ID                string     `json:"id"`
	AssignmentID      string     `json:"assignment_id"`
	StudentID         string     `json:"student_id"`
	StudentEmail      string     `json:"student_email,omitempty"`
	SubmittedAt       *time.Time `json:"submitted_at,omitempty"`
	IsLate            bool       `json:"is_late"`
	LateDays          int        `json:"late_days"`
	Status            string     `json:"status"`
	Grade             *int       `json:"grade,omitempty"`
	GradeAfterPenalty *int       `json:"grade_after_penalty,omitempty"`
	Feedback          *string    `json:"feedback,omitempty"`
	GradedAt          *time.Time `json:"graded_at,omitempty"`
	Files             []FileResponse `json:"files"`
	CreatedAt         time.Time  `json:"created_at"`
}

// FileResponse is the API response for an uploaded file.
type FileResponse struct {
	ID       string `json:"id"`
	FileURL  string `json:"file_url"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
}
