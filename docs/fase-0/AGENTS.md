# 🤖 AGENTS.md — Instruksi untuk AI Agent

> **Dokumen ini BUKAN untuk manusia.** Ini adalah instruksi eksplisit untuk AI coding agent (Cursor, Copilot, Gemini, dsb.) yang bekerja pada codebase AkuBelajar. Tanpa dokumen ini, AI agent bekerja secara "nebak" tentang konteksnya sendiri.

---

## Identitas Proyek

```yaml
nama: AkuBelajar
deskripsi: Platform edukasi digital AI-first, cross-platform, enterprise-grade
bahasa_utama: [Go 1.23+, TypeScript 5+, SQL (PostgreSQL)]
framework: [Gin (Go), Next.js 15 (React)]
database: PostgreSQL 16+ dengan Row-Level Security
cache: Redis 7+
ai: Google Gemini 2.0 Flash
style_css: Tailwind CSS v4 + Shadcn UI
```

---

## Aturan Global (WAJIB Dipatuhi)

### ✅ BOLEH

- Membuat file baru sesuai `FOLDER_STRUCTURE.md`
- Menambah endpoint baru mengikuti pattern di `API_SPEC.md`
- Membuat komponen UI baru mengikuti pattern Shadcn UI
- Menambah test baru mengikuti `TEST_CASES.md`
- Menambah migration baru (backward-compatible only)
- Menggunakan tools: `make`, `docker-compose`, `pnpm`

### ❌ TIDAK BOLEH

- **Mengubah skema database tanpa migration file** — tidak boleh edit langsung
- **Menghapus atau memodifikasi file `audit_logs`** — tabel ini immutable
- **Menyimpan secret/credential di source code** — gunakan `.env` atau Vault
- **Menggunakan `SELECT *`** — selalu explicit column selection
- **Menggunakan string concatenation untuk SQL** — selalu parameterized queries
- **Mengabaikan error di Go** — setiap `err` harus di-handle, tidak boleh `_ = err`
- **Menggunakan `any` di TypeScript** — strict typing wajib
- **Mengubah file di folder `migrations/` yang sudah dijalankan** — buat migration baru
- **Membuat endpoint tanpa auth middleware** — semua endpoint wajib diproteksi
- **Deploy langsung ke production** — harus melalui CI/CD pipeline

---

## Urutan Kerja yang Benar

### Saat Membuat Fitur Baru

```
1. Baca USER_STORIES.md → pahami kebutuhan user
2. Baca ACCEPTANCE_CRITERIA.md → pahami kapan fitur "selesai"
3. Baca API_SPEC.md → cek apakah endpoint sudah didefinisikan
4. Baca DATA_MODELS.md → pahami struct Go + TypeScript types
5. Baca BUSINESS_LOGIC.md → pahami aturan domain
6. Buat migration file (jika ada perubahan schema)
7. Buat/update Go service (backend)
8. Buat/update Next.js component (frontend)
9. Tulis test sesuai TEST_CASES.md
10. Verifikasi terhadap ACCEPTANCE_CRITERIA.md
```

### Saat Memperbaiki Bug

```
1. Reproduksi bug → pahami expected vs actual behavior
2. Baca BUSINESS_LOGIC.md → pastikan fix sesuai aturan domain
3. Tulis test yang menangkap bug (red)
4. Perbaiki kode (green)
5. Pastikan test lain tidak rusak
6. Update CHANGELOG.md jika bug user-facing
```

### Saat Berurusan dengan Security

```
1. Baca SECURE_CODING.md → ikuti semua aturan
2. Baca THREAT_MODEL.md → pahami attack surface
3. Validasi kode terhadap OWASP_CHECKLIST.md
4. Pastikan audit log ter-trigger untuk operasi kritis
```

---

## Scope per Sesi

AI agent harus membatasi scope perubahan per sesi kerja:

| Scope | Maksimal per Sesi | Contoh |
|:---|:---|:---|
| File yang diubah | ≤ 10 file | Satu fitur fokus |
| Migration file | ≤ 1 migration | Satu perubahan schema |
| Endpoint baru | ≤ 3 endpoint | CRUD satu entity |
| Komponen UI baru | ≤ 5 komponen | Satu halaman lengkap |

> **Prinsip:** Lebih baik perubahan kecil dan teruji daripada perubahan besar yang rapuh.

---

## Referensi File yang Harus Dibaca Sebelum Coding

| Prioritas | File | Kapan Dibaca |
|:---|:---|:---|
| 🔴 Wajib | `FOLDER_STRUCTURE.md` | Sebelum membuat file baru |
| 🔴 Wajib | `API_SPEC.md` | Sebelum membuat/mengubah endpoint |
| 🔴 Wajib | `DATA_MODELS.md` | Sebelum membuat struct/type baru |
| 🔴 Wajib | `SECURE_CODING.md` | Sebelum menulis kode apapun |
| 🟡 Penting | `BUSINESS_LOGIC.md` | Sebelum implementasi fitur domain |
| 🟡 Penting | `CODING_STANDARDS.md` | Untuk naming convention & style |
| 🟢 Referensi | `ACCEPTANCE_CRITERIA.md` | Untuk validasi fitur selesai |
| 🟢 Referensi | `TEST_CASES.md` | Untuk menulis test yang tepat |

---

## Context Window Management

Jika context window terbatas, prioritaskan membaca file dalam urutan ini:

```
1. AGENTS.md          (file ini — aturan kerja)
2. FOLDER_STRUCTURE.md (di mana menaruh file)
3. DATA_MODELS.md     (shape data di semua layer)
4. API_SPEC.md        (kontrak API)
5. SECURE_CODING.md   (aturan keamanan kode)
```

---

## Error Recovery

Jika AI agent mengalami kebingungan:

| Situasi | Tindakan |
|:---|:---|
| Tidak tahu file harus di mana | Baca `FOLDER_STRUCTURE.md` |
| Tidak tahu format request/response API | Baca `API_SPEC.md` |
| Tidak yakin field name di layer berbeda | Baca `DATA_MODELS.md` |
| Tidak tahu aturan bisnis | Baca `BUSINESS_LOGIC.md` |
| Tidak tahu apakah fitur sudah selesai | Baca `ACCEPTANCE_CRITERIA.md` |
| Kode terasa tidak aman | Baca `SECURE_CODING.md` |

---

*Terakhir diperbarui: 21 Maret 2026*
