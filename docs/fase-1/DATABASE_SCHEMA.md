# 🗄️ Database Schema — AkuBelajar

> Dokumentasi skema PostgreSQL lengkap termasuk RLS policy, indexing strategy, dan migration plan.

---

## Daftar Isi

1. [Prinsip Desain Database](#prinsip-desain-database)
2. [Entity Relationship Diagram](#entity-relationship-diagram)
3. [Tabel Inti](#tabel-inti)
4. [Row-Level Security (RLS) Policies](#row-level-security-rls-policies)
5. [Indexing Strategy](#indexing-strategy)
6. [Migration Guidelines](#migration-guidelines)

---

## Prinsip Desain Database

| Prinsip | Implementasi |
|:---|:---|
| **Normalisasi** | 3NF minimum — menghindari data redundancy |
| **UUID v7 Primary Key** | Anti-IDOR, time-sortable, distributed-friendly |
| **Soft Delete** | Semua tabel utama memiliki kolom `deleted_at` |
| **Audit Trail** | Tabel `audit_logs` immutable (no UPDATE/DELETE) |
| **RLS** | Row-Level Security di PostgreSQL untuk isolasi data |
| **Timezone** | Semua timestamp menggunakan `TIMESTAMPTZ` (UTC) |

---

## Entity Relationship Diagram

```
┌──────────┐     1:N     ┌──────────┐     1:N     ┌────────────┐
│  schools │─────────────│  users   │─────────────│  user_     │
│          │             │          │             │  profiles  │
└──────────┘             └────┬─────┘             └────────────┘
                              │
              ┌───────────────┼───────────────┐
              │ 1:N           │ 1:N           │ 1:N
        ┌─────▼─────┐  ┌─────▼─────┐  ┌──────▼──────┐
        │  classes   │  │  subjects │  │ attendances │
        └─────┬─────┘  └─────┬─────┘  └─────────────┘
              │               │
              │         ┌─────▼──────┐     1:N     ┌──────────────┐
              │         │  quizzes   │─────────────│ quiz_        │
              │         │            │             │ questions    │
              │         └─────┬──────┘             └──────────────┘
              │               │
              │         ┌─────▼──────┐     1:N     ┌──────────────┐
              │         │  quiz_     │─────────────│ quiz_        │
              │         │ submissions│             │ answers      │
              │         └────────────┘             └──────────────┘
              │
        ┌─────▼──────┐     1:N     ┌──────────────┐
        │ assignments│─────────────│ submissions  │
        └────────────┘             └──────────────┘
```

---

## Tabel Inti

### `users`

```sql
CREATE TYPE user_role_enum AS ENUM ('super_admin', 'teacher', 'class_leader', 'student');

CREATE TABLE users (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email               VARCHAR(255) UNIQUE NOT NULL,
    password            VARCHAR(255) NOT NULL,          -- Argon2id hash
    role                user_role_enum NOT NULL,
    school_id           UUID REFERENCES schools(id) ON DELETE CASCADE,
    is_active           BOOLEAN DEFAULT TRUE,

    -- Security fields
    failed_login_count  INTEGER DEFAULT 0,
    locked_until        TIMESTAMPTZ,
    last_login_at       TIMESTAMPTZ,
    last_login_ip       INET,

    -- Timestamps
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ                     -- Soft delete
);
```

### `schools`

```sql
CREATE TABLE schools (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    npsn            VARCHAR(20) UNIQUE,                 -- Nomor Pokok Sekolah Nasional
    address         TEXT,
    province        VARCHAR(100),
    city            VARCHAR(100),
    config          JSONB DEFAULT '{}'::JSONB,          -- Per-school settings
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);
```

### `quizzes`

```sql
CREATE TABLE quizzes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           VARCHAR(255) NOT NULL,
    subject_id      UUID REFERENCES subjects(id),
    class_id        UUID REFERENCES classes(id),
    teacher_id      UUID REFERENCES users(id),
    time_limit      INTEGER NOT NULL,                   -- Menit
    is_published     BOOLEAN DEFAULT FALSE,
    ai_generated    BOOLEAN DEFAULT FALSE,
    start_at        TIMESTAMPTZ,
    end_at          TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
```

### `audit_logs` (Immutable)

```sql
CREATE TABLE audit_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    action          VARCHAR(50) NOT NULL,               -- 'CREATE', 'UPDATE', 'DELETE'
    entity_type     VARCHAR(50) NOT NULL,               -- 'quiz', 'grade', 'user'
    entity_id       UUID NOT NULL,
    old_value       JSONB,
    new_value       JSONB,
    ip_address      INET,
    user_agent      TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Mencegah mutasi audit log
REVOKE UPDATE, DELETE ON audit_logs FROM app_user;
```

---

## Row-Level Security (RLS) Policies

### Prinsip

Setiap query secara otomatis difilter berdasarkan `current_setting('app.current_user_id')` yang di-set oleh middleware Go sebelum setiap transaksi database.

### Policy: Siswa hanya lihat nilainya sendiri

```sql
ALTER TABLE grades ENABLE ROW LEVEL SECURITY;

CREATE POLICY student_grades_isolation ON grades
    FOR SELECT
    USING (student_id = current_setting('app.current_user_id')::UUID);
```

### Policy: Guru hanya akses kelas yang diajar

```sql
ALTER TABLE classes ENABLE ROW LEVEL SECURITY;

CREATE POLICY teacher_class_isolation ON classes
    USING (
        school_id = current_setting('app.current_school_id')::UUID
        AND EXISTS (
            SELECT 1 FROM teacher_subjects ts
            WHERE ts.teacher_id = current_setting('app.current_user_id')::UUID
              AND ts.class_id = classes.id
        )
    );
```

---

## Indexing Strategy

| Tabel | Index | Tipe | Alasan |
|:---|:---|:---|:---|
| `users` | `email` | UNIQUE B-tree | Lookup login |
| `users` | `school_id, role` | Composite B-tree | Filter user per sekolah per role |
| `quizzes` | `class_id, is_published` | Composite B-tree | Daftar kuis aktif per kelas |
| `quiz_submissions` | `quiz_id, student_id` | UNIQUE Composite | Satu submission per siswa per kuis |
| `attendances` | `class_id, date` | Composite B-tree | Lookup absensi harian |
| `audit_logs` | `entity_type, entity_id` | Composite B-tree | Trace perubahan per entity |
| `audit_logs` | `created_at` | B-tree DESC | Query audit terbaru |

---

## Migration Guidelines

### Tools

- **golang-migrate** (`github.com/golang-migrate/migrate`)
- File naming: `YYYYMMDDHHMMSS_description.up.sql` / `.down.sql`

### Rules

1. **Setiap migration harus reversible** — selalu buat file `.down.sql`
2. **Tidak boleh ada data loss** di migration `.down.sql`
3. **Backward-compatible** — kolom baru harus `DEFAULT` atau `NULL` capable
4. **Test locally** sebelum push: `make migrate-up && make migrate-down && make migrate-up`
5. **Tidak boleh DROP COLUMN** yang masih digunakan oleh kode yang berjalan

```bash
# Buat migration baru
make migration-create name=add_quiz_feedback_column

# Jalankan migration
make migrate-up

# Rollback terakhir
make migrate-down-1
```

---

## Referensi Terkait

- [System Overview](SYSTEM_OVERVIEW.md) — Arsitektur keseluruhan
- [Audit Log Spec](../fase-3/AUDIT_LOG_SPEC.md) — Detail event yang di-log
- [Backup & Recovery](../fase-4/BACKUP_AND_RECOVERY.md) — Prosedur backup database

---

*Terakhir diperbarui: 21 Maret 2026*
