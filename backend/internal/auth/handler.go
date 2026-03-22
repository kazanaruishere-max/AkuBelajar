package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/response"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/validator"
)

// Handler handles auth HTTP requests.
type Handler struct {
	service   *Service
	validator *validator.Validator
}

// NewHandler creates a new auth handler.
func NewHandler(service *Service, v *validator.Validator) *Handler {
	return &Handler{service: service, validator: v}
}

// Login handles POST /auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "AUTH_001", "Format request tidak valid")
		return
	}

	// Validate
	if errs := h.validator.ValidateStruct(req); errs != nil {
		response.BadRequestWithDetails(c, "AUTH_002", "Validasi gagal", errs)
		return
	}

	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	result, err := h.service.Login(c.Request.Context(), req, ip, userAgent)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	response.OK(c, result)
}

// RefreshToken handles POST /auth/refresh
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "AUTH_001", "Format request tidak valid")
		return
	}

	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	result, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken, ip, userAgent)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	response.OK(c, result)
}

// Logout handles POST /auth/logout
func (h *Handler) Logout(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "AUTH_001", "Format request tidak valid")
		return
	}

	if err := h.service.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		response.InternalError(c, "AUTH_010", "Gagal logout")
		return
	}

	response.NoContent(c)
}

// ChangePassword handles POST /auth/change-password
func (h *Handler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "AUTH_001", "Format request tidak valid")
		return
	}

	if errs := h.validator.ValidateStruct(req); errs != nil {
		response.BadRequestWithDetails(c, "AUTH_002", "Validasi gagal", errs)
		return
	}

	// Get user ID from middleware context
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "AUTH_003", "Tidak terautentikasi")
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), userID.(string), req); err != nil {
		h.handleAuthError(c, err)
		return
	}

	response.OK(c, gin.H{"message": "Password berhasil diubah. Silakan login kembali."})
}

// GetMe handles GET /auth/me
func (h *Handler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "AUTH_003", "Tidak terautentikasi")
		return
	}

	result, err := h.service.GetMe(c.Request.Context(), userID.(string))
	if err != nil {
		response.InternalError(c, "AUTH_010", "Gagal mengambil data user")
		return
	}

	response.OK(c, result)
}

// handleAuthError maps domain errors to HTTP responses.
func (h *Handler) handleAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{"code": "AUTH_004", "message": err.Error()},
		})
	case errors.Is(err, ErrAccountInactive):
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{"code": "AUTH_005", "message": err.Error()},
		})
	case errors.Is(err, ErrAccountLocked):
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": gin.H{"code": "AUTH_006", "message": err.Error()},
		})
	case errors.Is(err, ErrInvalidRefreshToken):
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{"code": "AUTH_007", "message": err.Error()},
		})
	default:
		response.InternalError(c, "AUTH_010", "Terjadi kesalahan server")
	}
}
