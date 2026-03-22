package upload

import (
	"github.com/gin-gonic/gin"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/response"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/storage"
)

type Handler struct {
	store *storage.SupabaseStorage
}

func NewHandler(store *storage.SupabaseStorage) *Handler {
	return &Handler{store: store}
}

// UploadFile handles file uploads and returns the public URL.
func (h *Handler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "UPLOAD_001", "File tidak ditemukan dalam request")
		return
	}
	defer file.Close()

	folder := c.DefaultPostForm("folder", "uploads")

	// Validate
	if err := storage.ValidateFile(header, 20, nil); err != nil {
		response.BadRequest(c, "UPLOAD_002", err.Error())
		return
	}

	url, err := h.store.Upload(file, header, folder)
	if err != nil {
		response.InternalError(c, "UPLOAD_010", "Gagal mengupload file: "+err.Error())
		return
	}

	response.OK(c, gin.H{
		"url":       url,
		"file_name": header.Filename,
		"file_size": header.Size,
	})
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc) {
	rg.POST("/upload", authMW, h.UploadFile)
}
