# 🚫 FORBIDDEN_PATTERNS.md — Pola Kode yang DILARANG

> Daftar anti-pattern yang TIDAK BOLEH ditulis. AI agent harus cek file ini sebelum coding.

---

## 1. Keamanan

### ❌ DILARANG: JWT

```go
// ❌ SALAH — kita TIDAK pakai JWT
import "github.com/golang-jwt/jwt"
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// ✅ BENAR — pakai Paseto v4
import "aidanwoods.dev/go-paseto"
token := paseto.NewToken()
token.SetExpiration(time.Now().Add(15 * time.Minute))
```

### ❌ DILARANG: bcrypt

```go
// ❌ SALAH
import "golang.org/x/crypto/bcrypt"
hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

// ✅ BENAR — pakai Argon2id
import "github.com/alexedwards/argon2id"
hash, _ := argon2id.CreateHash(password, argon2id.DefaultParams)
```

### ❌ DILARANG: String concatenation SQL

```go
// ❌ SALAH — SQL INJECTION RISK!
query := "SELECT * FROM users WHERE email = '" + email + "'"

// ✅ BENAR — Parameterized query
query := "SELECT id, email, role FROM users WHERE email = $1"
row := db.QueryRow(ctx, query, email)
```

### ❌ DILARANG: Credentials di log

```go
// ❌ SALAH
logger.Info("login", zap.String("password", password))
logger.Info("token", zap.String("access_token", token))
logger.Info("otp", zap.String("code", otp))

// ✅ BENAR
logger.Info("login_attempt", zap.String("email", maskEmail(email)))
```

### ❌ DILARANG: Hardcoded secrets

```go
// ❌ SALAH
pasetoKey := "my-super-secret-key-12345678901"
dbURL := "postgres://admin:password123@localhost/db"

// ✅ BENAR
pasetoKey := os.Getenv("PASETO_KEY")
dbURL := os.Getenv("DATABASE_URL")
```

---

## 2. Database

### ❌ DILARANG: SELECT *

```go
// ❌ SALAH — ambil semua kolom termasuk password_hash
rows, _ := db.Query(ctx, "SELECT * FROM users WHERE id = $1", id)

// ✅ BENAR — list kolom eksplisit
rows, _ := db.Query(ctx, "SELECT id, email, role, is_active FROM users WHERE id = $1", id)
```

### ❌ DILARANG: Skip error handling

```go
// ❌ SALAH
result, _ := db.Exec(ctx, query, args...)
json.NewDecoder(r.Body).Decode(&input) // error ignored!

// ✅ BENAR
result, err := db.Exec(ctx, query, args...)
if err != nil {
    return fmt.Errorf("insert user: %w", err)
}
```

### ❌ DILARANG: Auto-increment ID

```sql
-- ❌ SALAH
CREATE TABLE users (
    id SERIAL PRIMARY KEY
);

-- ✅ BENAR — UUID v7
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
);
```

### ❌ DILARANG: Direct schema change tanpa migration

```go
// ❌ SALAH — ALTER TABLE langsung dari kode
db.Exec(ctx, "ALTER TABLE users ADD COLUMN phone TEXT")

// ✅ BENAR — buat migration file
// backend/migrations/000016_add_phone_to_users.up.sql
// backend/migrations/000016_add_phone_to_users.down.sql
```

---

## 3. Backend Pattern

### ❌ DILARANG: Business logic di handler

```go
// ❌ SALAH — handler harus tipis
func CreateAssignment(c *gin.Context) {
    // 50 baris business logic di sini...
    db.Exec(query)
}

// ✅ BENAR — handler → service → repository
func CreateAssignment(c *gin.Context) {
    var req CreateAssignmentRequest
    if err := c.ShouldBindJSON(&req); err != nil { ... }
    result, err := assignmentService.Create(ctx, req)
    c.JSON(200, result)
}
```

### ❌ DILARANG: Panic tanpa recover

```go
// ❌ SALAH
panic("something went wrong")

// ✅ BENAR — return error
return fmt.Errorf("something went wrong: %w", err)
```

---

## 4. Frontend Pattern

### ❌ DILARANG: useEffect untuk data fetching

```tsx
// ❌ SALAH
useEffect(() => {
  fetch('/api/users').then(r => r.json()).then(setUsers)
}, [])

// ✅ BENAR — TanStack Query
const { data: users } = useQuery({
  queryKey: ['users'],
  queryFn: () => api.getUsers()
})
```

### ❌ DILARANG: Global state untuk server data

```tsx
// ❌ SALAH — Zustand untuk API response
const useStore = create(set => ({
  users: [],
  fetchUsers: async () => { ... set({ users }) }
}))

// ✅ BENAR — TanStack Query untuk server state, Zustand untuk client-only state
// TanStack Query: users, assignments, grades (server data)
// Zustand: sidebar open, theme, form draft (client state)
```

### ❌ DILARANG: Inline styles untuk layout

```tsx
// ❌ SALAH
<div style={{ display: 'flex', padding: '16px' }}>

// ✅ BENAR — CSS modules / design tokens
<div className={styles.container}>
```

---

*Terakhir diperbarui: 21 Maret 2026*
