# 🔑 Password Reset Flow — AkuBelajar

> Semua skenario reset password: self-service via OTP dan admin reset.

---

## 1. Self-Service Reset (via Halaman Login)

```mermaid
sequenceDiagram
    participant User
    participant FE as Next.js
    participant API as Go Backend
    participant Redis
    participant WA as Fonnte WA API
    participant Email as SMTP

    User->>FE: Klik "Lupa Password"
    FE->>FE: Tampilkan form (input email/WA)
    User->>FE: Input identifier
    FE->>API: POST /auth/password-reset/request
    API->>API: Cari user by email ATAU phone_wa
    alt User tidak ditemukan
        API-->>FE: 200 "Jika terdaftar, OTP telah dikirim" (pesan generik)
    else User ditemukan
        API->>API: Generate OTP 6 digit
        API->>API: Hash OTP (SHA-256)
        API->>Redis: Store hash + rate limit check
        alt Punya nomor WA
            API->>WA: Kirim OTP via WA (prioritas)
        else Fallback email
            API->>Email: Kirim OTP via Email
        end
        API-->>FE: 200 "Jika terdaftar, OTP telah dikirim"
    end
    User->>FE: Input OTP + password baru
    FE->>API: POST /auth/password-reset/verify
    API->>API: Constant-time compare OTP hash
    API->>API: Hash password baru (Argon2id)
    API->>API: REVOKE semua active sessions
    API-->>FE: 200 "Password berhasil diubah"
    FE->>FE: Redirect ke /login
```

### Batasan OTP

| Parameter | Nilai |
|:---|:---|
| Panjang OTP | 6 digit angka |
| TTL | 5 menit |
| Rate limit | 3 request per jam per IP |
| Penyimpanan | SHA-256 hash (BUKAN plain text) |
| Validasi | Constant-time comparison |

---

## 2. Reset oleh Admin

```mermaid
sequenceDiagram
    participant Admin as SuperAdmin/Guru
    participant API as Go Backend
    participant DB as PostgreSQL

    Admin->>API: PUT /users/:id/reset-password
    API->>API: Cek: Admin punya izin atas user ini?
    API->>API: Generate temp password (12 char random)
    API->>DB: UPDATE password (Argon2id hash)
    API->>DB: SET is_first_login = TRUE
    API->>DB: REVOKE semua active sessions user
    API->>DB: INSERT audit_log
    API-->>Admin: 200 { "temporary_password": "xK9mR2pL4wQs" }
    Note over Admin: Password ditampilkan SEKALI di layar
    Admin->>Admin: Sampaikan manual ke siswa
```

- SuperAdmin bisa reset **semua user** di sekolahnya
- Guru bisa reset **siswa di kelasnya**
- Password sementara ditampilkan **sekali saja** (tidak dikirim ulang)
- `is_first_login = TRUE` → siswa dipaksa ganti saat login
- Audit log mencatat: admin ID, user ID, timestamp, IP

---

## 3. Schema: password_reset_tokens

```sql
CREATE TABLE password_reset_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id),
    token_hash  VARCHAR(255) NOT NULL,     -- SHA-256 hash of OTP
    expires_at  TIMESTAMPTZ NOT NULL,      -- NOW() + 5 menit
    used_at     TIMESTAMPTZ,               -- NULL = belum dipakai
    ip_address  INET,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_prt_user ON password_reset_tokens(user_id, created_at);
```

---

## 4. Security Constraints

| Constraint | Detail |
|:---|:---|
| OTP storage | SHA-256 hash, BUKAN plain text |
| Anti-enumeration | Response selalu generik ("Jika terdaftar...") |
| Constant-time comparison | Mencegah timing attack |
| Token cleanup | Scheduled job hapus token expired setiap 1 jam |
| Session revocation | Semua session aktif di-revoke setelah reset |
| Password history | Tidak boleh reuse 3 password terakhir |

---

## 5. Edge Cases

| Skenario | Penanganan |
|:---|:---|
| OTP yang sudah dipakai | Error `AUTH_009`: "Kode OTP salah" (pesan generik) |
| Request OTP berulang cepat | Rate limit: 3/jam per IP → `AUTH_011` |
| User ganti password lalu lupa lagi | Flow normal — bisa request OTP lagi (setelah cooldown habis) |
| Email/WA tidak ditemukan | Tetap tampilkan pesan sukses generik (anti-enumeration) |
| OTP dikirim via WA tapi nomor tidak valid | Fallback ke email. Jika email juga gagal → OTP tidak terkirim, user harus hubungi admin |

---

*Terakhir diperbarui: 21 Maret 2026*
