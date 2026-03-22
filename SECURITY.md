# 🔒 Kebijakan Keamanan

## Versi yang Didukung

| Versi | Status |
|:--|:--|
| 0.x (Demo) | ✅ Didukung |

## Melaporkan Kerentanan

Jika Anda menemukan kerentanan keamanan, **jangan** membuat issue publik. Sebagai gantinya:

1. **Email**: Kirim detail ke alamat email yang tercantum di profil pengelola
2. **Deskripsi**: Jelaskan kerentanan secara detail
3. **Reproduksi**: Sertakan langkah-langkah untuk mereproduksi masalah
4. **Dampak**: Jelaskan potensi dampak kerentanan

### Waktu Respons

- **Konfirmasi penerimaan**: Dalam 48 jam
- **Penilaian awal**: Dalam 1 minggu
- **Perbaikan**: Tergantung severity

## Praktik Keamanan

### Autentikasi

- ✅ **Paseto v4** — Tidak rentan terhadap serangan `alg: none` seperti JWT
- ✅ **Argon2id** — Password hashing terkuat saat ini (OWASP 2024 recommended)
- ✅ **Brute-force protection** — Account lockout setelah 5 kali gagal login
- ✅ **Refresh token rotation** — Token dirotasi setiap kali refresh

### Database

- ✅ **Parameterized queries** — Semua query menggunakan parameter (`$1, $2, ...`), bukan string concatenation
- ✅ **Row Level Security** — Isolasi data per sekolah di PostgreSQL
- ✅ **Soft delete** — Data tidak benar-benar dihapus, hanya ditandai `deleted_at`

### API

- ✅ **Rate limiting** — 120 request/menit per IP
- ✅ **CORS** — Konfigurasi strict per origin
- ✅ **Input validation** — Semua request body divalidasi dengan `go-playground/validator/v10`
- ✅ **Error sanitization** — Error internal tidak di-expose ke client

### Infrastructure

- ✅ **SSL/TLS** — Koneksi database dan Redis terenkripsi (`sslmode=require`, `rediss://`)
- ✅ **Environment variables** — Semua secret disimpan di environment variables, bukan di kode
- ✅ **Graceful shutdown** — Server menunggu request aktif sebelum shutdown

## Dependensi

Kami secara rutin memeriksa dependensi untuk kerentanan yang diketahui:

```bash
# Go
go list -m -json all | nancy sleuth

# Node.js
pnpm audit
```

---

Terima kasih telah membantu menjaga keamanan AkuBelajar! 🛡️
