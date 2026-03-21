# 🔐 Secure Coding — AkuBelajar

> Aturan kode yang WAJIB diikuti agar setiap baris kode yang dihasilkan secara default sudah aman. Berbeda dari OWASP_CHECKLIST.md (audit), ini adalah **rules saat menulis kode**.

---

## Prinsip

1. **Secure by Default** — Kode baru harus aman tanpa konfigurasi tambahan
2. **Fail Secure** — Jika ragu, tolak akses / reject input
3. **Defense in Depth** — Jangan bergantung pada satu layer keamanan saja
4. **Least Privilege** — Berikan akses minimum yang dibutuhkan

---

## 🔴 RULE KERAS (Tidak Boleh Dilanggar)

### SQL

```go
// ✅ WAJIB — Parameterized query
row := db.QueryRow(ctx, "SELECT * FROM users WHERE email = $1", email)

// ❌ DILARANG KERAS — String concatenation
row := db.QueryRow(ctx, "SELECT * FROM users WHERE email = '" + email + "'")

// ❌ DILARANG KERAS — fmt.Sprintf untuk query
query := fmt.Sprintf("SELECT * FROM users WHERE id = '%s'", id)
```

### Input Validation

```go
// ✅ WAJIB — Validasi sebelum diproses (Go)
type CreateQuizRequest struct {
    Title     string `json:"title"     validate:"required,min=3,max=255"`
    SubjectID string `json:"subject_id" validate:"required,uuid"`
    TimeLimit int    `json:"time_limit" validate:"required,min=5,max=180"`
}

if err := validator.Validate(req); err != nil {
    return BadRequest(c, err)
}
```

```typescript
// ✅ WAJIB — Zod schema sebelum diproses (TypeScript)
const CreateQuizSchema = z.object({
  title: z.string().min(3).max(255),
  subjectId: z.string().uuid(),
  timeLimit: z.number().min(5).max(180),
});

const data = CreateQuizSchema.parse(rawInput);

// ❌ DILARANG — Trust tanpa validasi
const data = req.body as CreateQuizInput; // type assertion bukan validasi!
```

### Password & Hashing

```go
// ✅ WAJIB — Argon2id untuk password
hash, _ := argon2id.CreateHash(password, argon2id.DefaultParams)

// ❌ DILARANG — MD5, SHA1, SHA256 untuk password
hash := md5.Sum([]byte(password))

// ❌ DILARANG — bcrypt (sudah dianggap kurang kuat)
hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
```

### Secrets

```go
// ✅ WAJIB — Dari environment variable
apiKey := os.Getenv("GEMINI_API_KEY")

// ❌ DILARANG KERAS — Hardcoded secret
apiKey := "AIzaSyB1234567890abcdef"

// ❌ DILARANG KERAS — Secret di struct tag, comment, atau test file
```

### Error Messages

```go
// ✅ WAJIB — Generic error ke client
c.JSON(401, gin.H{"error": "invalid_credentials"})

// ❌ DILARANG — Detail internal ke client
c.JSON(401, gin.H{"error": "user with email x@y.com not found in database"})
c.JSON(500, gin.H{"error": err.Error()}) // Leaks stack trace!
```

---

## 🟡 RULE PENTING

### Authentication

| Rule | Detail |
|:---|:---|
| Semua endpoint WAJIB auth middleware | Kecuali: `/health`, `/auth/login`, `/auth/register` |
| Access token TTL | Maksimal 15 menit |
| Refresh token | Rotate setiap kali digunakan (single-use) |
| Signed URL | TTL maksimal 15 menit |
| Session pada password change | Invalidasi SEMUA active session |

### File Upload

| Rule | Detail |
|:---|:---|
| MIME type check | Validasi content, bukan hanya extension |
| Max file size | 10 MB (configurable) |
| Storage path | Gunakan UUID, JANGAN nama file asli |
| Serve files | Via signed URL dengan TTL, BUKAN direct path |
| Antivirus | Scan sebelum simpan (ClamAV) |

### CORS

```go
// ✅ WAJIB — Restrictive CORS
cors.Config{
    AllowOrigins:     []string{"https://app.akubelajar.id"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Authorization", "Content-Type"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}

// ❌ DILARANG di production
AllowOrigins: []string{"*"}
```

### HTTP Headers

```go
// ✅ WAJIB — Security headers di setiap response
c.Header("X-Content-Type-Options", "nosniff")
c.Header("X-Frame-Options", "DENY")
c.Header("X-XSS-Protection", "0")  // Modern: rely on CSP instead
c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
c.Header("Content-Security-Policy", "default-src 'self'; ...")
c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
```

### Logging

```go
// ✅ WAJIB — Log security events
logger.Warn("login_failed", zap.String("email", email), zap.String("ip", ip))

// ❌ DILARANG — Log sensitive data
logger.Info("login", zap.String("password", password))
logger.Info("token", zap.String("jwt", token))
```

---

## Checklist Sebelum Commit

- [ ] Tidak ada SQL string concatenation
- [ ] Semua input melalui validasi (Go validator / Zod)
- [ ] Tidak ada `any` di TypeScript
- [ ] Error messages generik (tidak leak internals)
- [ ] Tidak ada secret/credential hardcoded
- [ ] Auth middleware di semua endpoint baru
- [ ] RBAC check sesuai matriks peran
- [ ] Audit log untuk operasi mutasi data kritis
- [ ] File upload divalidasi MIME + size
- [ ] `json:"-"` pada field sensitif di Go struct

---

## Referensi

- [OWASP Checklist](OWASP_CHECKLIST.md) — Audit keamanan per rilis
- [Threat Model](THREAT_MODEL.md) — Attack surface
- [Coding Standards](../fase-2/CODING_STANDARDS.md) — Style guide umum

---

*Terakhir diperbarui: 21 Maret 2026*
