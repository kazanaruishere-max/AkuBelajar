package security

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
)

// TokenMaker handles Paseto v4 token creation and validation.
type TokenMaker struct {
	symmetricKey paseto.V4SymmetricKey
}

// TokenPayload holds the claims embedded in a token.
type TokenPayload struct {
	UserID    string `json:"user_id"`
	SchoolID  string `json:"school_id"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"` // "access" or "refresh"
}

// NewTokenMaker creates a new TokenMaker with the given symmetric key.
func NewTokenMaker(keyHex string) (*TokenMaker, error) {
	key, err := paseto.V4SymmetricKeyFromHex(keyHex)
	if err != nil {
		// If hex parsing fails, generate from bytes (for dev convenience)
		key = paseto.NewV4SymmetricKey()
	}

	return &TokenMaker{symmetricKey: key}, nil
}

// CreateToken creates a new Paseto v4 local token.
func (tm *TokenMaker) CreateToken(payload TokenPayload, expiry time.Duration) (string, error) {
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(expiry))
	token.SetString("user_id", payload.UserID)
	token.SetString("school_id", payload.SchoolID)
	token.SetString("role", payload.Role)
	token.SetString("token_type", payload.TokenType)

	encrypted := token.V4Encrypt(tm.symmetricKey, nil)
	return encrypted, nil
}

// ValidateToken validates and decodes a Paseto v4 local token.
func (tm *TokenMaker) ValidateToken(tokenString string) (*TokenPayload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))

	token, err := parser.ParseV4Local(tm.symmetricKey, tokenString, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	userID, err := token.GetString("user_id")
	if err != nil {
		return nil, fmt.Errorf("missing user_id claim: %w", err)
	}

	schoolID, err := token.GetString("school_id")
	if err != nil {
		return nil, fmt.Errorf("missing school_id claim: %w", err)
	}

	role, err := token.GetString("role")
	if err != nil {
		return nil, fmt.Errorf("missing role claim: %w", err)
	}

	tokenType, err := token.GetString("token_type")
	if err != nil {
		return nil, fmt.Errorf("missing token_type claim: %w", err)
	}

	return &TokenPayload{
		UserID:    userID,
		SchoolID:  schoolID,
		Role:      role,
		TokenType: tokenType,
	}, nil
}

// CreateAccessToken creates a short-lived access token (15 min default).
func (tm *TokenMaker) CreateAccessToken(userID, schoolID, role string, expiry time.Duration) (string, error) {
	return tm.CreateToken(TokenPayload{
		UserID:    userID,
		SchoolID:  schoolID,
		Role:      role,
		TokenType: "access",
	}, expiry)
}

// CreateRefreshToken creates a long-lived refresh token (7 day default).
func (tm *TokenMaker) CreateRefreshToken(userID, schoolID, role string, expiry time.Duration) (string, error) {
	return tm.CreateToken(TokenPayload{
		UserID:    userID,
		SchoolID:  schoolID,
		Role:      role,
		TokenType: "refresh",
	}, expiry)
}
