# 📐 Business Logic — AkuBelajar

> Aturan bisnis akademik yang tidak bisa disimpulkan dari skema database saja. Setiap rule di sini adalah **keputusan domain yang eksplisit** — bukan teknis, melainkan kebijakan.

---

## 1. Penilaian & Kalkulasi Nilai

### Rumus Nilai Akhir

```
Nilai Akhir = (Bobot Tugas × Rata-rata Tugas) + (Bobot Kuis × Rata-rata Kuis)
```

| Komponen | Bobot Default | Configurable? |
|:---|:---|:---|
| Tugas (assignment) | 60% | ✅ Per sekolah via `schools.config` |
| Kuis / Ujian (quiz) | 40% | ✅ Per sekolah via `schools.config` |

**Contoh:**
```
Rata-rata Tugas = 80
Rata-rata Kuis  = 90
Nilai Akhir     = (0.6 × 80) + (0.4 × 90) = 48 + 36 = 84
```

### Konversi Huruf

| Rentang | Huruf | Predikat |
|:---|:---|:---|
| 90 - 100 | A | Sangat Baik |
| 80 - 89 | B | Baik |
| 70 - 79 | C | Cukup |
| 60 - 69 | D | Kurang |
| 0 - 59 | E | Sangat Kurang |

> ⚠️ Tabel konversi ini **configurable per sekolah**. Beberapa sekolah menggunakan skala A-E, lainnya A+/A/A- dll.

### Nilai KKM (Kriteria Ketuntasan Minimal)

- Default KKM: **70** (configurable per mata pelajaran per sekolah)
- Siswa dengan nilai akhir < KKM → ditandai **"Belum Tuntas"**
- Siswa belum tuntas bisa mengikuti remedial (max 1x per semester)

---

## 2. Kuis & CBT (Computer Based Test)

### Lifecycle Kuis

```
DRAFT → PUBLISHED → ACTIVE → ENDED → GRADED → LOCKED
  ↑                                              ↓
  └────── Guru bisa edit ──────┘     Nilai tidak bisa diubah
```

| State | Siapa Bisa Aksi | Apa yang Terjadi |
|:---|:---|:---|
| `DRAFT` | Guru | Edit soal, hapus, preview |
| `PUBLISHED` | Guru | Edit jadwal, tapi soal terkunci |
| `ACTIVE` | Siswa | Kerjakan kuis (saat `start_at` ≤ now ≤ `end_at`) |
| `ENDED` | Sistem | Auto-submit jika siswa belum submit |
| `GRADED` | Sistem | Semua jawaban sudah di-grade secara otomatis |
| `LOCKED` | — | Nilai final, tidak bisa diubah siapapun |

### Timer Rules

| Rule | Detail |
|:---|:---|
| Timer source | **Server-only** — waktu client diabaikan |
| Start time | Saat siswa klik "Start Quiz" (bukan saat publish) |
| Auto-submit | Jika waktu habis dan siswa belum submit |
| Grace period | 0 detik — waktu habis = langsung submit |
| Pause | **Tidak diizinkan** — timer berjalan terus |

### Anti-Cheating Rules

| Pelanggaran | Toleransi | Aksi |
|:---|:---|:---|
| Pindah tab (`visibilitychange`) | 3x peringatan | Ke-4: ujian auto-submit + flag |
| Minimize window (`blur`) | 2x peringatan | Ke-3: ujian auto-submit + flag |
| DevTools terbuka | 0 toleransi | Ujian langsung terkunci |
| WebSocket disconnect | 30 detik | >30 detik: session expired, auto-submit |
| IP address berubah mid-session | 0 toleransi | Session expired, perlu contact guru |

### Scoring Rules

| Tipe Soal | Cara Scoring |
|:---|:---|
| Pilihan Ganda | Hash comparison (Argon2id) — benar = 1, salah = 0 |
| Essay | AI grading (Gemini) + review manual guru |
| Soal tidak dijawab | Dianggap salah (skor 0) |

---

## 3. Tugas (Assignment)

### Deadline Rules

| Skenario | Apa yang Terjadi |
|:---|:---|
| Submit sebelum deadline | ✅ Diterima normal |
| Submit setelah deadline | ⚠️ Diterima dengan **penalti -10% per hari** (max -50%) |
| Submit > 5 hari setelah deadline | ❌ Ditolak otomatis |
| Belum submit sama sekali | Nilai = 0 setelah deadline + 5 hari |

