# ✅ Acceptance Criteria — AkuBelajar

> Kriteria penerimaan eksplisit per fitur — kapan sebuah fitur dianggap **selesai**. Tanpa ini, kode bisa "berjalan" tapi tidak memenuhi kebutuhan bisnis.

---

## Format

```
Fitur: [Nama Fitur]
User Story: US-XXX
Status: ✅ Done / 🟡 In Progress / ❌ Not Started

AC-1: [Deskripsi kriteria]
AC-2: [Deskripsi kriteria]
...
```

---

## AUTH — Login & Register

**User Story:** US-001  
**Status:** ✅ Done

| # | Kriteria | Status |
|:---|:---|:---|
| AC-1 | User bisa login dengan email + password yang valid | ✅ |
| AC-2 | Login gagal mengembalikan error message tanpa hint password | ✅ |
| AC-3 | Setelah 5x gagal login, akun terkunci selama 15 menit | ✅ |
| AC-4 | Access token expire dalam 15 menit | ✅ |
| AC-5 | Refresh token rotate setiap kali digunakan | ✅ |
| AC-6 | Logout menginvalidasi refresh token di server | ✅ |
| AC-7 | Password di-hash dengan Argon2id (bukan bcrypt/MD5) | ✅ |
| AC-8 | Rate limit: max 5 login attempts per 15 menit per IP | ✅ |

---

## QUIZ — AI Quiz Generation

**User Story:** US-002  
**Status:** ✅ Done

| # | Kriteria | Status |
|:---|:---|:---|
| AC-1 | Guru bisa input topik, jumlah soal (5-50), dan tingkat kesulitan | ✅ |
| AC-2 | AI menghasilkan soal pilihan ganda (4 opsi) yang valid dan relevan | ✅ |
| AC-3 | Output AI selalu valid JSON (Gemini 2.0 JSON Mode) | ✅ |
| AC-4 | Guru bisa preview, edit, dan hapus soal sebelum publish | ✅ |
| AC-5 | Jawaban benar di-hash (Argon2id) di database — tidak readable | ✅ |
| AC-6 | Input topic di-sanitasi (anti prompt injection) | ✅ |
| AC-7 | Jika AI gagal (timeout/error), tampilkan pesan error yang jelas | ✅ |
| AC-8 | Rate limit: max 10 AI generations per guru per jam | ✅ |

---

## CBT — Computer Based Test

**User Story:** US-003  
**Status:** ✅ Done

| # | Kriteria | Status |
|:---|:---|:---|
| AC-1 | Siswa hanya bisa mulai kuis saat `start_at ≤ now ≤ end_at` | ✅ |
| AC-2 | Timer dikontrol server — manipulasi jam client tidak berdampak | ✅ |
| AC-3 | Auto-submit saat waktu habis (0 detik grace period) | ✅ |
| AC-4 | Deteksi pindah tab: 3x peringatan, ke-4 auto-submit | ✅ |
| AC-5 | Deteksi minimize: 2x peringatan, ke-3 auto-submit | ✅ |
| AC-6 | WebSocket disconnect > 30 detik = session expired | ✅ |
| AC-7 | Nilai instan ditampilkan setelah submit | ✅ |
| AC-8 | Siswa tidak bisa submit kuis yang sama dua kali | ✅ |
| AC-9 | Pembahasan AI hanya tersedia setelah submit | ✅ |

---

## ATTENDANCE — Absensi Digital

**User Story:** US-004  
**Status:** ✅ Done

| # | Kriteria | Status |
|:---|:---|:---|
| AC-1 | Guru/Ketua Kelas bisa input absensi per kelas per hari | ✅ |
| AC-2 | Status: present, permission, sick, absent | ✅ |
| AC-3 | Bisa diakses dari mobile (PWA responsive) | ✅ |
| AC-4 | Input absensi hanya untuk T+0 (hari ini) atau T-1 (kemarin) | ✅ |
| AC-5 | Edit absensi lama max T-7, hanya oleh guru | ✅ |
| AC-6 | Setiap perubahan masuk audit log | ✅ |
| AC-7 | Offline support: data tersimpan dan sync saat online | ✅ |

---

## GRADES — Lihat Nilai

**User Story:** US-005  
**Status:** ✅ Done

| # | Kriteria | Status |
|:---|:---|:---|
| AC-1 | Siswa hanya bisa melihat nilai sendiri (RLS enforced) | ✅ |
| AC-2 | Guru hanya bisa melihat nilai kelas yang diajar (RLS) | ✅ |
| AC-3 | Nilai ditampilkan per mata pelajaran, per periode | ✅ |
| AC-4 | Grafik tren nilai tersedia | ✅ |
| AC-5 | Notifikasi saat nilai baru diinput | ✅ |
| AC-6 | Nilai akhir dihitung otomatis berdasarkan bobot configurable | ✅ |

---

## BULK IMPORT — Import User Massal

**User Story:** US-006  
**Status:** ✅ Done

| # | Kriteria | Status |
|:---|:---|:---|
| AC-1 | Admin bisa upload file CSV/Excel | ✅ |
| AC-2 | Validasi per baris dengan error report detail | ✅ |
| AC-3 | Proses berjalan async (tidak memblokir UI) | ✅ |
| AC-4 | Notifikasi saat proses selesai | ✅ |
| AC-5 | Auto-generate password + kirim via email | ✅ |
| AC-6 | Idempotent — aman dijalankan ulang tanpa duplikasi | ✅ |
| AC-7 | Mampu handle 100K+ baris tanpa crash | ✅ |

---

## NOTIFICATION — Multi-Channel

**User Story:** US-008  
**Status:** 🟡 In Progress

| # | Kriteria | Status |
|:---|:---|:---|
| AC-1 | In-app notification real-time (WebSocket) | 🟡 |
| AC-2 | Email notification | ✅ |
| AC-3 | WhatsApp notification | 🟡 |
| AC-4 | Push notification (PWA) | 🟡 |
| AC-5 | User bisa set preferensi channel per tipe notifikasi | ❌ |

---

## Template Kriteria Baru

```markdown
## [NAMA FITUR]

**User Story:** US-XXX
**Status:** ❌ Not Started

| # | Kriteria | Status |
|:---|:---|:---|
| AC-1 | [Deskripsi kriteria yang spesifik dan terukur] | ❌ |
| AC-2 | [Deskripsi kriteria] | ❌ |
```

---

## Referensi

- [User Stories](../fase-5/USER_STORIES.md) — Kebutuhan user
- [Business Logic](BUSINESS_LOGIC.md) — Aturan domain
- [Test Cases](TEST_CASES.md) — Test yang memvalidasi kriteria ini

---

*Terakhir diperbarui: 21 Maret 2026*
