# ⏱️ Zero Downtime Deploy — AkuBelajar

> Strategi rolling update dan migrasi database backward-compatible agar tidak ada downtime saat deploy.

---

## Prinsip

1. **Tidak ada saat di mana user tidak bisa akses** — transisi harus seamless
2. **Migrasi database harus backward-compatible** — kode lama bisa jalan dengan schema baru
3. **Rollback harus cepat** — ≤ 5 menit untuk kembali ke versi sebelumnya

---

## Strategi Rolling Update (Kubernetes)

```yaml
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1          # Tambah 1 pod baru sebelum hapus yang lama
      maxUnavailable: 0     # Tidak boleh ada pod yang tidak tersedia
```

### Alur

```
Pod v1  Pod v1  Pod v1     ← Awal (3 pods, semua v1)
    │
    ▼
Pod v1  Pod v1  Pod v1  Pod v2  ← Surge: 1 pod v2 ditambahkan
    │
    ▼
Pod v1  Pod v1  Pod v2  Pod v2  ← Rolling: 1 pod v1 dihapus, 1 v2 ditambah
    │
    ▼
Pod v2  Pod v2  Pod v2     ← Selesai (semua v2)
```

---

## Database Migration: Expand-Contract Pattern

### ❌ SALAH — Breaking migration

```sql
-- Ini akan break kode v1 yang masih jalan!
ALTER TABLE users RENAME COLUMN name TO full_name;
```

### ✅ BENAR — Expand-Contract (3 tahap)

**Tahap 1: Expand (deploy bersamaan dengan kode v2)**

```sql
-- Tambah kolom baru, isi dari kolom lama
ALTER TABLE users ADD COLUMN full_name VARCHAR(255);
UPDATE users SET full_name = name WHERE full_name IS NULL;
```

**Tahap 2: Migrate Code (kode v2 menulis ke kedua kolom)**

```go
// Kode v2 — tulis ke full_name DAN name (backward-compat)
func (r *UserRepo) UpdateName(ctx context.Context, id uuid.UUID, fullName string) error {
    _, err := r.pool.Exec(ctx,
        `UPDATE users SET full_name = $1, name = $1 WHERE id = $2`,
        fullName, id,
    )
    return err
}
```

**Tahap 3: Contract (setelah semua pod sudah v2, deploy terpisah)**

```sql
-- Hapus kolom lama
ALTER TABLE users DROP COLUMN name;
```

---

## Rollback Strategy

### Kubernetes Rollback

```bash
# Lihat history
kubectl rollout history deployment/akubelajar-api -n akubelajar-prod

# Rollback ke versi sebelumnya
kubectl rollout undo deployment/akubelajar-api -n akubelajar-prod

# Rollback ke versi spesifik
kubectl rollout undo deployment/akubelajar-api --to-revision=3
```

### Database Rollback

```bash
# Rollback 1 migration terakhir
make migrate-down-1

# Rollback ke versi spesifik
make migrate-goto VERSION=20260321001
```

---

## Pre-Deploy Checklist

- [ ] Migration backward-compatible? (kode lama bisa jalan?)
- [ ] Readiness probe terpasang?
- [ ] Rollback plan sudah disiapkan?
- [ ] Staging test passed?
- [ ] Feature flags ready untuk gradual rollout?
- [ ] Monitoring alert threshold sudah di-set?

---

## Referensi

- [Deployment Guide](DEPLOYMENT_GUIDE.md)
- [Feature Flag](../fase-2/FEATURE_FLAG.md)
- [Database Schema](../fase-1/DATABASE_SCHEMA.md)

---

*Terakhir diperbarui: 21 Maret 2026*
