# 🛡️ Input Validation Rules — AkuBelajar

> Aturan validasi untuk SETIAP field input. Digunakan bersama antara Zod (frontend) dan Go struct validator (backend).

---

## Prinsip

1. **Validate on both sides** — frontend (UX) + backend (security)
2. **Backend adalah sumber kebenaran** — frontend bisa di-bypass
3. **Whitelist, bukan blacklist** — hanya terima yang explisit diizinkan

---

## Field Identitas User

### email

| Rule | Detail |
|:---|:---|
| Format | RFC 5322 |
| Max length | 255 char |
| Transform | lowercase, trim whitespace |
| Unique | Per seluruh sistem |

```typescript
// Zod
const emailSchema = z.string().email().max(255).transform(v => v.toLowerCase().trim());
```
```go
// Go
Email string `json:"email" validate:"required,email,max=255"`
```

- ✅ Valid: `guru@akubelajar.id`
- ❌ Invalid: `guru akubelajar.id` → "Format email tidak valid"

### password

| Rule | Detail |
|:---|:---|
| Min | 8 char |
| Max | 72 char |
| Wajib | 1 uppercase, 1 lowercase, 1 angka, 1 simbol |
| Tidak boleh | Sama dengan 3 password terakhir, mengandung nama/email/NISN |
| Cek | Terhadap 10.000 common passwords |

```typescript
const passwordSchema = z.string().min(8).max(72)
  .regex(/[A-Z]/, 'Harus ada huruf besar')
  .regex(/[a-z]/, 'Harus ada huruf kecil')
  .regex(/[0-9]/, 'Harus ada angka')
  .regex(/[^A-Za-z0-9]/, 'Harus ada simbol');
```
```go
Password string `json:"password" validate:"required,min=8,max=72"`
```

### name

| Rule | Detail |
|:---|:---|
| Min | 2 char |
| Max | 100 char |
| Allowed | Huruf, spasi, titik, koma |
| Tidak boleh | Angka, karakter khusus lain |

```typescript
const nameSchema = z.string().min(2).max(100).regex(/^[a-zA-Z\s.,]+$/);
```

### nisn

| Rule | Detail |
|:---|:---|
| Format | Tepat 10 digit angka |
| Tidak boleh | `0000000000` |
| Unique | Per sekolah |

```typescript
const nisnSchema = z.string().regex(/^[0-9]{10}$/).refine(v => v !== '0000000000');
```
```go
NISN string `json:"nisn" validate:"len=10,numeric"`
```

### phone_wa

| Rule | Detail |
|:---|:---|
| Format | E.164: `+62xxxxxxxxxx` |
| Min digits | 10 |
| Max digits | 15 |

```typescript
const phoneSchema = z.string().regex(/^\+62\d{8,13}$/);
```

### birth_date

| Rule | Detail |
|:---|:---|
| Format | `YYYY-MM-DD` |
| Tidak boleh | Di masa depan |
| Tidak boleh | > 100 tahun lalu |

---

## Field Akademik

| Field | Min | Max | Format | Contoh Valid |
|:---|:---|:---|:---|:---|
| class_name | 1 | 20 | String | `8A`, `XII-IPA-1` |
| subject_name | 1 | 100 | String | `Matematika` |
| school_code | 6 | 6 | Uppercase alphanumeric | `AB0001` |
| assignment_title | 5 | 200 | String | `Laporan Praktikum Biologi` |
| assignment_description | 0 | 10000 | HTML sanitized (bluemonday) | — |
| grade_value | 0 | 100 | Integer | `85` |
| quiz_title | 5 | 200 | String | `Kuis Bab 3: Struktur Sel` |
| question_text | 10 | 2000 | String | `Organel yang berfungsi...` |
| answer_option | 1 | 500 | String | `Mitokondria` |

---

## File Upload

| Kategori | Max Size | Allowed Types | Dimensi |
|:---|:---|:---|:---|
| Avatar | 2 MB | JPG, PNG, WebP | 100×100 – 2000×2000 px |
| Assignment file | 20 MB/file, 3 file max | PDF, DOCX, PPTX, XLSX, JPG, PNG, ZIP | — |
| Bukti izin/sakit | 5 MB | PDF, JPG, PNG | — |

**Security:**
- Validate magic bytes, bukan hanya extension
- Strip EXIF metadata dari gambar
- Scan filename: **tidak boleh** mengandung `../` atau absolute path
- Rename ke `{uuid_v7}.{ext}`

---

## Query Parameters

| Param | Type | Min | Max | Default |
|:---|:---|:---|:---|:---|
| `per_page` / `limit` | Integer | 1 | 100 | 20 |
| `page` | Integer | 1 | — | 1 |
| `date` | String | — | — | — (ISO 8601) |
| `sort` | String | — | — | Whitelist only |
| `search` | String | 0 | 100 | — |

```go
// Go — sort whitelist
var allowedSorts = map[string]bool{
    "created_at": true, "name": true, "email": true, "score": true,
}
```

---

*Terakhir diperbarui: 21 Maret 2026*
