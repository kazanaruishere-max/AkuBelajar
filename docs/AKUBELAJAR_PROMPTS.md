# PROMPT BANK — AKUBELAJAR MISSING DOCS
## Kumpulan Prompt untuk Antigravity AI Agent
**Gunakan setiap prompt secara terpisah. Berikan AKUBELAJAR_JOURNAL.md sebagai konteks di setiap sesi.**

---

## INSTRUKSI PENGGUNAAN

Sebelum menjalankan prompt manapun, selalu awali sesi dengan kalimat ini:

```
Kamu adalah technical writer dan software architect senior untuk project AkuBelajar.
Berikut adalah dokumentasi utama project: [lampirkan AKUBELAJAR_JOURNAL.md]

Gunakan seluruh konteks dari dokumen tersebut sebagai referensi utama.
Jangan membuat asumsi di luar apa yang tertulis. Jika ada ambiguitas, tulis
sebagai [KEPUTUSAN DIPERLUKAN: ...] agar tim dapat mengisi sendiri.
```

---

# ════════════════════════════════════════
# KELOMPOK A — ALUR PENGGUNA (USER FLOWS)
# ════════════════════════════════════════

---

## PROMPT A-1
**File output:** `docs/flows/ACCOUNT_CREATION_FLOW.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama ACCOUNT_CREATION_FLOW.md untuk
project AkuBelajar.

Dokumen ini harus mendefinisikan SECARA EKSPLISIT bagaimana setiap role mendapatkan
akun di sistem, karena saat ini sama sekali tidak terdokumentasi.

Cakup SEMUA jalur berikut dengan flowchart mermaid dan penjelasan teknis:

1. SUPERADMIN
   - Dibuat saat instalasi awal (database seed)
   - Tidak bisa self-register
   - Credentials dikirim ke developer/owner sekolah

2. GURU
   - Jalur A: Bulk import Excel oleh SuperAdmin (format kolom wajib)
   - Jalur B: Invite token yang dikirim SuperAdmin via WA/Email
   - Validasi: email unik, NIP opsional
   - Akun langsung aktif atau perlu approval?

3. SISWA
   - Jalur A: Bulk import Excel awal tahun ajaran (format kolom wajib)
   - Jalur B: Invite token dari Guru/Admin untuk siswa pindahan
   - Jalur C: Self-register menggunakan kode kelas (school_code + class_code)
              dengan status PENDING sampai diapprove Guru/Admin
   - Validasi: NISN unik per sekolah, email opsional

4. KETUA KELAS
   - Bukan register baru — Guru melakukan role upgrade dari akun Siswa yang
     sudah ada
   - Maksimal 1 Ketua Kelas per kelas per tahun ajaran
   - Otomatis downgrade kembali ke Siswa saat tahun ajaran berakhir

5. ORANG TUA (roadmap Fase 2)
   - Mendapat invite link yang di-generate dari akun Siswa yang sudah ada
   - Satu akun Orang Tua bisa linked ke beberapa akun Siswa (kakak-adik)
   - Akses read-only terhadap data anak

Untuk setiap jalur sertakan:
- Flowchart mermaid (sequenceDiagram atau flowchart LR)
- Field data yang diperlukan saat pembuatan akun
- Siapa yang bisa memicu proses ini
- Status akun setelah dibuat (active / pending / suspended)
- Notifikasi apa yang dikirim dan ke mana
- Error scenario: duplikat email, NISN invalid, token expired
- Tabel schema untuk tabel invite_tokens yang dibutuhkan

Stack: Go 1.23+ backend, Next.js 15 + TypeScript frontend,
PostgreSQL 16 + UUID v7 primary key, Paseto v4 auth,
notifikasi via Fonnte WA API + Email.
```

---

## PROMPT A-2
**File output:** `docs/flows/FIRST_LOGIN_FLOW.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama FIRST_LOGIN_FLOW.md untuk
project AkuBelajar.

Dokumen ini mendefinisikan apa yang terjadi saat pengguna login PERTAMA KALI
setelah akun dibuat oleh admin/sistem.

Cakup semua hal berikut:

1. DETEKSI FIRST LOGIN
   - Kolom di tabel users: is_first_login BOOLEAN DEFAULT TRUE
   - Middleware Next.js yang redirect ke halaman /onboarding jika flag aktif
   - Flag di-set FALSE hanya setelah semua langkah onboarding selesai

2. ALUR FORCE PASSWORD CHANGE
   - Siswa/Guru menerima password sementara (generated random 12 karakter)
   - Saat login pertama, sistem WAJIB redirect ke halaman ganti password
   - Tidak bisa skip atau akses halaman lain sebelum selesai
   - Validasi password baru: minimal 8 karakter, 1 huruf besar, 1 angka, 1 simbol
   - Password lama (sementara) langsung di-invalidate setelah ganti
   - Sertakan flowchart mermaid

3. ALUR KELENGKAPAN PROFIL
   Setelah ganti password, arahkan ke wizard profil:
   - Step 1 (Semua role): Foto profil (opsional, max 2MB, JPG/PNG/WebP)
   - Step 2 (Siswa): NISN, tanggal lahir, nama orang tua, nomor WA orang tua
   - Step 2 (Guru): NIP (opsional), mata pelajaran yang diajarkan, nomor WA
   - Step 3 (Semua role): Preferensi notifikasi (WA / Email / In-app)
   - Siswa WAJIB memilih kelas aktif jika belum di-assign oleh admin
   - Sertakan flowchart wizard step-by-step

4. LANDING SETELAH ONBOARDING
   - Dashboard pertama kali menampilkan tour/tooltip singkat
   - Konten yang ditampilkan berbeda per role (role-based dashboard)
   - Sertakan mapping: role → halaman pertama yang dituju

5. EDGE CASES
   - Bagaimana jika user menutup browser di tengah onboarding?
   - Bagaimana jika foto profil gagal upload?
   - Bagaimana jika nomor WA orang tua sudah terdaftar di akun lain?

Stack: Next.js 15 App Router, Go backend, PostgreSQL, MinIO untuk file storage,
Argon2id untuk hashing password baru.
```

---

## PROMPT A-3
**File output:** `docs/flows/PASSWORD_RESET_FLOW.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama PASSWORD_RESET_FLOW.md untuk
project AkuBelajar.

Dokumen ini mendefinisikan SEMUA skenario reset password yang bisa terjadi.

Cakup hal-hal berikut:

1. SELF-SERVICE RESET (via halaman login)
   - Input: email ATAU nomor WA terdaftar
   - Sistem kirim OTP 6 digit
   - OTP TTL: 5 menit, maksimal 3x request per jam per IP (rate limit)
   - OTP dikirim via: WA (prioritas utama) atau Email (fallback)
   - Halaman input OTP + password baru
   - Setelah berhasil: semua session aktif di-revoke (logout paksa semua device)
   - Sertakan flowchart mermaid sequenceDiagram

2. RESET OLEH ADMIN (untuk siswa yang tidak punya email/WA)
   - SuperAdmin atau Guru bisa trigger reset untuk user di bawah mereka
   - Sistem generate password sementara baru (12 karakter random)
   - Password sementara ditampilkan SEKALI di layar admin (tidak dikirim ulang)
   - Admin wajib sampaikan secara manual ke siswa
   - Flag is_first_login di-set TRUE kembali agar siswa dipaksa ganti
   - Audit log mencatat: admin mana, user mana, timestamp

3. TABEL DATABASE YANG DIBUTUHKAN
   Definisikan tabel password_reset_tokens:
   - id UUID PK
   - user_id UUID FK
   - token_hash VARCHAR (hashed OTP, bukan plain text)
   - expires_at TIMESTAMPTZ
   - used_at TIMESTAMPTZ (nullable)
   - ip_address INET
   - created_at TIMESTAMPTZ

4. SECURITY CONSTRAINTS
   - OTP tidak boleh disimpan plain text di DB (hash dengan SHA-256)
   - Token expired harus di-cleanup via scheduled job
   - Jika email/WA tidak ditemukan: tampilkan pesan generik (jangan bocorkan
     apakah email terdaftar atau tidak — anti user enumeration)
   - Constant-time comparison untuk validasi OTP

5. EDGE CASES
   - User input OTP yang sudah dipakai
   - User request OTP berulang kali dalam waktu singkat
   - User ganti password sementara lalu lupa lagi

Stack: Go backend, PostgreSQL, Redis untuk rate limiting OTP,
Fonnte WA API, SMTP Email.
```

---

