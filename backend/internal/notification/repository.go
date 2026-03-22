package notification

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, req CreateNotificationRequest) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO notifications (user_id, title, message, type, link)
		VALUES ($1,$2,$3,$4,$5)
	`, req.UserID, req.Title, req.Message, req.Type, req.Link)
	return err
}

func (r *Repository) Broadcast(ctx context.Context, req BroadcastRequest) error {
	for _, uid := range req.UserIDs {
		_, err := r.db.Exec(ctx, `
			INSERT INTO notifications (user_id, title, message, type, link)
			VALUES ($1,$2,$3,$4,$5)
		`, uid, req.Title, req.Message, req.Type, req.Link)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) ListByUser(ctx context.Context, userID string, limit int) ([]NotificationResponse, error) {
	if limit <= 0 { limit = 50 }
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, title, message, type, COALESCE(link,''), is_read, read_at, created_at
		FROM notifications WHERE user_id = $1
		ORDER BY created_at DESC LIMIT $2
	`, userID, limit)
	if err != nil { return nil, err }
	defer rows.Close()

	var items []NotificationResponse
	for rows.Next() {
		var n NotificationResponse
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.Type, &n.Link, &n.IsRead, &n.ReadAt, &n.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, n)
	}
	return items, nil
}

func (r *Repository) MarkRead(ctx context.Context, userID, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notifications SET is_read = TRUE, read_at = NOW() WHERE id = $1 AND user_id = $2
	`, id, userID)
	return err
}

func (r *Repository) MarkAllRead(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notifications SET is_read = TRUE, read_at = NOW() WHERE user_id = $1 AND is_read = FALSE
	`, userID)
	return err
}

func (r *Repository) UnreadCount(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE`, userID).Scan(&count)
	return count, err
}
