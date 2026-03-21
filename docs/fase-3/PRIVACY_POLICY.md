# 🔏 Kebijakan Privasi — AkuBelajar

> Kepatuhan terhadap **UU PDP (Undang-Undang Pelindungan Data Pribadi)** No. 27 Tahun 2022.

---

## 1. Pengendali Data

| Informasi | Detail |
|:---|:---|
| **Nama** | AkuBelajar (PT. [Nama Perusahaan]) |
| **Alamat** | [Alamat Kantor] |
| **Email DPO** | `privacy@akubelajar.id` |
| **Telepon** | +62-xxx-xxx-xxxx |

---

## 2. Data yang Dikumpulkan

### Data Pribadi Umum

| Data | Tujuan | Dasar Hukum |
|:---|:---|:---|
| Nama lengkap | Identifikasi pengguna | Pelaksanaan perjanjian |
| Email | Autentikasi, notifikasi | Pelaksanaan perjanjian |
| Nomor telepon | Notifikasi WhatsApp | Persetujuan |
| NISN/NIP | Identifikasi akademik | Kewajiban hukum (Kemendikbud) |
| Alamat IP | Keamanan, audit log | Kepentingan sah |

### Data Pribadi Spesifik (Sensitif)

| Data | Tujuan | Perlindungan Khusus |
|:---|:---|:---|
| Nilai akademik | Fungsi inti platform | RLS, enkripsi, akses terbatas |
| Rekaman kehadiran | Fungsi inti platform | Audit log immutable |
| Hasil ujian | Fungsi inti platform | Hash jawaban, isolasi data |

---

## 3. Hak Subjek Data (Pasal 5-13 UU PDP)

| Hak | Implementasi |
|:---|:---|
| **Hak untuk tahu** | Privacy policy ini + notifikasi saat data dikumpulkan |
| **Hak akses** | Data export via dashboard (JSON/CSV) |
| **Hak koreksi** | Edit profil langsung di platform |
| **Hak hapus** | Request via email DPO, proses ≤ 3x24 jam |
| **Hak tarik persetujuan** | Fitur opt-out notifikasi, delete account |
| **Hak keberatan** | Saluran pengaduan ke DPO |
| **Hak portabilitas** | Export data dalam format standar (JSON) |

---

## 4. Penyimpanan & Retensi Data

| Tipe Data | Retensi | Setelah Retensi |
|:---|:---|:---|
| Data akun aktif | Selama akun aktif | — |
| Data akun dihapus | 30 hari (grace period) | Hard delete + anonymize |
| Audit log | 5 tahun | Archive terenkripsi |
| Backup | 90 hari | Auto-expire |
| Log server | 90 hari | Auto-rotate & delete |

---

## 5. Pihak Ketiga yang Menerima Data

| Pihak Ketiga | Data yang Dibagi | Tujuan |
|:---|:---|:---|
| Google (Gemini AI) | Topik soal (bukan data siswa) | Generasi soal AI |
| Cloudflare | IP address, request headers | CDN & WAF |
| SMTP Provider | Email address | Pengiriman notifikasi |
| Kemendikbud | NISN, data akademik | Kewajiban pelaporan |

---

## 6. Keamanan Data

Langkah teknis yang diterapkan:

- ✅ Enkripsi at-rest (AES-256) dan in-transit (TLS 1.3)
- ✅ Row-Level Security (RLS) di database
- ✅ Password hashing Argon2id
- ✅ Audit log immutable
- ✅ Akses terbatas berbasis RBAC
- ✅ Backup terenkripsi
- ✅ Vulnerability scanning berkala

---

## 7. Insiden Pelanggaran Data

Sesuai Pasal 46 UU PDP:

1. Notifikasi ke subjek data **≤ 3x24 jam** setelah terdeteksi
2. Notifikasi ke otoritas (BSSN/Kominfo) **≤ 3x24 jam**
3. Isi notifikasi: jenis data yang bocor, kronologi, langkah yang diambil

---

## 8. Kontak

Untuk pertanyaan atau permintaan terkait data pribadi:

| Kanal | Alamat |
|:---|:---|
| Email DPO | `privacy@akubelajar.id` |
| Formulir | `app.akubelajar.id/privacy-request` |

---

*Terakhir diperbarui: 21 Maret 2026*