## PROMPT A-4
**File output:** `docs/flows/ATTENDANCE_FLOW.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama ATTENDANCE_FLOW.md untuk
project AkuBelajar.

Dokumen ini mendefinisikan sistem absensi dari awal hingga akhir.

Cakup semua aspek berikut:

1. KODE STATUS ABSENSI
   Definisikan enum dengan jelas:
   - H = Hadir
   - I = Izin (dengan keterangan wajib, disetujui Guru)
   - S = Sakit (dengan keterangan, bisa upload surat dokter)
   - A = Alfa/Tanpa Keterangan
   - T = Terlambat (tandai Hadir tapi catat waktu masuk)

2. ALUR INPUT ABSENSI OLEH GURU
   - Guru buka halaman absensi di PWA/browser
   - Pilih kelas + mata pelajaran + tanggal
   - Sistem tampilkan daftar siswa kelas tersebut
   - Guru tap/klik status tiap siswa (default: Hadir)
   - Konfirmasi submit — setelah submit tidak bisa diedit sendiri
   - Sertakan flowchart mermaid

3. ALUR INPUT ABSENSI OLEH KETUA KELAS
   - Akses terbatas: hanya bisa input H atau A
   - Tidak bisa set I atau S (harus Guru)
   - Data yang diinput Ketua Kelas berstatus DRAFT
   - Guru harus review dan konfirmasi sebelum final
   - Sertakan flowchart mermaid

4. ALUR EDIT ABSENSI (KOREKSI)
   - Hanya Guru yang bisa edit absensi finalnya sendiri
   - SuperAdmin bisa edit semua
   - Edit absensi setelah 24 jam membutuhkan alasan (field reason wajib diisi)
   - Semua edit tercatat di audit_log: siapa, kapan, dari status apa ke apa

5. ALUR PENGAJUAN IZIN/SAKIT OLEH SISWA
   - Siswa submit form izin/sakit sebelum atau sesudah tanggal (max 3 hari setelah)
   - Upload bukti opsional (surat dokter, foto)
   - Status: PENDING → Guru approve/reject → FINAL
   - Notifikasi ke Guru saat ada pengajuan baru
   - Sertakan flowchart approval

6. REKAP & STATISTIK
   - Rekap harian: total H/I/S/A per kelas
   - Rekap mingguan dan bulanan per siswa
   - Alert otomatis jika siswa alfa lebih dari 3 kali berturut-turut
   - Persentase kehadiran minimum: [KEPUTUSAN DIPERLUKAN: berapa %?]
   - Dampak kehadiran rendah terhadap nilai: [KEPUTUSAN DIPERLUKAN]

7. SCHEMA TABEL
   Definisikan tabel attendances secara lengkap dengan semua kolom,
   tipe data, constraint, dan index yang dibutuhkan.

Stack: Go backend, PostgreSQL RLS (siswa hanya lihat data sendiri),
Next.js PWA mobile-first, WebSocket untuk update real-time.
```

---

## PROMPT A-5
**File output:** `docs/flows/ASSIGNMENT_FLOW.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama ASSIGNMENT_FLOW.md untuk
project AkuBelajar.

Dokumen ini mendefinisikan siklus hidup tugas dari pembuatan hingga penilaian.

Cakup semua hal berikut:

1. ALUR GURU MEMBUAT TUGAS
   - Form: judul, deskripsi (rich text/markdown), deadline, bobot nilai (%)
   - Lampiran: file materi opsional (max 50MB, format: PDF/DOCX/PPTX/MP4)
   - Target: seluruh kelas ATAU beberapa siswa tertentu
   - Opsi: izinkan submit terlambat (Y/N) + penalti nilai jika terlambat (%)
   - Draft → Published → Closed lifecycle
   - Notifikasi otomatis ke siswa saat published
   - Sertakan flowchart mermaid

2. ALUR SISWA MENGERJAKAN & SUBMIT TUGAS
   - Siswa lihat tugas di dashboard (sorted by deadline terdekat)
   - Download materi jika ada (via signed URL MinIO, TTL 15 menit)
   - Upload submission: max 3 file, max 20MB per file
   - Format yang diterima: PDF, DOCX, PPTX, JPG, PNG, ZIP
   - File di-rename otomatis: {student_id}_{assignment_id}_{timestamp}.ext
   - Konfirmasi submit — setelah submit bisa revisi jika guru izinkan
   - Status submission: DRAFT → SUBMITTED → GRADED
   - Sertakan flowchart mermaid

3. ALUR SUBMIT TERLAMBAT
   - Jika deadline lewat dan izin terlambat diaktifkan:
     sistem terima submission tapi tandai LATE
   - Sistem otomatis hitung penalti: nilai × (1 - penalti%)
   - Jika izin terlambat tidak aktif: tombol submit di-disable, tampilkan pesan
   - Sertakan flowchart decision

4. ALUR GURU MENGOREKSI & MEMBERI NILAI
   - Guru lihat daftar submission per tugas (filter: belum dinilai / sudah dinilai)
   - Download file submission siswa (signed URL)
   - Input nilai (0–100) + komentar teks wajib minimal 10 karakter
   - Opsi: minta revisi (siswa bisa upload ulang sekali)
   - Setelah dinilai: notifikasi ke siswa
   - Nilai masuk ke buku nilai otomatis
   - Sertakan flowchart mermaid

5. EDGE CASES YANG HARUS DITANGANI
   - Siswa upload file yang terinfeksi malware (bagaimana sistem merespons?)
   - Guru menghapus tugas setelah ada yang submit
   - Siswa submit file 0 byte
   - Deadline berubah setelah ada yang submit
   - Guru edit nilai setelah rapor dikunci

6. SCHEMA TABEL LENGKAP
   Definisikan tabel assignments dan submissions dengan semua kolom,
   constraint, index, dan RLS policy.

Stack: Go backend, MinIO untuk file storage (private bucket),
PostgreSQL, Next.js, notifikasi multi-channel.
```

---

## PROMPT A-6
**File output:** `docs/flows/ROLE_MANAGEMENT_FLOW.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama ROLE_MANAGEMENT_FLOW.md untuk
project AkuBelajar.

Dokumen ini mendefinisikan bagaimana role pengguna dikelola selama siklus hidup akun.

Cakup hal-hal berikut:

1. HIERARKI ROLE & SIAPA BISA MENGUBAH SIAPA
   Buat tabel matrix yang jelas:
   - SuperAdmin bisa ubah role Guru dan Siswa
   - Guru bisa upgrade Siswa → Ketua Kelas (di kelasnya sendiri)
   - Guru bisa downgrade Ketua Kelas → Siswa
   - Tidak ada self-upgrade (tidak bisa ubah role sendiri)
   - Sertakan diagram mermaid

2. ALUR UPGRADE SISWA → KETUA KELAS
   - Guru pilih siswa di kelasnya
   - Konfirmasi: "Jadikan Budi Santoso sebagai Ketua Kelas 8A?"
   - Jika sudah ada Ketua Kelas aktif: sistem peringatkan, guru konfirmasi
     apakah akan mengganti (otomatis downgrade yang lama)
   - Batasan: maks 1 Ketua Kelas per kelas per tahun ajaran aktif
   - Notifikasi ke siswa yang di-upgrade
   - Catat di audit_log
   - Sertakan flowchart mermaid

3. ALUR DOWNGRADE KETUA KELAS → SISWA
   - Manual oleh Guru kapan saja
   - Otomatis saat tahun ajaran berakhir (scheduled job)
   - Notifikasi ke pengguna yang bersangkutan
   - Data yang pernah dibuat sebagai Ketua Kelas tetap tersimpan

4. ALUR NONAKTIFKAN AKUN (SUSPEND)
   - Hanya SuperAdmin yang bisa suspend akun
   - Akun ter-suspend tidak bisa login, session aktif di-revoke
   - Data historis (nilai, absensi) tetap tersimpan — tidak dihapus
   - Bisa di-reaktivasi kapan saja

5. ALUR HAPUS AKUN (SOFT DELETE)
   - Hanya SuperAdmin
   - Soft delete: kolom deleted_at di-set, data tidak benar-benar dihapus
   - Data historis akademik siswa tetap bisa diakses SuperAdmin untuk keperluan
     akreditasi
   - Hard delete tidak diizinkan (data akademik wajib disimpan min. 5 tahun)
   - Sertakan flowchart mermaid

6. AUTO ROLE LIFECYCLE (SCHEDULED JOBS)
   - Akhir tahun ajaran: Ketua Kelas → Siswa otomatis
   - Siswa lulus: akun di-archive (read-only), bukan dihapus
   - Definisikan: kapan scheduled job ini berjalan, siapa yang trigger,
     bagaimana jika gagal

Stack: Go backend, PostgreSQL audit_logs, Redis untuk revoke token,
Paseto token management.
```

---

