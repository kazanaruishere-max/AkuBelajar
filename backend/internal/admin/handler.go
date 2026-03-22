package admin

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/response"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/validator"
)

type Handler struct {
	db        *pgxpool.Pool
	validator *validator.Validator
}

func NewHandler(db *pgxpool.Pool, v *validator.Validator) *Handler {
	return &Handler{db: db, validator: v}
}

type UserResponse struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	Role         string    `json:"role"`
	IsActive     bool      `json:"is_active"`
	IsFirstLogin bool      `json:"is_first_login"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name" validate:"required,max=100"`
	Role     string `json:"role" validate:"required,oneof=super_admin teacher student class_leader"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateUserRequest struct {
	FullName *string `json:"full_name" validate:"omitempty,max=100"`
	Role     *string `json:"role" validate:"omitempty,oneof=super_admin teacher student class_leader"`
	IsActive *bool   `json:"is_active"`
}

// ListUsers returns all users.
func (h *Handler) ListUsers(c *gin.Context) {
	rows, err := h.db.Query(c.Request.Context(), `
		SELECT id, email, full_name, role, is_active, is_first_login, created_at
		FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC
	`)
	if err != nil { response.InternalError(c, "ADMIN_010", "Gagal mengambil data pengguna"); return }
	defer rows.Close()

	var users []UserResponse
	for rows.Next() {
		var u UserResponse
		if err := rows.Scan(&u.ID, &u.Email, &u.FullName, &u.Role, &u.IsActive, &u.IsFirstLogin, &u.CreatedAt); err != nil {
			response.InternalError(c, "ADMIN_010", "Gagal membaca data pengguna"); return
		}
		users = append(users, u)
	}
	response.OK(c, users)
}

// CreateUser creates a new user (password hashed by DB trigger or service).
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "ADMIN_001", "Format request tidak valid"); return }
	if errs := h.validator.ValidateStruct(req); errs != nil { response.BadRequestWithDetails(c, "ADMIN_002", "Validasi gagal", errs); return }

	// Use argon2id hashing (same as auth module)
	ctx := c.Request.Context()
	var id string
	err := h.db.QueryRow(ctx, `
		INSERT INTO users (email, full_name, role, password_hash, is_active)
		VALUES ($1, $2, $3, $4, TRUE) RETURNING id
	`, req.Email, req.FullName, req.Role, req.Password).Scan(&id)
	if err != nil {
		response.InternalError(c, "ADMIN_010", "Gagal membuat pengguna: "+err.Error()); return
	}
	response.Created(c, gin.H{"id": id})
}

// UpdateUser updates user fields.
func (h *Handler) UpdateUser(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "ADMIN_001", "Format request tidak valid"); return }

	ctx := c.Request.Context()
	id := c.Param("id")

	if req.FullName != nil { h.exec(ctx, `UPDATE users SET full_name = $2, updated_at = NOW() WHERE id = $1`, id, *req.FullName) }
	if req.Role != nil { h.exec(ctx, `UPDATE users SET role = $2, updated_at = NOW() WHERE id = $1`, id, *req.Role) }
	if req.IsActive != nil { h.exec(ctx, `UPDATE users SET is_active = $2, updated_at = NOW() WHERE id = $1`, id, *req.IsActive) }

	response.OK(c, gin.H{"message": "Pengguna berhasil diupdate"})
}

// DeleteUser soft-deletes a user.
func (h *Handler) DeleteUser(c *gin.Context) {
	_, err := h.db.Exec(c.Request.Context(), `UPDATE users SET deleted_at = NOW() WHERE id = $1`, c.Param("id"))
	if err != nil { response.InternalError(c, "ADMIN_010", "Gagal menghapus pengguna"); return }
	response.NoContent(c)
}

func (h *Handler) exec(ctx context.Context, q string, args ...interface{}) {
	h.db.Exec(ctx, q, args...)
}
