package assignment

import "errors"

var (
	ErrNotFound         = errors.New("tugas tidak ditemukan")
	ErrNotAuthorized    = errors.New("Anda tidak memiliki akses ke tugas ini")
	ErrAlreadySubmitted = errors.New("tugas sudah disubmit")
	ErrDeadlinePassed   = errors.New("batas waktu pengumpulan sudah lewat")
	ErrNotPublished     = errors.New("tugas belum dipublish")
	ErrAlreadyGraded    = errors.New("tugas sudah dinilai")
)
