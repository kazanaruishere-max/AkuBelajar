package assignment

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles assignment database operations.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new assignment repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// ListByTeacher returns assignments created by a teacher.
func (r *Repository) ListByTeacher(ctx context.Context, teacherID string) ([]AssignmentResponse, error) {
	return r.queryAssignments(ctx, `
		WHERE a.teacher_id = $1 AND a.deleted_at IS NULL ORDER BY a.created_at DESC
	`, teacherID)
}

// ListByClass returns published assignments for a class.
func (r *Repository) ListByClass(ctx context.Context, classID string) ([]AssignmentResponse, error) {
	return r.queryAssignments(ctx, `
		WHERE a.class_id = $1 AND a.status = 'published' AND a.deleted_at IS NULL ORDER BY a.deadline_at DESC
	`, classID)
}

// GetByID returns an assignment by ID.
func (r *Repository) GetByID(ctx context.Context, id string) (*AssignmentResponse, error) {
	items, err := r.queryAssignments(ctx, `WHERE a.id = $1 AND a.deleted_at IS NULL`, id)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &items[0], nil
}

func (r *Repository) queryAssignments(ctx context.Context, where string, args ...interface{}) ([]AssignmentResponse, error) {
	q := `
		SELECT a.id, a.class_id, COALESCE(cl.name,''), a.subject_id, COALESCE(s.name,''),
		       a.teacher_id, COALESCE(u.email,''), a.title, a.description,
		       a.deadline_at, a.allow_late, a.late_penalty_pct, a.max_late_days,
		       a.max_file_count, a.max_file_size_mb, a.weight_pct, a.status,
		       (SELECT COUNT(*) FROM submissions sub WHERE sub.assignment_id = a.id AND sub.status != 'draft'),
		       (SELECT COUNT(*) FROM submissions sub WHERE sub.assignment_id = a.id AND sub.status = 'graded'),
		       a.created_at
		FROM assignments a
		LEFT JOIN classes cl ON cl.id = a.class_id
		LEFT JOIN subjects s ON s.id = a.subject_id
		LEFT JOIN users u ON u.id = a.teacher_id
	` + where

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []AssignmentResponse
	for rows.Next() {
		var a AssignmentResponse
		if err := rows.Scan(
			&a.ID, &a.ClassID, &a.ClassName, &a.SubjectID, &a.SubjectName,
			&a.TeacherID, &a.TeacherEmail, &a.Title, &a.Description,
			&a.DeadlineAt, &a.AllowLate, &a.LatePenaltyPct, &a.MaxLateDays,
			&a.MaxFileCount, &a.MaxFileSizeMB, &a.WeightPct, &a.Status,
			&a.SubmissionCount, &a.GradedCount, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, nil
}

// Create inserts a new assignment.
func (r *Repository) Create(ctx context.Context, teacherID string, req CreateAssignmentRequest) (*AssignmentResponse, error) {
	var id string
	var createdAt, deadlineAt time.Time
	err := r.db.QueryRow(ctx, `
		INSERT INTO assignments (
			class_id, subject_id, teacher_id, title, description,
			deadline_at, allow_late, late_penalty_pct, max_late_days,
			max_file_count, max_file_size_mb, weight_pct, status
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,'draft')
		RETURNING id, created_at, deadline_at
	`, req.ClassID, req.SubjectID, teacherID, req.Title, req.Description,
		req.DeadlineAt, req.AllowLate, req.LatePenaltyPct, req.MaxLateDays,
		req.MaxFileCount, req.MaxFileSizeMB, req.WeightPct,
	).Scan(&id, &createdAt, &deadlineAt)
	if err != nil {
		return nil, err
	}

	return &AssignmentResponse{
		ID: id, ClassID: req.ClassID, SubjectID: req.SubjectID,
		TeacherID: teacherID, Title: req.Title, Description: req.Description,
		DeadlineAt: deadlineAt, AllowLate: req.AllowLate,
		LatePenaltyPct: req.LatePenaltyPct, MaxLateDays: req.MaxLateDays,
		MaxFileCount: req.MaxFileCount, MaxFileSizeMB: req.MaxFileSizeMB,
		WeightPct: req.WeightPct, Status: "draft", CreatedAt: createdAt,
	}, nil
}

// UpdateStatus changes the assignment status (draft → published → closed).
func (r *Repository) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE assignments SET status = $2::assignment_status, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`, id, status)
	return err
}

// Delete soft-deletes an assignment.
func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE assignments SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL
	`, id)
	return err
}

