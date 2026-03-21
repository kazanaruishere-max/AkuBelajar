# 📋 Changelog — AkuBelajar

Semua perubahan signifikan pada proyek ini didokumentasikan dalam file ini.

Format mengikuti [Keep a Changelog](https://keepachangelog.com/id-ID/1.1.0/), dan proyek ini menggunakan [Semantic Versioning](https://semver.org/lang/id/).

---

## Panduan Versioning

```
MAJOR.MINOR.PATCH

MAJOR  → Breaking change (API berubah, migration required)
MINOR  → Fitur baru (backward-compatible)
PATCH  → Bug fix & security patch
```

### Kategori Perubahan

| Label | Arti |
|:---|:---|
| **Added** | Fitur atau kapabilitas baru |
| **Changed** | Perubahan pada fitur yang sudah ada |
| **Deprecated** | Fitur yang akan dihapus di versi mendatang |
| **Removed** | Fitur yang dihapus |
| **Fixed** | Perbaikan bug |
| **Security** | Perbaikan kerentanan keamanan |

---

## [Unreleased]

### Added
- _(Fitur yang sedang dikembangkan dan belum dirilis)_

---

## [2.0.0] — 2026-03-21

### ⚠️ Breaking Changes
- Migrasi arsitektur dari Laravel Monolith ke **Go Microservices + Next.js 15**
- Database schema direset total — gunakan migration script untuk upgrade
- API endpoint berubah dari `/api/*` ke `/api/v1/*`
- Autentikasi berubah dari Laravel Sanctum ke **JWT + Paseto**

### Added
- 🤖 **AI Quiz Generator** via Google Gemini 2.0 Flash dengan JSON Mode
- 🔒 **Zero-Trust Security Architecture** — RBAC, RLS, WAF, rate limiting
- 📱 **PWA Support** — installable, offline-capable via Service Worker
- 🖥️ **CBT Interface** — Computer Based Test dengan anti-cheat detection
- 📊 **Early Warning System** — Prediksi siswa berisiko via AI analytics
- 📄 **Digital Report Card** — PDF rapor dengan QR code verifikasi
- 👥 **Bulk Import** — Async import massal 100K+ baris tanpa blocking
- 🔔 **Multi-Channel Notification** — Email, WhatsApp, In-App, Push
- 🌙 **Dark Mode** — Otomatis via `prefers-color-scheme`

### Changed
- Backend engine: PHP 8.3/Laravel → **Go 1.23+/Gin**
- Frontend engine: Blade Template → **Next.js 15 + TypeScript**
- Database driver: Eloquent ORM → **pgx (native PostgreSQL driver)**
- Password hashing: bcrypt → **Argon2id**
- Primary keys: sequential integer → **UUID v7**
- State management: Session-based → **JWT/Paseto stateless**

### Security
- Implementasi **Row-Level Security (RLS)** di PostgreSQL
- **Immutable Audit Log** — tabel yang tidak bisa di-UPDATE/DELETE
- **Prompt Injection Protection** pada AI service
- **Rate Limiting** per endpoint dengan Redis sliding window
- **CORS** strict policy + **CSP** headers

---

## [1.0.0] — 2025-09-01

### Added
- Sistem manajemen sekolah dasar (CRUD user, kelas, mata pelajaran)
- Manajemen tugas dan pengumpulan online
- Sistem absensi digital
- Dashboard admin, guru, dan siswa
- Export nilai ke Excel

### Stack
- Laravel 10 · PHP 8.3 · MySQL 8 · Blade Template

---

## Template Entry Baru

```markdown
## [X.Y.Z] — YYYY-MM-DD

### Added
- Deskripsi fitur baru (#nomor-issue)

### Changed
- Deskripsi perubahan (#nomor-issue)

### Fixed
- Deskripsi perbaikan (#nomor-issue)

### Security
- Deskripsi patch keamanan (#nomor-issue)
```

---

*Terakhir diperbarui: 21 Maret 2026*
