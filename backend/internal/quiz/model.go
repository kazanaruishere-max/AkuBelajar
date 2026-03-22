package quiz

import "time"

// ── Request DTOs ──────────────────────────────────────────

type CreateQuizRequest struct {
	ClassID            string `json:"class_id" validate:"required,uuid"`
	SubjectID          string `json:"subject_id" validate:"required,uuid"`
	Title              string `json:"title" validate:"required,max=200"`
	TimeLimit          int    `json:"time_limit" validate:"required,min=1,max=300"` // minutes
	RandomizeQuestions bool   `json:"randomize_questions"`
	RandomizeOptions   bool   `json:"randomize_options"`
	MaxAttempts        int    `json:"max_attempts" validate:"min=1,max=5"`
	AllowReview        bool   `json:"allow_review"`
	StartAt            string `json:"start_at"`
	EndAt              string `json:"end_at"`
}

type QuestionRequest struct {
	QuestionText string  `json:"question_text" validate:"required"`
	QuestionType string  `json:"question_type" validate:"required,oneof=multiple_choice essay"`
	Options      []byte  `json:"options"`      // JSONB: [{key:"A",text:"...",is_correct:true}]
	AnswerHash   *string `json:"answer_hash"`  // correct answer key hash
	Explanation  *string `json:"explanation"`
	OrderNum     int     `json:"order_num" validate:"required,min=1"`
}

type AnswerRequest struct {
	QuestionID  string  `json:"question_id" validate:"required,uuid"`
	SelectedKey *string `json:"selected_key"` // for MC
	EssayText   *string `json:"essay_text"`   // for essay
}

// ── Response DTOs ─────────────────────────────────────────

type QuizResponse struct {
	ID                 string     `json:"id"`
	ClassID            string     `json:"class_id"`
	ClassName          string     `json:"class_name,omitempty"`
	SubjectID          string     `json:"subject_id"`
	SubjectName        string     `json:"subject_name,omitempty"`
	TeacherID          string     `json:"teacher_id"`
	Title              string     `json:"title"`
	TimeLimit          int        `json:"time_limit"`
	RandomizeQuestions bool       `json:"randomize_questions"`
	RandomizeOptions   bool       `json:"randomize_options"`
	MaxAttempts        int        `json:"max_attempts"`
	AllowReview        bool       `json:"allow_review"`
	Status             string     `json:"status"`
	QuestionCount      int        `json:"question_count"`
	SessionCount       int        `json:"session_count"`
	StartAt            *time.Time `json:"start_at,omitempty"`
	EndAt              *time.Time `json:"end_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
}

type QuestionResponse struct {
	ID           string  `json:"id"`
	QuestionText string  `json:"question_text"`
	QuestionType string  `json:"question_type"`
	Options      []byte  `json:"options,omitempty"`
	Explanation  *string `json:"explanation,omitempty"`
	OrderNum     int     `json:"order_num"`
}

type SessionResponse struct {
	ID           string     `json:"id"`
	QuizID       string     `json:"quiz_id"`
	StudentID    string     `json:"student_id"`
	StudentEmail string     `json:"student_email,omitempty"`
	StartedAt    time.Time  `json:"started_at"`
	SubmittedAt  *time.Time `json:"submitted_at,omitempty"`
	ExpiresAt    time.Time  `json:"expires_at"`
	Status       string     `json:"status"`
	Score        *int       `json:"score,omitempty"`
	CheatCount   int        `json:"cheat_count"`
	CreatedAt    time.Time  `json:"created_at"`
}

type AnswerResponse struct {
	ID          string  `json:"id"`
	QuestionID  string  `json:"question_id"`
	SelectedKey *string `json:"selected_key,omitempty"`
	EssayText   *string `json:"essay_text,omitempty"`
	IsCorrect   *bool   `json:"is_correct,omitempty"`
}
