# 🛡️ Threat Model — AkuBelajar

> Dokumentasi semua kemungkinan serangan beserta mitigasinya menggunakan metodologi **STRIDE**.

---

## Metodologi: STRIDE

| Kategori | Deskripsi | Pertanyaan Kunci |
|:---|:---|:---|
| **S**poofing | Menyamar sebagai entitas lain | Siapa bisa berpura-pura jadi orang lain? |
| **T**ampering | Memodifikasi data tanpa izin | Data apa yang bisa diubah secara ilegal? |
| **R**epudiation | Menyangkal telah melakukan tindakan | Aksi apa yang tidak bisa dibuktikan? |
| **I**nformation Disclosure | Kebocoran informasi | Data sensitif apa yang bisa bocor? |
| **D**enial of Service | Membuat layanan tidak tersedia | Apa yang bisa membuat sistem down? |
| **E**levation of Privilege | Mendapatkan akses lebih tinggi | Bagaimana user biasa bisa jadi admin? |

---

## Threat Matrix

### S — Spoofing (Pemalsuan Identitas)

| # | Ancaman | Severity | Mitigasi | Status |
|:---|:---|:---|:---|:---|
| S1 | Brute force login | 🔴 High | Rate limiter (5 attempts/15 min), account lockout, Argon2id | ✅ Implemented |
| S2 | Credential stuffing | 🔴 High | Rate limiter per IP, breached password check | ✅ Implemented |
| S3 | Session hijacking | 🔴 High | HttpOnly cookies, Secure flag, SameSite=Strict | ✅ Implemented |
| S4 | JWT token theft | 🟡 Medium | Short TTL (15 min), refresh token rotation | ✅ Implemented |
| S5 | Phishing via email | 🟡 Medium | DMARC/SPF/DKIM, user education | 🟡 Partial |

### T — Tampering (Manipulasi Data)

| # | Ancaman | Severity | Mitigasi | Status |
|:---|:---|:---|:---|:---|
| T1 | SQL Injection | 🔴 Critical | Parameterized queries (pgx), no string concat | ✅ Implemented |
| T2 | XSS (Cross-Site Scripting) | 🔴 High | CSP headers, React auto-escaping, DOMPurify | ✅ Implemented |
| T3 | Manipulasi nilai/grade | 🔴 High | Immutable audit log, optimistic locking, RLS | ✅ Implemented |
| T4 | Timer manipulation (CBT) | 🟡 Medium | Server-authoritative timer, WebSocket heartbeat | ✅ Implemented |
| T5 | File upload malware | 🟡 Medium | MIME check, size limit, antivirus scan, isolated storage | ✅ Implemented |

### R — Repudiation (Penyangkalan)

| # | Ancaman | Severity | Mitigasi | Status |
|:---|:---|:---|:---|:---|
| R1 | Guru menyangkal telah mengubah nilai | 🟡 Medium | Immutable audit log (no UPDATE/DELETE) | ✅ Implemented |
| R2 | Admin menyangkal menghapus user | 🟡 Medium | Audit log dengan IP address + user agent | ✅ Implemented |
| R3 | Siswa menyangkal telah submit kuis | 🟢 Low | Submission timestamp + hash + audit trail | ✅ Implemented |

### I — Information Disclosure (Kebocoran Data)

| # | Ancaman | Severity | Mitigasi | Status |
|:---|:---|:---|:---|:---|
| I1 | IDOR — siswa lihat nilai siswa lain | 🔴 High | UUID v7 PK + RLS policy per user | ✅ Implemented |
| I2 | API response leaking sensitive fields | 🟡 Medium | Explicit field selection (no SELECT *) | ✅ Implemented |
| I3 | Error message leaking stack trace | 🟡 Medium | Generic error di production, detail di log | ✅ Implemented |
| I4 | Kunci jawaban kuis terekspos | 🔴 High | Jawaban di-hash (Argon2id) di DB | ✅ Implemented |
| I5 | Secret di source code | 🔴 Critical | HashiCorp Vault, .env di .gitignore, git-secrets | ✅ Implemented |

### D — Denial of Service

| # | Ancaman | Severity | Mitigasi | Status |
|:---|:---|:---|:---|:---|
| D1 | DDoS attack | 🔴 High | Cloudflare WAF + DDoS protection | ✅ Implemented |
| D2 | API abuse (excessive requests) | 🟡 Medium | Per-user rate limiting (Redis sliding window) | ✅ Implemented |
| D3 | Large file upload DoS | 🟡 Medium | File size limit (10MB), timeout | ✅ Implemented |
| D4 | Slowloris attack | 🟡 Medium | Nginx timeout config, connection limits | ✅ Implemented |
| D5 | Database connection exhaustion | 🟡 Medium | PgBouncer, connection pool limits | ✅ Implemented |

### E — Elevation of Privilege

| # | Ancaman | Severity | Mitigasi | Status |
|:---|:---|:---|:---|:---|
| E1 | Siswa akses endpoint guru | 🔴 Critical | RBAC middleware per endpoint | ✅ Implemented |
| E2 | JWT role manipulation | 🔴 High | Server-side role validation, Paseto signed tokens | ✅ Implemented |
| E3 | Mass assignment | 🟡 Medium | Explicit struct binding (Gin ShouldBindJSON) | ✅ Implemented |
| E4 | Prompt injection ke AI | 🟡 Medium | Input sanitization, blocked phrases, length limit | ✅ Implemented |

---

## Attack Surface Map

```
Internet
    │
    ▼
[Cloudflare WAF] ──── Block malicious traffic
    │
    ▼
[Nginx] ──── TLS, rate limit, request size limit
    │
    ▼
[Go API] ──── Auth, RBAC, input validation, audit log
    │
    ▼
[PostgreSQL + RLS] ──── Row-level isolation, immutable audit
```

---

## Referensi

- [OWASP Checklist](OWASP_CHECKLIST.md)
- [Incident Response](INCIDENT_RESPONSE.md)
- [Security Policy](../fase-0/SECURITY.md)

---

*Terakhir diperbarui: 21 Maret 2026*
