package notification

import "time"

type CreateNotificationRequest struct {
	UserID  string `json:"user_id" validate:"required,uuid"`
	Title   string `json:"title" validate:"required,max=200"`
	Message string `json:"message" validate:"required"`
	Type    string `json:"type" validate:"required,oneof=info assignment quiz grade attendance system"`
	Link    string `json:"link"`
}

type BroadcastRequest struct {
	UserIDs []string `json:"user_ids" validate:"required,min=1,dive,uuid"`
	Title   string   `json:"title" validate:"required,max=200"`
	Message string   `json:"message" validate:"required"`
	Type    string   `json:"type" validate:"required,oneof=info assignment quiz grade attendance system"`
	Link    string   `json:"link"`
}

type NotificationResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	Link      string    `json:"link,omitempty"`
	IsRead    bool      `json:"is_read"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type UnreadCount struct {
	Count int `json:"count"`
}
