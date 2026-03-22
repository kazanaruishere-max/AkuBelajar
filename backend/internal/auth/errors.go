package auth

import "errors"

// Domain errors for auth operations.
var (
	ErrInvalidCredentials  = errors.New("email atau password salah")
	ErrAccountInactive     = errors.New("akun dinonaktifkan, hubungi admin")
	ErrAccountLocked       = errors.New("akun terkunci karena terlalu banyak percobaan login. Coba lagi dalam 15 menit")
	ErrInvalidRefreshToken = errors.New("sesi tidak valid, silakan login kembali")
)
