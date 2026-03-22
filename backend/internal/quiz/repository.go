package quiz

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// ── Quiz CRUD ─────────────────────────────────────────────

func (r *Repository) ListByTeacher(ctx context.Context, teacherID string) ([]QuizResponse, error) {
	return r.queryQuizzes(ctx, `WHERE q.teacher_id = $1 AND q.deleted_at IS NULL ORDER BY q.created_at DESC`, teacherID)
}

func (r *Repository) ListByClass(ctx context.Context, classID string) ([]QuizResponse, error) {
	return r.queryQuizzes(ctx, `WHERE q.class_id = $1 AND q.status IN ('published','active','ended','graded') AND q.deleted_at IS NULL ORDER BY q.created_at DESC`, classID)
}

func (r *Repository) GetByID(ctx context.Context, id string) (*QuizResponse, error) {
	items, err := r.queryQuizzes(ctx, `WHERE q.id = $1 AND q.deleted_at IS NULL`, id)
	if err != nil { return nil, err }
	if len(items) == 0 { return nil, nil }
	return &items[0], nil
}

func (r *Repository) queryQuizzes(ctx context.Context, where string, args ...interface{}) ([]QuizResponse, error) {
	q := `SELECT q.id, q.class_id, COALESCE(cl.name,''), q.subject_id, COALESCE(s.name,''),
	       q.teacher_id, q.title, q.time_limit, q.randomize_questions, q.randomize_options,
	       q.max_attempts, q.allow_review, q.status,
	       (SELECT COUNT(*) FROM quiz_questions qq WHERE qq.quiz_id = q.id),
	       (SELECT COUNT(*) FROM quiz_sessions qs WHERE qs.quiz_id = q.id),
	       q.start_at, q.end_at, q.created_at
	FROM quizzes q
	LEFT JOIN classes cl ON cl.id = q.class_id
	LEFT JOIN subjects s ON s.id = q.subject_id ` + where

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var items []QuizResponse
	for rows.Next() {
		var qz QuizResponse
		if err := rows.Scan(&qz.ID, &qz.ClassID, &qz.ClassName, &qz.SubjectID, &qz.SubjectName,
			&qz.TeacherID, &qz.Title, &qz.TimeLimit, &qz.RandomizeQuestions, &qz.RandomizeOptions,
			&qz.MaxAttempts, &qz.AllowReview, &qz.Status, &qz.QuestionCount, &qz.SessionCount,
			&qz.StartAt, &qz.EndAt, &qz.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, qz)
	}
	return items, nil
}

func (r *Repository) Create(ctx context.Context, teacherID string, req CreateQuizRequest) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		INSERT INTO quizzes (class_id, subject_id, teacher_id, title, time_limit,
		    randomize_questions, randomize_options, max_attempts, allow_review, start_at, end_at, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,NULLIF($10,'')::TIMESTAMPTZ,NULLIF($11,'')::TIMESTAMPTZ,'draft')
		RETURNING id
	`, req.ClassID, req.SubjectID, teacherID, req.Title, req.TimeLimit,
		req.RandomizeQuestions, req.RandomizeOptions, req.MaxAttempts, req.AllowReview,
		req.StartAt, req.EndAt).Scan(&id)
	return id, err
}

func (r *Repository) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Exec(ctx, `UPDATE quizzes SET status = $2::quiz_status, updated_at = NOW() WHERE id = $1`, id, status)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `UPDATE quizzes SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}

// ── Questions ─────────────────────────────────────────────

func (r *Repository) AddQuestion(ctx context.Context, quizID string, req QuestionRequest) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		INSERT INTO quiz_questions (quiz_id, question_text, question_type, options, answer_hash, explanation, order_num)
		VALUES ($1,$2,$3::question_type,$4,$5,$6,$7) RETURNING id
	`, quizID, req.QuestionText, req.QuestionType, req.Options, req.AnswerHash, req.Explanation, req.OrderNum).Scan(&id)
	return id, err
}

func (r *Repository) ListQuestions(ctx context.Context, quizID string) ([]QuestionResponse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, question_text, question_type, options, explanation, order_num
		FROM quiz_questions WHERE quiz_id = $1 ORDER BY order_num
	`, quizID)
	if err != nil { return nil, err }
	defer rows.Close()

	var items []QuestionResponse
	for rows.Next() {
		var q QuestionResponse
		if err := rows.Scan(&q.ID, &q.QuestionText, &q.QuestionType, &q.Options, &q.Explanation, &q.OrderNum); err != nil {
			return nil, err
		}
		items = append(items, q)
	}
	return items, nil
}

