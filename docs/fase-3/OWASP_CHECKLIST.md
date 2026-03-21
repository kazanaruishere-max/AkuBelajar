# ✅ OWASP Top 10 Checklist — AkuBelajar

> Daftar centang OWASP Top 10 (2021) yang **wajib dilewati** sebelum setiap rilis.

---

## Cara Menggunakan

- Sebelum setiap rilis, jalankan checklist ini
- Setiap item harus di-check oleh minimal 1 security reviewer
- Item yang gagal = **rilis diblokir** sampai diperbaiki

---

## A01:2021 — Broken Access Control

- [x] RBAC middleware aktif di semua endpoint
- [x] RLS (Row-Level Security) aktif di PostgreSQL
- [x] UUID v7 sebagai primary key (anti-IDOR)
- [x] CORS policy restrictive (hanya domain yang diizinkan)
- [x] Directory listing disabled di Nginx
- [x] Force HTTPS (redirect HTTP → HTTPS)

## A02:2021 — Cryptographic Failures

- [x] Password hashing: Argon2id (bukan bcrypt/MD5/SHA)
- [x] Kunci jawaban kuis di-hash di database
- [x] TLS 1.3 untuk semua komunikasi
- [x] Secrets di HashiCorp Vault (bukan di .env production)
- [x] Sensitive data tidak di-log (password, token, API key)

## A03:2021 — Injection

- [x] Parameterized queries di semua SQL (pgx)
- [x] Input sanitization di Go (validator library)
- [x] CSP headers untuk mencegah XSS
- [x] React auto-escaping aktif (default)
- [x] AI prompt sanitization (blocked phrases, length limit)

## A04:2021 — Insecure Design

- [x] Threat model terdokumentasi (STRIDE)
- [x] Rate limiting per endpoint
- [x] Account lockout setelah 5 failed login
- [x] Server-authoritative logic (timer, grading)
- [x] Input validation di backend (bukan hanya frontend)

## A05:2021 — Security Misconfiguration

- [x] Default credentials diganti sebelum production
- [x] Debug mode off di production
- [x] Error messages generik di production
- [x] Unnecessary ports/services ditutup
- [x] Security headers (X-Frame-Options, X-Content-Type-Options, HSTS)

## A06:2021 — Vulnerable and Outdated Components

- [x] `go mod tidy` + vulnerability check (`govulncheck`)
- [x] `pnpm audit` di CI pipeline
- [x] Dependabot / Renovate aktif di GitHub
- [x] Docker base image terbaru (alpine-based)
- [ ] Jadwal review dependensi bulanan

## A07:2021 — Identification and Authentication Failures

- [x] Multi-factor authentication (roadmap)
- [x] Paseto/JWT dengan short expiry (15 min access, 7 day refresh)
- [x] Refresh token rotation
- [x] Session invalidation saat password change
- [x] Brute force protection (rate limit + lockout)

## A08:2021 — Software and Data Integrity Failures

- [x] CI/CD pipeline integrity (signed commits)
- [x] Immutable audit log (REVOKE UPDATE/DELETE)
- [x] Optimistic locking untuk concurrent edits
- [x] Checksum pada file upload
- [x] Docker image vulnerability scanning

## A09:2021 — Security Logging and Monitoring Failures

- [x] Structured logging (zap) dengan correlation ID
- [x] Audit log untuk semua mutasi data kritis
- [x] Grafana + Prometheus monitoring
- [x] Alert untuk anomali (spike in 401/403, unusual IP)
- [ ] SIEM integration (roadmap)

## A10:2021 — Server-Side Request Forgery (SSRF)

- [x] AI service hanya call Gemini API (whitelist)
- [x] User-supplied URL validation
- [x] Internal service communication via private network
- [x] Metadata endpoint blocked (169.254.169.254)

---

## Summary

| Category | Status | Score |
|:---|:---|:---|
| A01 — Broken Access Control | ✅ Pass | 6/6 |
| A02 — Cryptographic Failures | ✅ Pass | 5/5 |
| A03 — Injection | ✅ Pass | 5/5 |
| A04 — Insecure Design | ✅ Pass | 5/5 |
| A05 — Security Misconfiguration | ✅ Pass | 5/5 |
| A06 — Outdated Components | ⚠️ Partial | 4/5 |
| A07 — Auth Failures | ✅ Pass | 5/5 |
| A08 — Data Integrity | ✅ Pass | 5/5 |
| A09 — Logging Failures | ⚠️ Partial | 4/5 |
| A10 — SSRF | ✅ Pass | 4/4 |

---

*Terakhir diperbarui: 21 Maret 2026*
