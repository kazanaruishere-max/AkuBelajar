# 🤝 Panduan Kontribusi — AkuBelajar

Terima kasih telah tertarik berkontribusi! Dokumen ini mendefinisikan aturan dan alur kerja agar semua kontributor bekerja dengan standar yang sama.

---

## Daftar Isi

1. [Prinsip Kontribusi](#prinsip-kontribusi)
2. [Alur Kontribusi (Workflow)](#alur-kontribusi-workflow)
3. [Standar Branch & Commit](#standar-branch--commit)
4. [Pull Request (PR) Guidelines](#pull-request-pr-guidelines)
5. [Code Review Checklist](#code-review-checklist)
6. [Issue & Bug Report](#issue--bug-report)
7. [Lingkungan Pengembangan](#lingkungan-pengembangan)

---

## Prinsip Kontribusi

1. **Satu PR, Satu Tujuan** — Setiap PR hanya menyelesaikan satu issue atau satu fitur.
2. **Test Wajib** — Tidak ada PR yang di-merge tanpa test yang memadai.
3. **Review Wajib** — Minimal 1 reviewer untuk PR biasa, 2 reviewer untuk perubahan yang menyentuh keamanan atau database.
4. **Breaking Change = ADR** — Setiap breaking change harus disertai Architecture Decision Record.
5. **Bahasa** — Kode dan comment dalam Bahasa Inggris. Dokumentasi dalam Bahasa Indonesia.

---

## Alur Kontribusi (Workflow)

```
1. Fork & Clone repository
         ↓
2. Buat branch dari `develop`
         ↓
3. Tulis kode + test
         ↓
4. Pastikan `make lint` dan `make test` lolos
         ↓
5. Commit dengan format Conventional Commits
         ↓
6. Push ke fork & buat Pull Request ke `develop`
         ↓
7. Tunggu review & perbaiki feedback
         ↓
8. Merge setelah approved ✅
```

---

## Standar Branch & Commit

### Format Nama Branch

```
<tipe>/<nomor-issue>-<deskripsi-singkat>
```

| Tipe | Contoh | Keterangan |
|:---|:---|:---|
| `feature/` | `feature/42-ai-quiz-generator` | Fitur baru |
| `fix/` | `fix/78-login-rate-limit` | Perbaikan bug |
| `hotfix/` | `hotfix/99-xss-vulnerability` | Fix kritis di production |
| `docs/` | `docs/15-setup-local` | Perubahan dokumentasi |
| `refactor/` | `refactor/33-auth-middleware` | Refaktor tanpa perubahan fungsional |
| `test/` | `test/55-quiz-service-unit` | Penambahan/perbaikan test |

### Format Commit Message (Conventional Commits)

```
<tipe>(<scope>): <deskripsi singkat>

[body opsional — jelaskan MENGAPA, bukan APA]

[footer opsional — referensi issue]
```

**Tipe yang diizinkan:**

| Tipe | Deskripsi |
|:---|:---|
| `feat` | Fitur baru |
| `fix` | Perbaikan bug |
| `docs` | Perubahan dokumentasi saja |
| `style` | Perubahan formatting (bukan logic) |
| `refactor` | Perubahan kode tanpa perbaikan bug atau fitur baru |
| `perf` | Peningkatan performa |
| `test` | Penambahan atau perbaikan test |
| `chore` | Perubahan build, CI, atau dependency |
| `security` | Perbaikan terkait keamanan |

**Contoh:**

```
feat(quiz): add AI-powered quiz generation via Gemini 2.0

Implement structured output prompt for guaranteed JSON parsing.
Uses Gemini 2.0 Flash JSON mode to eliminate parsing failures.

Closes #42
```

---

## Pull Request (PR) Guidelines

### Template PR

```markdown
## Deskripsi
<!-- Jelaskan apa yang berubah dan MENGAPA -->

## Tipe Perubahan
- [ ] 🆕 Fitur baru
- [ ] 🐛 Perbaikan bug
- [ ] 📝 Dokumentasi
- [ ] ♻️ Refactor
- [ ] 🔒 Keamanan

## Checklist
- [ ] Test baru ditambahkan/diperbarui
- [ ] `make lint` lolos tanpa error
- [ ] `make test` lolos semua
- [ ] Dokumentasi diperbarui (jika mengubah API/behavior)
- [ ] CHANGELOG.md diperbarui
- [ ] Tidak ada secret/credential yang ter-commit

## Screenshot (jika UI berubah)
<!-- Lampirkan screenshot before/after -->

## Issue Terkait
Closes #___
```

### Aturan Merge

| Kondisi | Persyaratan |
|:---|:---|
| Kode biasa | ≥ 1 approval |
| Perubahan database/schema | ≥ 2 approval + DBA review |
| Perubahan security/auth | ≥ 2 approval + security review |
| Breaking change | ADR wajib + ≥ 2 approval |

---

## Code Review Checklist

Reviewer harus memeriksa:

- [ ] **Kebenaran logika** — Apakah kode menyelesaikan masalah yang dimaksud?
- [ ] **Keamanan** — Ada SQL injection, XSS, IDOR, atau kerentanan lain?
- [ ] **Performa** — Ada N+1 query, memory leak, atau operasi O(n²) yang tidak perlu?
- [ ] **Test** — Apakah test mencakup happy path AND edge cases?
- [ ] **Naming** — Apakah nama variabel, fungsi, dan file konsisten dengan codebase?
- [ ] **Error handling** — Apakah semua error di-handle, bukan di-ignore?
- [ ] **Dokumentasi** — Apakah perubahan API terdokumentasi?

---

## Issue & Bug Report

### Template Bug Report

```markdown
## Deskripsi Bug
<!-- Jelaskan bug secara singkat -->

## Langkah Reproduksi
1. Buka halaman '...'
2. Klik '...'
3. Perhatikan bahwa '...'

## Expected Behavior
<!-- Apa yang seharusnya terjadi -->

## Actual Behavior
<!-- Apa yang sebenarnya terjadi -->

## Environment
- OS: [e.g., Windows 11, macOS 15]
- Browser: [e.g., Chrome 130]
- Versi AkuBelajar: [e.g., v2.0.0]

## Screenshot / Log
<!-- Lampirkan jika ada -->
```

---

## Lingkungan Pengembangan

Lihat **[fase-2/SETUP_LOCAL.md](../fase-2/SETUP_LOCAL.md)** untuk panduan setup lengkap.

**Minimum yang diperlukan:**

| Tool | Versi Minimum |
|:---|:---|
| Go | 1.23+ |
| Node.js | 22+ |
| pnpm | 9+ |
| Docker | 25+ |
| PostgreSQL | 16+ |
| Redis | 7+ |

---

*Terakhir diperbarui: 21 Maret 2026*