// ── Sessions ──────────────────────────────────────────────

func (r *Repository) CountStudentSessions(ctx context.Context, quizID, studentID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM quiz_sessions WHERE quiz_id = $1 AND student_id = $2`, quizID, studentID).Scan(&count)
	return count, err
}

func (r *Repository) CreateSession(ctx context.Context, quizID, studentID, ip, ua string, expiresAt time.Time, questionOrder []int) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		INSERT INTO quiz_sessions (quiz_id, student_id, expires_at, ip_address, user_agent, question_order)
		VALUES ($1,$2,$3,$4::INET,$5,$6) RETURNING id
	`, quizID, studentID, expiresAt, ip, ua, questionOrder).Scan(&id)
	return id, err
}

func (r *Repository) GetSession(ctx context.Context, sessionID string) (*SessionResponse, error) {
	var s SessionResponse
	err := r.db.QueryRow(ctx, `
		SELECT id, quiz_id, student_id, started_at, submitted_at, expires_at, status, score, cheat_count, created_at
		FROM quiz_sessions WHERE id = $1
	`, sessionID).Scan(&s.ID, &s.QuizID, &s.StudentID, &s.StartedAt, &s.SubmittedAt, &s.ExpiresAt, &s.Status, &s.Score, &s.CheatCount, &s.CreatedAt)
	if err == pgx.ErrNoRows { return nil, nil }
	if err != nil { return nil, err }
	return &s, nil
}

func (r *Repository) ListSessions(ctx context.Context, quizID string) ([]SessionResponse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT qs.id, qs.quiz_id, qs.student_id, u.email, qs.started_at, qs.submitted_at,
		       qs.expires_at, qs.status, qs.score, qs.cheat_count, qs.created_at
		FROM quiz_sessions qs JOIN users u ON u.id = qs.student_id
		WHERE qs.quiz_id = $1 ORDER BY u.email
	`, quizID)
	if err != nil { return nil, err }
	defer rows.Close()

	var items []SessionResponse
	for rows.Next() {
		var s SessionResponse
		if err := rows.Scan(&s.ID, &s.QuizID, &s.StudentID, &s.StudentEmail, &s.StartedAt, &s.SubmittedAt,
			&s.ExpiresAt, &s.Status, &s.Score, &s.CheatCount, &s.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, nil
}

func (r *Repository) SubmitSession(ctx context.Context, sessionID string, score int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE quiz_sessions SET submitted_at = NOW(), status = 'submitted', score = $2, updated_at = NOW()
		WHERE id = $1
	`, sessionID, score)
	return err
}

// ── Answers ───────────────────────────────────────────────

func (r *Repository) SaveAnswer(ctx context.Context, sessionID string, req AnswerRequest, isCorrect *bool) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO quiz_answers (session_id, question_id, selected_key, essay_text, is_correct, answered_at)
		VALUES ($1,$2,$3,$4,$5,NOW())
		ON CONFLICT (session_id, question_id) DO UPDATE
		SET selected_key = $3, essay_text = $4, is_correct = $5, answered_at = NOW()
	`, sessionID, req.QuestionID, req.SelectedKey, req.EssayText, isCorrect)
	return err
}

func (r *Repository) GetAnswers(ctx context.Context, sessionID string) ([]AnswerResponse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, question_id, selected_key, essay_text, is_correct
		FROM quiz_answers WHERE session_id = $1
	`, sessionID)
	if err != nil { return nil, err }
	defer rows.Close()

	var items []AnswerResponse
	for rows.Next() {
		var a AnswerResponse
		if err := rows.Scan(&a.ID, &a.QuestionID, &a.SelectedKey, &a.EssayText, &a.IsCorrect); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, nil
}

func (r *Repository) CountCorrectAnswers(ctx context.Context, sessionID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM quiz_answers WHERE session_id = $1 AND is_correct = TRUE`, sessionID).Scan(&count)
	return count, err
}