## PROMPT A-7
**File output:** `docs/flows/CBT_DETAIL_FLOW.md`
**Prioritas:** PENTING

```
Buatkan dokumen markdown lengkap dengan nama CBT_DETAIL_FLOW.md untuk
project AkuBelajar.

Dokumen ini melengkapi sequence diagram CBT yang sudah ada di AKUBELAJAR_JOURNAL.md
dengan detail edge case dan skenario yang belum terdokumentasi.

Cakup hal-hal berikut:

1. PRE-EXAM CHECKLIST (sebelum siswa mulai)
   - Sistem cek: browser support fullscreen API?
   - Sistem cek: koneksi WebSocket berhasil?
   - Tampilkan aturan ujian dan minta konfirmasi siswa
   - Countdown timer sebelum soal muncul
   - Sertakan flowchart mermaid

2. SISTEM PENGACAKAN SOAL
   - Urutan soal diacak per siswa (seed berbeda tiap siswa)
   - Pilihan jawaban juga diacak (kecuali tipe "semua benar" / "tidak ada yang benar")
   - Pengacakan dilakukan di server, bukan client
   - Sertakan penjelasan algoritma shuffle

3. SKENARIO KONEKSI PUTUS DI TENGAH UJIAN
   - Jawaban yang sudah diisi di-autosave ke server setiap 30 detik
   - Jika WebSocket putus: tampilkan warning, timer tetap berjalan di server
   - Siswa bisa reconnect dan lanjutkan dari posisi terakhir
   - Jika tidak reconnect sampai waktu habis: submit otomatis dengan jawaban
     yang sudah tersimpan
   - Sertakan flowchart mermaid

4. DETEKSI KECURANGAN (ANTI-CHEAT DETAIL)
   - visibilitychange event: toleransi berapa kali sebelum dikunci?
   - Devtools detection: apa yang dilakukan sistem jika terdeteksi?
   - Multiple device detection: siswa coba login di 2 device sekaligus
   - Sertakan tabel: EVENT → AKSI SISTEM → NOTIFIKASI KE GURU

5. GURU RESET SESI UJIAN SISWA
   - Alur: Guru lihat daftar siswa yang session-nya terkunci
   - Guru klik "Reset Sesi" → sistem buka kembali akses siswa
   - Waktu ujian siswa: dilanjutkan dari sisa waktu (tidak di-reset ke awal)
   - Sertakan flowchart mermaid

6. HASIL & REVIEW PASCA UJIAN
   - Kapan siswa bisa lihat jawaban benar? (langsung / setelah semua submit / setting guru)
   - Format halaman review: soal + jawaban siswa + kunci + penjelasan AI
   - Guru bisa lihat distribusi jawaban per soal (analitik)
   - Sertakan schema data untuk quiz_sessions tabel

Stack: Go WebSocket (gorilla/websocket), Redis untuk session state,
Next.js CBT interface, Gemini AI untuk penjelasan jawaban.
```

---

## PROMPT A-8
**File output:** `docs/flows/NOTIFICATION_FLOW.md`
**Prioritas:** PENTING

```
Buatkan dokumen markdown lengkap dengan nama NOTIFICATION_FLOW.md untuk
project AkuBelajar.

Dokumen ini mendefinisikan sistem notifikasi end-to-end yang lengkap.

Cakup hal-hal berikut:

1. ARSITEKTUR QUEUE NOTIFIKASI
   - Sertakan diagram: event trigger → Redis queue → Go worker → channel
   - Prioritas queue: HIGH (security alert) → MEDIUM (deadline) → LOW (info)
   - Retry logic: berapa kali retry jika gagal kirim? Backoff strategy?
   - Dead letter queue untuk notifikasi yang gagal terus

2. PREFERENSI NOTIFIKASI PER USER
   - User bisa aktifkan/nonaktifkan tiap channel (WA / Email / In-app)
   - User bisa set quiet hours (misal: tidak terima notif jam 22:00–06:00)
   - Admin sekolah bisa set default preferensi untuk semua user baru
   - Schema tabel notification_preferences

3. KATALOG LENGKAP SEMUA EVENT NOTIFIKASI
   Buat tabel komprehensif dengan kolom:
   EVENT | TRIGGER | CHANNEL | PENERIMA | TEMPLATE | PRIORITAS
   Sertakan SEMUA event: tugas baru, deadline 24 jam, kuis publish, nilai keluar,
   rapor tersedia, akun baru, login device baru, sesi ujian dikunci, izin disetujui,
   peringatan kehadiran rendah, dll.

4. TEMPLATE PESAN PER CHANNEL
   Contoh template untuk setidaknya 5 event penting:
   - Format WA (plain text, max 1000 karakter)
   - Format Email (HTML sederhana)
   - Format In-app (judul max 50 char + body max 200 char)

5. UNSUBSCRIBE & COMPLIANCE
   - Cara user berhenti dari notifikasi WA
   - Cara user unsubscribe dari email
   - Compliance: tidak spam, tidak kirim di luar quiet hours

6. FAILURE HANDLING
   - Nomor WA tidak valid / tidak terdaftar di WhatsApp
   - Email bounce (hard bounce vs soft bounce)
   - Apa yang ditampilkan di in-app jika WA dan Email gagal semua

Stack: Go worker, Redis queue, Fonnte WA API, SMTP/Resend Email API,
WebSocket untuk in-app real-time notification.
```

---

# ════════════════════════════════════════════════════
# KELOMPOK B — KONTRAK API & SPESIFIKASI TEKNIS
# ════════════════════════════════════════════════════

---

## PROMPT B-1
**File output:** `docs/api/API_SPEC.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama API_SPEC.md untuk project AkuBelajar.

Dokumen ini adalah kontrak resmi antara Go backend dan Next.js frontend.
Tanpa dokumen ini, kedua sisi tidak bisa sinkron.

Ikuti konvensi berikut untuk SEMUA endpoint:
- Base URL: /api/v1/
- Auth header: Authorization: Bearer <paseto_token>
- Content-Type: application/json
- Response sukses: { "data": {...}, "meta": {...} }
- Response error: { "error": { "code": "string", "message": "string", "details": [...] } }

Dokumentasikan SEMUA endpoint berikut dengan format lengkap
(method, path, auth required, request body/params, response 200/400/401/403/404/422/500):

AUTH MODULE:
- POST /auth/login
- POST /auth/refresh
- POST /auth/logout
- POST /auth/password-reset/request
- POST /auth/password-reset/verify
- POST /auth/change-password

USER MODULE:
- GET /users/me
- PUT /users/me
- PUT /users/me/avatar
- GET /users/:id (admin only)
- GET /users (admin only, dengan pagination + filter by role/school)
- POST /users (admin: create single user)
- PUT /users/:id/role (admin/guru)
- PUT /users/:id/status (admin: activate/suspend)
- DELETE /users/:id (admin: soft delete)

INVITE TOKEN MODULE:
- POST /invite-tokens (admin/guru generate token)
- POST /invite-tokens/claim (user claim token untuk buat akun)
- GET /invite-tokens (admin: list semua token)
- DELETE /invite-tokens/:id (admin: revoke token)

ATTENDANCE MODULE:
- POST /attendances (guru/ketua kelas: input absensi)
- GET /attendances?class_id=&date= (guru: lihat absensi)
- PUT /attendances/:id (guru: edit absensi + reason)
- GET /attendances/my (siswa: lihat absensi sendiri)
- GET /attendances/summary?student_id=&month= (rekap)

ASSIGNMENT MODULE:
- POST /assignments (guru: buat tugas)
- GET /assignments?class_id= (list tugas)
- GET /assignments/:id
- PUT /assignments/:id (guru: edit tugas)
- DELETE /assignments/:id (guru: hapus tugas)
- POST /assignments/:id/submissions (siswa: submit)
- GET /assignments/:id/submissions (guru: lihat semua submission)
- PUT /assignments/:id/submissions/:sub_id/grade (guru: beri nilai)

QUIZ MODULE:
- POST /quizzes (guru: buat kuis)
- POST /quizzes/generate-ai (guru: generate via AI)
- GET /quizzes?class_id=
- GET /quizzes/:id
- PUT /quizzes/:id/publish
- POST /quizzes/:id/sessions (siswa: mulai ujian)
- POST /quizzes/:id/sessions/submit (siswa: submit jawaban)
- GET /quizzes/:id/results (guru: lihat hasil)

GRADE MODULE:
- GET /grades?student_id=&subject_id=&semester= (buku nilai)
- GET /grades/report-card/:student_id (data rapor)
- POST /grades/report-card/:student_id/generate (generate PDF)

NOTIFICATION MODULE:
- GET /notifications (in-app notif list)
- PUT /notifications/:id/read
- PUT /notifications/read-all
- GET /notifications/preferences
- PUT /notifications/preferences

Untuk setiap endpoint, sertakan contoh request body dan response JSON yang nyata.
Definisikan juga semua error code string yang mungkin muncul.
```

