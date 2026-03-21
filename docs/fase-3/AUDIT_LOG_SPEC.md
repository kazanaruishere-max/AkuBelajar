# 📋 Audit Log Specification — AkuBelajar

> Semua event yang wajib masuk audit log immutable — tidak bisa di-UPDATE atau di-DELETE.

---

## Prinsip Audit Log

1. **Immutable** — Sekali ditulis, tidak bisa diubah atau dihapus
2. **Complete** — Mencakup who, what, when, where, old/new value
3. **Tamper-evident** — Perubahan terdeteksi via checksum chain
4. **Queryable** — Indexed untuk pencarian cepat

---

## Schema

```sql
CREATE TABLE audit_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,                       -- Siapa
    action          VARCHAR(50) NOT NULL,                 -- Apa (CREATE/UPDATE/DELETE)
    entity_type     VARCHAR(50) NOT NULL,                 -- Tipe entity
    entity_id       UUID NOT NULL,                        -- ID entity
    old_value       JSONB,                                -- Nilai sebelum (untuk UPDATE/DELETE)
    new_value       JSONB,                                -- Nilai sesudah (untuk CREATE/UPDATE)
    ip_address      INET,                                 -- Dari mana
    user_agent      TEXT,                                  -- Browser/device info
    request_id      UUID,                                  -- Correlation ID
    created_at      TIMESTAMPTZ DEFAULT NOW()             -- Kapan
);

-- IMMUTABILITY ENFORCEMENT
REVOKE UPDATE, DELETE ON audit_logs FROM app_user;
REVOKE TRUNCATE ON audit_logs FROM app_user;
```

---

## Event yang WAJIB Di-Log

### 🔴 Critical Events (Selalu di-log)

| Entity | Create | Update | Delete |
|:---|:---|:---|:---|
| `users` | ✅ | ✅ (role change, status) | ✅ |
| `grades` | ✅ | ✅ (nilai berubah) | ✅ |
| `quiz_submissions` | ✅ | ❌ (immutable) | ❌ |
| `audit_logs` | ❌ (auto) | ❌ (blocked) | ❌ (blocked) |

### 🟡 Important Events

| Entity | Create | Update | Delete |
|:---|:---|:---|:---|
| `quizzes` | ✅ | ✅ (publish/unpublish) | ✅ |
| `assignments` | ✅ | ✅ | ✅ |
| `attendances` | ✅ | ✅ | ✅ |
| `classes` | ✅ | ✅ | ✅ |

### 🟢 Security Events (Selalu di-log)

| Event | Di-log |
|:---|:---|
| Login success | ✅ (user_id, IP, user_agent) |
| Login failed | ✅ (email attempted, IP, reason) |
| Account locked | ✅ |
| Password changed | ✅ (tanpa value lama/baru) |
| Role changed | ✅ (old role → new role) |
| Token refresh | ✅ |
| Permission denied (403) | ✅ |

---

## Retention Policy

| Tier | Umur Data | Storage | Akses |
|:---|:---|:---|:---|
| Hot | 0-90 hari | PostgreSQL (primary) | Direct query |
| Warm | 90 hari - 1 tahun | PostgreSQL (partitioned) | Query via archive view |
| Cold | 1-5 tahun | Object Storage (MinIO) | Export/restore on demand |
| Archive | 5+ tahun | Compressed backup | Legal/compliance only |

---

## Query Examples

```sql
-- Siapa yang mengubah nilai siswa X?
SELECT * FROM audit_logs
WHERE entity_type = 'grade'
  AND entity_id = 'uuid-grade-xxx'
ORDER BY created_at DESC;

-- Semua login gagal dari IP tertentu
SELECT * FROM audit_logs
WHERE action = 'LOGIN_FAILED'
  AND ip_address = '203.0.113.50'
  AND created_at > NOW() - INTERVAL '24 hours';

-- Semua perubahan role user
SELECT * FROM audit_logs
WHERE entity_type = 'user'
  AND action = 'UPDATE'
  AND new_value->>'role' IS DISTINCT FROM old_value->>'role';
```

---

*Terakhir diperbarui: 21 Maret 2026*
