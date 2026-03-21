# 🔄 Migration Strategy — AkuBelajar

> Strategi migrasi database: tool, naming, versioning, rollback, dan zero-downtime migration.

---

## 1. Tool & Setup

| Item | Detail |
|:---|:---|
| Tool | `golang-migrate/migrate` v4 |
| Driver | `pgx5` |
| Migration dir | `backend/migrations/` |
| Run | `make migrate-up` / `make migrate-down` |

### File Structure

```
backend/migrations/
├── 000001_create_schools.up.sql
├── 000001_create_schools.down.sql
├── 000002_create_users.up.sql
├── 000002_create_users.down.sql
├── 000003_create_academic_tables.up.sql
├── 000003_create_academic_tables.down.sql
└── ...
```

---

## 2. Naming Convention

```
{VERSION}_{DESCRIPTION}.{DIRECTION}.sql
```

| Part | Format | Contoh |
|:---|:---|:---|
| Version | 6-digit zero-padded | `000015` |
| Description | snake_case, verb_noun | `add_cheat_log_to_sessions` |
| Direction | `up` atau `down` | `up.sql` / `down.sql` |

### Contoh

```
000015_add_cheat_log_to_quiz_sessions.up.sql
000015_add_cheat_log_to_quiz_sessions.down.sql
```

---

## 3. Rules Wajib

### Setiap migration HARUS:

1. **Idempotent** — `IF NOT EXISTS` untuk CREATE, `IF EXISTS` untuk DROP
2. **Punya down migration** — rollback harus selalu mungkin
3. **Atomic** — wrapped dalam `BEGIN; ... COMMIT;`
4. **Backward-compatible** — lihat bagian Zero-Downtime di bawah
5. **Tested** — jalankan up + down + up sebelum commit

### Template

```sql
-- 000015_add_cheat_log_to_quiz_sessions.up.sql
BEGIN;

ALTER TABLE quiz_sessions
    ADD COLUMN IF NOT EXISTS cheat_log JSONB DEFAULT '[]';

COMMENT ON COLUMN quiz_sessions.cheat_log IS 'Array of cheat detection events';

COMMIT;
```

```sql
-- 000015_add_cheat_log_to_quiz_sessions.down.sql
BEGIN;

ALTER TABLE quiz_sessions
    DROP COLUMN IF EXISTS cheat_log;

COMMIT;
```

---

## 4. Zero-Downtime Migration (Expand-Contract)

Saat menambah/mengubah kolom di production, gunakan pola **3 fase**:

### Fase 1: EXPAND (deploy 1)

```sql
-- Tambah kolom baru, nullable, tanpa NOT NULL
ALTER TABLE users ADD COLUMN phone_verified BOOLEAN;
```

Backend: tulis ke kolom baru DAN kolom lama.

### Fase 2: MIGRATE DATA (deploy 2)

```sql
-- Backfill data lama
UPDATE users SET phone_verified = FALSE WHERE phone_verified IS NULL;
```

### Fase 3: CONTRACT (deploy 3)

```sql
-- Sekarang aman untuk set NOT NULL dan drop kolom lama
ALTER TABLE users ALTER COLUMN phone_verified SET NOT NULL;
ALTER TABLE users ALTER COLUMN phone_verified SET DEFAULT FALSE;
```

**JANGAN pernah** dalam satu deploy:
- ❌ Rename kolom (break existing queries)
- ❌ Change type (break inserts)
- ❌ Add NOT NULL tanpa default (break inserts)
- ❌ Drop kolom yang masih direferensikan

---

## 5. Data Migration

Untuk migrasi data besar (> 10.000 rows):

```sql
-- Batch update, bukan satu UPDATE besar
DO $$
DECLARE
    batch_size INT := 1000;
    affected INT;
BEGIN
    LOOP
        UPDATE users
        SET status = 'active'
        WHERE status IS NULL
        AND id IN (SELECT id FROM users WHERE status IS NULL LIMIT batch_size);

        GET DIAGNOSTICS affected = ROW_COUNT;
        EXIT WHEN affected = 0;
        
        PERFORM pg_sleep(0.1);  -- Rate limit agar tidak lock
    END LOOP;
END $$;
```

---

## 6. Commands

```makefile
# Create migration baru
migrate-create:
	migrate create -ext sql -dir backend/migrations -seq $(name)

# Run semua migration up
migrate-up:
	migrate -path backend/migrations -database "$(DB_URL)" up

# Rollback 1 step
migrate-down:
	migrate -path backend/migrations -database "$(DB_URL)" down 1

# Force version (jika dirty)
migrate-force:
	migrate -path backend/migrations -database "$(DB_URL)" force $(version)

# Status
migrate-version:
	migrate -path backend/migrations -database "$(DB_URL)" version
```

---

## 7. Pre-Commit Checklist

- [ ] `up.sql` dan `down.sql` keduanya ada
- [ ] `down.sql` benar-benar reverse dari `up.sql`
- [ ] Jalankan: up → down → up tanpa error
- [ ] Tidak ada `DROP COLUMN` atau `ALTER TYPE` tanpa expand-contract
- [ ] Data migration menggunakan batch (bukan single UPDATE)
- [ ] Index baru menggunakan `CONCURRENTLY` di production

---

*Terakhir diperbarui: 21 Maret 2026*