---

## PROMPT B-2
**File output:** `docs/api/WEBSOCKET_SPEC.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama WEBSOCKET_SPEC.md untuk
project AkuBelajar.

Saat ini hanya disebutkan "WebSocket heartbeat" di journal tanpa detail apapun.
Dokumen ini mendefinisikan kontrak lengkap WebSocket.

Cakup semua hal berikut:

1. ENDPOINT & AUTENTIKASI
   - URL: wss://api.akubelajar.id/ws
   - Auth: Paseto token dikirim sebagai query param ?token= atau header
     (jelaskan mana yang lebih aman dan mengapa)
   - Koneksi per user: maksimal berapa koneksi simultan?
   - Timeout idle connection: berapa menit sebelum disconnect?

2. FORMAT PESAN (MESSAGE ENVELOPE)
   Definisikan standar envelope JSON untuk SEMUA pesan:
   ```json
   {
     "type": "event_name",
     "payload": {},
     "timestamp": "2026-03-21T10:00:00Z",
     "id": "uuid"
   }
   ```

3. KATALOG SEMUA EVENT (CLIENT → SERVER)
   Dokumentasikan setiap event yang dikirim client:
   - quiz:heartbeat (setiap 5 detik selama ujian)
   - quiz:autosave (kirim jawaban sementara)
   - notification:mark_read

4. KATALOG SEMUA EVENT (SERVER → CLIENT)
   Dokumentasikan setiap event yang dikirim server:
   - quiz:time_update (sisa waktu)
   - quiz:force_submit (waktu habis)
   - quiz:session_locked (anti-cheat trigger)
   - notification:new (notifikasi baru masuk)
   - attendance:updated (absensi berubah real-time)
   - grade:published (nilai baru keluar)

5. RECONNECT STRATEGY (CLIENT SIDE)
   - Exponential backoff: 1s, 2s, 4s, 8s, max 30s
   - Maksimal berapa kali retry sebelum tampilkan error ke user?
   - Bagaimana handle pesan yang missed selama disconnect?
   - Sertakan pseudocode atau TypeScript snippet

6. ERROR HANDLING
   - Close codes yang digunakan dan artinya
   - Skenario: token expired saat WebSocket aktif
   - Skenario: server restart, semua koneksi putus

Stack: Go gorilla/websocket, Redis untuk broadcast antar instance (pub/sub),
Next.js client dengan native WebSocket API.
```

---

## PROMPT B-3
**File output:** `docs/api/ERROR_CODE_CATALOGUE.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama ERROR_CODE_CATALOGUE.md untuk
project AkuBelajar.

Dokumen ini adalah katalog SEMUA error yang mungkin terjadi di sistem,
dengan kode standar yang konsisten antara backend Go dan frontend Next.js.

Struktur error response yang harus diikuti:
```json
{
  "error": {
    "code": "AUTH_001",
    "message": "Pesan ramah pengguna dalam Bahasa Indonesia",
    "details": ["field-level error jika ada"],
    "request_id": "uuid untuk debugging"
  }
}
```

Definisikan error code untuk SEMUA kategori berikut:

AUTH (AUTH_001 dst):
- Invalid credentials, account locked, token expired, token invalid,
  insufficient permissions, account suspended, first login required,
  password change required, OTP invalid, OTP expired, OTP rate limited

USER (USER_001 dst):
- User not found, email already exists, NISN already exists,
  invalid role assignment, cannot delete own account, school not found

VALIDATION (VAL_001 dst):
- Field required, invalid format (email, NISN, phone), value too long,
  value too short, invalid file type, file too large, invalid date range

ASSIGNMENT (ASSIGN_001 dst):
- Assignment not found, deadline passed, submission not allowed,
  already submitted, grading period closed, file upload failed

QUIZ (QUIZ_001 dst):
- Quiz not found, session not found, session expired, session locked,
  already submitted, answer not valid, ai generation failed, time up

ATTENDANCE (ATT_001 dst):
- Already recorded today, invalid status, cannot edit after 24h (without reason),
  class not found, date in future not allowed

SYSTEM (SYS_001 dst):
- Internal server error, service unavailable, database error,
  external API error (Gemini, Fonnte), rate limit exceeded, maintenance mode

Untuk setiap error code sertakan:
- Kode string (misal: AUTH_001)
- HTTP status code yang sesuai
- Pesan default Bahasa Indonesia untuk user
- Kapan error ini muncul
- Apa yang harus dilakukan developer/user untuk resolve

Stack: Go backend (gin framework), TypeScript frontend (type-safe error handling).
```

---

# ════════════════════════════════════════════════
# KELOMPOK C — DATABASE & LOGIKA BISNIS
# ════════════════════════════════════════════════

---

## PROMPT C-1
**File output:** `docs/database/DATABASE_SCHEMA_FULL.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama DATABASE_SCHEMA_FULL.md untuk
project AkuBelajar.

Dokumen yang ada saat ini hanya mendefinisikan 3 tabel (users, quizzes, quiz_questions).
Dokumen ini mendefinisikan SEMUA tabel secara lengkap.

Prinsip yang sudah ditetapkan di project:
- Primary key: UUID v7 di semua tabel
- Soft delete: kolom deleted_at TIMESTAMPTZ
- Timestamps: created_at, updated_at di semua tabel
- Security: school_id di semua tabel yang berkaitan dengan data sekolah
- Password: Argon2id hash
- Kunci jawaban kuis: Argon2id hash (tidak boleh plain text)

Definisikan DDL SQL lengkap + penjelasan untuk SEMUA tabel berikut:

CORE TABLES:
- schools (id, name, code, address, logo_url, theme_color, is_active)
- users (sudah ada, lengkapi)
- user_profiles (id, user_id FK, nisn, nip, birth_date, phone_wa,
  parent_name, parent_phone, photo_url, bio)

ACADEMIC TABLES:
- academic_years (id, school_id, name, start_date, end_date, is_active)
- classes (id, school_id, academic_year_id, name, grade_level, homeroom_teacher_id)
- subjects (id, school_id, name, code, description)
- class_subjects (id, class_id, subject_id, teacher_id, schedule_json)
- student_classes (id, student_id, class_id, academic_year_id, joined_at)

ATTENDANCE TABLES:
- attendances (id, student_id, class_id, subject_id, date, status ENUM,
  reason, proof_url, recorded_by, is_draft, approved_by, approved_at)

ASSIGNMENT TABLES:
- assignments (id, class_id, subject_id, teacher_id, title, description,
  deadline_at, allow_late, late_penalty_pct, max_file_count, max_file_size_mb,
  allowed_extensions, status ENUM)
- assignment_attachments (id, assignment_id, file_url, file_name, file_size)
- submissions (id, assignment_id, student_id, submitted_at, is_late,
  status ENUM, grade, feedback, graded_by, graded_at)
- submission_files (id, submission_id, file_url, file_name, file_size)

QUIZ TABLES:
- quizzes (sudah ada, lengkapi dengan: randomize_questions, randomize_options,
  max_attempts, allow_review, review_mode ENUM)
- quiz_questions (sudah ada, lengkapi)
- quiz_sessions (id, quiz_id, student_id, started_at, submitted_at,
  expires_at, status ENUM, score, ip_address, cheat_count)
- quiz_answers (id, session_id, question_id, answer_index, is_correct,
  answered_at)

GRADE TABLES:
- grades (id, student_id, subject_id, academic_year_id, semester,
  task_score, quiz_score, final_score, grade_letter, is_locked, locked_at)
- report_cards (id, student_id, academic_year_id, semester,
  pdf_url, generated_at, digital_signature_hash, is_published)

NOTIFICATION & SYSTEM TABLES:
- notifications (id, user_id, type, title, body, is_read, read_at,
  related_entity_type, related_entity_id)
- notification_preferences (id, user_id, email_enabled, wa_enabled,
  in_app_enabled, quiet_start, quiet_end)
- invite_tokens (id, school_id, created_by, token_hash, role, class_id,
  max_uses, uses_count, expires_at, used_at)
- password_reset_tokens (id, user_id, token_hash, expires_at, used_at, ip_address)
- audit_logs (id, actor_id, action, entity_type, entity_id, old_value JSONB,
  new_value JSONB, ip_address, user_agent, created_at)

Untuk setiap tabel sertakan:
1. DDL SQL CREATE TABLE lengkap
2. Penjelasan tujuan setiap kolom
3. Index yang direkomendasikan
4. RLS policy jika relevan
5. Constraint dan business rule di tingkat DB
```

