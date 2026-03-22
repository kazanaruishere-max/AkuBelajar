package academic

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles academic database operations.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new academic repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// ── Academic Years ────────────────────────────────────────

// ListAcademicYears returns all academic years for a school.
func (r *Repository) ListAcademicYears(ctx context.Context, schoolID string) ([]AcademicYearResponse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, start_date, end_date, is_active, created_at
		FROM academic_years
		WHERE school_id = $1
		ORDER BY start_date DESC
	`, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []AcademicYearResponse
	for rows.Next() {
		var ay AcademicYearResponse
		var startDate, endDate time.Time
		if err := rows.Scan(&ay.ID, &ay.Name, &startDate, &endDate, &ay.IsActive, &ay.CreatedAt); err != nil {
			return nil, err
		}
		ay.StartDate = startDate.Format("2006-01-02")
		ay.EndDate = endDate.Format("2006-01-02")
		items = append(items, ay)
	}
	return items, nil
}

// CreateAcademicYear inserts a new academic year.
func (r *Repository) CreateAcademicYear(ctx context.Context, schoolID string, req AcademicYearRequest) (*AcademicYearResponse, error) {
	isActive := false
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// If setting active, deactivate others first
	if isActive {
		_, _ = r.db.Exec(ctx, `UPDATE academic_years SET is_active = FALSE WHERE school_id = $1`, schoolID)
	}

	var ay AcademicYearResponse
	var startDate, endDate time.Time
	err := r.db.QueryRow(ctx, `
		INSERT INTO academic_years (school_id, name, start_date, end_date, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, start_date, end_date, is_active, created_at
	`, schoolID, req.Name, req.StartDate, req.EndDate, isActive).Scan(
		&ay.ID, &ay.Name, &startDate, &endDate, &ay.IsActive, &ay.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	ay.StartDate = startDate.Format("2006-01-02")
	ay.EndDate = endDate.Format("2006-01-02")
	return &ay, nil
}

// UpdateAcademicYear updates an existing academic year.
func (r *Repository) UpdateAcademicYear(ctx context.Context, schoolID, id string, req AcademicYearRequest) (*AcademicYearResponse, error) {
	isActive := false
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	if isActive {
		_, _ = r.db.Exec(ctx, `UPDATE academic_years SET is_active = FALSE WHERE school_id = $1 AND id != $2`, schoolID, id)
	}

	var ay AcademicYearResponse
	var startDate, endDate time.Time
	err := r.db.QueryRow(ctx, `
		UPDATE academic_years
		SET name = $3, start_date = $4, end_date = $5, is_active = $6, updated_at = NOW()
		WHERE id = $1 AND school_id = $2
		RETURNING id, name, start_date, end_date, is_active, created_at
	`, id, schoolID, req.Name, req.StartDate, req.EndDate, isActive).Scan(
		&ay.ID, &ay.Name, &startDate, &endDate, &ay.IsActive, &ay.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	ay.StartDate = startDate.Format("2006-01-02")
	ay.EndDate = endDate.Format("2006-01-02")
	return &ay, nil
}

// DeleteAcademicYear soft-deletes an academic year (only if not active).
func (r *Repository) DeleteAcademicYear(ctx context.Context, schoolID, id string) error {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM academic_years WHERE id = $1 AND school_id = $2 AND is_active = FALSE
	`, id, schoolID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrCannotDeleteActive
	}
	return nil
}

// ── Classes ───────────────────────────────────────────────

