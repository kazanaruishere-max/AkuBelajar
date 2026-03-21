# 🧪 Testing Strategy — AkuBelajar

> Apa yang wajib di-test, target coverage, tool yang digunakan, dan hierarki testing.

---

## Testing Pyramid

```
         ╱╲
        ╱  ╲         E2E Tests (5%)
       ╱    ╲        Playwright — critical user flows
      ╱──────╲
     ╱        ╲      Integration Tests (25%)
    ╱          ╲     API tests, DB tests, service tests
   ╱────────────╲
  ╱              ╲   Unit Tests (70%)
 ╱                ╲  Pure functions, business logic, utilities
╱──────────────────╲
```

---

## Target Coverage

| Layer | Minimum | Target | Blocking |
|:---|:---|:---|:---|
| Backend (Go) — Unit | 70% | 85% | < 60% = PR blocked |
| Backend (Go) — Integration | 50% | 70% | < 40% = PR blocked |
| Frontend (TypeScript) — Unit | 60% | 80% | < 50% = PR blocked |
| Frontend (React) — Component | 50% | 70% | < 40% = PR blocked |
| E2E — Critical paths | 100% | 100% | Any failure = deploy blocked |

---

## Tools per Layer

### Backend (Go)

| Tool | Purpose |
|:---|:---|
| `testing` (stdlib) | Unit tests |
| `testify/assert` | Assertions |
| `testify/mock` | Mocking interfaces |
| `testcontainers-go` | Integration tests dengan PostgreSQL & Redis containers |
| `httptest` | HTTP handler testing |
| `goleak` | Goroutine leak detection |

### Frontend (TypeScript/React)

| Tool | Purpose |
|:---|:---|
| `vitest` | Unit test runner (faster than Jest) |
| `@testing-library/react` | Component testing (user-centric) |
| `msw` (Mock Service Worker) | API mocking |
| `@playwright/test` | E2E testing |

---

## Apa yang Wajib Di-Test

### ✅ WAJIB

| Area | Contoh |
|:---|:---|
| **Auth flows** | Login, register, token refresh, logout |
| **RBAC** | Setiap role dapat/tidak dapat akses endpoint tertentu |
| **Data validation** | Input sanitasi, boundary values |
| **Business rules** | Quiz scoring, grade calculation, attendance rules |
| **Error handling** | 404, 403, 500, timeout |
| **Database constraints** | Unique violations, FK constraints, RLS |
| **AI output parsing** | Malformed JSON, empty response, timeout |

### ⚠️ OPSIONAL (Nice to have)

| Area | Contoh |
|:---|:---|
| UI visual regression | Screenshot comparison |
| Performance tests | Load testing with k6 |
| Accessibility tests | axe-core automation |

### ❌ TIDAK PERLU Di-Test

| Area | Alasan |
|:---|:---|
| Third-party library internals | Sudah di-test oleh maintainer |
| Getter/setter sederhana | Tidak ada logic |
| Framework boilerplate | Next.js routing, Gin binding |

---

## Naming Convention

```go
// Go test function naming
func TestQuizService_CreateQuiz_Success(t *testing.T) { ... }
func TestQuizService_CreateQuiz_InvalidInput(t *testing.T) { ... }
func TestQuizService_CreateQuiz_Unauthorized(t *testing.T) { ... }

// Pattern: Test{Service}_{Method}_{Scenario}
```

```typescript
// TypeScript test naming
describe('QuizService', () => {
  describe('createQuiz', () => {
    it('should create quiz with valid input', () => { ... });
    it('should throw error when title is empty', () => { ... });
    it('should reject unauthorized user', () => { ... });
  });
});
```

---

## CI Integration

```yaml
# .github/workflows/test.yml
test-backend:
  steps:
    - run: go test ./... -race -coverprofile=coverage.out
    - run: go tool cover -func=coverage.out
    - uses: codecov/codecov-action@v4

test-frontend:
  steps:
    - run: pnpm test --coverage
    - run: pnpm test:e2e
```

---

## Referensi

- [Coding Standards](CODING_STANDARDS.md)
- [Backend Guide](BACKEND_GUIDE.md)
- [Frontend Guide](FRONTEND_GUIDE.md)

---

*Terakhir diperbarui: 21 Maret 2026*
