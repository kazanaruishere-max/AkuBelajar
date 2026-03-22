# 🖼️ Wireframes & Layout Reference — AkuBelajar

> Referensi visual layout per halaman. AI agent gunakan ini saat menulis UI code.

---

## 1. Layout Master

### Desktop (≥ 1024px)

```
┌──────────────────────────────────────────────────┐
│  TopBar: [Breadcrumb]            [🔔] [Avatar ▼] │
├──────────┬───────────────────────────────────────┤
│          │                                       │
│  Sidebar │         Main Content Area             │
│  240px   │         (fluid width)                 │
│          │                                       │
│  [Logo]  │  ┌─────────────────────────────┐      │
│  ──────  │  │  PageHeader                 │      │
│  📊 Dash │  │  Title + Description + CTAs │      │
│  📝 Tugas│  └─────────────────────────────┘      │
│  📋 Kuis │                                       │
│  📅 Absen│  ┌──────┐ ┌──────┐ ┌──────┐          │
│  📊 Nilai│  │StatA │ │StatB │ │StatC │          │
│  👥 Users│  └──────┘ └──────┘ └──────┘          │
│  ⚙️ Sett │                                       │
│          │  ┌─────────────────────────────┐      │
│          │  │  DataTable / Content        │      │
│          │  │  ...                        │      │
│          │  └─────────────────────────────┘      │
└──────────┴───────────────────────────────────────┘
```

### Mobile (< 768px)

```
┌─────────────────────────┐
│ TopBar: [☰] Title [🔔]  │
├─────────────────────────┤
│                         │
│   Main Content Area     │
│   (full width)          │
│                         │
│   ┌───────────────┐     │
│   │  StatA        │     │
│   └───────────────┘     │
│   ┌───────────────┐     │
│   │  StatB        │     │
│   └───────────────┘     │
│                         │
│   ┌───────────────┐     │
│   │  Content      │     │
│   │  (scrollable) │     │
│   └───────────────┘     │
│                         │
├─────────────────────────┤
│ [📊] [📝] [📋] [📅] [👤]│
│      Bottom Tab Bar      │
└─────────────────────────┘
```

---

## 2. Halaman per Role

### Dashboard — SuperAdmin

```
┌─────────────────────────────────────┐
│  📊 Dashboard Overview              │
├─────────────────────────────────────┤
│                                     │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐
│  │Total │ │Total │ │Total │ │Active│
│  │Siswa │ │Guru  │ │Kelas │ │Users │
│  │ 342  │ │  28  │ │  12  │ │  45  │
│  └──────┘ └──────┘ └──────┘ └──────┘
│                                     │
│  ┌─────────────┐ ┌────────────────┐ │
│  │ Guru belum  │ │ Recent Audit   │ │
│  │ input nilai │ │ Log            │ │
│  │ ──────────  │ │ ──────────     │ │
│  │ Bu Ani: IPA │ │ User created.. │ │
│  │ Pak Budi: M │ │ Grade locked.. │ │
│  │ ...         │ │ Login from..   │ │
│  └─────────────┘ └────────────────┘ │
└─────────────────────────────────────┘
```

### Dashboard — Guru

```
┌─────────────────────────────────────┐
│  📊 Dashboard Guru                   │
├─────────────────────────────────────┤
│                                     │
│  ┌──────┐ ┌──────┐ ┌──────┐        │
│  │Tugas │ │Kuis  │ │Belum │        │
│  │Aktif │ │Aktif │ │Dinilai│       │
│  │  5   │ │  2   │ │  12  │        │
│  └──────┘ └──────┘ └──────┘        │
│                                     │
│  Kelas Saya:                        │
│  ┌───────────────────────────┐      │
│  │ 8A — Matematika    [Buka]│      │
│  │ 8B — Matematika    [Buka]│      │
│  │ 9A — Matematika    [Buka]│      │
│  └───────────────────────────┘      │
│                                     │
│  [+ Buat Tugas]  [🤖 Buat Kuis AI] │
└─────────────────────────────────────┘
```

### Dashboard — Siswa

```
┌─────────────────────────────────────┐
│  Halo, Rina! 👋                      │
├─────────────────────────────────────┤
│                                     │
│  ⚡ Yang Perlu Dikerjakan:           │
│  ┌───────────────────────────┐      │
│  │ 📝 Laporan IPA            │      │
│  │    Deadline: 2 hari lagi  │      │
│  │    [Kerjakan →]           │      │
│  ├───────────────────────────┤      │
│  │ 📋 Kuis Matematika Bab 3  │      │
│  │    30 menit, 10 soal      │      │
│  │    [Mulai Kuis →]         │      │
│  └───────────────────────────┘      │
│                                     │
│  📊 Nilai Terbaru:                   │
│  ┌───────────────────────────┐      │
│  │ Matematika: 85 (B)        │      │
│  │ B. Indonesia: 78 (C)      │      │
│  │ IPA: 90 (A) ⭐             │      │
│  └───────────────────────────┘      │
└─────────────────────────────────────┘
```

---

## 3. Flow Halaman: Buat Kuis AI