---

## PROMPT C-2
**File output:** `docs/database/ERD.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama ERD.md untuk project AkuBelajar.

Dokumen ini berisi Entity Relationship Diagram yang menggambarkan relasi
antar semua tabel di database AkuBelajar.

Gunakan format Mermaid erDiagram untuk semua relasi.

Buat ERD dalam 4 bagian terpisah agar tidak terlalu padat:

BAGIAN 1 — Core & User Management:
Tampilkan relasi antara: schools, users, user_profiles,
academic_years, classes, subjects, class_subjects, student_classes

BAGIAN 2 — Attendance & Assignment:
Tampilkan relasi antara: attendances, assignments, assignment_attachments,
submissions, submission_files

BAGIAN 3 — Quiz & Grades:
Tampilkan relasi antara: quizzes, quiz_questions, quiz_sessions,
quiz_answers, grades, report_cards

BAGIAN 4 — Notifications & System:
Tampilkan relasi antara: notifications, notification_preferences,
invite_tokens, password_reset_tokens, audit_logs

Untuk setiap relasi, gunakan notasi crow's foot yang benar:
- ||--|| (one to one)
- ||--o{ (one to many)
- }o--o{ (many to many)

Di bawah setiap diagram, tambahkan penjelasan singkat tentang
relasi-relasi yang tidak obvious (yang butuh konteks bisnis untuk dipahami).

Sertakan juga tabel ringkasan: nama tabel | jumlah kolom utama | relasi ke tabel mana.

Stack: PostgreSQL 16, UUID v7 primary keys.
```

---

## PROMPT C-3
**File output:** `docs/business/BUSINESS_LOGIC.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama BUSINESS_LOGIC.md untuk
project AkuBelajar.

Dokumen ini mendefinisikan SEMUA aturan bisnis akademik yang tidak bisa
disimpulkan dari skema database saja. Ini adalah "otak" dari sistem.

Cakup semua hal berikut:

1. KALKULASI NILAI AKHIR
   - Formula lengkap: Nilai Akhir = (Rata-rata Tugas × 60%) + (Rata-rata Kuis × 40%)
   - Bagaimana jika siswa tidak punya nilai tugas sama sekali?
   - Bagaimana jika siswa tidak punya nilai kuis sama sekali?
   - Apakah bobot 60/40 bisa dikonfigurasi per mapel? Per sekolah?
   - Nilai remedial: apakah menggantikan nilai asli atau di-average?
   - Sertakan tabel: kondisi → formula yang digunakan

2. PREDIKAT & KKM
   - Definisikan predikat berdasarkan nilai:
     A = 90–100, B = 80–89, C = 70–79, D = 60–69, E = < 60
   - KKM default: [KEPUTUSAN DIPERLUKAN: 70?] — bisa dikonfigurasi per mapel
   - Siswa dengan nilai < KKM pada mapel apa saja yang perlu ditandai
   - Predikat di rapor: A (Sangat Baik), B (Baik), C (Cukup), D (Perlu Bimbingan)

3. ATURAN KEHADIRAN
   - Persentase kehadiran minimum: [KEPUTUSAN DIPERLUKAN: 75%?]
   - Dampak jika kehadiran di bawah minimum:
     siswa masih bisa dinilai? atau otomatis tidak lulus?
   - Status Izin dan Sakit: dihitung hadir atau tidak dalam persentase?
   - Terlambat: dihitung penuh hadir atau setengah hadir?
   - Sertakan rumus kalkulasi persentase kehadiran

4. ATURAN DEADLINE & KETERLAMBATAN
   - Deadline dihitung sampai jam berapa? (23:59:59 WIB hari H?)
   - Toleransi keterlambatan sistem (grace period): [KEPUTUSAN: 0 menit?]
   - Penalti nilai keterlambatan: contoh 10% per hari terlambat
   - Maksimal hari keterlambatan yang masih diterima: [KEPUTUSAN: 3 hari?]
   - Submit setelah batas maksimal: tidak diterima sama sekali

5. ATURAN KUIS & CBT
   - Satu kuis bisa dikerjakan berapa kali? (default: 1, bisa dikonfigurasi)
   - Jika multi-attempt: nilai mana yang dipakai? Tertinggi? Terakhir? Rata-rata?
   - Toleransi cheat detection sebelum sesi dikunci: [KEPUTUSAN: 3 kali?]
   - Waktu ujian server-authoritative: jika siswa submit setelah waktu habis → tidak diterima

6. ATURAN RAPOR & LOCK DATA
   - Rapor diterbitkan per semester (Semester 1: Desember, Semester 2: Juni)
   - Setelah rapor di-publish: nilai tidak bisa diubah kecuali SuperAdmin
   - Data yang di-lock: grades, report_cards
   - Alur unlock (jika ada kesalahan): SuperAdmin unlock → Guru edit → SuperAdmin lock ulang
   - Sertakan flowchart alur lock/unlock

7. EARLY WARNING SYSTEM RULES
   - Trigger peringatan ke Guru/Admin jika:
     * Nilai rata-rata siswa < 60 selama 2 minggu berturut-turut
     * Absensi Alfa ≥ 3 kali dalam 1 minggu
     * Tidak submit tugas ≥ 3 kali berturut-turut
   - Siapa yang menerima warning? Guru kelas? Wali kelas? Admin?
   - Sertakan tabel: kondisi → threshold → penerima notifikasi

Stack: PostgreSQL untuk constraint, Go untuk business logic layer,
Gemini AI untuk early warning analysis.
```

---

# ═══════════════════════════════════════════
# KELOMPOK D — KEAMANAN (SECURITY GAPS)
# ═══════════════════════════════════════════

---

## PROMPT D-1
**File output:** `docs/security/INPUT_VALIDATION_RULES.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama INPUT_VALIDATION_RULES.md untuk
project AkuBelajar.

Dokumen ini mendefinisikan aturan validasi untuk SETIAP field input di sistem.
Ini digunakan bersama antara: schema Zod di Next.js frontend DAN struct validator
di Go backend. Kedua sisi harus menggunakan aturan yang identik.

Definisikan aturan untuk SEMUA field berikut:

FIELD IDENTITAS USER:
- email: format RFC 5322, max 255 char, lowercase, trim whitespace
- password: min 8 char, max 72 char (bcrypt limit), wajib 1 huruf besar +
  1 angka + 1 simbol, tidak boleh sama dengan 3 password terakhir,
  cek terhadap list 10.000 common passwords
- name: min 2 char, max 100 char, hanya huruf + spasi + titik + koma,
  tidak boleh angka, tidak boleh karakter khusus lain
- nisn: tepat 10 digit angka, tidak boleh 0000000000
- nip: 18 digit angka (format ASN Indonesia), opsional
- phone_wa: format E.164 (+62xxxxxxxxxx), min 10 digit, max 15 digit
- birth_date: tidak boleh di masa depan, tidak boleh lebih dari 100 tahun lalu

FIELD AKADEMIK:
- class name: max 20 char, contoh "8A", "XII-IPA-1"
- subject name: max 100 char
- school code: 6 karakter alphanumeric uppercase
- assignment title: min 5 char, max 200 char
- assignment description: max 10.000 char, HTML di-sanitasi dengan bluemonday
- grade value: integer 0–100
- quiz title: min 5 char, max 200 char
- question text: min 10 char, max 2.000 char
- answer option: min 1 char, max 500 char

FILE UPLOAD:
- avatar: max 2MB, hanya JPG/PNG/WebP, dimensi min 100×100px, max 2000×2000px
- assignment file: max 20MB per file, max 3 file per submission,
  allowed: PDF/DOCX/PPTX/XLSX/JPG/PNG/ZIP
- proof/evidence: max 5MB, hanya PDF/JPG/PNG
- semua file: scan filename untuk path traversal attack
  (tidak boleh ada ../ atau absolute path)

QUERY PARAMETERS:
- pagination limit: integer, min 1, max 100, default 20
- page/cursor: integer positif, atau opaque cursor string max 200 char
- date filter: format ISO 8601 (YYYY-MM-DD)
- sort field: hanya whitelist nilai yang diizinkan (tidak boleh arbitrary string)
- search query: max 100 char, strip SQL special characters

Untuk setiap field sertakan:
1. Aturan validasi lengkap
2. Contoh value VALID
3. Contoh value INVALID + error message yang ditampilkan
4. Kode Zod schema (TypeScript)
5. Kode Go struct tag dengan validator

Stack: Go go-playground/validator, TypeScript Zod v3,
Next.js react-hook-form + zodResolver.
```

---

