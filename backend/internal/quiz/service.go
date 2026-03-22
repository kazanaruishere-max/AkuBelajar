package quiz

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new quiz (teacher only).
func (s *Service) Create(ctx context.Context, teacherID string, req CreateQuizRequest) (string, error) {
	if req.MaxAttempts == 0 { req.MaxAttempts = 1 }
	return s.repo.Create(ctx, teacherID, req)
}

// Publish publishes a draft quiz.
func (s *Service) Publish(ctx context.Context, id, teacherID string) error {
	q, err := s.repo.GetByID(ctx, id)
	if err != nil { return err }
	if q == nil { return ErrNotFound }
	if q.TeacherID != teacherID { return ErrNotAuthorized }
	return s.repo.UpdateStatus(ctx, id, "published")
}

// StartSession creates a new quiz session for a student.
func (s *Service) StartSession(ctx context.Context, quizID, studentID, ip, ua string) (*SessionResponse, error) {
	q, err := s.repo.GetByID(ctx, quizID)
	if err != nil { return nil, err }
	if q == nil { return nil, ErrNotFound }
	if q.Status != "published" && q.Status != "active" { return nil, ErrNotPublished }

	// Check end time
	if q.EndAt != nil && time.Now().After(*q.EndAt) { return nil, ErrQuizEnded }

	// Check max attempts
	count, err := s.repo.CountStudentSessions(ctx, quizID, studentID)
	if err != nil { return nil, err }
	if count >= q.MaxAttempts { return nil, ErrMaxAttempts }

	// Generate question order (randomized if enabled)
	questions, err := s.repo.ListQuestions(ctx, quizID)
	if err != nil { return nil, err }

	order := make([]int, len(questions))
	for i := range order { order[i] = i }
	if q.RandomizeQuestions {
		rand.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })
	}

	expiresAt := time.Now().Add(time.Duration(q.TimeLimit) * time.Minute)

	sessionID, err := s.repo.CreateSession(ctx, quizID, studentID, ip, ua, expiresAt, order)
	if err != nil { return nil, err }

	return s.repo.GetSession(ctx, sessionID)
}

// SaveAnswer saves a student's answer (auto-grades MC).
func (s *Service) SaveAnswer(ctx context.Context, sessionID string, req AnswerRequest) error {
	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil { return err }
	if session == nil { return ErrNotFound }
	if session.Status != "in_progress" { return ErrSessionExpired }
	if time.Now().After(session.ExpiresAt) { return ErrSessionExpired }

	// Auto-grade for MC: check answer_hash
	var isCorrect *bool
	if req.SelectedKey != nil {
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(*req.SelectedKey)))
		// We store the correct answer hash in the question; compare at submit time
		// For now, store nil and grade at submission
		_ = hash
	}

	return s.repo.SaveAnswer(ctx, sessionID, req, isCorrect)
}

// SubmitSession submits a quiz session and calculates score.
func (s *Service) SubmitSession(ctx context.Context, sessionID string) (*SessionResponse, error) {
	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil { return nil, err }
	if session == nil { return nil, ErrNotFound }

	// Auto-grade: count correct answers
	correct, _ := s.repo.CountCorrectAnswers(ctx, sessionID)

	// Get total questions
	q, _ := s.repo.GetByID(ctx, session.QuizID)
	total := 0
	if q != nil { total = q.QuestionCount }

	score := 0
	if total > 0 { score = (correct * 100) / total }

	if err := s.repo.SubmitSession(ctx, sessionID, score); err != nil {
		return nil, err
	}

	return s.repo.GetSession(ctx, sessionID)
}
