# 🔥 Edge Cases — AkuBelajar

> Skenario ekstrem yang HARUS ditangani. Kasus yang paling sering diabaikan tapi paling sering menyebabkan bug di production.

---

## Format

- **Skenario:** Deskripsi singkat
- **Expected:** Apa yang SEHARUSNYA terjadi
- **Layer:** Frontend / Backend / Database / Semua

---

## Autentikasi & Session

| # | Skenario | Expected | Layer |
|:---|:---|:---|:---|
| E-01 | Login di device A, ganti password di device B | Device A: redirect ke login saat access token expired (refresh gagal karena session revoked) | Backend + FE |
| E-02 | Token expired tepat saat request dikirim | Server return 401 → client auto-refresh → retry request | FE |
| E-03 | User klik login 5× cepat (double submit) | Debounce di FE. Backend: idempotent — hanya 1 session dibuat | Semua |
| E-04 | Redis down tapi PostgreSQL up | Fallback: skip blocklist check, accept all tokens. Log warning | Backend |

---

## Manajemen File

| # | Skenario | Expected | Layer |
|:---|:---|:---|:---|
| E-05 | File .jpg tapi isinya executable | Reject: magic byte mismatch → `VAL_007` | Backend |
| E-06 | Filename: `../../../etc/passwd` | Reject: path traversal detected. Rename ke UUID | Backend |
| E-07 | Upload file 0 byte | Reject: min size 1 byte → `VAL_008` | Semua |
| E-08 | Upload tepat 20MB vs 20MB+1 | 20MB: accept. 20MB+1: reject → `VAL_008` | Semua |
| E-09 | 2 user upload nama file sama bersamaan | Aman: semua di-rename ke UUID unik | Backend |
| E-10 | MinIO down saat upload | Return `SYS_002` + retry button. File tersimpan di temp | Backend + FE |

---

## Absensi

| # | Skenario | Expected | Layer |
|:---|:---|:---|:---|
| E-11 | Guru input absensi tanggal sama 2× | Unique constraint prevent. Return `ATT_001` | DB |
| E-12 | Guru input absensi tanggal libur | Accept (guru tau jadwalnya). Bukan concern backend | — |
| E-13 | Siswa dikeluarkan dari kelas setelah absensi diinput | Absensi historis tetap ada. Siswa tidak muncul di list kelas baru | Backend |
| E-14 | Koneksi putus saat guru input 30 siswa | FE: simpan di local state. Tampilkan retry. Backend: atomic — semua masuk atau semua gagal | Semua |

---

## Tugas & Submission

| # | Skenario | Expected | Layer |
|:---|:---|:---|:---|
| E-15 | Siswa submit 2× dalam 1 detik (double click) | Unique constraint (assignment_id, student_id). Return `ASSIGN_004` pada submit ke-2 | DB + FE |
| E-16 | Guru hapus tugas setelah 10 submission | Soft delete. Submissions tetap ada. Nilai tetap valid | Backend |
| E-17 | Deadline dimundurkan setelah ada yang submit | Recalculate `is_late` dan `late_days` untuk semua submission yang ada | Backend |
| E-18 | Guru ubah nilai setelah rapor locked | Block: return 403. Rapor harus di-unlock dulu | Backend |

---

## Kuis / CBT

| # | Skenario | Expected | Layer |
|:---|:---|:---|:---|
| E-19 | Siswa buka kuis di 2 tab sekaligus | Tab pertama yang mulai session = valid. Tab kedua = `QUIZ_004` | Backend |
| E-20 | Waktu habis tepat saat siswa klik submit | Server time wins. Jika `NOW() > expires_at` → auto-submit (bukan manual submit) | Backend |
| E-21 | AI gagal generate soal → kuis kosong | Rollback: hapus kuis yang baru dibuat. Return `QUIZ_007` | Backend |
| E-22 | Koneksi putus 1 detik sebelum waktu habis | Jawaban terakhir autosave 30 detik lalu = tersimpan. Auto-submit oleh server | Backend |
| E-23 | Guru perpanjang waktu saat siswa mengerjakan | Update `quiz_sessions.expires_at` via WebSocket push. Siswa lihat timer bertambah | Semua |

---

## Nilai & Kalkulasi

| # | Skenario | Expected | Layer |
|:---|:---|:---|:---|
| E-24 | Nilai 0 untuk semua siswa | Valid. Guru bisa memberikan 0 dengan sengaja (tugas tidak dikerjakan) | — |
| E-25 | Bobot 60/40 tapi tidak ada nilai kuis | Formula fallback: `final_score = avg_tugas × 1.0` | Backend |
| E-26 | Nilai > 100 (bonus?) | Reject: constraint `grade >= 0 AND grade <= 100`. Jika bonus dibutuhkan → fitur terpisah | DB |

---

## Database & Concurrency

| # | Skenario | Expected | Layer |
|:---|:---|:---|:---|
| E-27 | 2 guru input nilai siswa sama bersamaan | Optimistic locking: `updated_at` check. Yang kedua mendapat conflict error | DB + Backend |
| E-28 | Bulk import 1000 siswa saat load tinggi | Batch processing: 50 per batch, 500ms delay. Progress bar di FE | Semua |
| E-29 | Database migration gagal di tengah jalan | Transaction rollback otomatis. Migration di-mark failed. Retry manual | DB |

---

*Terakhir diperbarui: 21 Maret 2026*
