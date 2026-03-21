# 🔑 Environment Variables — AkuBelajar

> Daftar semua env var yang dibutuhkan beserta referensi ke HashiCorp Vault. **Tidak ada secret yang tersimpan di repo.**

---

## Prinsip

1. **Tidak ada secret di source code** — semua secret di HashiCorp Vault
2. **`.env` hanya untuk development** — production menggunakan K8s Secrets dari Vault
3. **Validasi saat startup** — aplikasi gagal start jika env var kritis tidak ada

---

## Daftar Environment Variables

### Application Core

| Variable | Required | Default | Vault Path | Deskripsi |
|:---|:---|:---|:---|:---|
| `APP_ENV` | ✅ | `development` | — | Environment: development, staging, production |
| `APP_PORT` | ✅ | `8080` | — | Port API server |
| `APP_NAME` | ❌ | `akubelajar` | — | Nama aplikasi untuk logging |
| `FRONTEND_URL` | ✅ | `http://localhost:3000` | — | URL frontend (untuk CORS) |
| `LOG_LEVEL` | ❌ | `info` | — | debug, info, warn, error |

### Database (PostgreSQL)

| Variable | Required | Vault Path | Deskripsi |
|:---|:---|:---|:---|
| `DB_HOST` | ✅ | `secret/akubelajar/db#host` | PostgreSQL host |
| `DB_PORT` | ✅ | `secret/akubelajar/db#port` | PostgreSQL port |
| `DB_USER` | ✅ | `secret/akubelajar/db#user` | Database username |
| `DB_PASSWORD` | ✅ | `secret/akubelajar/db#password` | Database password |
| `DB_NAME` | ✅ | `secret/akubelajar/db#name` | Database name |
| `DB_SSL_MODE` | ❌ | — | disable (dev), require (prod) |
| `DB_MAX_CONNECTIONS` | ❌ | — | Connection pool max (default: 50) |

### Cache (Redis)

| Variable | Required | Vault Path | Deskripsi |
|:---|:---|:---|:---|
| `REDIS_HOST` | ✅ | `secret/akubelajar/redis#host` | Redis host |
| `REDIS_PORT` | ✅ | `secret/akubelajar/redis#port` | Redis port |
| `REDIS_PASSWORD` | ✅ | `secret/akubelajar/redis#password` | Redis password |
| `REDIS_DB` | ❌ | — | Redis database number (default: 0) |

### Authentication

| Variable | Required | Vault Path | Deskripsi |
|:---|:---|:---|:---|
| `JWT_SECRET` | ✅ | `secret/akubelajar/auth#jwt_secret` | JWT signing key (min 32 chars) |
| `JWT_EXPIRY` | ❌ | — | Access token TTL (default: 15m) |
| `REFRESH_TOKEN_EXPIRY` | ❌ | — | Refresh token TTL (default: 7d) |
| `PASETO_PRIVATE_KEY` | ✅ | `secret/akubelajar/auth#paseto_key` | Paseto v4 private key |

### AI Service (Gemini)

| Variable | Required | Vault Path | Deskripsi |
|:---|:---|:---|:---|
| `GEMINI_API_KEY` | ✅ | `secret/akubelajar/ai#gemini_key` | Google Gemini API key |
| `GEMINI_MODEL` | ❌ | — | Model name (default: gemini-2.0-flash) |
| `GEMINI_MAX_TOKENS` | ❌ | — | Max output tokens (default: 4096) |

### Object Storage (MinIO)

| Variable | Required | Vault Path | Deskripsi |
|:---|:---|:---|:---|
| `MINIO_ENDPOINT` | ✅ | `secret/akubelajar/storage#endpoint` | MinIO endpoint |
| `MINIO_ACCESS_KEY` | ✅ | `secret/akubelajar/storage#access_key` | Access key |
| `MINIO_SECRET_KEY` | ✅ | `secret/akubelajar/storage#secret_key` | Secret key |
| `MINIO_BUCKET` | ✅ | — | Bucket name (default: akubelajar) |

### Email / Notification

| Variable | Required | Vault Path | Deskripsi |
|:---|:---|:---|:---|
| `SMTP_HOST` | ✅ | `secret/akubelajar/email#host` | SMTP server host |
| `SMTP_PORT` | ✅ | — | SMTP port (587 TLS) |
| `SMTP_USER` | ✅ | `secret/akubelajar/email#user` | SMTP username |
| `SMTP_PASSWORD` | ✅ | `secret/akubelajar/email#password` | SMTP password |

---

## Vault Integration (Production)

```bash
# Login ke Vault
vault login -method=oidc

# Baca secret
vault kv get secret/akubelajar/db

# K8s: inject secrets via Vault Agent
# Lihat deployment manifest di fase-4/DEPLOYMENT_GUIDE.md
```

---

## Validasi Startup

```go
// config/config.go
func LoadConfig() (*Config, error) {
    cfg := &Config{}
    if err := envconfig.Process("", cfg); err != nil {
        return nil, fmt.Errorf("missing required env vars: %w", err)
    }
    // Validasi tambahan
    if len(cfg.JWTSecret) < 32 {
        return nil, errors.New("JWT_SECRET must be at least 32 characters")
    }
    return cfg, nil
}
```

---

*Terakhir diperbarui: 21 Maret 2026*
