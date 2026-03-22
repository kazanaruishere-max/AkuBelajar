package security

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

// HashPassword creates an Argon2id hash from a plain-text password.
// Argon2id is chosen over bcrypt for memory-hard GPU-resistance.
func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return hash, nil
}

// VerifyPassword checks if a plain-text password matches an Argon2id hash.
func VerifyPassword(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("verify password: %w", err)
	}
	return match, nil
}