## PROMPT D-2
**File output:** `docs/security/SESSION_AND_TOKEN_SPEC.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama SESSION_AND_TOKEN_SPEC.md untuk
project AkuBelajar.

Dokumen ini mendefinisikan secara eksplisit bagaimana session dan token
dikelola di seluruh sistem.

Cakup semua hal berikut:

1. TOKEN LIFECYCLE
   - Access token (Paseto v4 local):
     * TTL: 15 menit
     * Payload: user_id, role, school_id, session_id, iat, exp
     * Disimpan: memory only di client (BUKAN localStorage/sessionStorage)
   - Refresh token:
     * TTL: 7 hari (persistent login) atau 24 jam (non-persistent)
     * Disimpan: httpOnly + Secure + SameSite=Strict cookie
     * Di-rotate setiap kali digunakan (refresh token rotation)
   - Sertakan diagram lifecycle token

2. CONCURRENT SESSION MANAGEMENT
   - Maksimal berapa device bisa login bersamaan? [KEPUTUSAN: 3 device?]
   - Apa yang terjadi jika melebihi batas? (oldest session di-revoke?)
   - Bagaimana user bisa lihat dan revoke session aktif?
   - Schema tabel active_sessions

3. TOKEN REVOCATION
   - Kapan token di-revoke: logout, ganti password, suspend akun, role berubah
   - Strategi: blocklist di Redis dengan key = session_id, TTL = sisa TTL token
   - Bagaimana middleware Go cek apakah token di-revoke sebelum proses request?
   - Sertakan pseudocode flow validasi token

4. REMEMBER ME FEATURE
   - Checkbox "Ingat saya" di halaman login
   - Perbedaan TTL refresh token: 24 jam vs 7 hari
   - Bagaimana cara identify "trusted device" setelah remember me?

5. SECURITY EVENTS YANG TRIGGER AUTO-REVOKE
   Daftar semua kondisi yang menyebabkan semua session user di-revoke:
   - Ganti password (manual atau reset)
   - Admin suspend akun
   - Admin ubah role
   - Deteksi login dari lokasi sangat berbeda dalam waktu singkat
   - User minta "logout semua device" dari settings

6. PASETO TOKEN IMPLEMENTATION
   - Sertakan contoh kode Go untuk generate dan verify Paseto v4 local token
   - Key management: private key disimpan di HashiCorp Vault
   - Key rotation: prosedur rotate key tanpa downtime

Stack: Go o1ecc8o/paseto, Redis untuk token blocklist,
Next.js cookies (httpOnly), HashiCorp Vault.
```

---

## PROMPT D-3
**File output:** `docs/security/FILE_UPLOAD_SECURITY.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama FILE_UPLOAD_SECURITY.md untuk
project AkuBelajar.

File upload adalah salah satu attack vector paling berbahaya.
Dokumen ini mendefinisikan semua proteksi yang harus diterapkan.

Cakup semua hal berikut:

1. VALIDASI FILE DI BACKEND (Go)
   - Jangan percaya Content-Type dari client — selalu detect dari magic bytes
   - Daftar magic bytes untuk setiap tipe file yang diizinkan (PDF, DOCX, JPG, PNG, WebP, ZIP)
   - Sertakan contoh kode Go untuk validate magic bytes
   - Reject file jika extension tidak cocok dengan actual content
   - Strip semua metadata EXIF dari gambar sebelum disimpan
   - Limit: max file size per endpoint (avatar: 2MB, submission: 20MB, dll.)

2. PENAMAAN & PENYIMPANAN FILE
   - Rename SEMUA file yang diupload — jangan pernah simpan dengan nama asli
   - Format nama baru: {uuid_v7}.{ext_lowercase}
   - Path di MinIO: {bucket}/{school_id}/{entity_type}/{entity_id}/{filename}
   - Tidak boleh ada file yang publicly accessible — semua via signed URL
   - TTL signed URL: 15 menit untuk download, 5 menit untuk upload presigned URL

3. SIGNED URL WORKFLOW
   - Sertakan sequence diagram: client request download → server validate auth
     → server generate signed URL → client download langsung dari MinIO
   - Sertakan sequence diagram: client request upload → server validate →
     server generate presigned PUT URL → client upload langsung ke MinIO →
     client notify server bahwa upload selesai → server verify file ada

4. RATE LIMITING UNTUK UPLOAD
   - Max upload per user per jam: 20 file
   - Max storage per sekolah: [KEPUTUSAN: 10GB?]
   - Alert admin jika storage sekolah > 80%

5. MALWARE SCANNING
   - Apakah perlu antivirus scan? (ClamAV integration)
   - Jika ya: async scan setelah upload, file berstatus QUARANTINE sampai clear
   - Jika scan gagal/timeout: file di-reject atau diterima?
   - Sertakan flowchart decision

6. CLEANUP
   - File orphan (upload tapi tidak pernah attach ke entity): hapus setelah 24 jam
   - File yang entity-nya di-soft-delete: kapan file fisik dihapus?
   - Scheduled job untuk cleanup

Stack: Go, MinIO S3-compatible, net/http magic byte detection,
disintegration/imaging untuk image processing.
```

---

## PROMPT D-4
**File output:** `docs/security/RATE_LIMIT_POLICY.md`
**Prioritas:** PENTING

```
Buatkan dokumen markdown lengkap dengan nama RATE_LIMIT_POLICY.md untuk
project AkuBelajar.

Saat ini hanya ada contoh kode rate limiter tanpa policy yang jelas.
Dokumen ini mendefinisikan limit konkret untuk setiap endpoint.

Cakup hal-hal berikut:

1. STRATEGI RATE LIMITING
   - Algoritma: Token Bucket (cocok untuk burst traffic)
   - Koordinasi antar instance: Redis sebagai shared counter
   - Key: kombinasi IP + user_id (authenticated) atau IP saja (unauthenticated)
   - Response jika terkena limit: HTTP 429 + header Retry-After

2. TABEL LIMIT PER ENDPOINT
   Buat tabel komprehensif dengan kolom:
   ENDPOINT | LIMIT | WINDOW | KEY | BURST ALLOWANCE | NOTES

   Sertakan limit untuk:
   - POST /auth/login: 5/menit per IP
   - POST /auth/password-reset/request: 3/jam per IP
   - POST /auth/password-reset/verify: 5/menit per token
   - POST /invite-tokens/claim: 10/jam per IP
   - POST /quizzes/generate-ai: 10/jam per user (Gemini API cost)
   - POST /assignments/:id/submissions (file upload): 20/jam per user
   - GET /api/* (general): 120/menit per user
   - WebSocket connections: 5 koneksi simultan per user

3. LAYER RATE LIMITING
   Jelaskan setiap layer dan bagaimana mereka bekerja bersama:
   - Layer 1: Cloudflare (edge, per IP global)
   - Layer 2: Nginx (per IP, sebelum masuk app)
   - Layer 3: Go middleware (per user/IP, per endpoint)
   - Sertakan diagram alur request melewati semua layer

4. RESPONSE FORMAT
   Sertakan contoh response HTTP 429 yang lengkap:
   - Status code: 429
   - Header: Retry-After, X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset
   - Body JSON dengan error code standar

5. WHITELIST & BYPASS
   - IP whitelist untuk monitoring/health check endpoints
   - Cara SuperAdmin bypass limit untuk operasi massal (bulk import)
   - Internal service-to-service calls (tidak perlu rate limit)

Stack: Go middleware dengan Redis, Nginx limit_req_zone, Cloudflare Rate Limiting rules.
```

---

## PROMPT D-5
**File output:** `docs/security/PASSWORD_POLICY.md`
**Prioritas:** PENTING

```
Buatkan dokumen markdown lengkap dengan nama PASSWORD_POLICY.md untuk
project AkuBelajar.

Cakup semua hal berikut:

1. ATURAN KEKUATAN PASSWORD
   - Panjang minimum: 8 karakter
   - Panjang maksimum: 72 karakter (batas bcrypt, meski kita pakai Argon2id)
   - Wajib mengandung: 1 huruf besar, 1 huruf kecil, 1 angka, 1 simbol
   - Tidak boleh mengandung: nama pengguna, email, NISN
   - Cek terhadap daftar 10.000 common passwords (haveibeenpwned list)
   - Sertakan contoh password: VALID vs INVALID dengan penjelasan

2. HASHING DENGAN ARGON2ID
   - Parameter yang digunakan:
     * Memory: 65536 KB (64 MB)
     * Iterations: 3
     * Parallelism: 4
     * Salt: 16 byte random per password
     * Output length: 32 byte
   - Sertakan contoh kode Go untuk hash dan verify
   - Kapan parameter perlu di-upgrade? (hardware lebih cepat)

3. PASSWORD HISTORY
   - Tidak boleh reuse 3 password terakhir
   - Cara menyimpan history dengan aman (hash, bukan plain text)
   - Schema tabel password_histories

4. TEMPORARY PASSWORD
   - Format: 12 karakter random (uppercase + lowercase + angka)
   - Contoh generator di Go: crypto/rand
   - TTL temporary password: [KEPUTUSAN: 7 hari?] sebelum akun di-lock
   - Sudah expired: admin harus generate ulang

5. PASSWORD CHANGE POLICY
   - Force change: saat first login dan setelah reset
   - Voluntary change: kapan saja dari settings
   - Setelah ganti: semua session lain di-revoke (kecuali session saat ini)
   - Cooldown sebelum ganti password lagi: [KEPUTUSAN: 24 jam?]

Stack: Go golang.org/x/crypto/argon2, crypto/rand, PostgreSQL.
```

