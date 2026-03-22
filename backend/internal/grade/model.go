package grade

import "time"

type GradeResponse struct {
	ID             string    `json:"id"`
	StudentID      string    `json:"student_id"`
	StudentEmail   string    `json:"student_email,omitempty"`
	ClassID        string    `json:"class_id"`
	SubjectID      string    `json:"subject_id"`
	SubjectName    string    `json:"subject_name,omitempty"`
	AcademicYearID string    `json:"academic_year_id"`
	Category       string    `json:"category"` // assignment, quiz, exam, midterm, final
	Label          string    `json:"label"`
	Score          int       `json:"score"`
	WeightPct      int       `json:"weight_pct"`
	WeightedScore  float64   `json:"weighted_score"`
	CreatedAt      time.Time `json:"created_at"`
}

type GradeSummary struct {
	StudentID    string  `json:"student_id"`
	StudentEmail string  `json:"student_email"`
	SubjectID    string  `json:"subject_id"`
	SubjectName  string  `json:"subject_name"`
	Average      float64 `json:"average"`
	WeightedAvg  float64 `json:"weighted_average"`
	LetterGrade  string  `json:"letter_grade"`
}

type InsertGradeRequest struct {
	StudentID      string `json:"student_id" validate:"required,uuid"`
	ClassID        string `json:"class_id" validate:"required,uuid"`
	SubjectID      string `json:"subject_id" validate:"required,uuid"`
	AcademicYearID string `json:"academic_year_id" validate:"required,uuid"`
	Category       string `json:"category" validate:"required,oneof=assignment quiz exam midterm final"`
	Label          string `json:"label" validate:"required,max=100"`
	Score          int    `json:"score" validate:"required,min=0,max=100"`
	WeightPct      int    `json:"weight_pct" validate:"required,min=1,max=100"`
}
