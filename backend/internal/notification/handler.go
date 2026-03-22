package notification

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/response"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/validator"
)

type Handler struct {
	repo      *Repository
	validator *validator.Validator
}

func NewHandler(repo *Repository, v *validator.Validator) *Handler {
	return &Handler{repo: repo, validator: v}
}

func uid(c *gin.Context) string { v, _ := c.Get("user_id"); return v.(string) }

// ListNotifications returns notifications for the current user.
func (h *Handler) ListNotifications(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	items, err := h.repo.ListByUser(c.Request.Context(), uid(c), limit)
	if err != nil { response.InternalError(c, "NOTIF_010", "Gagal mengambil notifikasi"); return }
	response.OK(c, items)
}

// UnreadCount returns the count of unread notifications.
func (h *Handler) UnreadCount(c *gin.Context) {
	count, err := h.repo.UnreadCount(c.Request.Context(), uid(c))
	if err != nil { response.InternalError(c, "NOTIF_010", "Gagal mengambil jumlah notifikasi"); return }
	response.OK(c, UnreadCount{Count: count})
}

// MarkRead marks a single notification as read.
func (h *Handler) MarkRead(c *gin.Context) {
	if err := h.repo.MarkRead(c.Request.Context(), uid(c), c.Param("id")); err != nil {
		response.InternalError(c, "NOTIF_010", "Gagal menandai notifikasi"); return
	}
	response.OK(c, gin.H{"message": "Notifikasi ditandai sudah dibaca"})
}

// MarkAllRead marks all notifications as read.
func (h *Handler) MarkAllRead(c *gin.Context) {
	if err := h.repo.MarkAllRead(c.Request.Context(), uid(c)); err != nil {
		response.InternalError(c, "NOTIF_010", "Gagal menandai semua notifikasi"); return
	}
	response.OK(c, gin.H{"message": "Semua notifikasi ditandai sudah dibaca"})
}

// Send creates a notification (teacher/admin use).
func (h *Handler) Send(c *gin.Context) {
	var req CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "NOTIF_001", "Format request tidak valid"); return }
	if errs := h.validator.ValidateStruct(req); errs != nil { response.BadRequestWithDetails(c, "NOTIF_002", "Validasi gagal", errs); return }
	if err := h.repo.Create(c.Request.Context(), req); err != nil {
		response.InternalError(c, "NOTIF_010", "Gagal mengirim notifikasi"); return
	}
	response.Created(c, gin.H{"message": "Notifikasi berhasil dikirim"})
}

// Broadcast sends notifications to multiple users.
func (h *Handler) Broadcast(c *gin.Context) {
	var req BroadcastRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "NOTIF_001", "Format request tidak valid"); return }
	if errs := h.validator.ValidateStruct(req); errs != nil { response.BadRequestWithDetails(c, "NOTIF_002", "Validasi gagal", errs); return }
	if err := h.repo.Broadcast(c.Request.Context(), req); err != nil {
		response.InternalError(c, "NOTIF_010", "Gagal mengirim broadcast"); return
	}
	response.Created(c, gin.H{"message": "Broadcast berhasil dikirim"})
}
