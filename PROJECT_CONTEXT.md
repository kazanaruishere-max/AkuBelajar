# 📋 PROJECT_CONTEXT.md — Ringkasan 1 Halaman

> Baca ini di awal setiap sesi. Ringkasan lengkap agar agent tidak perlu baca 1600+ baris journal.

---

## Apa Ini?

**AkuBelajar** adalah platform edukasi digital untuk SMP/SMA di Indonesia. Fitur utama: manajemen akademik (kelas, mapel, nilai), ujian online (CBT) dengan anti-cheat, absensi digital, dan AI quiz generation via Gemini.

## Siapa Penggunanya?

| Role | Akses |
|:---|:---|
| **SuperAdmin** | Kelola seluruh sekolah: user, kelas, mapel, tahun ajaran |
| **Guru (Teacher)** | Buat tugas, kuis, input nilai, input absensi |
| **Ketua Kelas (ClassLeader)** | Input draft absensi, submit tugas |
| **Siswa (Student)** | Lihat nilai, ikut kuis, submit tugas |

## Stack Teknis

| Layer | Teknologi |
|:---|:---|
| Frontend | Next.js 15 (App Router), TypeScript, Zustand, TanStack Query |
| Backend | Go 1.23, Gin, `pgx` (PostgreSQL driver) |
| Database | PostgreSQL 16 + RLS, UUID v7 PKs |
| Cache | Redis 7 (session, rate limit, queue) |
| Auth | Paseto v4 access token (15 min) + httpOnly refresh token (7 hari) |
| AI | Gemini 2.0 Flash (quiz generation, content safety filter) |
| File Storage | Supabase Storage (demo) / MinIO (production) |
| Email | Resend (demo) |
| Hosting | Vercel (FE) + Render (BE) + Supabase (DB) — **Rp 0/bulan** |

## Keputusan Kunci (Sudah Final)

| Keputusan | Nilai |
|:---|:---|
| KKM default | 70 (configurable per sekolah & mapel) |
| Kehadiran minimum | 75% (dampak: tidak boleh ikut UAS) |
| Max session per user | 3 device |
| Bobot nilai | Configurable (default 60% tugas / 40% kuis) |
| Role orang tua di MVP | Tidak (roadmap Fase 2) |
| Multi-tenant | Ya (soft — `school_id` di semua tabel + RLS) |
| Password hashing | Argon2id (BUKAN bcrypt) |
| Remedial cap | Ya (max = KKM = 70) |
| Token type | Paseto v4 (BUKAN JWT) |

## Arsitektur 30 Detik

```
Browser → Vercel (Next.js SSR) → Render (Go API) → Supabase (PostgreSQL)
                                       ↕
                                 Upstash (Redis)
                                       ↕
                              Gemini AI / Resend Email
```

## Dokumentasi: 77 File

Lihat `docs/INDEX.md` untuk navigasi lengkap. Highlights:
- **API contracts**: `docs/api/API_SPEC_FULL.md`
- **DB Schema**: `docs/database/DATABASE_SCHEMA_FULL.md` (26 tabel)
- **RBAC Matrix**: `docs/security/RBAC_MATRIX.md`
- **Business Logic**: `docs/database/BUSINESS_LOGIC_FULL.md`

---

*Terakhir diperbarui: 21 Maret 2026*
