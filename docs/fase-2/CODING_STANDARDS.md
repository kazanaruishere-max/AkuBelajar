# 📏 Coding Standards — AkuBelajar

> Satu gaya kode untuk Go, TypeScript, dan SQL agar codebase konsisten di seluruh layer.

---

## Prinsip Umum

1. **Konsistensi di atas preferensi** — Ikuti standar yang ada, bukan selera pribadi
2. **Readability first** — Kode dibaca 10x lebih sering daripada ditulis
3. **Explicit over implicit** — Lebih baik verbose tapi jelas daripada singkat tapi membingungkan
4. **Autoformat wajib** — Tidak ada perdebatan formatting; serahkan ke tools

---

## Go (Backend)

### Formatter & Linter

```bash
# Wajib dijalankan sebelum commit
gofmt -w .          # Format otomatis
golangci-lint run   # Lint komprehensif
```

### Naming Convention

| Elemen | Convention | Contoh |
|:---|:---|:---|
| Package | lowercase, singular | `auth`, `quiz`, `notification` |
| Exported function | PascalCase | `CreateQuiz`, `GetStudentGrades` |
| Unexported function | camelCase | `hashPassword`, `validateInput` |
| Constant | PascalCase | `MaxLoginAttempts`, `DefaultPageSize` |
| Interface | PascalCase + `-er` suffix | `QuizService`, `UserRepository` |
| Struct | PascalCase | `QuizQuestion`, `AuthClaims` |
| File | snake_case | `quiz_service.go`, `auth_middleware.go` |

### Error Handling

```go
// ✅ BENAR — Wrap error dengan konteks
if err != nil {
    return fmt.Errorf("failed to create quiz: %w", err)
}

// ❌ SALAH — Error tanpa konteks
if err != nil {
    return err
}

// ❌ SALAH — Ignore error
result, _ := db.Query(...)
```

### Struct Tags

```go
type CreateQuizRequest struct {
    Title      string `json:"title"      validate:"required,min=3,max=255"`
    SubjectID  string `json:"subject_id" validate:"required,uuid"`
    TimeLimit  int    `json:"time_limit" validate:"required,min=5,max=180"`
}
```

---

## TypeScript (Frontend)

### Formatter & Linter

```bash
pnpm lint          # ESLint
pnpm format        # Prettier
```

### Naming Convention

| Elemen | Convention | Contoh |
|:---|:---|:---|
| Component | PascalCase | `QuizCard`, `StudentDashboard` |
| File (component) | PascalCase | `QuizCard.tsx`, `StudentDashboard.tsx` |
| File (utility) | camelCase | `apiClient.ts`, `formatDate.ts` |
| Function | camelCase | `fetchQuizzes`, `handleSubmit` |
| Constant | SCREAMING_SNAKE | `MAX_FILE_SIZE`, `API_BASE_URL` |
| Type/Interface | PascalCase | `Quiz`, `UserProfile`, `ApiResponse<T>` |
| Enum | PascalCase | `UserRole.Teacher`, `QuizStatus.Published` |
| CSS class | kebab-case | `quiz-card`, `student-dashboard` |

### TypeScript Rules

```typescript
// ✅ BENAR — Strict typing, no `any`
function getUser(id: string): Promise<User> { ... }

// ❌ SALAH — Using `any`
function getUser(id: any): Promise<any> { ... }

// ✅ BENAR — Zod validation pada API response
const user = UserSchema.parse(await response.json());

// ❌ SALAH — Type assertion tanpa validasi
const user = (await response.json()) as User;
```

---

## SQL (Database)

### Naming Convention

| Elemen | Convention | Contoh |
|:---|:---|:---|
| Table | snake_case, plural | `users`, `quiz_questions` |
| Column | snake_case | `created_at`, `school_id` |
| Index | `idx_{table}_{columns}` | `idx_users_email` |
| Foreign Key | `fk_{table}_{ref_table}` | `fk_quizzes_teacher` |
| Constraint | `chk_{table}_{column}` | `chk_users_email_format` |

### Query Rules

```sql
-- ✅ BENAR — Parameterized query
SELECT * FROM users WHERE email = $1 AND school_id = $2;

-- ❌ SALAH — String concatenation (SQL Injection!)
SELECT * FROM users WHERE email = '" + email + "';

-- ✅ BENAR — Explicit column selection
SELECT id, email, role FROM users WHERE school_id = $1;

-- ❌ SALAH — SELECT * di production
SELECT * FROM users;
```

---

## File Organization

```
// Go — Group by domain
internal/
├── auth/
│   ├── handler.go      # HTTP handlers
│   ├── service.go      # Business logic
│   ├── repository.go   # Database queries
│   └── model.go        # Structs & types

// TypeScript — Group by feature
components/
├── features/
│   ├── quiz/
│   │   ├── QuizCard.tsx
│   │   ├── QuizList.tsx
│   │   └── useQuiz.ts
│   └── attendance/
│       ├── AttendanceForm.tsx
│       └── useAttendance.ts
```

---

## Referensi

- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- [Effective Go](https://go.dev/doc/effective_go)
- [TypeScript Style Guide (Google)](https://google.github.io/styleguide/tsguide.html)

---

*Terakhir diperbarui: 21 Maret 2026*
