package academic

import "errors"

var (
	ErrCannotDeleteActive = errors.New("tidak dapat menghapus tahun ajaran yang sedang aktif")
	ErrNotFound           = errors.New("data tidak ditemukan")
)