// ListClasses returns all classes for a school and academic year.
func (r *Repository) ListClasses(ctx context.Context, schoolID, academicYearID string) ([]ClassResponse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT c.id, c.name, c.grade_level, c.academic_year_id, ay.name,
		       c.homeroom_teacher_id, u.email,
		       (SELECT COUNT(*) FROM student_classes sc WHERE sc.class_id = c.id),
		       c.created_at
		FROM classes c
		JOIN academic_years ay ON ay.id = c.academic_year_id
		LEFT JOIN users u ON u.id = c.homeroom_teacher_id
		WHERE c.school_id = $1 AND c.academic_year_id = $2 AND c.deleted_at IS NULL
		ORDER BY c.grade_level, c.name
	`, schoolID, academicYearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ClassResponse
	for rows.Next() {
		var c ClassResponse
		var createdAt time.Time
		if err := rows.Scan(
			&c.ID, &c.Name, &c.GradeLevel, &c.AcademicYearID, &c.AcademicYearName,
			&c.HomeroomTeacherID, &c.HomeroomTeacher, &c.StudentCount, &createdAt,
		); err != nil {
			return nil, err
		}
		c.CreatedAt = createdAt.Format(time.RFC3339)
		items = append(items, c)
	}
	return items, nil
}

// CreateClass inserts a new class.
func (r *Repository) CreateClass(ctx context.Context, schoolID string, req ClassRequest) (*ClassResponse, error) {
	var c ClassResponse
	var createdAt time.Time
	err := r.db.QueryRow(ctx, `
		INSERT INTO classes (school_id, academic_year_id, name, grade_level, homeroom_teacher_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, grade_level, academic_year_id, created_at
	`, schoolID, req.AcademicYearID, req.Name, req.GradeLevel, req.HomeroomTeacherID).Scan(
		&c.ID, &c.Name, &c.GradeLevel, &c.AcademicYearID, &createdAt,
	)
	if err != nil {
		return nil, err
	}
	c.CreatedAt = createdAt.Format(time.RFC3339)
	return &c, nil
}

// UpdateClass updates an existing class.
func (r *Repository) UpdateClass(ctx context.Context, schoolID, id string, req ClassRequest) (*ClassResponse, error) {
	var c ClassResponse
	var createdAt time.Time
	err := r.db.QueryRow(ctx, `
		UPDATE classes
		SET name = $3, grade_level = $4, homeroom_teacher_id = $5, updated_at = NOW()
		WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL
		RETURNING id, name, grade_level, academic_year_id, created_at
	`, id, schoolID, req.Name, req.GradeLevel, req.HomeroomTeacherID).Scan(
		&c.ID, &c.Name, &c.GradeLevel, &c.AcademicYearID, &createdAt,
	)
	if err != nil {
		return nil, err
	}
	c.CreatedAt = createdAt.Format(time.RFC3339)
	return &c, nil
}

// DeleteClass soft-deletes a class.
func (r *Repository) DeleteClass(ctx context.Context, schoolID, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE classes SET deleted_at = NOW() WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL
	`, id, schoolID)
	return err
}

