# 🛠️ Setup Lokal — AkuBelajar

> Panduan setup environment pengembangan lokal. Target: developer baru bisa berjalan dalam **< 30 menit** tanpa bertanya.

---

## Prerequisites

| Tool | Versi Minimum | Cara Install |
|:---|:---|:---|
| **Go** | 1.23+ | [go.dev/dl](https://go.dev/dl/) |
| **Node.js** | 22+ | [nodejs.org](https://nodejs.org/) atau `nvm install 22` |
| **pnpm** | 9+ | `npm install -g pnpm` |
| **Docker** | 25+ | [docker.com](https://www.docker.com/) |
| **Docker Compose** | 2.20+ | Bundled dengan Docker Desktop |
| **Git** | 2.40+ | [git-scm.com](https://git-scm.com/) |
| **Make** | Any | Windows: `choco install make` / Linux: preinstalled |
| **VS Code** | Latest | [code.visualstudio.com](https://code.visualstudio.com/) |

---

## Langkah 1: Clone & Setup (5 menit)

```bash
# Clone repository
git clone https://github.com/kazanaru/akubelajar.git
cd akubelajar

# Copy environment file
cp .env.example .env
```

### Edit `.env`

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=akubelajar
DB_PASSWORD=localdev123
DB_NAME=akubelajar_dev
DB_SSL_MODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT & Paseto
JWT_SECRET=your-dev-jwt-secret-min-32-chars-long
PASETO_PRIVATE_KEY=your-dev-paseto-private-key

# Gemini AI (opsional untuk development)
GEMINI_API_KEY=your-gemini-api-key

# Environment
APP_ENV=development
APP_PORT=8080
FRONTEND_URL=http://localhost:3000
```

---

## Langkah 2: Start Infrastruktur (3 menit)

```bash
# Start PostgreSQL + Redis via Docker
docker-compose up -d postgres redis

# Verifikasi berjalan
docker-compose ps
```

---

## Langkah 3: Setup Database (3 menit)

```bash
# Jalankan migrasi database
make migrate-up

# Seed data awal (admin, guru, siswa default)
make seed

# Verifikasi
make db-check
```

---

## Langkah 4: Start Backend (3 menit)

```bash
# Install Go dependencies
go mod download

# Jalankan backend dengan hot-reload
make dev-api
# Backend berjalan di http://localhost:8080

# Atau tanpa hot-reload
make run-api
```

---

## Langkah 5: Start Frontend (5 menit)

```bash
cd frontend

# Install dependencies
pnpm install

# Jalankan development server
pnpm dev
# Frontend berjalan di http://localhost:3000
```

---

## Langkah 6: Verifikasi (2 menit)

| Check | URL / Command | Expected |
|:---|:---|:---|
| Backend Health | `curl http://localhost:8080/health` | `{"status":"ok"}` |
| Frontend | `http://localhost:3000` | Login page |
| PostgreSQL | `make db-check` | Connected |
| Redis | `docker exec akubelajar-redis redis-cli ping` | `PONG` |

### Login Default

| Role | Email | Password |
|:---|:---|:---|
| Super Admin | `admin@akubelajar.id` | `Admin@123!` |
| Guru | `guru@akubelajar.id` | `Guru@123!` |
| Siswa | `siswa@akubelajar.id` | `Siswa@123!` |

---

## VS Code Extensions (Recommended)

```json
// .vscode/extensions.json
{
  "recommendations": [
    "golang.go",
    "dbaeumer.vscode-eslint",
    "bradlc.vscode-tailwindcss",
    "prisma.prisma",
    "ms-vscode.vscode-typescript-next",
    "esbenp.prettier-vscode",
    "redhat.vscode-yaml",
    "ms-azuretools.vscode-docker"
  ]
}
```

---

## Troubleshooting

| Masalah | Solusi |
|:---|:---|
| Port 5432 sudah digunakan | Stop PostgreSQL lokal: `sudo systemctl stop postgresql` |
| Port 3000 sudah digunakan | Ubah port di `frontend/.env.local`: `PORT=3001` |
| `permission denied` saat Docker | Tambahkan user ke grup docker: `sudo usermod -aG docker $USER` |
| Go module error | Run `go clean -modcache && go mod download` |
| pnpm install gagal | Hapus `node_modules` dan `pnpm-lock.yaml`, lalu `pnpm install` |

---

## Referensi Terkait

- [Coding Standards](CODING_STANDARDS.md)
- [Git Workflow](GIT_WORKFLOW.md)
- [Backend Guide](BACKEND_GUIDE.md)
- [Frontend Guide](FRONTEND_GUIDE.md)

---

*Terakhir diperbarui: 21 Maret 2026*
