# 🔀 Git Workflow — AkuBelajar

> Branching strategy, format commit message, dan aturan merge yang wajib diikuti.

---

## Branching Strategy (Git Flow Simplified)

```
main          ─────────●─────────────────●───────────── (production releases)
                       ↑                 ↑
release/2.1   ─────────┤     ────────────┘
                       ↑
develop       ───●───●───●───●───●───●───●──────────── (integration branch)
                 ↑   ↑       ↑       ↑
feature/42   ────┘   │       │       │
fix/78       ────────┘       │       │
feature/55   ────────────────┘       │
hotfix/99    ────────────────────────┘ (langsung ke main + develop)
```

### Branch Utama

| Branch | Tujuan | Deploy ke |
|:---|:---|:---|
| `main` | Kode production yang stabil | Production |
| `develop` | Integrasi semua fitur terbaru | Staging |
| `release/x.y` | Persiapan rilis (freeze fitur, bug fix only) | Pre-production |

### Branch Kerja

| Prefix | Tujuan | Base Branch | Merge ke |
|:---|:---|:---|:---|
| `feature/` | Fitur baru | `develop` | `develop` |
| `fix/` | Bug fix | `develop` | `develop` |
| `hotfix/` | Fix kritis production | `main` | `main` + `develop` |
| `docs/` | Dokumentasi | `develop` | `develop` |
| `refactor/` | Refactor tanpa perubahan fungsional | `develop` | `develop` |

---

## Commit Message Format

Mengikuti **[Conventional Commits](https://www.conventionalcommits.org/)** v1.0.0:

```
<tipe>(<scope>): <deskripsi imperatif>

[body — jelaskan MENGAPA, bukan APA]

[footer — referensi issue, breaking change]
```

### Contoh

```bash
# Fitur baru
feat(quiz): add AI-powered quiz generation via Gemini 2.0

# Bug fix
fix(auth): prevent race condition on concurrent login attempts

Closes #78

# Breaking change
feat(api)!: rename /grades endpoint to /assessments

BREAKING CHANGE: All clients must update endpoint from /grades to /assessments.
Migration guide: docs/fase-5/UPGRADE_GUIDE.md

# Dokumentasi
docs(setup): add troubleshooting section for Docker on Windows
```

---

## Pull Request (PR) Rules

### Ukuran PR

| Ukuran | Baris Berubah | Review Time Target |
|:---|:---|:---|
| 🟢 Small | < 100 | < 30 menit |
| 🟡 Medium | 100-400 | < 2 jam |
| 🔴 Large | 400+ | **Wajib dipecah** menjadi PR lebih kecil |

### Checklist Sebelum PR

- [ ] Branch sudah rebase dari `develop` terbaru
- [ ] `make lint` lolos
- [ ] `make test` lolos
- [ ] Tidak ada `console.log` / `fmt.Println` debug tersisa
- [ ] Tidak ada secret/credential ter-commit
- [ ] CHANGELOG.md diperbarui (jika perubahan user-facing)

---

## Protected Branch Rules

| Branch | Rules |
|:---|:---|
| `main` | No direct push, ≥ 2 approvals, CI must pass, admin only merge |
| `develop` | No direct push, ≥ 1 approval, CI must pass |
| `release/*` | No direct push, ≥ 2 approvals, CI must pass |

---

## Release Process

```
1. Buat branch release/x.y dari develop
2. Bump version di package.json dan version.go
3. Update CHANGELOG.md
4. Fix bug yang ditemukan (commit ke release branch)
5. Merge ke main (create tag vx.y.z)
6. Merge balik ke develop
7. Delete release branch
```

---

*Terakhir diperbarui: 21 Maret 2026*