// ── Submissions ───────────────────────────────────────────

// GetSubmission returns a student's submission for an assignment.
func (r *Repository) GetSubmission(ctx context.Context, assignmentID, studentID string) (*SubmissionResponse, error) {
	var s SubmissionResponse
	err := r.db.QueryRow(ctx, `
		SELECT id, assignment_id, student_id, submitted_at, is_late, late_days,
		       status, grade, grade_after_penalty, feedback, graded_at, created_at
		FROM submissions
		WHERE assignment_id = $1 AND student_id = $2
	`, assignmentID, studentID).Scan(
		&s.ID, &s.AssignmentID, &s.StudentID, &s.SubmittedAt, &s.IsLate, &s.LateDays,
		&s.Status, &s.Grade, &s.GradeAfterPenalty, &s.Feedback, &s.GradedAt, &s.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Get files
	s.Files, _ = r.getSubmissionFiles(ctx, s.ID)
	return &s, nil
}

// CreateSubmission creates or gets an existing submission.
func (r *Repository) CreateSubmission(ctx context.Context, assignmentID, studentID string) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		INSERT INTO submissions (assignment_id, student_id, status)
		VALUES ($1, $2, 'draft')
		ON CONFLICT (assignment_id, student_id) DO UPDATE SET updated_at = NOW()
		RETURNING id
	`, assignmentID, studentID).Scan(&id)
	return id, err
}

// SubmitSubmission marks a submission as submitted.
func (r *Repository) SubmitSubmission(ctx context.Context, submissionID string, isLate bool, lateDays int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE submissions
		SET submitted_at = NOW(), status = 'submitted', is_late = $2, late_days = $3, updated_at = NOW()
		WHERE id = $1
	`, submissionID, isLate, lateDays)
	return err
}

// AddSubmissionFile adds a file to a submission.
func (r *Repository) AddSubmissionFile(ctx context.Context, submissionID, fileURL, fileName string, fileSize int64) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO submission_files (submission_id, file_url, file_name, file_size)
		VALUES ($1, $2, $3, $4)
	`, submissionID, fileURL, fileName, fileSize)
	return err
}

// ListSubmissions returns all submissions for an assignment (teacher view).
func (r *Repository) ListSubmissions(ctx context.Context, assignmentID string) ([]SubmissionResponse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT s.id, s.assignment_id, s.student_id, u.email, s.submitted_at,
		       s.is_late, s.late_days, s.status, s.grade, s.grade_after_penalty,
		       s.feedback, s.graded_at, s.created_at
		FROM submissions s
		JOIN users u ON u.id = s.student_id
		WHERE s.assignment_id = $1
		ORDER BY u.email
	`, assignmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []SubmissionResponse
	for rows.Next() {
		var s SubmissionResponse
		if err := rows.Scan(
			&s.ID, &s.AssignmentID, &s.StudentID, &s.StudentEmail, &s.SubmittedAt,
			&s.IsLate, &s.LateDays, &s.Status, &s.Grade, &s.GradeAfterPenalty,
			&s.Feedback, &s.GradedAt, &s.CreatedAt,
		); err != nil {
			return nil, err
		}
		s.Files, _ = r.getSubmissionFiles(ctx, s.ID)
		items = append(items, s)
	}
	return items, nil
}

// GradeSubmission sets the grade for a submission.
func (r *Repository) GradeSubmission(ctx context.Context, submissionID, gradedBy string, grade, gradeAfterPenalty int, feedback *string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE submissions
		SET grade = $2, grade_after_penalty = $3, feedback = $4, status = 'graded',
		    graded_by = $5, graded_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, submissionID, grade, gradeAfterPenalty, feedback, gradedBy)
	return err
}

func (r *Repository) getSubmissionFiles(ctx context.Context, submissionID string) ([]FileResponse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, file_url, file_name, file_size FROM submission_files WHERE submission_id = $1
	`, submissionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []FileResponse
	for rows.Next() {
		var f FileResponse
		if err := rows.Scan(&f.ID, &f.FileURL, &f.FileName, &f.FileSize); err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}
