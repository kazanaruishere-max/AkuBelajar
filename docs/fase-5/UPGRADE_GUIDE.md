# ⬆️ Upgrade Guide — AkuBelajar

> Panduan upgrade dari satu versi ke versi berikutnya — termasuk breaking changes dan langkah migrasi data.

---

## Cara Membaca Guide Ini

Setiap versi memiliki bagian tersendiri yang berisi:
- **Breaking Changes** — Perubahan yang memerlukan aksi manual
- **Migration Steps** — Langkah demi langkah untuk upgrade
- **Rollback** — Cara kembali ke versi sebelumnya jika diperlukan

---

## Upgrade v1.x → v2.0.0 (Major Upgrade)

> ⚠️ **BREAKING CHANGE** — Ini adalah migrasi arsitektur total. Backup database sebelum memulai!

### Apa yang Berubah

| Aspek | v1.x | v2.0.0 |
|:---|:---|:---|
| Backend | PHP/Laravel | Go 1.23+ |
| Frontend | Blade Template | Next.js 15 |
| Database | MySQL 8 | PostgreSQL 16+ |
| Auth | Laravel Sanctum | JWT + Paseto |
| Primary Keys | Auto-increment (INT) | UUID v7 |
| API Format | `/api/*` | `/api/v1/*` |

### Migration Steps

```
1. BACKUP database MySQL lama
   $ mysqldump -u root akubelajar > backup_v1.sql

2. Install PostgreSQL 16+
   $ docker-compose up -d postgres

3. Jalankan migration script (MySQL → PostgreSQL)
   $ make migrate-from-v1 --source=backup_v1.sql
   Script ini akan:
   - Konversi schema MySQL → PostgreSQL
   - Konversi auto-increment ID → UUID v7
   - Re-hash password dari bcrypt → Argon2id
   - Buat mapping table ID lama → UUID baru

4. Verifikasi migrasi data
   $ make verify-migration
   - Cek row count per tabel
   - Cek foreign key integrity
   - Test login dengan user yang dimigrasikan

5. Deploy kode v2.0.0
   $ docker-compose up -d

6. Update client-side bookmarks/API calls
   - /api/* → /api/v1/*
   - Session cookie → Bearer token
```

### Rollback

```bash
# Jika upgrade gagal:
# 1. Stop v2 services
docker-compose down

# 2. Restore MySQL backup
mysql -u root akubelajar < backup_v1.sql

# 3. Redeploy v1 code
git checkout v1.0.0
docker-compose up -d
```

---

## Upgrade v2.0.x → v2.1.0 (Minor Upgrade)

> ✅ Backward-compatible — tidak ada breaking change.

### New Features

- Bulk import v2 (async, improved error reporting)
- Dark mode support
- Early Warning System (AI risk prediction)

### Migration Steps

```bash
# 1. Pull latest
git pull origin main

# 2. Run migrations
make migrate-up

# 3. Rebuild & deploy
docker-compose build && docker-compose up -d

# 4. Verify
curl http://localhost:8080/health
```

---

## Template untuk Versi Baru

```markdown
## Upgrade vX.Y.Z → vX.Y.Z

### Breaking Changes
- [ daftar perubahan yang memerlukan aksi ]

### Migration Steps
1. Step 1
2. Step 2
3. ...

### Rollback
- Cara kembali ke versi sebelumnya

### Post-Upgrade Verification
- [ ] Health check passed
- [ ] Smoke test passed
- [ ] Data integrity verified
```

---

## Referensi

- [Changelog](../fase-0/CHANGELOG.md)
- [Zero Downtime Deploy](../fase-4/ZERO_DOWNTIME_DEPLOY.md)
- [Database Schema](../fase-1/DATABASE_SCHEMA.md)

---

*Terakhir diperbarui: 21 Maret 2026*