---

# ══════════════════════════════════════════
# KELOMPOK E — FRONTEND & UI/UX
# ══════════════════════════════════════════

---

## PROMPT E-1
**File output:** `docs/frontend/SITEMAP.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama SITEMAP.md untuk
project AkuBelajar.

Dokumen ini mendefinisikan SEMUA halaman yang ada di aplikasi Next.js,
URL path-nya, siapa yang bisa akses, dan apa yang ada di halaman tersebut.

Format tabel untuk setiap halaman:
PATH | NAMA HALAMAN | ROLE YANG BISA AKSES | KOMPONEN UTAMA | AUTH REQUIRED

Definisikan halaman untuk semua route berikut:

PUBLIC ROUTES (tidak perlu login):
- / (landing page)
- /login
- /forgot-password
- /reset-password/[token]
- /invite/[token] (claim invite token)

ONBOARDING ROUTES (sudah login tapi belum lengkap):
- /onboarding/change-password
- /onboarding/complete-profile
- /onboarding/select-class (khusus siswa)

SHARED ROUTES (semua role yang sudah login):
- /dashboard
- /notifications
- /settings
- /settings/profile
- /settings/security
- /settings/notifications

ADMIN ROUTES:
- /admin/users
- /admin/users/[id]
- /admin/users/import
- /admin/schools
- /admin/academic-years
- /admin/classes
- /admin/subjects
- /admin/invite-tokens
- /admin/audit-logs
- /admin/system

TEACHER ROUTES:
- /teacher/dashboard
- /teacher/attendance/[class_id]
- /teacher/assignments
- /teacher/assignments/new
- /teacher/assignments/[id]
- /teacher/assignments/[id]/submissions
- /teacher/quizzes
- /teacher/quizzes/new
- /teacher/quizzes/[id]
- /teacher/quizzes/[id]/results
- /teacher/grades/[class_id]
- /teacher/report-cards/[class_id]

STUDENT ROUTES:
- /student/dashboard
- /student/assignments
- /student/assignments/[id]
- /student/quizzes
- /student/quizzes/[id]/take (CBT interface)
- /student/quizzes/[id]/review
- /student/grades
- /student/attendance
- /student/report-cards

Untuk setiap route sertakan:
- Middleware yang berlaku (auth check, role check, onboarding check)
- Apakah di-render SSR, ISR, atau CSR, dan mengapa
- Data fetching strategy (server component atau client + react-query)
- Layout yang digunakan

Sertakan juga diagram pohon navigasi (tree diagram menggunakan ASCII art atau mermaid).

Stack: Next.js 15 App Router, TypeScript, middleware.ts untuk auth redirect.
```

---

## PROMPT E-2
**File output:** `docs/frontend/STATE_MANAGEMENT.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama STATE_MANAGEMENT.md untuk
project AkuBelajar.

Saat ini hanya disebutkan "menggunakan Zustand" tanpa detail apapun.
Dokumen ini mendefinisikan arsitektur state management secara menyeluruh.

Cakup semua hal berikut:

1. PRINSIP STATE MANAGEMENT
   - Server state (data dari API): dikelola TanStack Query (react-query)
   - Client state (UI state): dikelola Zustand
   - Form state: dikelola react-hook-form
   - URL state (filter, pagination): dikelola Next.js searchParams
   - Kapan masing-masing digunakan? Sertakan decision tree

2. ZUSTAND STORES (definisikan setiap store)

   authStore:
   - State: user (profile + role), accessToken, isAuthenticated, isOnboarding
   - Actions: setUser, clearAuth, setToken
   - Persistence: tidak (token di cookie httpOnly, tidak bisa diakses JS)

   uiStore:
   - State: sidebarOpen, theme (light/dark), activeNotificationCount
   - Actions: toggleSidebar, setTheme, setNotificationCount
   - Persistence: localStorage untuk theme preference

   quizSessionStore:
   - State: sessionId, quizId, currentQuestionIndex, answers Map,
     timeRemaining, sessionStatus, wsConnection
   - Actions: initSession, answerQuestion, nextQuestion, submitSession
   - Persistence: sessionStorage (survive refresh, tidak survive close tab)

   Definisikan juga: notificationStore, attendanceStore (draft absensi)

3. TANSTACK QUERY PATTERNS
   - Query keys convention: ['entity', 'action', { params }]
   - Cache time per entity: berapa lama data dianggap fresh?
   - Stale time: kapan refetch otomatis dilakukan?
   - Optimistic updates: untuk operasi mana yang perlu optimistic update?
   - Mutation + cache invalidation pattern
   - Sertakan contoh kode TypeScript untuk setiap pattern

4. LOADING, ERROR, EMPTY STATES
   Definisikan komponen standar dan kapan digunakan:
   - <PageSkeleton /> — full page loading
   - <TableSkeleton rows={n} /> — loading tabel
   - <CardSkeleton /> — loading card
   - <ErrorBoundary /> — unexpected error
   - <ApiError error={e} onRetry={fn} /> — API error dengan retry
   - <EmptyState icon title description action /> — data kosong

5. REAL-TIME STATE UPDATES (WebSocket)
   - Bagaimana WebSocket event mengupdate state Zustand dan TanStack Query cache?
   - Sertakan contoh: saat server kirim "notification:new", update notificationStore
     DAN invalidate query ['notifications']

Stack: Zustand v5, TanStack Query v5, Next.js App Router, TypeScript strict mode.
```

---

# ══════════════════════════════════════════════════
# KELOMPOK F — TESTING & QUALITY GATES
# ══════════════════════════════════════════════════

---

## PROMPT F-1
**File output:** `docs/testing/ACCEPTANCE_CRITERIA.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama ACCEPTANCE_CRITERIA.md untuk
project AkuBelajar.

Dokumen ini mendefinisikan kapan setiap fitur dianggap SELESAI dan SIAP RILIS.
Tanpa ini, AI agent tidak tahu batas "done".

Format untuk setiap fitur:
**Fitur: [Nama Fitur]**
Given [konteks awal]
When [aksi yang dilakukan]
Then [hasil yang diharapkan]
And [kondisi tambahan]

Tulis acceptance criteria lengkap untuk SEMUA fitur berikut:

1. AUTENTIKASI
   - Login berhasil dengan kredensial valid
   - Login gagal dengan password salah (max 5x sebelum lock)
   - Akun terkunci otomatis dan bisa di-unlock
   - Logout membersihkan semua session data
   - Refresh token rotation bekerja

2. MANAJEMEN AKUN
   - Admin bisa bulk import 100 siswa dari Excel tanpa error
   - Invite token bisa di-generate, dikirim, diklaim, dan expired dengan benar
   - First login redirect ke onboarding yang benar
   - Force password change tidak bisa di-skip

3. ABSENSI
   - Guru bisa input absensi untuk semua siswa dalam 1 kelas
   - Status absensi tersimpan dengan benar dan bisa dilihat siswa
   - Edit absensi setelah 24 jam memerlukan alasan
   - Rekap kehadiran menghitung persentase dengan benar

4. TUGAS
   - Siswa bisa upload file dan submit sebelum deadline
   - Sistem menolak submit setelah deadline jika tidak ada grace period
   - Guru bisa download semua submission dan beri nilai
   - Nilai tugas masuk ke buku nilai otomatis

5. KUIS / CBT
   - Siswa hanya bisa mulai kuis yang sudah dipublish dalam window waktu
   - Timer berjalan di server, tidak bisa dimanipulasi client
   - Cheat detection memicu warning setelah threshold terlampaui
   - Auto-submit berjalan saat waktu habis
   - Nilai langsung muncul setelah submit

6. NILAI & RAPOR
   - Formula nilai akhir dihitung dengan benar
   - Rapor PDF bisa di-generate dengan data yang benar
   - Setelah rapor di-lock, nilai tidak bisa diubah oleh Guru
   - Digital signature pada PDF bisa diverifikasi

7. NOTIFIKASI
   - Siswa menerima notifikasi WA saat tugas baru dipublish
   - Push notification in-app muncul real-time via WebSocket
   - Preferensi notifikasi user direspek

Sertakan juga: daftar browser yang harus disupport, ukuran layar yang harus
ditest (375px, 768px, 1280px, 1920px), dan performance criteria (LCP < 2.5s).
```

