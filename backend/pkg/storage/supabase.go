package storage

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SupabaseStorage handles file uploads to Supabase Storage.
type SupabaseStorage struct {
	projectURL string
	serviceKey string
	bucket     string
}

// NewSupabaseStorage creates a new Supabase Storage client.
func NewSupabaseStorage() *SupabaseStorage {
	return &SupabaseStorage{
		projectURL: os.Getenv("SUPABASE_URL"),
		serviceKey: os.Getenv("SUPABASE_SERVICE_KEY"),
		bucket:     os.Getenv("SUPABASE_BUCKET"),
	}
}

// Upload uploads a file to Supabase Storage and returns the public URL.
func (s *SupabaseStorage) Upload(file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	uniqueName := fmt.Sprintf("%s/%s_%s%s", folder, time.Now().Format("20060102"), uuid.New().String()[:8], ext)

	// Read file content
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, file); err != nil {
		return "", fmt.Errorf("gagal membaca file: %w", err)
	}

	// Upload to Supabase Storage
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.projectURL, s.bucket, uniqueName)
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return "", fmt.Errorf("gagal membuat request: %w", err)
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+s.serviceKey)
	req.Header.Set("x-upsert", "true")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal upload ke storage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload gagal (status %d): %s", resp.StatusCode, string(body))
	}

	// Return public URL
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.projectURL, s.bucket, uniqueName)
	return publicURL, nil
}

// ValidateFile checks file size and extension.
func ValidateFile(header *multipart.FileHeader, maxSizeMB int, allowedExts []string) error {
	// Check size
	if header.Size > int64(maxSizeMB)*1024*1024 {
		return fmt.Errorf("ukuran file melebihi %dMB", maxSizeMB)
	}

	// Check extension
	if len(allowedExts) > 0 {
		ext := strings.ToLower(filepath.Ext(header.Filename))
		allowed := false
		for _, e := range allowedExts {
			if ext == "."+strings.ToLower(e) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("ekstensi file %s tidak diizinkan", ext)
		}
	}

	return nil
}