// AssignStudents adds students to a class.
func (r *Repository) AssignStudents(ctx context.Context, classID, academicYearID string, studentIDs []string) error {
	for _, sid := range studentIDs {
		_, err := r.db.Exec(ctx, `
			INSERT INTO student_classes (student_id, class_id, academic_year_id)
			VALUES ($1, $2, $3)
			ON CONFLICT (student_id, class_id, academic_year_id) DO NOTHING
		`, sid, classID, academicYearID)
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveStudent removes a student from a class.
func (r *Repository) RemoveStudent(ctx context.Context, classID, studentID string) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM student_classes WHERE class_id = $1 AND student_id = $2
	`, classID, studentID)
	return err
}

// ListStudentsInClass returns all students in a class.
func (r *Repository) ListStudentsInClass(ctx context.Context, classID string) ([]StudentListItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT u.id, u.email, COALESCE(p.parent_name, '') AS name
		FROM student_classes sc
		JOIN users u ON u.id = sc.student_id
		LEFT JOIN user_profiles p ON p.user_id = u.id
		WHERE sc.class_id = $1
		ORDER BY u.email
	`, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []StudentListItem
	for rows.Next() {
		var s StudentListItem
		if err := rows.Scan(&s.ID, &s.Email, &s.Name); err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, nil
}

// ── Subjects ──────────────────────────────────────────────

// ListSubjects returns all subjects for a school.
func (r *Repository) ListSubjects(ctx context.Context, schoolID string) ([]SubjectResponse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, code, description, created_at
		FROM subjects
		WHERE school_id = $1 AND deleted_at IS NULL
		ORDER BY name
	`, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []SubjectResponse
	for rows.Next() {
		var s SubjectResponse
		var createdAt time.Time
		if err := rows.Scan(&s.ID, &s.Name, &s.Code, &s.Description, &createdAt); err != nil {
			return nil, err
		}
		s.CreatedAt = createdAt.Format(time.RFC3339)
		items = append(items, s)
	}
	return items, nil
}

// CreateSubject inserts a new subject.
func (r *Repository) CreateSubject(ctx context.Context, schoolID string, req SubjectRequest) (*SubjectResponse, error) {
	var s SubjectResponse
	var createdAt time.Time
	err := r.db.QueryRow(ctx, `
		INSERT INTO subjects (school_id, name, code, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, code, description, created_at
	`, schoolID, req.Name, req.Code, req.Description).Scan(
		&s.ID, &s.Name, &s.Code, &s.Description, &createdAt,
	)
	if err != nil {
		return nil, err
	}
	s.CreatedAt = createdAt.Format(time.RFC3339)
	return &s, nil
}

// UpdateSubject updates an existing subject.
func (r *Repository) UpdateSubject(ctx context.Context, schoolID, id string, req SubjectRequest) (*SubjectResponse, error) {
	var s SubjectResponse
	var createdAt time.Time
	err := r.db.QueryRow(ctx, `
		UPDATE subjects
		SET name = $3, code = $4, description = $5, updated_at = NOW()
		WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL
		RETURNING id, name, code, description, created_at
	`, id, schoolID, req.Name, req.Code, req.Description).Scan(
		&s.ID, &s.Name, &s.Code, &s.Description, &createdAt,
	)
	if err != nil {
		return nil, err
	}
	s.CreatedAt = createdAt.Format(time.RFC3339)
	return &s, nil
}

// DeleteSubject soft-deletes a subject.
func (r *Repository) DeleteSubject(ctx context.Context, schoolID, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE subjects SET deleted_at = NOW() WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL
	`, id, schoolID)
	return err
}

// ── Class-Subject Assignments ─────────────────────────────

// AssignTeacherToSubject assigns a teacher to teach a subject in a class.
func (r *Repository) AssignTeacherToSubject(ctx context.Context, classID string, req ClassSubjectRequest) (*ClassSubjectResponse, error) {
	var cs ClassSubjectResponse
	err := r.db.QueryRow(ctx, `
		INSERT INTO class_subjects (class_id, subject_id, teacher_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (class_id, subject_id) DO UPDATE SET teacher_id = $3
		RETURNING id, subject_id, teacher_id
	`, classID, req.SubjectID, req.TeacherID).Scan(&cs.ID, &cs.SubjectID, &cs.TeacherID)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// ListClassSubjects returns all subject-teacher assignments for a class.
func (r *Repository) ListClassSubjects(ctx context.Context, classID string) ([]ClassSubjectResponse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT cs.id, cs.subject_id, s.name, cs.teacher_id, u.email
		FROM class_subjects cs
		JOIN subjects s ON s.id = cs.subject_id
		JOIN users u ON u.id = cs.teacher_id
		WHERE cs.class_id = $1
		ORDER BY s.name
	`, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ClassSubjectResponse
	for rows.Next() {
		var cs ClassSubjectResponse
		if err := rows.Scan(&cs.ID, &cs.SubjectID, &cs.SubjectName, &cs.TeacherID, &cs.TeacherEmail); err != nil {
			return nil, err
		}
		items = append(items, cs)
	}
	return items, nil
}

// RemoveClassSubject removes a subject-teacher assignment from a class.
func (r *Repository) RemoveClassSubject(ctx context.Context, classSubjectID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM class_subjects WHERE id = $1`, classSubjectID)
	return err
}
