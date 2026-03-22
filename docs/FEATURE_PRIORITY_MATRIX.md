# 🎯 Feature Priority Matrix — AkuBelajar

> Dari 40+ fitur, mana yang MVP, mana yang Fase 2. AI agent mulai dari sini.

---

## Priority Legend

| Label | Arti | Target |
|:---|:---|:---|
| 🔴 **P0** | Must-have MVP — tanpa ini app tidak berfungsi | Sprint 1-3 |
| 🟡 **P1** | Should-have — penting tapi app bisa jalan tanpa ini | Sprint 4-6 |
| 🟢 **P2** | Nice-to-have — enhancement setelah MVP | Post-launch |
| ⚪ **P3** | Future — long-term roadmap | Fase 2+ |

---

## Authentication & User Management

| Fitur | Priority | Sprint | Notes |
|:---|:---|:---|:---|
| Login (email + password) | 🔴 P0 | 1 | Paseto v4 |
| Logout (single + all devices) | 🔴 P0 | 1 | |
| Register via invite token | 🔴 P0 | 1 | |
| First login force password change | 🔴 P0 | 1 | |
| Profile edit (nama, avatar) | 🟡 P1 | 4 | |
| Password reset (OTP via email) | 🟡 P1 | 4 | Via Resend |
| Bulk import user (Excel) | 🟡 P1 | 5 | |
| OAuth (Google) | 🟢 P2 | — | Post-MVP |
| Role Orang Tua | ⚪ P3 | — | Fase 2 |

## Academic Management

| Fitur | Priority | Sprint | Notes |
|:---|:---|:---|:---|
| CRUD Tahun Ajaran | 🔴 P0 | 1 | |
| CRUD Kelas | 🔴 P0 | 1 | |
| CRUD Mata Pelajaran | 🔴 P0 | 1 | |
| Assign guru ke kelas+mapel | 🔴 P0 | 1 | |
| Assign siswa ke kelas | 🔴 P0 | 1 | |
| Transisi tahun ajaran | 🟡 P1 | 6 | Naik kelas, lulus |
| Kenaikan/kelulusan batch | 🟡 P1 | 6 | |

## Assignments (Tugas)

| Fitur | Priority | Sprint | Notes |
|:---|:---|:---|:---|
| Guru buat tugas | 🔴 P0 | 2 | |
| Siswa submit tugas (file) | 🔴 P0 | 2 | |
| 📷 Foto tugas dari kamera | 🔴 P0 | 2 | No-mirror, kamera belakang default |
| 🖼️ Upload dari galeri | 🔴 P0 | 2 | Alternatif jika kamera error |
| Guru beri nilai + feedback | 🔴 P0 | 2 | |
| Deadline enforcement | 🔴 P0 | 2 | Late = rejected |
| Resubmit (jika guru allow) | 🟡 P1 | 5 | |

## Quiz / CBT

| Fitur | Priority | Sprint | Notes |
|:---|:---|:---|:---|
| Guru buat kuis manual | 🔴 P0 | 2 | |
| Siswa ikut kuis | 🔴 P0 | 2 | |
| Timer + auto-submit | 🔴 P0 | 2 | |
| AI generate soal (Gemini) | 🔴 P0 | 3 | Fitur utama |
| Pengacakan urutan soal | 🟡 P1 | 4 | |
| Anti-cheat (tab detection) | 🟡 P1 | 4 | |
| Review kuis setelah selesai | 🟡 P1 | 5 | |

## Attendance (Absensi)

| Fitur | Priority | Sprint | Notes |
|:---|:---|:---|:---|
| Guru input absensi | 🔴 P0 | 3 | |
| 📱 QR scan absensi (siswa) | 🟡 P1 | 4 | QR rotate 30 detik, fallback manual |
| Ketua Kelas input draft | 🟡 P1 | 4 | Perlu approval guru |
| Izin/sakit dengan lampiran | 🟡 P1 | 5 | File upload |
| Dashboard attendance rate | 🟡 P1 | 5 | |
| Early warning (alfa streak) | 🟢 P2 | — | Notifikasi otomatis |

## Grades & Report Card

| Fitur | Priority | Sprint | Notes |
|:---|:---|:---|:---|
| Auto hitung nilai akhir | 🔴 P0 | 3 | Bobot configurable |
| Dashboard nilai per siswa | 🔴 P0 | 3 | |
| Lock rapor oleh admin | 🟡 P1 | 6 | |
| Generate PDF rapor | 🟢 P2 | — | Chromium headless |
| QR verification rapor | 🟢 P2 | — | |

## Notifications

| Fitur | Priority | Sprint | Notes |
|:---|:---|:---|:---|
| In-app notification | 🔴 P0 | 3 | WebSocket real-time |
| Email notification | 🟡 P1 | 4 | Resend API |
| WA notification | ⚪ P3 | — | Butuh Fonnte (berbayar) |
| Notification preferences | 🟡 P1 | 5 | |

## Infrastructure

| Fitur | Priority | Sprint | Notes |
|:---|:---|:---|:---|
| Audit log | 🟡 P1 | 4 | |
| Rate limiting | 🟡 P1 | 4 | |
| File upload (Supabase Storage) | 🔴 P0 | 2 | Tugas + avatar |
| Dark mode | 🟢 P2 | — | |
| PWA offline | ⚪ P3 | — | Fase 2 |
| Mobile app (React Native) | ⚪ P3 | — | Fase 2 |

---

## Sprint Roadmap (MVP = Sprint 1-3)

| Sprint | Fokus | Duration |
|:---|:---|:---|
| **Sprint 1** | Auth + Academic CRUD + Skeleton UI | 2 minggu |
| **Sprint 2** | Assignments + Quiz + File Upload | 2 minggu |
| **Sprint 3** | AI Quiz + Attendance + Grades + Notifications | 2 minggu |
| Sprint 4 | Polish: password reset, bulk import, anti-cheat | 2 minggu |
| Sprint 5 | Review, resubmit, preferences, dashboard | 2 minggu |
| Sprint 6 | Tahun ajaran transition, lock rapor, audit | 2 minggu |

**MVP target: 6 minggu development**

---

*Terakhir diperbarui: 22 Maret 2026*
