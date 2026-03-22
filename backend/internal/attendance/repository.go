package attendance

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// MarkBatch inserts or updates attendance records for a date.
func (r *Repository) MarkBatch(ctx context.Context, teacherID string, req MarkAttendanceRequest) error {
	for _, rec := range req.Records {
		_, err := r.db.Exec(ctx, `
			INSERT INTO attendances (student_id, class_id, subject_id, academic_year_id, date, status, note, marked_by)
			VALUES ($1,$2,$3,$4,$5,$6::attendance_status,$7,$8)
			ON CONFLICT (student_id, class_id, subject_id, date)
			DO UPDATE SET status = $6::attendance_status, note = $7, marked_by = $8, updated_at = NOW()
		`, rec.StudentID, req.ClassID, req.SubjectID, req.AcademicYearID, req.Date, rec.Status, rec.Note, teacherID)
		if err != nil {
			return err
		}
	}
	return nil
}

// ListByDate returns attendance for a class on a specific date.
func (r *Repository) ListByDate(ctx context.Context, classID, subjectID, date string) ([]AttendanceResponse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT a.id, a.student_id, u.email, a.class_id, a.subject_id, a.date, a.status, COALESCE(a.note,''), a.marked_by, a.created_at
		FROM attendances a
		JOIN users u ON u.id = a.student_id
		WHERE a.class_id = $1 AND a.subject_id = $2 AND a.date = $3
		ORDER BY u.email
	`, classID, subjectID, date)
	if err != nil { return nil, err }
	defer rows.Close()

	var items []AttendanceResponse
	for rows.Next() {
		var a AttendanceResponse
		if err := rows.Scan(&a.ID, &a.StudentID, &a.Email, &a.ClassID, &a.SubjectID, &a.Date, &a.Status, &a.Note, &a.MarkedBy, &a.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, nil
}

// Summary returns attendance summary per student for a class/subject.
func (r *Repository) Summary(ctx context.Context, classID, subjectID, ayID string) ([]AttendanceSummary, error) {
	rows, err := r.db.Query(ctx, `
		SELECT a.student_id, u.email,
		    COUNT(*) FILTER (WHERE a.status = 'present') AS present,
		    COUNT(*) FILTER (WHERE a.status = 'absent') AS absent,
		    COUNT(*) FILTER (WHERE a.status = 'late') AS late,
		    COUNT(*) FILTER (WHERE a.status = 'excused') AS excused,
		    COUNT(*) AS total
		FROM attendances a
		JOIN users u ON u.id = a.student_id
		WHERE a.class_id = $1 AND a.subject_id = $2 AND a.academic_year_id = $3
		GROUP BY a.student_id, u.email
		ORDER BY u.email
	`, classID, subjectID, ayID)
	if err != nil { return nil, err }
	defer rows.Close()

	var items []AttendanceSummary
	for rows.Next() {
		var s AttendanceSummary
		if err := rows.Scan(&s.StudentID, &s.Email, &s.Present, &s.Absent, &s.Late, &s.Excused, &s.Total); err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, nil
}

// StudentHistory returns a student's attendance history.
func (r *Repository) StudentHistory(ctx context.Context, studentID, classID, subjectID string) ([]AttendanceResponse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT a.id, a.student_id, '', a.class_id, a.subject_id, a.date, a.status, COALESCE(a.note,''), a.marked_by, a.created_at
		FROM attendances a
		WHERE a.student_id = $1 AND a.class_id = $2 AND a.subject_id = $3
		ORDER BY a.date DESC
	`, studentID, classID, subjectID)
	if err != nil { return nil, err }
	defer rows.Close()

	var items []AttendanceResponse
	for rows.Next() {
		var a AttendanceResponse
		if err := rows.Scan(&a.ID, &a.StudentID, &a.Email, &a.ClassID, &a.SubjectID, &a.Date, &a.Status, &a.Note, &a.MarkedBy, &a.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, nil
}
