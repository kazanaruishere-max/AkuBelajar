# 👥 User Personas — AkuBelajar

> Siapa pengguna sebenarnya? Guru daerah vs admin kota punya kebutuhan berbeda.

---

## Persona 1: Pak Budi — Guru Matematika (Daerah)

| Detail | Nilai |
|:---|:---|
| **Usia** | 45 tahun |
| **Lokasi** | Kabupaten Garut, Jawa Barat |
| **Sekolah** | SMP Negeri, 600 siswa |
| **Device** | HP Android mid-range (Redmi Note 12), laptop sekolah (shared) |
| **Internet** | 4G (sering tidak stabil, kadang 3G) |
| **Tech savvy** | ⭐⭐ (bisa WA, Google, YouTube, belum pernah pakai LMS) |

### Pain Points
- Menulis soal ujian manual butuh 2-3 jam per kuis
- Nilai dicatat di buku → sering salah hitung
- Absensi pakai kertas → rekap akhir semester menyiksa
- Takut salah klik, merusak data

### Kebutuhan
- ✅ **UI sederhana** — tombol besar, label jelas, sedikit input
- ✅ **AI bantu bikin soal** — hemat 2 jam per kuis
- ✅ **Auto hitung nilai** — tidak perlu kalkulator
- ✅ **Offline-friendly** — halaman tetap bisa dibuka saat sinyal lemah
- ✅ **Bahasa Indonesia** — tidak ada istilah Inggris yang membingungkan

### Skenario Penggunaan
> Pak Budi buka AkuBelajar di HP sambil di kelas. Input absensi (tap-tap cepat), lalu bikin kuis AI untuk Bab 3. Sore hari, cek submission tugas di laptop sekolah.

---

## Persona 2: Bu Ani — SuperAdmin (Kota)

| Detail | Nilai |
|:---|:---|
| **Usia** | 35 tahun |
| **Lokasi** | Jakarta Selatan |
| **Sekolah** | SMA Swasta, 400 siswa, 30 guru |
| **Device** | Laptop Windows (sekolah), iPhone 15 (pribadi) |
| **Internet** | WiFi sekolah (stabil), 4G |
| **Tech savvy** | ⭐⭐⭐⭐ (pakai Google Workspace, Canva, pernah pakai Moodle) |

### Pain Points
- Daftar 400 siswa + 30 guru satu per satu = butuh berminggu-minggu
- Koordinasi dengan 30 guru soal nilai akhir kacau
- Sekolah sebelumnya pakai Google Sheets, data tersebar di mana-mana
- Sering diminta rapor mendadak oleh yayasan

### Kebutuhan
- ✅ **Bulk import** (Excel → semua siswa sekaligus)
- ✅ **Dashboard overview** — berapa guru sudah input nilai, berapa belum
- ✅ **Role management** — assign wali kelas, ketua kelas
- ✅ **Report card generation** — 1 klik, semua rapor jadi PDF

### Skenario Penggunaan
> Bu Ani buka dashboard di laptop. Upload Excel 400 siswa, assign ke 12 kelas. Setiap akhir semester, cek dashboard: 5 guru belum finalisasi nilai → kirim reminder. Generate rapor batch → download ZIP.

---

## Persona 3: Rina — Siswi Kelas 9 (Kota)

| Detail | Nilai |
|:---|:---|
| **Usia** | 15 tahun |
| **Lokasi** | Surabaya |
| **Sekolah** | SMP Negeri, kelas 9A |
| **Device** | HP Android (Oppo A78), tidak punya laptop |
| **Internet** | 4G (paket data terbatas, sering habis) |
| **Tech savvy** | ⭐⭐⭐⭐⭐ (Instagram, TikTok, game, cepat belajar app baru) |

### Pain Points
- Tidak tahu ada tugas baru sampai teman kasih tahu di WA grup
- Lupa deadline → nilai 0
- Takut CBT karena kalau HP mati / sinyal hilang, jawaban hilang
- Bosan dengan soal monoton (pilihan ganda terus)

### Kebutuhan
- ✅ **Notifikasi tugas** — push notification langsung ke HP
- ✅ **Mobile-first** — semua fitur harus nyaman di HP
- ✅ **CBT reliable** — auto-save jawaban, reconnect otomatis
- ✅ **Light mode/dark mode** — belajar malam pakai dark mode
- ✅ **Hemat data** — halaman ringan, gambar di-compress

### Skenario Penggunaan
> Rina buka AkuBelajar di HP jam 8 malam. Ada notif: "Tugas Bahasa Indonesia deadline besok!" Submit file foto tugas dari kamera. Besoknya, ikut kuis IPA 20 menit di sekolah — sinyal agak lemah tapi jawaban auto-save.

---

## Persona 4: Dimas — Ketua Kelas 8A (Daerah)

| Detail | Nilai |
|:---|:---|
| **Usia** | 14 tahun |
| **Lokasi** | Magelang, Jawa Tengah |
| **Sekolah** | SMP Negeri, kelas 8A, 32 siswa |
| **Device** | HP Android bekas kakak (Samsung A22) |
| **Tech savvy** | ⭐⭐⭐ (bisa WA, YouTube, game mobile) |

### Kebutuhan
- ✅ **Input absensi cepat** — checklist 32 nama dalam < 3 menit
- ✅ **Simple UI** — tidak banyak menu, langsung ke absensi
- ✅ **Tanggung jawab jelas** — draft saja, guru yang approve

---

## Implikasi Desain

| Insight | Desain Decision |
|:---|:---|
| 60%+ user pakai HP Android | **Mobile-first** responsive design |
| Internet tidak selalu stabil | **Auto-save**, reconnect, minimal data transfer |
| Banyak user belum tech-savvy | **Label jelas**, tombol besar, onboarding wizard |
| Data terbatas | **Image compression**, lazy loading, skeleton UI |
| Guru takut salah | **Confirmation dialog** untuk aksi destructive |

---

*Terakhir diperbarui: 21 Maret 2026*
