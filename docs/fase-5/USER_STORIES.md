# 👤 User Stories — AkuBelajar

> Kebutuhan nyata pengguna sebagai anchor keputusan teknis — setiap fitur harus bisa ditelusuri kembali ke user story.

---

## Format User Story

```
Sebagai [ROLE],
saya ingin [TINDAKAN],
agar [MANFAAT].

Kriteria Penerimaan:
- [ ] Syarat 1
- [ ] Syarat 2
```

---

## 🔴 Prioritas Tinggi (MVP)

### US-001: Login Aman

```
Sebagai PENGGUNA (semua role),
saya ingin login dengan email dan password secara aman,
agar data saya terlindungi dari akses tidak sah.

Kriteria Penerimaan:
- [x] Login dengan email + password (Argon2id)
- [x] Rate limit: maks 5 attempts per 15 menit
- [x] Account lockout setelah 5 gagal berturut-turut
- [x] JWT/Paseto token dengan expiry 15 menit
- [x] Refresh token rotation
```

### US-002: Buat Kuis AI

```
Sebagai GURU,
saya ingin membuat soal kuis secara otomatis menggunakan AI,
agar saya menghemat waktu dan bisa fokus mengajar.

Kriteria Penerimaan:
- [x] Input: topik, jumlah soal, tingkat kesulitan
- [x] AI menghasilkan soal pilihan ganda yang valid
- [x] Guru bisa preview dan edit sebelum publish
- [x] Jawaban ter-hash di database (tidak readable)
- [x] Soal sesuai kurikulum (K-13/Merdeka)
```

### US-003: Kerjakan Kuis (CBT)

```
Sebagai SISWA,
saya ingin mengerjakan kuis online dengan timer yang adil,
agar saya bisa menunjukkan kemampuan saya tanpa kecurangan.

Kriteria Penerimaan:
- [x] Timer dikontrol server (anti-manipulasi)
- [x] Deteksi pindah tab / minimize browser
- [x] Jawaban auto-submit saat waktu habis
- [x] Nilai instan setelah submit
- [x] Pembahasan AI tersedia setelah submit
```

### US-004: Input Absensi

```
Sebagai GURU atau KETUA KELAS,
saya ingin mencatat kehadiran siswa dengan cepat,
agar data absensi selalu akurat dan real-time.

Kriteria Penerimaan:
- [x] Input absensi per kelas per hari
- [x] Status: Hadir, Izin, Sakit, Alpha
- [x] Bisa diakses dari mobile (PWA)
- [x] Tersimpan meski offline (sync saat online)
- [x] Audit trail untuk setiap perubahan
```

### US-005: Lihat Nilai

```
Sebagai SISWA,
saya ingin melihat nilai saya secara real-time,
agar saya bisa memantau perkembangan akademik saya.

Kriteria Penerimaan:
- [x] Hanya bisa melihat nilai sendiri (RLS)
- [x] Nilai per mata pelajaran, per periode
- [x] Grafik tren nilai dari waktu ke waktu
- [x] Notifikasi saat nilai baru diinput
```

---

## 🟡 Prioritas Sedang (Post-MVP)

### US-006: Bulk Import User

```
Sebagai SUPER ADMIN,
saya ingin mengimport data guru dan siswa secara massal dari file Excel/CSV,
agar onboarding sekolah baru bisa dilakukan dalam hitungan menit.

Kriteria Penerimaan:
- [x] Upload file CSV/Excel
- [x] Validasi per baris dengan error report
- [x] Proses async (tidak memblokir UI)
- [x] Auto-generate password + kirim via email
- [x] Idempotent (aman dijalankan ulang)
```

### US-007: Rapor Digital

```
Sebagai SEKOLAH,
saya ingin mencetak rapor dalam format PDF yang sah dan terverifikasi,
agar rapor tidak bisa dipalsukan.

Kriteria Penerimaan:
- [x] PDF rapor dengan layout sesuai template sekolah
- [x] QR code untuk verifikasi keaslian online
- [x] Timestamp kriptografis
- [x] Bobot nilai bisa dikonfigurasi per sekolah
- [x] Nilai di-lock setelah rapor diterbitkan
```

### US-008: Notifikasi Multi-Channel

```
Sebagai PENGGUNA (semua role),
saya ingin menerima notifikasi penting melalui berbagai channel,
agar saya tidak melewatkan informasi penting.

Kriteria Penerimaan:
- [ ] In-app notification (real-time via WebSocket)
- [ ] Email notification
- [ ] WhatsApp notification
- [ ] Push notification (PWA)
- [ ] Pengaturan preferensi notifikasi per user
```

---

## 🟢 Prioritas Rendah (Future)

### US-009: Early Warning System

```
Sebagai GURU,
saya ingin mengetahui siswa mana yang berisiko tidak naik kelas,
agar saya bisa melakukan intervensi lebih awal.

Kriteria Penerimaan:
- [ ] AI analisis tren nilai, kehadiran, keterlambatan tugas
- [ ] Risk level: Low / Medium / High
- [ ] Rekomendasi intervensi
- [ ] Dashboard visualisasi per kelas
```

### US-010: Multi-Sekolah (SaaS)

```
Sebagai PLATFORM OWNER,
saya ingin beberapa sekolah bisa menggunakan platform secara mandiri,
agar bisnis bisa di-scale tanpa deploy ulang per sekolah.

Kriteria Penerimaan:
- [ ] Self-service onboarding sekolah baru
- [ ] Data isolasi antar sekolah (RLS)
- [ ] Custom branding per sekolah
- [ ] Tier pricing (Free, Basic, Premium)
```

---

## Traceability Matrix

| User Story | Fitur | ADR | Test |
|:---|:---|:---|:---|
| US-001 | Auth Service | ADR-001 (Paseto vs JWT) | `auth_test.go` |
| US-002 | AI Quiz Generator | — | `ai_service_test.go` |
| US-003 | CBT Interface | — | `quiz_handler_test.go` |
| US-004 | Attendance Module | — | `attendance_test.go` |
| US-005 | Grades Module | — | `grade_test.go` |
| US-006 | Bulk Import | — | `bulk_import_test.go` |
| US-007 | PDF Rapor | — | `report_test.go` |
| US-010 | Multitenancy | — | [Multitenancy Plan](MULTITENANCY_PLAN.md) |

---

*Terakhir diperbarui: 21 Maret 2026*