```
Step 1: Form                Step 2: Loading           Step 3: Review
┌─────────────────┐        ┌─────────────────┐       ┌─────────────────┐
│ 🤖 Generate AI   │        │                 │       │ 10 Soal Generated│
│                 │        │   ⏳ Generating   │       │ ─────────────── │
│ Mapel: [▼ MTK]  │        │   soal dengan   │       │ 1. Berapakah... │
│ Topik: [_______]│   →    │   AI...         │  →    │    ○ A. 42      │
│ Jumlah: [10   ] │        │                 │       │    ● B. 56 ✓    │
│ Level: [▼ Camp] │        │   ~15-30 detik  │       │    ○ C. 28      │
│                 │        │                 │       │    ○ D. 70      │
│ [Generate Soal] │        │   [Batal]       │       │  [✏️Edit] [🗑️]   │
└─────────────────┘        └─────────────────┘       │                 │
                                                     │ [Save Draft]    │
                                                     │ [Publish]       │
                                                     └─────────────────┘
```

---

## 4. Flow Halaman: CBT (Siswa)

```
Pre-Exam                    During Exam                Post-Exam
┌─────────────────┐        ┌─────────────────┐       ┌─────────────────┐
│ Kuis: Bab 3     │        │ ⏱️ 24:35          │       │ ✅ Kuis Selesai   │
│ ──────────────  │        │ Soal 3 dari 10  │       │                 │
│ Mapel: MTK      │        │ ─────────────── │       │ Skor: 80/100    │
│ Soal: 10        │   →    │                 │  →    │ Benar: 8/10     │
│ Waktu: 30 menit │        │ Jika x + 5 = 12│       │ Waktu: 22:15    │
│ ──────────────  │        │ maka x = ?      │       │                 │
│ ⚠️ Tab detection │        │                 │       │ [Lihat Review]  │
│    aktif        │        │ ○ A. 5          │       │ [Kembali]       │
│                 │        │ ○ B. 7 ←selected│       └─────────────────┘
│ [Mulai Kuis]    │        │ ○ C. 12         │
└─────────────────┘        │ ○ D. 17         │
                           │                 │
                           │ [← Prev] [Next→]│
                           │ ○○●○○○○○○○      │
                           └─────────────────┘
```

---

## 5. Flow Halaman: Input Absensi (Guru)

```
┌─────────────────────────────────────┐
│ 📅 Absensi — 8A — Matematika        │
│ Tanggal: 21 Maret 2026              │
├─────────────────────────────────────┤
│                                     │
│  #  Nama           Status           │
│  ── ─────────────  ──────           │
│  1  Ahmad Fauzi    [✅ H]            │
│  2  Budi Santoso   [✅ H]            │
│  3  Citra Dewi     [🟡 I]  ← tap    │
│  4  Dimas Pratama  [🔴 A]            │
│  5  Eka Putri      [✅ H]            │
│  ...                                │
│  32 Zara Amelia    [✅ H]            │
│                                     │
│  Summary: H:28  I:2  S:1  A:1      │
│                                     │
│  [Simpan] [Batal]                   │
└─────────────────────────────────────┘

Legend: H=Hadir, I=Izin, S=Sakit, A=Alfa
Tap status untuk toggle: H → I → S → A → H
```

---

## 6. Flow Halaman: Foto Tugas (Kamera)

```
Step 1: Pilih             Step 2: Kamera           Step 3: Preview
┌─────────────────┐        ┌─────────────────┐       ┌─────────────────┐
│ Submit Tugas    │        │ 📷 Foto Tugas    │       │ Preview         │
│                 │        │                 │       │                 │
│ [📁 Pilih File] │        │  ┌───────────┐  │       │  ┌───────────┐  │
│                 │   →    │  │ Viewfinder│  │  →    │  │ [Foto 1]  │  │
│ [📷 Foto Tugas] │        │  │ NO MIRROR │  │       │  │ Tugas.jpg │  │
│                 │        │  │ ✓ teks    │  │       │  └───────────┘  │
│  atau           │        │  │ terbaca  │  │       │                 │
│                 │        │  └───────────┘  │       │  [📷 Tambah Foto]│
│ [🖼️ Dari Galeri]│        │ [🔄] [⭕ Jepret] │       │  [Submit Tugas] │
└─────────────────┘        └─────────────────┘       └─────────────────┘

🔄 = toggle kamera depan/belakang
Kamera belakang = default (untuk foto dokumen)
Preview TIDAK di-mirror → teks terbaca normal
```

---

## 7. Flow Halaman: QR Scan Absensi

```
GURU (Layar/Proyektor):       SISWA (HP):
┌─────────────────┐               ┌─────────────────┐
│ 📱 QR Absensi  │               │ 📷 Scan QR     │
│ 8A - MTK       │               │                 │
│                 │               │  ┌───────────┐  │
│  ┌─────────┐   │               │  │ Kamera    │  │
│  │█████████│   │  ← scan ←    │  │ [■] aim  │  │
│  │█ QR    █│   │               │  │ NO MIRROR │  │
│  │█ CODE  █│   │               │  └───────────┘  │
│  │█████████│   │               │                 │
│  └─────────┘   │               │  ✅ Hadir!       │
│ Refresh: 00:28 │               │  09:01 WIB      │
│ Hadir: 28/32   │               │  Kelas 8A - MTK │
│ [Manual Input] │               └─────────────────┘
└─────────────────┘

QR berubah otomatis setiap 30 detik
Guru bisa fallback ke input manual kapan saja
```

---

## 8. Design Principles

| Principle | Implementasi |
|:---|:---|
| **Mobile-first** | Desain HP dulu, scale up ke desktop |
| **Tap-friendly** | Min touch target: 44×44px |
| **Forgiving** | Confirmation dialog untuk aksi destructive |
| **Scannable** | Status menggunakan warna + icon (colorblind-safe) |
| **Fast** | Skeleton loading, optimistic updates |
| **Accessible** | Focus states, screen reader labels, high contrast |

---

*Terakhir diperbarui: 22 Maret 2026*
