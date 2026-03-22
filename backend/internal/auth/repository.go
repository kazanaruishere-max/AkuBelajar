package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles all auth-related database operations.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new auth repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// FindUserByEmail looks up a user by email address.
func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT u.id, u.school_id, u.email, u.password, u.role,
		       u.is_active, u.is_first_login, u.failed_login_count,
		       u.locked_until, u.last_login_at, u.last_login_ip,
		       u.created_at, u.updated_at
		FROM users u
		WHERE u.email = $1 AND u.deleted_at IS NULL`

	var user User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.SchoolID, &user.Email, &user.Password, &user.Role,
		&user.IsActive, &user.IsFirstLogin, &user.FailedLoginCount,
		&user.LockedUntil, &user.LastLoginAt, &user.LastLoginIP,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &user, nil
}

// FindUserByID looks up a user by ID.
func (r *Repository) FindUserByID(ctx context.Context, userID string) (*User, error) {
	query := `
		SELECT u.id, u.school_id, u.email, u.password, u.role,
		       u.is_active, u.is_first_login, u.failed_login_count,
		       u.locked_until, u.last_login_at, u.last_login_ip,
		       u.created_at, u.updated_at
		FROM users u
		WHERE u.id = $1 AND u.deleted_at IS NULL`

	var user User
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID, &user.SchoolID, &user.Email, &user.Password, &user.Role,
		&user.IsActive, &user.IsFirstLogin, &user.FailedLoginCount,
		&user.LockedUntil, &user.LastLoginAt, &user.LastLoginIP,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &user, nil
}

// GetSchoolName retrieves the school name for a given school_id.
func (r *Repository) GetSchoolName(ctx context.Context, schoolID string) (string, error) {
	var name string
	err := r.db.QueryRow(ctx, "SELECT name FROM schools WHERE id = $1", schoolID).Scan(&name)
	if err != nil {
		return "", fmt.Errorf("get school name: %w", err)
	}
	return name, nil
}

// GetUserProfile retrieves a brief profile for a user.
func (r *Repository) GetUserProfile(ctx context.Context, userID string) (*ProfileBrief, error) {
	query := `SELECT nisn, nip, phone_wa, photo_url FROM user_profiles WHERE user_id = $1`

	var p ProfileBrief
	var nisn, nip, phone, photo *string
	err := r.db.QueryRow(ctx, query, userID).Scan(&nisn, &nip, &phone, &photo)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user profile: %w", err)
	}

	if nisn != nil { p.NISN = *nisn }
	if nip != nil { p.NIP = *nip }
	if phone != nil { p.PhoneWA = *phone }
	if photo != nil { p.PhotoURL = *photo }

	return &p, nil
}

// IncrementFailedLogin increments the failed login counter and locks if threshold reached.
func (r *Repository) IncrementFailedLogin(ctx context.Context, userID string) error {
	query := `
		UPDATE users SET
			failed_login_count = failed_login_count + 1,
			locked_until = CASE
				WHEN failed_login_count + 1 >= 5
				THEN NOW() + INTERVAL '15 minutes'
				ELSE locked_until
			END,
			updated_at = NOW()
		WHERE id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

// ResetFailedLogin resets the counter after successful login.
func (r *Repository) ResetFailedLogin(ctx context.Context, userID string, ip string) error {
	query := `
		UPDATE users SET
			failed_login_count = 0,
			locked_until = NULL,
			last_login_at = NOW(),
			last_login_ip = $2,
			updated_at = NOW()
		WHERE id = $1`
	_, err := r.db.Exec(ctx, query, userID, ip)
	return err
}

// CreateSession stores a new active session.
func (r *Repository) CreateSession(ctx context.Context, userID, refreshToken, ip, userAgent string, isRemember bool, expiresAt time.Time) (string, error) {
	hash := hashToken(refreshToken)
	var sessionID string
	query := `
		INSERT INTO active_sessions (user_id, refresh_token_hash, ip_address, device_info, is_remember_me, expires_at)
		VALUES ($1, $2, $3::inet, $4::jsonb, $5, $6)
		RETURNING id`

	deviceJSON := fmt.Sprintf(`{"user_agent": "%s"}`, userAgent)
	err := r.db.QueryRow(ctx, query, userID, hash, ip, deviceJSON, isRemember, expiresAt).Scan(&sessionID)
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}
	return sessionID, nil
}

// FindSessionByToken finds a session by its hashed refresh token.
func (r *Repository) FindSessionByToken(ctx context.Context, refreshToken string) (*Session, error) {
	hash := hashToken(refreshToken)
	query := `
		SELECT id, user_id, refresh_token_hash, ip_address, is_remember_me, expires_at, created_at
		FROM active_sessions
		WHERE refresh_token_hash = $1 AND expires_at > NOW()`

	var s Session
	var ip *string
	err := r.db.QueryRow(ctx, query, hash).Scan(
		&s.ID, &s.UserID, &s.RefreshTokenHash, &ip, &s.IsRememberMe, &s.ExpiresAt, &s.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find session: %w", err)
	}
	if ip != nil { s.IPAddress = *ip }
	return &s, nil
}

// DeleteSession removes a single session (logout).
func (r *Repository) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM active_sessions WHERE id = $1", sessionID)
	return err
}

// DeleteAllUserSessions removes all sessions for a user (force logout everywhere).
func (r *Repository) DeleteAllUserSessions(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM active_sessions WHERE user_id = $1", userID)
	return err
}

// UpdatePassword changes the user's password and records history.
func (r *Repository) UpdatePassword(ctx context.Context, userID, newHash string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Save to password history
	_, err = tx.Exec(ctx,
		"INSERT INTO password_histories (user_id, password_hash) VALUES ($1, $2)",
		userID, newHash,
	)
	if err != nil {
		return fmt.Errorf("save password history: %w", err)
	}

	// Update user password and clear first_login flag
	_, err = tx.Exec(ctx,
		"UPDATE users SET password = $2, is_first_login = FALSE, updated_at = NOW() WHERE id = $1",
		userID, newHash,
	)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	return tx.Commit(ctx)
}

// hashToken creates a SHA-256 hash of a token for storage comparison.
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
