# 🔗 Data Models — AkuBelajar

> Menjembatani gap antara skema database PostgreSQL, struct Go, dan TypeScript types. Dokumen ini memastikan nama field konsisten di semua layer — sumber bug terbesar dalam proyek multi-layer.

---

## Konvensi Penamaan per Layer

| Layer | Convention | Contoh |
|:---|:---|:---|
| **PostgreSQL** | `snake_case` | `student_id`, `created_at`, `is_active` |
| **Go (struct)** | `PascalCase` | `StudentID`, `CreatedAt`, `IsActive` |
| **Go (JSON tag)** | `snake_case` | `json:"student_id"` |
| **TypeScript** | `camelCase` | `studentId`, `createdAt`, `isActive` |
| **API Request/Response** | `snake_case` (JSON) | `"student_id": "..."` |

> ⚠️ **API JSON menggunakan `snake_case`**, bukan camelCase. Frontend TypeScript harus melakukan mapping saat consume.

---

## User

### PostgreSQL

```sql
CREATE TABLE users (
    id                  UUID PRIMARY KEY,
    email               VARCHAR(255) UNIQUE NOT NULL,
    password            VARCHAR(255) NOT NULL,
    role                user_role_enum NOT NULL,
    school_id           UUID REFERENCES schools(id),
    is_active           BOOLEAN DEFAULT TRUE,
    failed_login_count  INTEGER DEFAULT 0,
    locked_until        TIMESTAMPTZ,
    last_login_at       TIMESTAMPTZ,
    last_login_ip       INET,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);
```

### Go Struct

```go
// internal/auth/model.go
type User struct {
    ID               uuid.UUID  `json:"id"                db:"id"`
    Email            string     `json:"email"             db:"email"`
    Password         string     `json:"-"                 db:"password"`          // NEVER exposed in JSON
    Role             UserRole   `json:"role"              db:"role"`
    SchoolID         uuid.UUID  `json:"school_id"         db:"school_id"`
    IsActive         bool       `json:"is_active"         db:"is_active"`
    FailedLoginCount int        `json:"-"                 db:"failed_login_count"` // Internal only
    LockedUntil      *time.Time `json:"-"                 db:"locked_until"`       // Internal only
    LastLoginAt      *time.Time `json:"last_login_at"     db:"last_login_at"`
    LastLoginIP      *string    `json:"-"                 db:"last_login_ip"`      // Internal only
    CreatedAt        time.Time  `json:"created_at"        db:"created_at"`
    UpdatedAt        time.Time  `json:"updated_at"        db:"updated_at"`
    DeletedAt        *time.Time `json:"-"                 db:"deleted_at"`         // Internal only
}

type UserRole string
const (
    RoleSuperAdmin  UserRole = "super_admin"
    RoleTeacher     UserRole = "teacher"
    RoleClassLeader UserRole = "class_leader"
    RoleStudent     UserRole = "student"
)
```

### TypeScript Interface

```typescript
// types/user.ts
export interface User {
  id: string;             // UUID
  email: string;
  role: UserRole;
  schoolId: string;       // ← camelCase, mapped from API's snake_case
  isActive: boolean;
  lastLoginAt: string | null;  // ISO 8601
  createdAt: string;
  updatedAt: string;
}

export type UserRole = 'super_admin' | 'teacher' | 'class_leader' | 'student';

// Zod schema for runtime validation
export const UserSchema = z.object({
  id: z.string().uuid(),
  email: z.string().email(),
  role: z.enum(['super_admin', 'teacher', 'class_leader', 'student']),
  school_id: z.string().uuid(),          // API returns snake_case
  is_active: z.boolean(),
  last_login_at: z.string().nullable(),
  created_at: z.string(),
  updated_at: z.string(),
});
```

---

## Quiz

### PostgreSQL

```sql
CREATE TABLE quizzes (
    id              UUID PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    subject_id      UUID REFERENCES subjects(id),
    class_id        UUID REFERENCES classes(id),
    teacher_id      UUID REFERENCES users(id),
    time_limit      INTEGER NOT NULL,
    is_published    BOOLEAN DEFAULT FALSE,
    ai_generated    BOOLEAN DEFAULT FALSE,
    start_at        TIMESTAMPTZ,
    end_at          TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
```

### Go Struct

