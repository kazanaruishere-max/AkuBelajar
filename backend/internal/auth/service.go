package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/kazanaruishere-max/akubelajar/backend/pkg/security"
)

// Service handles authentication business logic.
type Service struct {
	repo       *Repository
	tokenMaker *security.TokenMaker
	accessExp  time.Duration
	refreshExp time.Duration
}

// NewService creates a new auth service.
func NewService(repo *Repository, tokenMaker *security.TokenMaker, accessExp, refreshExp time.Duration) *Service {
	return &Service{
		repo:       repo,
		tokenMaker: tokenMaker,
		accessExp:  accessExp,
		refreshExp: refreshExp,
	}
}

// Login authenticates a user and returns tokens.
func (s *Service) Login(ctx context.Context, req LoginRequest, ip, userAgent string) (*TokenResponse, error) {
	// Find user by email
	user, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Check if account is active
	if !user.IsActive {
		return nil, ErrAccountInactive
	}

	// Check if account is locked
	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		return nil, ErrAccountLocked
	}

	// Verify password (Argon2id)
	match, err := security.VerifyPassword(req.Password, user.Password)
	if err != nil {
		return nil, fmt.Errorf("verify password: %w", err)
	}
	if !match {
		// Increment failed login counter (auto-locks after 5)
		_ = s.repo.IncrementFailedLogin(ctx, user.ID)
		return nil, ErrInvalidCredentials
	}

	// Success — reset failed login counter and update last_login
	if err := s.repo.ResetFailedLogin(ctx, user.ID, ip); err != nil {
		return nil, fmt.Errorf("reset failed login: %w", err)
	}

	// Generate tokens
	accessToken, err := s.tokenMaker.CreateAccessToken(user.ID, user.SchoolID, user.Role, s.accessExp)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	refreshExp := s.refreshExp
	if req.RememberMe {
		refreshExp = 30 * 24 * time.Hour // 30 days for "Remember Me"
	}

	refreshToken, err := s.tokenMaker.CreateRefreshToken(user.ID, user.SchoolID, user.Role, refreshExp)
	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}

	// Store session
	_, err = s.repo.CreateSession(ctx, user.ID, refreshToken, ip, userAgent, req.RememberMe, time.Now().Add(refreshExp))
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	// Build response
	userResp := s.buildUserResponse(ctx, user)

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.accessExp.Seconds()),
		User:         userResp,
	}, nil
}

// RefreshToken validates a refresh token and returns new token pair.
func (s *Service) RefreshToken(ctx context.Context, refreshToken, ip, userAgent string) (*TokenResponse, error) {
	// Find session by hashed refresh token
	session, err := s.repo.FindSessionByToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("find session: %w", err)
	}
	if session == nil {
		return nil, ErrInvalidRefreshToken
	}

	// Delete old session (single-use refresh token rotation)
	if err := s.repo.DeleteSession(ctx, session.ID); err != nil {
		return nil, fmt.Errorf("delete old session: %w", err)
	}

	// Get the user
	user, err := s.repo.FindUserByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil || !user.IsActive {
		return nil, ErrAccountInactive
	}

	// Generate new tokens
	accessToken, err := s.tokenMaker.CreateAccessToken(user.ID, user.SchoolID, user.Role, s.accessExp)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	refreshExp := s.refreshExp
	if session.IsRememberMe {
		refreshExp = 30 * 24 * time.Hour
	}

	newRefreshToken, err := s.tokenMaker.CreateRefreshToken(user.ID, user.SchoolID, user.Role, refreshExp)
	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}

	// Store new session
	_, err = s.repo.CreateSession(ctx, user.ID, newRefreshToken, ip, userAgent, session.IsRememberMe, time.Now().Add(refreshExp))
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	userResp := s.buildUserResponse(ctx, user)

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int(s.accessExp.Seconds()),
		User:         userResp,
	}, nil
}

// Logout invalidates a refresh token session.
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	session, err := s.repo.FindSessionByToken(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("find session: %w", err)
	}
	if session == nil {
		return nil // Already logged out — idempotent
	}
	return s.repo.DeleteSession(ctx, session.ID)
}

// ChangePassword changes the user's password.
func (s *Service) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return ErrInvalidCredentials
	}

	// Verify old password
	match, err := security.VerifyPassword(req.OldPassword, user.Password)
	if err != nil {
		return fmt.Errorf("verify old password: %w", err)
	}
	if !match {
		return ErrInvalidCredentials
	}

	// Hash new password
	newHash, err := security.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("hash new password: %w", err)
	}

	// Update password + save to history
	if err := s.repo.UpdatePassword(ctx, userID, newHash); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	// Invalidate all sessions (force re-login everywhere)
	if err := s.repo.DeleteAllUserSessions(ctx, userID); err != nil {
		return fmt.Errorf("delete sessions: %w", err)
	}

	return nil
}

// GetMe returns the current authenticated user's info.
func (s *Service) GetMe(ctx context.Context, userID string) (*UserResponse, error) {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}
	resp := s.buildUserResponse(ctx, user)
	return &resp, nil
}

// buildUserResponse constructs a UserResponse from a User model.
func (s *Service) buildUserResponse(ctx context.Context, user *User) UserResponse {
	resp := UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		Role:         user.Role,
		SchoolID:     user.SchoolID,
		IsFirstLogin: user.IsFirstLogin,
		LastLoginAt:  user.LastLoginAt,
	}

	// Get school name (best-effort)
	if name, err := s.repo.GetSchoolName(ctx, user.SchoolID); err == nil {
		resp.SchoolName = name
	}

	// Get profile (best-effort)
	if profile, err := s.repo.GetUserProfile(ctx, user.ID); err == nil {
		resp.Profile = profile
	}

	return resp
}