---

## PROMPT F-2
**File output:** `docs/testing/EDGE_CASES.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama EDGE_CASES.md untuk
project AkuBelajar.

Dokumen ini mendefinisikan skenario-skenario ekstrem dan negatif yang HARUS
ditangani oleh sistem. Ini adalah kasus yang paling sering diabaikan
tapi paling sering menyebabkan bug di production.

Format untuk setiap edge case:
**Skenario:** [Deskripsi singkat]
**Konteks:** [Kondisi saat ini]
**Aksi:** [Yang dilakukan pengguna / sistem]
**Expected behavior:** [Apa yang SEHARUSNYA terjadi]
**Risiko jika tidak ditangani:** [Dampak ke user/data]

Definisikan edge case untuk semua kategori:

AUTENTIKASI & SESSION:
- User login di device A, ganti password di device B → apa yang terjadi di device A?
- Token expired tepat saat request dikirim (race condition)
- User klik login berulang kali sebelum response kembali (double submit)
- Session Redis habis sebelum token JWT expired

MANAJEMEN FILE:
- User upload file dengan ekstensi .jpg tapi isinya executable
- User upload file dengan nama mengandung ../../../etc/passwd
- Upload file 0 byte
- Upload file tepat di batas limit (20MB - 1 byte vs 20MB + 1 byte)
- Dua user upload file dengan nama sama secara bersamaan
- MinIO tidak tersedia saat upload

ABSENSI:
- Guru input absensi untuk tanggal yang sama dua kali
- Guru input absensi untuk tanggal libur nasional
- Siswa dikeluarkan dari kelas setelah absensi sudah diinput
- Koneksi putus saat guru sedang input absensi 30 siswa

TUGAS & SUBMISSION:
- Siswa submit dua kali dalam hitungan detik (double submit)
- Guru menghapus tugas setelah ada 10 submission masuk
- Deadline diubah mundur setelah ada yang submit
- Guru mengubah nilai setelah rapor di-lock

KUIS / CBT:
- Siswa buka kuis di dua tab browser secara bersamaan
- Waktu kuis berakhir tepat saat siswa menekan submit
- Soal AI gagal di-generate → kuis sudah dibuat tapi kosong
- Siswa kehilangan koneksi 1 detik sebelum waktu habis
- Guru memperpanjang waktu saat siswa sedang mengerjakan

NILAI & KALKULASI:
- Nilai tugas diinput 0 untuk semua siswa — apakah ini valid atau error?
- Bobot tugas 60% + kuis 40% = 100%, tapi tidak ada nilai kuis
- Nilai di atas 100 (karena bonus?) — sistem allow atau reject?

DATABASE & CONCURRENCY:
- Dua guru menginput nilai untuk siswa yang sama secara bersamaan
- Bulk import 1.000 siswa saat server sedang tinggi load
- Database migration gagal di tengah jalan

Untuk setiap edge case, tentukan juga: apakah ini harus dicegah di frontend,
backend, database, atau semua lapisan?
```

---

# ══════════════════════════════════════════
# KELOMPOK G — OPERASIONAL
# ══════════════════════════════════════════

---

## PROMPT G-1
**File output:** `docs/ops/DOCKER_SETUP.md`
**Prioritas:** KRITIS

```
Buatkan dokumen markdown lengkap dengan nama DOCKER_SETUP.md untuk
project AkuBelajar.

Dokumen ini menyertakan file konfigurasi lengkap yang langsung bisa digunakan
oleh developer untuk menjalankan seluruh stack di local machine.

Sertakan file-file berikut secara lengkap dan siap pakai:

1. docker-compose.yml (development)
   Harus mencakup semua service:
   - postgres:16 dengan health check
   - redis:7-alpine dengan persistence
   - minio/minio dengan bucket auto-create
   - mailhog (SMTP catcher untuk dev, bukan email asli)
   - go-api (build dari Dockerfile, hot reload dengan air)
   - nextjs (build dari Dockerfile, hot reload)
   - nginx (reverse proxy semua service)
   Tambahkan: volume mounts, network, environment variables dari .env

2. docker-compose.prod.yml (production override)
   Perbedaan dari dev:
   - No volume mount source code
   - Build dari image registry
   - Tidak ada mailhog
   - Resource limits (CPU, memory) per container
   - Restart policy: always

3. Dockerfile untuk Go API
   - Multi-stage build: builder (golang:1.23) → runner (alpine:3.19)
   - Copy hanya binary, bukan source code
   - Non-root user (uid 1000)
   - Health check endpoint /health

4. Dockerfile untuk Next.js
   - Multi-stage: deps → builder → runner
   - Standalone output mode
   - Non-root user
   - Proper COPY untuk .next/standalone

5. .env.example
   Semua environment variable yang dibutuhkan dengan nilai contoh atau placeholder:
   - Database URL
   - Redis URL
   - MinIO credentials
   - Paseto private key placeholder
   - Gemini API key placeholder
   - Fonnte WA API key placeholder
   - SMTP settings
   - Domain/URL settings

6. Makefile
   Commands yang tersedia:
   - make dev (docker-compose up development)
   - make build (build semua image)
   - make migrate-up / migrate-down
   - make seed
   - make test
   - make lint
   - make logs service=api

Sertakan juga instruksi setup lengkap step-by-step untuk developer baru.

Stack: Docker 24+, Docker Compose v2, Go 1.23, Node.js 22, PostgreSQL 16, Redis 7.
```

---

## PROMPT G-2
**File output:** `docs/ops/SEED_DATA_SPEC.md`
**Prioritas:** PENTING

```
Buatkan dokumen markdown lengkap dengan nama SEED_DATA_SPEC.md untuk
project AkuBelajar.

Dokumen ini mendefinisikan SEMUA data yang harus ada saat sistem pertama kali
dijalankan, baik untuk development maupun production.

Cakup hal-hal berikut:

1. PRODUCTION SEED (data wajib ada di production)
   - 1 SuperAdmin default (credentials dikirim ke owner, wajib ganti saat first login)
   - Enum values / lookup tables jika ada
   - Konfigurasi default sistem (notification settings, dll.)
   - Tidak boleh ada data dummy di production seed

2. DEVELOPMENT SEED (data untuk mempermudah development & testing)
   Buat seed data realistis:

   Sekolah:
   - 1 sekolah: "SMP AkuBelajar" dengan school_code "AB0001"

   User (dengan password default "Test@12345"):
   - 1 SuperAdmin: admin@akubelajar.test
   - 3 Guru: guru1@, guru2@, guru3@akubelajar.test
   - 1 Ketua Kelas: ketua@akubelajar.test
   - 10 Siswa: siswa01@ sampai siswa10@akubelajar.test

   Academic data:
   - 1 Tahun Ajaran: 2025/2026 (aktif)
   - 2 Kelas: 8A (10 siswa), 8B (kosong)
   - 5 Mata Pelajaran: Matematika, Bahasa Indonesia, IPA, IPS, Bahasa Inggris
   - Jadwal mapel untuk 8A

   Sample content:
   - 3 Tugas (1 aktif, 1 sudah deadline, 1 sudah dinilai)
   - 2 Kuis (1 published, 1 draft)
   - Data absensi 2 minggu terakhir untuk 8A

3. IDEMPOTENT SEED
   - Seed harus aman dijalankan berkali-kali (INSERT ... ON CONFLICT DO NOTHING)
   - Cara run: make seed atau go run cmd/seed/main.go
   - Flag untuk reset: make seed RESET=true (hapus semua data dulu)

4. FACTORY FUNCTIONS (untuk test)
   Definisikan helper functions yang bisa dipakai di unit test:
   - CreateTestUser(role, overrides)
   - CreateTestClass(schoolId, overrides)
   - CreateTestAssignment(classId, teacherId, overrides)
   - CreateTestQuiz(classId, teacherId, overrides)
   Sertakan contoh penggunaan dalam Go test file.

Stack: Go, PostgreSQL, pgx untuk bulk insert.
```

---

*— End of AKUBELAJAR_PROMPTS.md —*
*Total: 18 prompt untuk 28 gap kritis + 14 gap penting*
*Urutan pengerjaan yang disarankan: C-1 → C-2 → C-3 → B-1 → A-1 → A-2 → A-3 → sisanya*
