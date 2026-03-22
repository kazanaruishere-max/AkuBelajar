package grade

import (
	"context"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Insert adds a new grade entry.
func (r *Repository) Insert(ctx context.Context, req InsertGradeRequest) error {
	weightedScore := float64(req.Score) * float64(req.WeightPct) / 100.0
	_, err := r.db.Exec(ctx, `
		INSERT INTO grades (student_id, class_id, subject_id, academic_year_id, category, label, score, weight_pct, weighted_score)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, req.StudentID, req.ClassID, req.SubjectID, req.AcademicYearID, req.Category, req.Label, req.Score, req.WeightPct, weightedScore)
	return err
}

// ListByStudent returns all grades for a student in a class/subject.
func (r *Repository) ListByStudent(ctx context.Context, studentID, classID, subjectID string) ([]GradeResponse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT g.id, g.student_id, '', g.class_id, g.subject_id, COALESCE(s.name,''),
		       g.academic_year_id, g.category, g.label, g.score, g.weight_pct, g.weighted_score, g.created_at
		FROM grades g
		LEFT JOIN subjects s ON s.id = g.subject_id
		WHERE g.student_id = $1 AND g.class_id = $2 AND g.subject_id = $3
		ORDER BY g.created_at DESC
	`, studentID, classID, subjectID)
	if err != nil { return nil, err }
	defer rows.Close()

	var items []GradeResponse
	for rows.Next() {
		var g GradeResponse
		if err := rows.Scan(&g.ID, &g.StudentID, &g.StudentEmail, &g.ClassID, &g.SubjectID, &g.SubjectName,
			&g.AcademicYearID, &g.Category, &g.Label, &g.Score, &g.WeightPct, &g.WeightedScore, &g.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, g)
	}
	return items, nil
}

// Summary returns weighted averages per student per subject.
func (r *Repository) Summary(ctx context.Context, classID, ayID string) ([]GradeSummary, error) {
	rows, err := r.db.Query(ctx, `
		SELECT g.student_id, u.email, g.subject_id, s.name,
		    AVG(g.score)::FLOAT AS average,
		    SUM(g.weighted_score) / NULLIF(SUM(g.weight_pct),0) * 100 AS weighted_avg
		FROM grades g
		JOIN users u ON u.id = g.student_id
		JOIN subjects s ON s.id = g.subject_id
		WHERE g.class_id = $1 AND g.academic_year_id = $2
		GROUP BY g.student_id, u.email, g.subject_id, s.name
		ORDER BY s.name, u.email
	`, classID, ayID)
	if err != nil { return nil, err }
	defer rows.Close()

	var items []GradeSummary
	for rows.Next() {
		var gs GradeSummary
		if err := rows.Scan(&gs.StudentID, &gs.StudentEmail, &gs.SubjectID, &gs.SubjectName, &gs.Average, &gs.WeightedAvg); err != nil {
			return nil, err
		}
		gs.Average = math.Round(gs.Average*100) / 100
		gs.WeightedAvg = math.Round(gs.WeightedAvg*100) / 100
		gs.LetterGrade = letterGrade(gs.WeightedAvg)
		items = append(items, gs)
	}
	return items, nil
}

func letterGrade(score float64) string {
	switch {
	case score >= 90: return "A"
	case score >= 80: return "B"
	case score >= 70: return "C"
	case score >= 60: return "D"
	default:          return "E"
	}
}
