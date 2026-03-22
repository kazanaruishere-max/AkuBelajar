package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// GeminiClient calls the Google Gemini API.
type GeminiClient struct {
	apiKey string
	model  string
}

// NewGeminiClient creates a new Gemini client.
func NewGeminiClient() *GeminiClient {
	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = "gemini-2.0-flash"
	}
	return &GeminiClient{
		apiKey: os.Getenv("GEMINI_API_KEY"),
		model:  model,
	}
}

// IsAvailable checks if the API key is configured.
func (g *GeminiClient) IsAvailable() bool {
	return g.apiKey != ""
}

// GenerateContent sends a prompt and returns the generated text.
func (g *GeminiClient) GenerateContent(prompt string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", g.model, g.apiKey)

	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": prompt}}},
		},
		"generationConfig": map[string]interface{}{
			"temperature":   0.7,
			"maxOutputTokens": 4096,
		},
	}

	jsonBody, _ := json.Marshal(body)
	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("gagal memanggil Gemini API: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("gagal parsing response: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("tidak ada response dari Gemini")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
