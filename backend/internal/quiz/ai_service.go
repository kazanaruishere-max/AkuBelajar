package quiz

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kazanaruishere-max/akubelajar/backend/pkg/ai"
)

// AIService generates quiz questions using AI.
type AIService struct {
	gemini *ai.GeminiClient
	repo   *Repository
}

// NewAIService creates a new AI quiz service.
func NewAIService(gemini *ai.GeminiClient, repo *Repository) *AIService {
	return &AIService{gemini: gemini, repo: repo}
}

// GenerateQuestionsRequest is the payload to generate questions.
type GenerateQuestionsRequest struct {
	QuizID      string `json:"quiz_id" validate:"required,uuid"`
	Topic       string `json:"topic" validate:"required,max=200"`
	Count       int    `json:"count" validate:"required,min=1,max=20"`
	Difficulty  string `json:"difficulty" validate:"required,oneof=easy medium hard"`
	Language    string `json:"language"`
}

// GenerateQuestions generates quiz questions using AI and saves them.
func (s *AIService) GenerateQuestions(ctx context.Context, req GenerateQuestionsRequest) (int, error) {
	if !s.gemini.IsAvailable() {
		return 0, fmt.Errorf("GEMINI_API_KEY belum dikonfigurasi")
	}

	lang := req.Language
	if lang == "" { lang = "Bahasa Indonesia" }

	prompt := fmt.Sprintf(`Buatkan %d soal pilihan ganda tentang "%s" dengan tingkat kesulitan %s.
Gunakan bahasa: %s.

Format respons HARUS berupa JSON array tanpa markdown, seperti ini:
[
  {
    "question_text": "Pertanyaan di sini?",
    "options": [
      {"key": "A", "text": "Opsi A", "is_correct": false},
      {"key": "B", "text": "Opsi B", "is_correct": true},
      {"key": "C", "text": "Opsi C", "is_correct": false},
      {"key": "D", "text": "Opsi D", "is_correct": false}
    ],
    "explanation": "Penjelasan jawaban benar"
  }
]

PENTING: Hanya kembalikan JSON array, tanpa teks tambahan atau markdown.`, req.Count, req.Topic, req.Difficulty, lang)

	text, err := s.gemini.GenerateContent(prompt)
	if err != nil {
		return 0, fmt.Errorf("gagal generate soal: %w", err)
	}

	// Clean response (remove markdown code blocks if present)
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	// Parse questions
	var questions []struct {
		QuestionText string `json:"question_text"`
		Options      []struct {
			Key       string `json:"key"`
			Text      string `json:"text"`
			IsCorrect bool   `json:"is_correct"`
		} `json:"options"`
		Explanation string `json:"explanation"`
	}

	if err := json.Unmarshal([]byte(text), &questions); err != nil {
		return 0, fmt.Errorf("gagal parsing respons AI: %w (response: %s)", err, text[:min(200, len(text))])
	}

	// Save questions to database
	saved := 0
	for i, q := range questions {
		optionsJSON, _ := json.Marshal(q.Options)
		explanation := q.Explanation
		_, err := s.repo.AddQuestion(ctx, req.QuizID, QuestionRequest{
			QuestionText: q.QuestionText,
			QuestionType: "multiple_choice",
			Options:      optionsJSON,
			Explanation:  &explanation,
			OrderNum:     i + 1,
		})
		if err == nil { saved++ }
	}

	return saved, nil
}

func min(a, b int) int {
	if a < b { return a }
	return b
}
