package quiz

import "errors"

var (
	ErrNotFound       = errors.New("kuis tidak ditemukan")
	ErrNotAuthorized  = errors.New("Anda tidak memiliki akses ke kuis ini")
	ErrNotPublished   = errors.New("kuis belum dipublish")
	ErrQuizEnded      = errors.New("kuis sudah berakhir")
	ErrAlreadyStarted = errors.New("Anda sudah mengerjakan kuis ini")
	ErrSessionExpired = errors.New("waktu mengerjakan sudah habis")
	ErrMaxAttempts    = errors.New("batas percobaan sudah tercapai")
)
