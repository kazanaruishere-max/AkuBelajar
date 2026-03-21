# ✅ Acceptance Criteria Full — AkuBelajar

> Kriteria kapan setiap fitur dianggap SELESAI dan SIAP RILIS. Format: Given/When/Then.

---

## 1. Autentikasi

### AC-AUTH-01: Login Berhasil
```
Given   User dengan akun aktif dan password benar
When    User submit form login
Then    Redirect ke dashboard sesuai role
And     Access token dan refresh token diterbitkan
And     last_login_at dan last_login_ip di-update
```

### AC-AUTH-02: Login Gagal — Account Lock
```
Given   User gagal login 5× berturut-turut
When    User mencoba login ke-6
Then    Response 429, akun terkunci 15 menit
And     Pesan: "Akun terkunci. Coba lagi dalam 15 menit"
And     Audit log: login_failed + account_locked
```

### AC-AUTH-03: Logout
```
Given   User sudah login
When    User klik logout
Then    Refresh token di-revoke dari DB
And     Session dihapus dari active_sessions
And     Redirect ke /login
And     Akses ke halaman protected → redirect /login
```

### AC-AUTH-04: Refresh Token Rotation
```
Given   Access token expired, refresh token valid
When    Client POST /auth/refresh
Then    Access token baru diterbitkan
And     Refresh token baru diterbitkan (old invalidated)
And     Jika refresh token lama digunakan lagi → semua session revoked (token theft detection)
```

---

## 2. Manajemen Akun

### AC-ACCOUNT-01: Bulk Import 100 Siswa
```
Given   Admin upload Excel dengan 100 baris data siswa valid
When    Admin submit import
Then    100 akun siswa dibuat dengan is_first_login=true
And     Temp password generated per siswa
And     Report: "100 berhasil, 0 gagal"
And     Total waktu < 30 detik
```

### AC-ACCOUNT-02: Invite Token Lifecycle
```
Given   Admin generate invite token (max_uses=5, expires 72 jam)
When    3 user claim token
Then    uses_count = 3, token masih valid
When    Token expired setelah 72 jam
Then    Claim gagal: "Token kedaluwarsa"
```

### AC-ACCOUNT-03: First Login Force Password Change
```
Given   User login pertama kali (is_first_login=true)
When    Login berhasil
Then    Redirect ke /onboarding/change-password
And     Semua route lain blocked
And     Setelah ganti password → lanjut ke profil wizard
And     Setelah wizard selesai → is_first_login=false
```

---

## 3. Absensi

### AC-ATT-01: Input Absensi Kelas
```
Given   Guru membuka halaman absensi untuk kelas 8A
When    Guru set status 30 siswa dan submit
Then    30 record absensi tersimpan (status: FINAL)
And     Waktu input < 2 menit untuk 30 siswa
And     Audit log tercatat
```

### AC-ATT-02: Edit Absensi Setelah 24 Jam
```
Given   Absensi sudah diinput > 24 jam lalu
When    Guru edit status siswa tanpa mengisi alasan
Then    Error: "Edit setelah 24 jam memerlukan alasan"
When    Guru mengisi alasan dan submit ulang
Then    Update berhasil, audit log mencatat old→new value
```

### AC-ATT-03: Rekap Kehadiran
```
Given   Siswa punya 20 hari hadir, 2 izin, 1 sakit, 3 alfa dalam 1 bulan
When    Sistem hitung persentase
Then    Hadir efektif = (20+2+1) = 23 dari 26 = 88.5%
```

---

## 4. Tugas

### AC-ASSIGN-01: Submit Sebelum Deadline
```
Given   Tugas aktif dengan deadline 25 Maret 23:59
When    Siswa upload 2 file PDF dan submit pada 24 Maret
Then    Status: SUBMITTED, is_late: false
And     File di-rename dan disimpan di MinIO
And     Guru menerima notifikasi
```

### AC-ASSIGN-02: Submit Setelah Deadline — Not Allowed
```
Given   Tugas dengan allow_late=false, deadline sudah lewat
When    Siswa mencoba submit
Then    Tombol submit disabled
And     Pesan: "Deadline sudah lewat"
```

### AC-ASSIGN-03: Nilai Masuk Buku Nilai
```
Given   Guru memberi nilai 85 pada submission siswa
When    Nilai disimpan
Then    grades.assignment_avg di-recalculate otomatis
And     Siswa menerima notifikasi
```

---

## 5. Kuis / CBT

### AC-QUIZ-01: Timer Server-Side
```
Given   Kuis 30 menit dimulai
When    Siswa manipulasi clock client
Then    Timer server tetap berjalan normal
And     Auto-submit pada waktu server, bukan client
```

### AC-QUIZ-02: Anti-Cheat Threshold
```
Given   Siswa mengerjakan kuis
When    Siswa switch tab 3 kali
Then    Warning 1, 2, 3 ditampilkan
When    Tab switch ke-4
Then    Session LOCKED, soal hilang
And     Guru menerima real-time notification
And     Jawaban yang sudah ada tetap tersimpan
```

### AC-QUIZ-03: Auto-Submit Saat Waktu Habis
```
Given   Siswa sedang mengerjakan, 15 dari 20 soal dijawab
When    Timer mencapai 0
Then    Server auto-submit 15 jawaban
And     5 soal tidak dijawab = skor 0
And     Nilai langsung muncul (jika review_mode=immediately)
```

---

## 6. Nilai & Rapor

### AC-GRADE-01: Formula Nilai Akhir
```
Given   Siswa punya rata-rata tugas 80, rata-rata kuis 90
And     Bobot: tugas 60%, kuis 40%
When    Sistem hitung
Then    final_score = (80×0.6) + (90×0.4) = 84.0
And     grade_letter = "B" (80-89)
```

### AC-GRADE-02: Lock Rapor
```
Given   Admin klik "Lock Rapor" untuk semester 1
When    Lock diproses
Then    grades.is_locked = true untuk semua siswa
And     Guru tidak bisa edit nilai (API return 403)
And     PDF rapor bisa di-generate
```

---

## 7. Notifikasi

### AC-NOTIF-01: WA Notification
```
Given   Tugas baru dipublish untuk kelas 8A
When    Notifikasi dikirim
Then    Semua siswa 8A yang WA enabled menerima pesan WA
And     Pesan sesuai template (judul, mapel, deadline, link)
And     Delivery time < 5 menit
```

### AC-NOTIF-02: Quiet Hours
```
Given   User set quiet hours 22:00-06:00
When    Notifikasi LOW/MEDIUM di-trigger jam 23:00
Then    Notifikasi di-queue, dikirim jam 06:00
When    Notifikasi HIGH (security) di-trigger jam 23:00
Then    Dikirim langsung (ignore quiet hours)
```

---

## Browser & Performance Requirements

| Criteria | Target |
|:---|:---|
| Browser support | Chrome 90+, Firefox 90+, Safari 15+, Edge 90+ |
| Screen sizes | 375px, 768px, 1280px, 1920px |
| LCP | < 2.5 detik |
| FID | < 100ms |
| CLS | < 0.1 |
| TTI | < 3.5 detik |
| Lighthouse score | ≥ 90 (performance, accessibility) |

---

*Terakhir diperbarui: 21 Maret 2026*
