package assignment

import (
	"context"
	"math"
	"time"
)

// Service handles assignment business logic.
type Service struct {
	repo *Repository
}

// NewService creates a new assignment service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new assignment (teacher only).
func (s *Service) Create(ctx context.Context, teacherID string, req CreateAssignmentRequest) (*AssignmentResponse, error) {
	return s.repo.Create(ctx, teacherID, req)
}

// Publish publishes a draft assignment.
func (s *Service) Publish(ctx context.Context, id, teacherID string) error {
	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if a == nil {
		return ErrNotFound
	}
	if a.TeacherID != teacherID {
		return ErrNotAuthorized
	}
	return s.repo.UpdateStatus(ctx, id, "published")
}

// Close closes an assignment (no more submissions).
func (s *Service) Close(ctx context.Context, id, teacherID string) error {
	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if a == nil {
		return ErrNotFound
	}
	if a.TeacherID != teacherID {
		return ErrNotAuthorized
	}
	return s.repo.UpdateStatus(ctx, id, "closed")
}

// Submit submits a student's work for an assignment.
func (s *Service) Submit(ctx context.Context, assignmentID, studentID string) (*SubmissionResponse, error) {
	a, err := s.repo.GetByID(ctx, assignmentID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, ErrNotFound
	}
	if a.Status != "published" {
		return nil, ErrNotPublished
	}

	// Calculate late status
	now := time.Now()
	isLate := now.After(a.DeadlineAt)
	lateDays := 0

	if isLate {
		if !a.AllowLate {
			return nil, ErrDeadlinePassed
		}
		lateDays = int(math.Ceil(now.Sub(a.DeadlineAt).Hours() / 24))
		if lateDays > a.MaxLateDays {
			return nil, ErrDeadlinePassed
		}
	}

	// Create or get existing submission
	subID, err := s.repo.CreateSubmission(ctx, assignmentID, studentID)
	if err != nil {
		return nil, err
	}

	// Mark as submitted
	if err := s.repo.SubmitSubmission(ctx, subID, isLate, lateDays); err != nil {
		return nil, err
	}

	return s.repo.GetSubmission(ctx, assignmentID, studentID)
}

// Grade grades a submission with late penalty calculation.
func (s *Service) Grade(ctx context.Context, submissionID, assignmentID, teacherID string, req GradeSubmissionRequest) error {
	a, err := s.repo.GetByID(ctx, assignmentID)
	if err != nil {
		return err
	}
	if a == nil {
		return ErrNotFound
	}
	if a.TeacherID != teacherID {
		return ErrNotAuthorized
	}

	// Get submission to check late status
	sub, err := s.repo.GetSubmission(ctx, assignmentID, "")
	if err != nil {
		return err
	}

	gradeAfterPenalty := req.Grade
	if sub != nil && sub.IsLate && sub.LateDays > 0 {
		penalty := sub.LateDays * a.LatePenaltyPct
		gradeAfterPenalty = req.Grade - (req.Grade * penalty / 100)
		if gradeAfterPenalty < 0 {
			gradeAfterPenalty = 0
		}
	}

	return s.repo.GradeSubmission(ctx, submissionID, teacherID, req.Grade, gradeAfterPenalty, req.Feedback)
}
