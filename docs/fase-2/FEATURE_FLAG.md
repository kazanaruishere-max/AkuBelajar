# 🚩 Feature Flag — AkuBelajar

> Cara rollout fitur secara bertahap tanpa deploy ulang — esensial untuk upgrade yang aman.

---

## Mengapa Feature Flag?

| Tanpa Feature Flag | Dengan Feature Flag |
|:---|:---|
| Deploy = semua user langsung terdampak | Rollout bertahap: 5% → 25% → 100% |
| Bug di fitur baru = rollback deploy | Bug = matikan flag, instant off |
| Perlu deploy ulang untuk disable fitur | Config update, no deploy needed |
| Testing hanya di staging | Testing di production (canary users) |

---

## Arsitektur

```
┌──────────────┐     ┌──────────────────┐     ┌───────────────┐
│   Dashboard  │────▶│  Redis Cache     │────▶│  PostgreSQL   │
│   Admin UI   │     │  (TTL: 30s)      │     │  (Source of   │
│              │     │                  │     │   Truth)      │
└──────────────┘     └───────┬──────────┘     └───────────────┘
                             │
                    ┌────────▼────────┐
                    │  Go Middleware   │
                    │  isEnabled(flag) │
                    └─────────────────┘
```

---

## Database Schema

```sql
CREATE TABLE feature_flags (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) UNIQUE NOT NULL,    -- 'ai_quiz_v2', 'dark_mode'
    description TEXT,
    is_enabled  BOOLEAN DEFAULT FALSE,           -- Global toggle
    rollout_pct INTEGER DEFAULT 0 CHECK (rollout_pct BETWEEN 0 AND 100),
    rules       JSONB DEFAULT '{}'::JSONB,       -- Targeting rules
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);
```

### Targeting Rules (JSONB)

```json
{
  "school_ids": ["uuid-sekolah-a", "uuid-sekolah-b"],
  "roles": ["teacher", "super_admin"],
  "user_ids": ["uuid-tester-1"],
  "min_version": "2.1.0"
}
```

---

## Go Implementation

```go
// internal/featureflag/service.go
type FeatureFlagService struct {
    repo  FeatureFlagRepo
    cache *redis.Client
}

func (s *FeatureFlagService) IsEnabled(ctx context.Context, flagName string, user *User) bool {
    flag, err := s.getFlag(ctx, flagName)
    if err != nil || !flag.IsEnabled {
        return false
    }

    // Check targeting rules
    if len(flag.Rules.UserIDs) > 0 && contains(flag.Rules.UserIDs, user.ID) {
        return true
    }
    if len(flag.Rules.Roles) > 0 && contains(flag.Rules.Roles, user.Role) {
        return true
    }
    if len(flag.Rules.SchoolIDs) > 0 && contains(flag.Rules.SchoolIDs, user.SchoolID) {
        return true
    }

    // Percentage rollout (deterministic based on user ID)
    if flag.RolloutPct > 0 {
        hash := crc32.ChecksumIEEE([]byte(user.ID.String()))
        return (hash % 100) < uint32(flag.RolloutPct)
    }

    return false
}
```

### Usage di Handler

```go
func (h *QuizHandler) CreateQuiz(c *gin.Context) {
    user := auth.GetUser(c)

    if h.flags.IsEnabled(c, "ai_quiz_v2", user) {
        // Gunakan AI Quiz v2 (Gemini 2.0)
        return h.createQuizV2(c)
    }

    // Fallback ke versi lama
    return h.createQuizV1(c)
}
```

---

## Frontend Usage

```tsx
// hooks/useFeatureFlag.ts
export function useFeatureFlag(flagName: string): boolean {
  const { data } = useQuery({
    queryKey: ['feature-flags', flagName],
    queryFn: () => api.flags.check(flagName),
    staleTime: 30_000, // Cache 30 detik
  });
  return data?.enabled ?? false;
}

// Component usage
function QuizPage() {
  const isV2Enabled = useFeatureFlag('ai_quiz_v2');

  return isV2Enabled ? <QuizV2 /> : <QuizV1 />;
}
```

---

## Rollout Playbook

```
1. Buat flag → is_enabled: false, rollout_pct: 0
2. Enable untuk internal team → rules.roles: ["super_admin"]
3. Testing di production → pastikan tidak ada error
4. Rollout 5% → rollout_pct: 5, monitor 24 jam
5. Rollout 25% → monitor 48 jam
6. Rollout 100% → monitor 1 minggu
7. Cleanup → hapus flag, hapus code path lama
```

---

## Flag Management

### Flag yang Aktif

| Flag Name | Status | Rollout | Keterangan |
|:---|:---|:---|:---|
| `ai_quiz_v2` | ✅ Active | 100% | Gemini 2.0 Flash quiz generator |
| `dark_mode` | ✅ Active | 100% | Dark mode via prefers-color-scheme |
| `bulk_import_v2` | 🟡 Rolling | 25% | Improved async import |
| `early_warning` | 🔴 Internal | 0% | AI risk prediction (testing) |

---

*Terakhir diperbarui: 21 Maret 2026*
