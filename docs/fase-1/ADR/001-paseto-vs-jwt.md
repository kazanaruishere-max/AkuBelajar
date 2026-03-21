# ADR-001: Paseto vs JWT untuk Token Autentikasi

| Metadata | Nilai |
|:---|:---|
| **Status** | ✅ Accepted |
| **Tanggal** | 2026-02-15 |
| **Pembuat** | Kazanaru |
| **Reviewer** | Tim Security |

---

## Konteks

AkuBelajar membutuhkan mekanisme token untuk autentikasi stateless antar service. Dua kandidat utama: **JWT (JSON Web Token)** dan **Paseto (Platform-Agnostic Security Tokens)**.

---

## Keputusan

**Menggunakan Paseto v4 sebagai token utama**, dengan JWT sebagai fallback untuk integrasi third-party.

---

## Alasan

| Aspek | JWT | Paseto |
|:---|:---|:---|
| **Algoritma** | Bisa memilih alg (termasuk `none`) | Algoritma di-lock per versi — tidak bisa mis-konfigurasi |
| **Kerentanan `alg:none`** | Rentan jika tidak di-hardcode | Tidak mungkin terjadi by design |
| **Key Confusion Attack** | Memungkinkan jika `alg` tidak divalidasi | Tidak mungkin — satu versi = satu algoritma |
| **Parsing Complexity** | Header + Payload + Sig, banyak edge case | Simpler structure, less attack surface |
| **Ecosystem** | Sangat luas (semua bahasa) | Cukup matang (Go, Python, JS, dll.) |

### Mengapa tetap menyimpan JWT?

- Beberapa third-party API (OAuth providers) mengembalikan JWT
- Backward compatibility untuk integrasi eksternal
- JWT digunakan **hanya** untuk menerima token dari pihak luar, **bukan** untuk internal issuance

---

## Konsekuensi

- ✅ Risiko mis-konfigurasi kriptografi berkurang drastis
- ✅ Attack surface lebih kecil
- ⚠️ Perlu maintain dua library token (paseto + jwt)
- ⚠️ Developer baru perlu mempelajari Paseto (kurva belajar minimal)

---

## Referensi

- [Paseto.io — Official Docs](https://paseto.io/)
- [Critical JWT Vulnerabilities](https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/)

---

# Template ADR Baru

Untuk membuat ADR baru, buat file dengan format `NNN-judul-kebab-case.md`:

```markdown
# ADR-NNN: Judul Keputusan

| Metadata | Nilai |
|:---|:---|
| **Status** | 🟡 Proposed / ✅ Accepted / ❌ Rejected / 🔄 Superseded by ADR-XXX |
| **Tanggal** | YYYY-MM-DD |
| **Pembuat** | Nama |
| **Reviewer** | Nama |

## Konteks
<!-- Mengapa keputusan ini perlu diambil? -->

## Keputusan
<!-- Apa yang diputuskan? -->

## Alasan
<!-- Mengapa opsi ini dipilih dibanding alternatif? -->

## Konsekuensi
<!-- Dampak positif dan negatif dari keputusan ini -->

## Referensi
<!-- Link ke sumber yang relevan -->
```

---

*Terakhir diperbarui: 21 Maret 2026*
