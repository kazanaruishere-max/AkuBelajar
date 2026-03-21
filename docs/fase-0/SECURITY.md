# 🔐 Kebijakan Keamanan — AkuBelajar

Keamanan data pendidikan adalah prioritas utama AkuBelajar. Dokumen ini mendefinisikan cara melaporkan celah keamanan secara **responsible disclosure** — agar researcher tidak langsung mempublikasikan exploit ke publik.

---

## Daftar Isi

1. [Versi yang Didukung](#versi-yang-didukung)
2. [Melaporkan Kerentanan](#melaporkan-kerentanan)
3. [Proses Respons](#proses-respons)
4. [Scope](#scope)
5. [Out of Scope](#out-of-scope)
6. [Safe Harbor](#safe-harbor)
7. [Penghargaan](#penghargaan)

---

## Versi yang Didukung

| Versi | Status Keamanan |
|:---|:---|
| 2.x.x (latest) | ✅ Aktif didukung — menerima security patch |
| 1.x.x (legacy) | ⚠️ End of Life — **tidak menerima patch lagi** |

> **Selalu gunakan versi terbaru** untuk mendapatkan perlindungan keamanan terkini.

---

## Melaporkan Kerentanan

> ⚠️ **JANGAN** membuat GitHub Issue publik untuk kerentanan keamanan.

### Kirim Laporan ke:

| Kanal | Alamat |
|:---|:---|
| **Email (Utama)** | `security@akubelajar.id` |
| **PGP Key** | [Unduh PGP Public Key](https://akubelajar.id/.well-known/pgp-key.asc) |
| **GitHub Security Advisory** | Via tab "Security" → "Report a vulnerability" |

### Informasi yang Diperlukan

Sertakan sebanyak mungkin informasi berikut:

```
1. Tipe kerentanan (XSS, SQL Injection, IDOR, RCE, dll.)
2. Lokasi (file/endpoint yang terdampak)
3. Langkah reproduksi (step-by-step)
4. Proof of Concept (PoC) — kode/request/screenshot
5. Dampak potensial (apa yang bisa dilakukan attacker)
6. Saran perbaikan (opsional, sangat dihargai)
```

### Enkripsi Laporan (Opsional tapi Dianjurkan)

Untuk laporan yang sangat sensitif, gunakan PGP key kami:

```
Fingerprint: XXXX XXXX XXXX XXXX XXXX  XXXX XXXX XXXX XXXX XXXX
```

---

## Proses Respons

```
Hari 0    : Laporan diterima → acknowledgment otomatis
Hari 1-2  : Tim keamanan melakukan triase & konfirmasi
Hari 3-5  : Penilaian dampak & severity (CVSS scoring)
Hari 5-14 : Pengembangan patch
Hari 14   : Patch dirilis + advisory dipublikasikan
Hari 30   : Reporter boleh melakukan public disclosure
```

### Severity Levels (CVSS v3.1)

| Level | CVSS Score | SLA Patch | Contoh |
|:---|:---|:---|:---|
| **Critical** | 9.0 - 10.0 | ≤ 24 jam | RCE, full data breach |
| **High** | 7.0 - 8.9 | ≤ 3 hari | SQL Injection, auth bypass |
| **Medium** | 4.0 - 6.9 | ≤ 7 hari | Stored XSS, IDOR terbatas |
| **Low** | 0.1 - 3.9 | ≤ 14 hari | Info disclosure minor |

---

## Scope

Berikut area yang termasuk dalam program responsible disclosure:

| Area | Domain/Target |
|:---|:---|
| Web Application | `app.akubelajar.id` |
| API Backend | `api.akubelajar.id/v1/*` |
| Autentikasi | JWT/Paseto, OAuth, login/register flows |
| Database | RLS policies, SQL injection vectors |
| File Upload | Validasi tipe, ukuran, dan path traversal |
| AI Service | Prompt injection, output manipulation |

---

## Out of Scope

Hal berikut **tidak termasuk** dalam scope:

- ❌ Serangan DDoS/DoS terhadap infrastruktur production
- ❌ Social engineering terhadap karyawan/pengguna
- ❌ Physical security testing
- ❌ Serangan terhadap layanan third-party (Gemini API, Cloudflare, dll.)
- ❌ Automated scanning tanpa koordinasi terlebih dahulu
- ❌ Akses ke data pengguna asli (gunakan akun test yang disediakan)

---

## Safe Harbor

AkuBelajar berkomitmen untuk:

1. **Tidak melakukan tuntutan hukum** terhadap researcher yang mengikuti panduan ini
2. **Tidak melaporkan ke penegak hukum** selama laporan dilakukan dengan itikad baik
3. **Bekerja sama** dengan reporter untuk memahami dan memperbaiki kerentanan
4. **Mengakui kontribusi** reporter (jika diizinkan) dalam security advisory

### Syarat Safe Harbor

- Tidak mengeksploitasi kerentanan lebih dari yang diperlukan untuk PoC
- Tidak mengakses, memodifikasi, atau menghapus data pengguna lain
- Tidak mengganggu ketersediaan layanan
- Melaporkan secara langsung ke tim keamanan, bukan ke publik

---

## Penghargaan

Reporter yang membantu meningkatkan keamanan AkuBelajar akan mendapatkan:

- 📝 Pencantuman nama di **Security Hall of Fame** (jika diizinkan)
- 🏆 Sertifikat digital penghargaan
- 🎁 Merchandise AkuBelajar (untuk temuan Critical/High)

---

## Referensi

- [OWASP Vulnerability Disclosure Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Vulnerability_Disclosure_Cheat_Sheet.html)
- [FIRST CVSS Calculator](https://www.first.org/cvss/calculator/3.1)
- [Threat Model AkuBelajar](../fase-3/THREAT_MODEL.md)

---

*Terakhir diperbarui: 21 Maret 2026*
