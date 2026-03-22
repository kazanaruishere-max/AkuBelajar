package auth

import "time"

// ---- Request DTOs ----

// LoginRequest holds login credentials.
type LoginRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	RememberMe bool   `json:"remember_me"`
}

// RefreshRequest holds the refresh token for renewal.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ChangePasswordRequest holds data for changing password.
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,strong_password"`
}

// ---- Response DTOs ----

// TokenResponse is returned after successful login/refresh.
type TokenResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"` // seconds
	User         UserResponse `json:"user"`
}

// UserResponse holds public user info returned in responses.
type UserResponse struct {
	ID           string         `json:"id"`
	Email        string         `json:"email"`
	Role         string         `json:"role"`
	SchoolID     string         `json:"school_id"`
	SchoolName   string         `json:"school_name,omitempty"`
	IsFirstLogin bool           `json:"is_first_login"`
	Profile      *ProfileBrief  `json:"profile,omitempty"`
	LastLoginAt  *time.Time     `json:"last_login_at,omitempty"`
}

// ProfileBrief is a subset of user_profiles for auth responses.
type ProfileBrief struct {
	NISN      string `json:"nisn,omitempty"`
	NIP       string `json:"nip,omitempty"`
	PhoneWA   string `json:"phone_wa,omitempty"`
	PhotoURL  string `json:"photo_url,omitempty"`
}

// ---- DB Models ----

// User represents a row in the users table.
type User struct {
	ID               string     `json:"id"`
	SchoolID         string     `json:"school_id"`
	Email            string     `json:"email"`
	Password         string     `json:"-"` // never expose hash
	Role             string     `json:"role"`
	IsActive         bool       `json:"is_active"`
	IsFirstLogin     bool       `json:"is_first_login"`
	FailedLoginCount int        `json:"failed_login_count"`
	LockedUntil      *time.Time `json:"locked_until,omitempty"`
	LastLoginAt      *time.Time `json:"last_login_at,omitempty"`
	LastLoginIP      *string    `json:"last_login_ip,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// Session represents a row in active_sessions table.
type Session struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	RefreshTokenHash string    `json:"-"`
	DeviceInfo       string    `json:"device_info,omitempty"`
	IPAddress        string    `json:"ip_address,omitempty"`
	IsRememberMe     bool      `json:"is_remember_me"`
	ExpiresAt        time.Time `json:"expires_at"`
	CreatedAt        time.Time `json:"created_at"`
}
