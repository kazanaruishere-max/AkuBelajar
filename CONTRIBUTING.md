# 🤝 Panduan Kontribusi

Terima kasih atas ketertarikan Anda untuk berkontribusi pada **AkuBelajar**! Berikut panduan untuk membantu Anda memulai.

## 📋 Daftar Isi

- [Kode Etik](#kode-etik)
- [Cara Berkontribusi](#cara-berkontribusi)
- [Setup Development](#setup-development)
- [Coding Standards](#coding-standards)
- [Commit Convention](#commit-convention)
- [Pull Request](#pull-request)

## Kode Etik

Proyek ini mengikuti [Kode Etik Kontributor](CODE_OF_CONDUCT.md). Dengan berpartisipasi, Anda diharapkan untuk mematuhi kode etik ini.

## Cara Berkontribusi

### 🐛 Melaporkan Bug

1. Cek apakah bug sudah pernah dilaporkan di [Issues](../../issues)
2. Jika belum, buat issue baru dengan template **Bug Report**
3. Sertakan langkah reproduksi, expected behavior, dan screenshots jika ada

### 💡 Mengusulkan Fitur

1. Cek apakah fitur sudah pernah diusulkan di [Issues](../../issues)
2. Buat issue baru dengan template **Feature Request**
3. Jelaskan use case dan manfaat fitur tersebut

### 🔧 Kontribusi Kode

1. Fork repository ini
2. Buat branch baru: `git checkout -b feature/nama-fitur`
3. Lakukan perubahan
4. Commit dengan pesan yang deskriptif
5. Push ke fork Anda: `git push origin feature/nama-fitur`
6. Buat Pull Request

## Setup Development

### Prasyarat

- Go 1.23+
- Node.js 20+ & pnpm 9+
- PostgreSQL 16 (via Supabase)

### Backend

```bash
cd backend
cp .env.example .env
# Edit .env sesuai konfigurasi Anda
go mod tidy
go run cmd/seed/main.go
go run cmd/api/main.go
```

### Frontend

```bash
cd frontend
pnpm install
echo "NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1" > .env.local
pnpm dev
```

## Coding Standards

### Go (Backend)

- Ikuti [Effective Go](https://golang.org/doc/effective_go) guidelines
- Gunakan `gofmt` untuk formatting
- Gunakan error handling yang eksplisit (jangan `panic`)
- Setiap module baru harus mengikuti struktur: `handler.go`, `service.go`, `repository.go`, `model.go`, `routes.go`
- Parameterized queries wajib (tanpa string concatenation di SQL)

### TypeScript (Frontend)

- Strict mode enabled
- Gunakan `type` untuk props, `interface` untuk API responses
- Zustand untuk state management (bukan Context API)
- API calls harus melalui `lib/api/client.ts`

## Commit Convention

Kami mengikuti [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]
```

### Types

| Type | Deskripsi |
|:--|:--|
| `feat` | Fitur baru |
| `fix` | Perbaikan bug |
| `docs` | Perubahan dokumentasi |
| `style` | Formatting (tanpa perubahan logika) |
| `refactor` | Refactoring kode |
| `test` | Menambah/memperbaiki tes |
| `chore` | Maintenance (dependencies, config) |

### Contoh

```
feat(quiz): add AI quiz generation via Gemini API
fix(auth): handle PgBouncer prepared statement error
docs(readme): add architecture diagram and flowcharts
chore(deps): update pgx to v5.7.2
```

## Pull Request

### Checklist

- [ ] Kode sudah di-format (`gofmt` / `prettier`)
- [ ] Tidak ada error saat build (`go build ./...` / `pnpm build`)
- [ ] Commit message mengikuti convention
- [ ] Perubahan sudah ditest secara lokal
- [ ] Dokumentasi sudah diupdate jika perlu

### Review Process

1. Minimal 1 approval diperlukan
2. CI checks harus passing
3. Tidak boleh ada conflict dengan branch utama

---

Terima kasih telah berkontribusi! 🎉
