# 🔏 Password Policy — AkuBelajar

> Aturan kekuatan password, hashing Argon2id, password history, dan change policy.

---

## 1. Aturan Kekuatan Password

| Rule | Nilai |
|:---|:---|
| Min length | 8 karakter |
| Max length | 72 karakter |
| Wajib | 1 huruf besar, 1 huruf kecil, 1 angka, 1 simbol |
| Tidak boleh | Mengandung nama, email, atau NISN user |
| Cek | Terhadap 10.000 common passwords |

### Contoh

| Password | Valid? | Alasan |
|:---|:---|:---|
| `Guru@2026!` | ✅ | Memenuhi semua kriteria |
| `P@ssw0rd!` | ❌ | Di daftar common passwords |
| `gurubudisantoso1!` | ❌ | Tidak ada huruf besar |
| `ABCDEFGH` | ❌ | Tidak ada lowercase, angka, simbol |
| `Gu@1` | ❌ | Kurang dari 8 karakter |

---

## 2. Hashing dengan Argon2id

### Parameter

| Parameter | Nilai | Alasan |
|:---|:---|:---|
| Memory | 65536 KB (64 MB) | Mempersulit GPU attack |
| Iterations | 3 | Balance antara keamanan dan speed |
| Parallelism | 4 | Manfaatkan multi-core |
| Salt | 16 byte random per password | Unique per user |
| Output length | 32 byte | Standard security |

### Go Implementation

```go
import "golang.org/x/crypto/argon2"

type Argon2Params struct {
    Memory      uint32 // 65536
    Iterations  uint32 // 3
    Parallelism uint8  // 4
    SaltLength  uint32 // 16
    KeyLength   uint32 // 32
}

func HashPassword(password string) (string, error) {
    salt := make([]byte, 16)
    _, _ = rand.Read(salt)
    
    hash := argon2.IDKey([]byte(password), salt, 3, 65536, 4, 32)
    
    // Encode as: $argon2id$v=19$m=65536,t=3,p=4$<salt>$<hash>
    return fmt.Sprintf(
        "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
        argon2.Version, 65536, 3, 4,
        base64.RawStdEncoding.EncodeToString(salt),
        base64.RawStdEncoding.EncodeToString(hash),
    ), nil
}

func VerifyPassword(password, encoded string) (bool, error) {
    // Parse encoded string → extract params, salt, hash
    // Recompute hash with same params + salt
    // Constant-time compare
    return subtle.ConstantTimeCompare(hash, recomputed) == 1, nil
}
```

### Kapan Upgrade Parameter?

- Jika hash time < 0.5 detik pada server production → naikkan iterations
- Review setiap 2 tahun atau saat upgrade hardware
- Rehash pada login berikutnya (transparent upgrade)

---

## 3. Password History

- Tidak boleh reuse **3 password terakhir**
- History disimpan sebagai **hash** (bukan plain text)

```sql
CREATE TABLE password_histories (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    password_hash   VARCHAR(255) NOT NULL,  -- Argon2id hash
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

### Check Flow

```
1. User input password baru
2. Backend hash password baru
3. Verify terhadap hash di password_histories (max 3 terbaru)
4. Jika match salah satu → REJECT: "Tidak boleh menggunakan password yang sama"
5. Jika tidak match → simpan hash baru ke password_histories
```

---

## 4. Temporary Password

| Property | Nilai |
|:---|:---|
| Format | 12 karakter: uppercase + lowercase + angka |
| Generator | `crypto/rand` (Go) |
| TTL | **7 hari** atau sampai digunakan (mana yang lebih dulu) |
| Setelah expired | Akun di-lock, admin harus generate ulang |

```go
func GenerateTempPassword() string {
    const charset = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789"
    b := make([]byte, 12)
    for i := range b {
        n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
        b[i] = charset[n.Int64()]
    }
    return string(b)
}
// Contoh output: "xK9mR2pL4wQs"
// Note: excluded ambiguous chars (0, O, 1, l, I)
```

---

## 5. Password Change Policy

| Scenario | Force? | Revoke Sessions? |
|:---|:---|:---|
| First login | ✅ Force | ❌ (belum ada session lain) |
| Setelah reset | ✅ Force | ✅ Semua session |
| Voluntary (settings) | ❌ | ✅ Semua session kecuali current |

- **Cooldown** sebelum ganti password lagi: **24 jam** — mencegah abuse jika akun compromised
- Tujuan: mencegah abuse jika akun compromised → attacker ganti password berulang

---

*Terakhir diperbarui: 21 Maret 2026*
