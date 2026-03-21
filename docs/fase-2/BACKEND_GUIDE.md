# ⚙️ Backend Guide — AkuBelajar

> Panduan khusus lapisan backend Go: arsitektur, pattern, error handling, dan best practices.

---

## Arsitektur Service

```
internal/
├── auth/              # Autentikasi & otorisasi
│   ├── handler.go     # HTTP handler (input/output)
│   ├── service.go     # Business logic
│   ├── repository.go  # Database queries
│   └── model.go       # Request/Response structs
├── academic/          # Modul akademik
├── quiz/              # Quiz & CBT engine
├── ai/                # Gemini AI integration
├── notification/      # Multi-channel notifikasi
└── middleware/        # Auth, RBAC, rate limit, logging
```

### Layer Pattern

```
Handler (HTTP) → Service (Logic) → Repository (DB)
      ↓                ↓                  ↓
 Validate input    Apply rules       Execute query
 Parse request     Orchestrate       Return data
 Return response   Call external
```

**Rules:**
- Handler **tidak boleh** query database langsung
- Service **tidak boleh** akses `*gin.Context`
- Repository **hanya** berisi SQL queries, tidak ada logic

---

## Error Handling Pattern

```go
// Definisikan custom errors per domain
var (
    ErrQuizNotFound     = errors.New("quiz not found")
    ErrQuizAlreadyTaken = errors.New("quiz already submitted")
    ErrUnauthorized     = errors.New("unauthorized access")
)

// Service layer — return domain error
func (s *QuizService) GetQuiz(ctx context.Context, id uuid.UUID) (*Quiz, error) {
    quiz, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("get quiz %s: %w", id, err)
    }
    if quiz == nil {
        return nil, ErrQuizNotFound
    }
    return quiz, nil
}

// Handler layer — translate ke HTTP status
func (h *QuizHandler) GetQuiz(c *gin.Context) {
    quiz, err := h.service.GetQuiz(c.Request.Context(), id)
    if err != nil {
        switch {
        case errors.Is(err, ErrQuizNotFound):
            c.JSON(404, gin.H{"error": "quiz_not_found"})
        case errors.Is(err, ErrUnauthorized):
            c.JSON(403, gin.H{"error": "forbidden"})
        default:
            c.JSON(500, gin.H{"error": "internal_error"})
            h.logger.Error("unexpected error", zap.Error(err))
        }
        return
    }
    c.JSON(200, quiz)
}
```

---

## Database Patterns

### Transaction

```go
func (r *QuizRepo) SubmitAnswers(ctx context.Context, submission *Submission) error {
    tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    defer tx.Rollback(ctx) // safe — no-op jika sudah commit

    // Insert submission
    _, err = tx.Exec(ctx, `INSERT INTO quiz_submissions ...`, ...)
    if err != nil {
        return fmt.Errorf("insert submission: %w", err)
    }

    // Insert answers
    for _, answer := range submission.Answers {
        _, err = tx.Exec(ctx, `INSERT INTO quiz_answers ...`, ...)
        if err != nil {
            return fmt.Errorf("insert answer: %w", err)
        }
    }

    return tx.Commit(ctx)
}
```

### Pagination

```go
// Cursor-based pagination (lebih performan dari OFFSET)
func (r *UserRepo) ListUsers(ctx context.Context, cursor *uuid.UUID, limit int) ([]User, error) {
    query := `SELECT id, email, role FROM users WHERE deleted_at IS NULL`
    args := []interface{}{}

    if cursor != nil {
        query += ` AND id > $1`
        args = append(args, *cursor)
    }
    query += fmt.Sprintf(` ORDER BY id LIMIT %d`, limit)

    rows, err := r.pool.Query(ctx, query, args...)
    // ...
}
```

---

## Middleware Stack

```go
router.Use(
    middleware.RequestID(),        // Tambah X-Request-ID
    middleware.Logger(logger),     // Structured logging
    middleware.Recovery(),         // Panic recovery
    middleware.CORS(corsConfig),   // CORS policy
    middleware.RateLimiter(redis), // Rate limiting
    middleware.Authenticate(),     // JWT/Paseto validation
)
```

---

## Logging Standard

```go
// ✅ Structured logging dengan zap
logger.Info("quiz created",
    zap.String("quiz_id", quiz.ID.String()),
    zap.String("teacher_id", teacherID.String()),
    zap.Int("question_count", len(quiz.Questions)),
)

// ❌ JANGAN gunakan fmt.Println di production
fmt.Println("quiz created:", quiz.ID)
```

---

## Referensi

- [Coding Standards](CODING_STANDARDS.md)
- [Testing Strategy](TESTING_STRATEGY.md)
- [Database Schema](../fase-1/DATABASE_SCHEMA.md)

---

*Terakhir diperbarui: 21 Maret 2026*