**Contoh penalti:**
```
Nilai mentah: 85
Terlambat 2 hari → Penalti: -20%
Nilai final: 85 × (1 - 0.20) = 68
```

> ⚠️ Guru bisa **override penalti** secara manual per kasus (misal: sakit). Override dicatat di audit log.

### File Upload Rules

| Rule | Nilai |
|:---|:---|
| Max file size | 10 MB |
| Tipe file yang diizinkan | `.pdf`, `.doc`, `.docx`, `.jpg`, `.png`, `.txt` |
| Max files per submission | 5 |
| Resubmit | Diizinkan sebelum deadline (replace, bukan append) |

---

## 4. Absensi (Attendance)

### Status yang Diizinkan

| Status | Kode | Dihitung Hadir? |
|:---|:---|:---|
| Hadir | `present` | ✅ |
| Izin | `permission` | ✅ (dihitung kehadiran) |
| Sakit | `sick` | ✅ (dihitung kehadiran) |
| Alpha | `absent` | ❌ |

### Aturan Input

| Rule | Detail |
|:---|:---|
| Siapa yang input | Guru pengampu atau Ketua Kelas |
| Kapan bisa input | Hanya untuk hari ini (T+0) atau kemarin (T-1) |
| Edit absensi lama | Hanya oleh guru, max T-7 (seminggu ke belakang) |
| Perubahan dicatat | Setiap edit masuk audit log |

### Minimum Kehadiran

- Minimum kehadiran per semester: **75%**
- Siswa < 75% → **Tidak bisa mengikuti UAS**
- Perhitungan: `(present + permission + sick) / total_hari × 100%`

---

## 5. Rapor & Pelaporan

### Lock Rules

| Langkah | Detail |
|:---|:---|
| 1. Guru finalisasi nilai | Klik "Finalisasi" per mata pelajaran |
| 2. Sistem hitung nilai akhir | Berdasarkan bobot yang dikonfigurasi |
| 3. Wali kelas review | Melihat semua nilai per siswa |
| 4. Admin lock rapor | Setelah lock: **semua nilai immutable** |
| 5. Generate PDF | PDF dengan QR code verifikasi |
| 6. Distribusi | Signed URL dikirim ke orang tua/siswa |

### Post-Lock Rules

- ❌ Nilai **tidak bisa diubah** setelah rapor di-lock
- ❌ Absensi **tidak bisa diubah** setelah rapor di-lock
- ✅ Jika ada kesalahan → Admin harus **unlock** (dicatat di audit log) → perbaiki → lock ulang

---

## 6. Role & Permission Edge Cases

| Skenario | Keputusan |
|:---|:---|
| Guru pindah sekolah | Akun di-deactivate di sekolah lama, buat baru di sekolah baru |
| Siswa pindah kelas | Histori nilai tetap ada, assignment baru mengikuti kelas baru |
| Ketua kelas ganti | Admin update role, ketua lama kembali menjadi student biasa |
| Super admin mengedit data sekolah lain | ❌ Tidak diizinkan — RLS berdasarkan school_id berlaku untuk semua |
| Guru yang bukan pengajar mata pelajaran mencoba akses kuis | ❌ RLS memblokir — hanya guru yang linked ke subject tersebut |

---

## 7. AI Service Rules

| Rule | Detail |
|:---|:---|
| Input sanitization | Semua input ke Gemini melalui `sanitizeInput()` |
| Max prompt length | 150 karakter untuk topik |
| Max questions per call | 50 soal |
| Retry policy | 3x retry dengan exponential backoff (1s, 2s, 4s) |
| Fallback | Jika AI gagal → beri pesan error, guru buat manual |
| Cost control | Rate limit: 10 AI calls per guru per jam |

---

## Referensi

- [User Stories](../fase-5/USER_STORIES.md) — Kebutuhan user
- [Acceptance Criteria](../fase-2/ACCEPTANCE_CRITERIA.md) — Kapan fitur selesai
- [Data Models](../fase-1/DATA_MODELS.md) — Mapping field antar layer

---

*Terakhir diperbarui: 21 Maret 2026*
