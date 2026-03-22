# 📜 DECISIONS_LOG.md — Catatan Keputusan Final

> Keputusan yang sudah diambil dan TIDAK BOLEH diubah oleh AI agent tanpa konfirmasi user.

---

## Cara Membaca

- ✅ **FINAL** — Sudah diputuskan, jangan ubah
- ⏳ **TENTATIVE** — Bisa berubah jika ada info baru
- Status tanggal = kapan keputusan dibuat

---

## Arsitektur & Stack

| # | Keputusan | Nilai | Alasan | Status | Tanggal |
|:-:|:---|:---|:---|:---|:---|
| 1 | Backend language | Go 1.23 | Performance, concurrency, cross-compile | ✅ FINAL | 2026-03-21 |
| 2 | Frontend framework | Next.js 15 (App Router) | SSR, RSC, Vercel native | ✅ FINAL | 2026-03-21 |
| 3 | Database | PostgreSQL 16 | RLS, JSON support, reliability | ✅ FINAL | 2026-03-21 |
| 4 | Cache | Redis 7 (Upstash) | Session, rate limit, notification queue | ✅ FINAL | 2026-03-21 |
| 5 | Primary key type | UUID v7 | Sortable by time, no collision | ✅ FINAL | 2026-03-21 |
| 6 | API framework (Go) | Gin | Battle-tested, fast, middleware ecosystem | ✅ FINAL | 2026-03-21 |
| 7 | DB driver (Go) | pgx v5 | Native PostgreSQL, prepared statements | ✅ FINAL | 2026-03-21 |

## Keamanan

| # | Keputusan | Nilai | Alasan | Status | Tanggal |
|:-:|:---|:---|:---|:---|:---|
| 8 | Auth token | Paseto v4 (BUKAN JWT) | No algorithm confusion, simpler, safer | ✅ FINAL | 2026-03-21 |
| 9 | Password hashing | Argon2id | Memory-hard, resists GPU cracking | ✅ FINAL | 2026-03-21 |
| 10 | Access token TTL | 15 menit | Short-lived, minimizes stolen token risk | ✅ FINAL | 2026-03-21 |
| 11 | Refresh token TTL | 7 hari | Balance UX vs security | ✅ FINAL | 2026-03-21 |
| 12 | Max concurrent sessions | 3 device | HP + laptop + device lain | ✅ FINAL | 2026-03-21 |
| 13 | Refresh token rotation | Wajib | Detect token theft via reuse detection | ✅ FINAL | 2026-03-21 |
| 14 | IP anomaly handling | Warning (bukan auto-revoke) | False positive tinggi di jaringan seluler ID | ✅ FINAL | 2026-03-21 |

## Bisnis & Akademik

| # | Keputusan | Nilai | Alasan | Status | Tanggal |
|:-:|:---|:---|:---|:---|:---|
| 15 | KKM default | 70 (configurable) | Standar nasional SMP/SMA | ✅ FINAL | 2026-03-21 |
| 16 | Kehadiran minimum | 75% | Standar Kemendikbud | ✅ FINAL | 2026-03-21 |
| 17 | Dampak kehadiran rendah | Hanya larangan UAS (bukan ke nilai) | Simplicity, sesuai aturan | ✅ FINAL | 2026-03-21 |
| 18 | Remedial cap | Ya (max = KKM) | Standar mayoritas sekolah | ✅ FINAL | 2026-03-21 |
| 19 | Bobot nilai | Configurable (default 60/40) | Tiap sekolah beda kebijakan | ✅ FINAL | 2026-03-21 |
| 20 | Role orang tua di MVP | Tidak ada | Complexity, postpone ke Fase 2 | ✅ FINAL | 2026-03-21 |
| 21 | Multi-tenant strategy | Soft (school_id + RLS) | Arsitektur siap tanpa re-architect | ✅ FINAL | 2026-03-21 |

## Hosting & Deployment

| # | Keputusan | Nilai | Alasan | Status | Tanggal |
|:-:|:---|:---|:---|:---|:---|
| 22 | Frontend hosting | Vercel (free) | Native Next.js, auto-deploy | ✅ FINAL | 2026-03-21 |
| 23 | Backend hosting | Render (free) | Docker support, WebSocket | ✅ FINAL | 2026-03-21 |
| 24 | Database hosting | Supabase (free) | PostgreSQL managed, 500MB | ✅ FINAL | 2026-03-21 |
| 25 | Redis hosting | Upstash (free) | Serverless Redis, 10K cmd/day | ✅ FINAL | 2026-03-21 |
| 26 | Email provider | Resend (free) | 100 email/day, simple API | ✅ FINAL | 2026-03-21 |
| 27 | WA notifications | Skip (demo) | Butuh bayar Fonnte, bukan prioritas demo | ✅ FINAL | 2026-03-21 |
| 28 | Custom domain | Skip (demo) | Pakai *.vercel.app + *.onrender.com | ✅ FINAL | 2026-03-21 |

## Lainnya

| # | Keputusan | Nilai | Alasan | Status | Tanggal |
|:-:|:---|:---|:---|:---|:---|
| 29 | Temp password TTL | 7 hari | Cukup waktu untuk user login pertama | ✅ FINAL | 2026-03-21 |
| 30 | Password change cooldown | 24 jam | Anti-abuse jika akun compromised | ✅ FINAL | 2026-03-21 |
| 31 | Storage per sekolah | 1 GB (demo) | Supabase free tier = 1GB | ✅ FINAL | 2026-03-21 |

---

*Terakhir diperbarui: 21 Maret 2026*
