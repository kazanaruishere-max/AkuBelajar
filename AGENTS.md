# 🤖 AGENTS.md — Instruksi untuk AI Agent

> File ini adalah **sumber kebenaran pertama** yang harus dibaca AI agent sebelum menyentuh codebase.

---

## Identitas Proyek

| Key | Value |
|:---|:---|
| Nama | AkuBelajar |
| Deskripsi | Platform edukasi digital untuk SMP/SMA di Indonesia |
| Backend | Go 1.23, Gin, PostgreSQL 16, Redis 7 |
| Frontend | Next.js 15 (App Router), TypeScript, Zustand, TanStack Query |
| Auth | Paseto v4 (access token) + httpOnly cookie (refresh token) |
| AI | Google Gemini 2.0 Flash (quiz generation) |
| Hosting (demo) | Vercel (FE) + Render (BE) + Supabase (DB) |

---

## Sumber Kebenaran — Baca Ini Dulu

| Urutan | File | Mengapa |
|:---|:---|:---|
| 1 | `PROJECT_CONTEXT.md` | Ringkasan 1 halaman — baca sebelum apa pun |
| 2 | `DECISIONS_LOG.md` | Keputusan yang sudah final — JANGAN ubah |
| 3 | `FORBIDDEN_PATTERNS.md` | Pola kode yang DILARANG |
| 4 | `docs/INDEX.md` | Navigasi ke 77+ dokumen detail |
| 5 | `docs/fase-0/AGENTS.md` | Instruksi detail scope & rules |

---

## Rules Wajib

### ✅ HARUS

1. Baca `PROJECT_CONTEXT.md` di awal setiap sesi
2. Ikuti naming convention: `snake_case` (DB/Go), `camelCase` (TypeScript)
3. Semua endpoint gunakan prefix `/api/v1/`
4. Log setiap error dengan `request_id` (structured JSON)
5. Validasi input di **2 layer**: frontend (Zod) DAN backend (Go validator)
6. Hash password dengan Argon2id — TIDAK ADA pengecualian
7. Gunakan `UUID v7` untuk semua primary key
8. Wrap operasi multi-tabel dalam transaction
9. Return error dalam format standar: `{ "code": "XXX_001", "message": "..." }`
10. Tulis test untuk setiap handler baru (minimum: happy path + auth + validation error)

### ❌ DILARANG

1. JANGAN gunakan JWT — kita pakai Paseto v4
2. JANGAN simpan password/token/OTP di log
3. JANGAN raw SQL tanpa parameterized query (SQL injection risk)
4. JANGAN hardcode credentials — semua dari environment variable
5. JANGAN `SELECT *` — selalu list kolom eksplisit
6. JANGAN skip error handling (`if err != nil` wajib di-handle)
7. JANGAN buat file di luar folder structure yang sudah didefinisikan
8. JANGAN ubah keputusan yang ada di `DECISIONS_LOG.md`
9. JANGAN install dependency baru tanpa minta konfirmasi user
10. JANGAN bypass RLS — semua query harus melalui session context

---

## Urutan Kerja per Sesi

```
1. Baca PROJECT_CONTEXT.md
2. Baca DECISIONS_LOG.md (cek update terakhir)
3. Baca FORBIDDEN_PATTERNS.md
4. Baca file yang relevan dengan task saat ini
5. Implementasi sesuai docs
6. Tulis/update test
7. Verifikasi: go build, go test, npm run lint
```

---

## Folder Structure Ringkas

```
akubelajar/
├── backend/
│   ├── cmd/api/             ← entry point
│   ├── internal/
│   │   ├── handler/         ← HTTP handlers
│   │   ├── service/         ← business logic
│   │   ├── repository/      ← DB queries
│   │   ├── model/           ← Go structs
│   │   ├── middleware/       ← auth, RBAC, rate limit
│   │   └── pkg/             ← shared utilities
│   └── migrations/          ← SQL migration files
├── frontend/
│   ├── src/app/             ← Next.js App Router pages
│   ├── src/components/      ← React components
│   ├── src/lib/             ← API client, auth, utils
│   ├── src/hooks/           ← custom hooks
│   └── src/stores/          ← Zustand stores
└── docs/                    ← 77 dokumentasi files
```

Referensi lengkap: `docs/fase-2/FOLDER_STRUCTURE.md`

---

## Kontak & Eskalasi

Jika menemukan ambiguitas yang tidak bisa dijawab oleh dokumentasi, **JANGAN asumsi** — tanyakan ke user.

---

*Terakhir diperbarui: 21 Maret 2026*
