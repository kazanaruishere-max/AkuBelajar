# 📖 README — AkuBelajar

> **Platform Edukasi Digital Generasi Berikutnya**  
> AI-First • Cross-Platform • Enterprise-Grade Security

---

## Apa itu AkuBelajar?

**AkuBelajar** adalah platform manajemen edukasi digital yang dirancang untuk era **"Zero Paper"**, **"AI-First"**, dan **"Mobile-Native"**. Platform ini menjawab kebutuhan nyata institusi pendidikan di Indonesia dengan menghadirkan:

- 🤖 **Quiz Generator AI** — Soal otomatis via Google Gemini 2.0 Flash
- 🔒 **Zero-Trust Security** — RBAC, RLS, WAF, rate limiting
- 📱 **Cross-Platform** — PWA + React Native (Expo)
- ⚡ **High Performance** — Go backend, 50K req/s per instance

---

## Tech Stack

| Layer | Teknologi |
|:---|:---|
| **Frontend** | Next.js 15 · TypeScript 5+ · Tailwind CSS v4 · Shadcn UI |
| **Backend** | Go 1.23+ · Gin Framework · pgx (PostgreSQL driver) |
| **Database** | PostgreSQL 16+ · Redis 7+ |
| **AI** | Google Gemini 2.0 Flash |
| **Infra** | Docker · Kubernetes · GitHub Actions · Cloudflare |
| **Monitoring** | Grafana · Prometheus · Loki |

---

## Quick Start (5 Menit)

### Prerequisites

- Docker & Docker Compose
- Go 1.23+
- Node.js 22+ & pnpm
- PostgreSQL 16+ (atau via Docker)

### Instalasi

```bash
# 1. Clone repository
git clone https://github.com/kazanaru/akubelajar.git
cd akubelajar

# 2. Copy environment
cp .env.example .env
# Edit .env: isi DB_*, REDIS_*, GEMINI_API_KEY

# 3. Start infrastruktur
docker-compose up -d postgres redis

# 4. Jalankan migrasi database
make migrate-up

# 5. Seed data awal
make seed

# 6. Jalankan backend
make run-api

# 7. Jalankan frontend (terminal terpisah)
cd frontend && pnpm install && pnpm dev
```

### Akses Default

| Role | Email | Password |
|:---|:---|:---|
| Super Admin | `admin@akubelajar.id` | `Admin@123!` |
| Guru | `guru@akubelajar.id` | `Guru@123!` |
| Siswa | `siswa@akubelajar.id` | `Siswa@123!` |

> ⚠️ **Ganti semua password default sebelum production!**

---

## Struktur Project

```
akubelajar/
├── backend/               # Go API services
├── frontend/              # Next.js 15 web app
├── mobile/                # React Native (Expo) - Roadmap
├── docs/                  # 📚 Dokumentasi lengkap (Anda di sini)
│   ├── INDEX.md           # Indeks semua dokumen
│   ├── fase-0/            # Pondasi & Perencanaan
│   ├── fase-1/            # Arsitektur & Desain
│   ├── fase-2/            # Development & Standar
│   ├── fase-3/            # Keamanan & Compliance
│   ├── fase-4/            # Deployment & Ops
│   └── fase-5/            # Pertumbuhan Jangka Panjang
├── docker-compose.yml
├── Makefile
└── AKUBELAJAR.md          # Jurnal teknis utama
```

---

## Dokumentasi Lengkap

Lihat **[docs/INDEX.md](../INDEX.md)** untuk navigasi lengkap semua dokumen yang disusun per fase pengembangan.

---

## Kontribusi

Baca **[CONTRIBUTING.md](CONTRIBUTING.md)** sebelum membuat Pull Request.

## Keamanan

Temukan celah keamanan? Baca **[SECURITY.md](SECURITY.md)** untuk responsible disclosure.

## Lisensi

Proyek ini dilisensikan di bawah **MIT License**. Lihat **[LICENSE.md](LICENSE.md)**.

---

## Kontak

| Kanal | Alamat |
|:---|:---|
| Email | kazanaru@akubelajar.id |
| GitHub | [@kazanaru](https://github.com/kazanaru) |

---

*Terakhir diperbarui: 21 Maret 2026*
