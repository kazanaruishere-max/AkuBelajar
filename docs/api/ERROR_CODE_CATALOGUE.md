# ❌ Error Code Catalogue — AkuBelajar

> Katalog semua error code yang konsisten antara Go backend dan Next.js frontend. Setiap error punya kode unik, HTTP status, pesan BI, dan resolusi.

---

## Format Response Error

```json
{
  "error": {
    "code": "AUTH_001",
    "message": "Email atau password salah",
    "details": [],
    "request_id": "019516a2-uuid-v7"
  }
}
```

---

## AUTH — Autentikasi & Otorisasi

| Code | HTTP | Pesan (Bahasa Indonesia) | Kapan Muncul | Resolusi |
|:---|:---|:---|:---|:---|
| `AUTH_001` | 401 | Email atau password salah | Login gagal | Cek email dan password |
| `AUTH_002` | 429 | Akun terkunci. Coba lagi dalam 15 menit | 5× gagal login | Tunggu 15 menit atau reset password |
| `AUTH_003` | 401 | Sesi telah berakhir. Silakan login kembali | Access token expired | Client auto-refresh |
| `AUTH_004` | 401 | Token tidak valid | Token corrupt/tampered | Login ulang |
| `AUTH_005` | 403 | Anda tidak memiliki izin untuk aksi ini | Role insufficient | Hubungi admin |
| `AUTH_006` | 403 | Akun Anda telah dinonaktifkan | Account suspended | Hubungi admin |
| `AUTH_007` | 403 | Silakan selesaikan proses pendaftaran terlebih dahulu | First login belum selesai | Selesaikan onboarding |
| `AUTH_008` | 403 | Anda harus mengubah password terlebih dahulu | Force password change | Ganti password |
| `AUTH_009` | 400 | Kode OTP salah | OTP invalid | Masukkan OTP yang benar |
| `AUTH_010` | 400 | Kode OTP telah kedaluwarsa | OTP expired | Request OTP baru |
| `AUTH_011` | 429 | Terlalu banyak permintaan OTP. Coba lagi nanti | Rate limited | Tunggu cooldown |

---

## USER — Manajemen User

| Code | HTTP | Pesan | Kapan | Resolusi |
|:---|:---|:---|:---|:---|
| `USER_001` | 404 | Pengguna tidak ditemukan | GET user by ID gagal | Cek ID |
| `USER_002` | 409 | Email sudah terdaftar | Create/update duplikat email | Gunakan email lain |
| `USER_003` | 409 | NISN sudah terdaftar di sekolah ini | Duplikat NISN | Cek NISN |
| `USER_004` | 400 | Perubahan role tidak diizinkan | Invalid role transition | Cek hierarki role |
| `USER_005` | 400 | Tidak dapat menghapus akun Anda sendiri | Self-delete attempt | — |
| `USER_006` | 404 | Sekolah tidak ditemukan | Invalid school_id | Cek school_id |

---

## VAL — Validasi Input

| Code | HTTP | Pesan | Kapan | Resolusi |
|:---|:---|:---|:---|:---|
| `VAL_001` | 400 | Field %s wajib diisi | Required field kosong | Isi field |
| `VAL_002` | 400 | Format email tidak valid | Email format salah | Cek format email |
| `VAL_003` | 400 | NISN harus tepat 10 digit angka | NISN invalid | Cek NISN |
| `VAL_004` | 400 | Nomor WA harus format +62xxxxxxxxxx | Phone format salah | Format E.164 |
| `VAL_005` | 400 | Nilai harus antara 0 dan 100 | Grade out of range | Input 0-100 |
| `VAL_006` | 400 | Teks terlalu panjang (maks %d karakter) | Exceeds max length | Persingkat teks |
| `VAL_007` | 400 | Tipe file tidak diizinkan | Invalid file type | Cek allowed types |
| `VAL_008` | 400 | Ukuran file terlalu besar (maks %dMB) | File too large | Kompres file |
| `VAL_009` | 400 | Rentang tanggal tidak valid | Invalid date range | end > start |
| `VAL_010` | 400 | Password harus minimal 8 karakter dengan huruf besar, angka, dan simbol | Weak password | Perkuat password |

---

## ASSIGN — Tugas