```go
// internal/quiz/model.go
type Quiz struct {
    ID          uuid.UUID  `json:"id"           db:"id"`
    Title       string     `json:"title"        db:"title"`
    SubjectID   uuid.UUID  `json:"subject_id"   db:"subject_id"`
    ClassID     uuid.UUID  `json:"class_id"     db:"class_id"`
    TeacherID   uuid.UUID  `json:"teacher_id"   db:"teacher_id"`
    TimeLimit   int        `json:"time_limit"   db:"time_limit"`     // Menit
    IsPublished bool       `json:"is_published" db:"is_published"`
    AIGenerated bool       `json:"ai_generated" db:"ai_generated"`
    StartAt     *time.Time `json:"start_at"     db:"start_at"`
    EndAt       *time.Time `json:"end_at"       db:"end_at"`
    CreatedAt   time.Time  `json:"created_at"   db:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"   db:"updated_at"`
    DeletedAt   *time.Time `json:"-"            db:"deleted_at"`
}
```

### TypeScript Interface

```typescript
// types/quiz.ts
export interface Quiz {
  id: string;
  title: string;
  subjectId: string;
  classId: string;
  teacherId: string;
  timeLimit: number;       // Menit
  isPublished: boolean;
  aiGenerated: boolean;
  startAt: string | null;  // ISO 8601
  endAt: string | null;
  createdAt: string;
  updatedAt: string;
}
```

---

## Quiz Question

### PostgreSQL

```sql
CREATE TABLE quiz_questions (
    id              UUID PRIMARY KEY,
    quiz_id         UUID REFERENCES quizzes(id) ON DELETE CASCADE,
    question_text   TEXT NOT NULL,
    question_type   question_type_enum NOT NULL,  -- 'multiple_choice', 'essay'
    options         JSONB,                         -- Array of {key, text}
    answer_hash     VARCHAR(255),                  -- Argon2id hash of correct answer
    explanation     TEXT,
    order_num       INTEGER NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

### Go Struct

```go
type QuizQuestion struct {
    ID           uuid.UUID    `json:"id"            db:"id"`
    QuizID       uuid.UUID    `json:"quiz_id"       db:"quiz_id"`
    QuestionText string       `json:"question"      db:"question_text"`
    QuestionType QuestionType `json:"question_type" db:"question_type"`
    Options      []Option     `json:"options"       db:"options"`        // JSONB
    AnswerHash   string       `json:"-"             db:"answer_hash"`   // NEVER exposed
    Explanation  string       `json:"explanation"   db:"explanation"`    // Only after submit
    OrderNum     int          `json:"order"         db:"order_num"`
    CreatedAt    time.Time    `json:"created_at"    db:"created_at"`
}

type Option struct {
    Key  string `json:"key"`   // "A", "B", "C", "D"
    Text string `json:"text"`
}
```

### TypeScript Interface

```typescript
export interface QuizQuestion {
  id: string;
  quizId: string;
  question: string;         // mapped from question_text
  questionType: 'multiple_choice' | 'essay';
  options: Option[];
  explanation?: string;      // Only available after submission
  order: number;
}

export interface Option {
  key: string;  // "A", "B", "C", "D"
  text: string;
}
```

---

## Field Mapping Quick Reference

| PostgreSQL | Go Struct | Go JSON Tag | TypeScript |
|:---|:---|:---|:---|
| `id` | `ID` | `"id"` | `id` |
| `email` | `Email` | `"email"` | `email` |
| `school_id` | `SchoolID` | `"school_id"` | `schoolId` |
| `is_active` | `IsActive` | `"is_active"` | `isActive` |
| `created_at` | `CreatedAt` | `"created_at"` | `createdAt` |
| `time_limit` | `TimeLimit` | `"time_limit"` | `timeLimit` |
| `question_text` | `QuestionText` | `"question"` | `question` |
| `answer_hash` | `AnswerHash` | `"-"` (hidden) | _(not exposed)_ |
| `failed_login_count` | `FailedLoginCount` | `"-"` (hidden) | _(not exposed)_ |
| `deleted_at` | `DeletedAt` | `"-"` (hidden) | _(not exposed)_ |

---

## Rules

1. **`json:"-"` untuk field sensitif** → Password, hash, IP, deleted_at tidak boleh terekspos
2. **Nullable field = pointer di Go** → `*time.Time`, `*string` untuk field yang bisa NULL
3. **Zod schema di frontend** → Validasi runtime saat consume API response
4. **JSON menggunakan `snake_case`** → Frontend melakukan mapping ke camelCase

---

## Referensi

- [API Spec](API_SPEC.md) — Kontrak endpoint
- [Database Schema](DATABASE_SCHEMA.md) — DDL lengkap
- [Coding Standards](../fase-2/CODING_STANDARDS.md) — Naming convention

---

*Terakhir diperbarui: 21 Maret 2026*