| Code | HTTP | Pesan | Kapan | Resolusi |
|:---|:---|:---|:---|:---|
| `ASSIGN_001` | 404 | Tugas tidak ditemukan | Invalid assignment ID | Cek ID |
| `ASSIGN_002` | 400 | Deadline sudah lewat | Submit after deadline (no late allowed) | Hubungi guru |
| `ASSIGN_003` | 400 | Pengumpulan tugas tidak diizinkan | Assignment closed | — |
| `ASSIGN_004` | 409 | Anda sudah mengumpulkan tugas ini | Double submit | — |
| `ASSIGN_005` | 400 | Tidak dapat menghapus tugas yang sudah ada pengumpulan | Delete with submissions | Archive instead |
| `ASSIGN_006` | 500 | Gagal mengunggah file | Upload error | Coba lagi |

---

## QUIZ — Kuis & CBT

| Code | HTTP | Pesan | Kapan | Resolusi |
|:---|:---|:---|:---|:---|
| `QUIZ_001` | 404 | Kuis tidak ditemukan | Invalid quiz ID | Cek ID |
| `QUIZ_002` | 400 | Kuis belum dimulai atau sudah berakhir | Outside time window | Tunggu jadwal |
| `QUIZ_003` | 404 | Sesi ujian tidak ditemukan | Invalid session ID | — |
| `QUIZ_004` | 409 | Anda sudah mengerjakan kuis ini | Already submitted | — |
| `QUIZ_005` | 403 | Sesi ujian Anda telah dikunci | Anti-cheat triggered | Hubungi guru |
| `QUIZ_006` | 400 | Waktu ujian telah habis | Time expired | — |
| `QUIZ_007` | 503 | Gagal membuat soal AI. Silakan coba lagi | Gemini API error | Retry |
| `QUIZ_008` | 400 | Jawaban tidak valid | Invalid answer key | Pilih A/B/C/D |

---

## ATT — Absensi

| Code | HTTP | Pesan | Kapan | Resolusi |
|:---|:---|:---|:---|:---|
| `ATT_001` | 409 | Absensi untuk tanggal ini sudah diinput | Duplicate date | Edit yang ada |
| `ATT_002` | 400 | Status absensi tidak valid | Invalid enum value | Cek enum |
| `ATT_003` | 403 | Edit absensi setelah 24 jam memerlukan alasan | No reason provided | Isi field reason |
| `ATT_004` | 404 | Kelas tidak ditemukan | Invalid class_id | Cek ID |
| `ATT_005` | 400 | Tidak dapat input absensi untuk tanggal di masa depan | Future date | Pilih hari ini/kemarin |

---

## SYS — Sistem

| Code | HTTP | Pesan | Kapan | Resolusi |
|:---|:---|:---|:---|:---|
| `SYS_001` | 500 | Terjadi kesalahan pada server | Unhandled error | Report bug |
| `SYS_002` | 503 | Layanan sedang tidak tersedia | Service down | Coba lagi nanti |
| `SYS_003` | 500 | Kesalahan database | DB error | Report bug |
| `SYS_004` | 503 | Layanan AI sedang tidak tersedia | Gemini down | Coba lagi nanti |
| `SYS_005` | 503 | Layanan notifikasi gagal | Fonnte/SMTP down | Retry otomatis |
| `SYS_006` | 429 | Terlalu banyak permintaan | Global rate limit | Tunggu |
| `SYS_007` | 503 | Sistem sedang dalam pemeliharaan | Maintenance mode | Tunggu |

---

## Go Implementation

```go
// pkg/apperror/codes.go
type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Status  int    `json:"-"`
}

var (
    ErrInvalidCredentials = &AppError{"AUTH_001", "Email atau password salah", 401}
    ErrAccountLocked      = &AppError{"AUTH_002", "Akun terkunci. Coba lagi dalam 15 menit", 429}
    ErrTokenExpired       = &AppError{"AUTH_003", "Sesi telah berakhir. Silakan login kembali", 401}
    // ... etc
)
```

## TypeScript Usage

```typescript
// lib/errors.ts
export function getErrorMessage(code: string): string {
  const messages: Record<string, string> = {
    AUTH_001: 'Email atau password salah',
    AUTH_002: 'Akun terkunci. Coba lagi dalam 15 menit',
    // ... etc
  };
  return messages[code] ?? 'Terjadi kesalahan yang tidak diketahui';
}
```

---

*Terakhir diperbarui: 21 Maret 2026*
