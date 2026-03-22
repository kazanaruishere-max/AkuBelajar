<div align="center">

# 📚 AkuBelajar

### Sistem Manajemen Pembelajaran (LMS) Modern untuk Sekolah Indonesia

[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![Next.js](https://img.shields.io/badge/Next.js-16-black?style=for-the-badge&logo=next.js&logoColor=white)](https://nextjs.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![Supabase](https://img.shields.io/badge/Supabase-Cloud-3ECF8E?style=for-the-badge&logo=supabase&logoColor=white)](https://supabase.com)
[![Redis](https://img.shields.io/badge/Redis-Upstash-DC382D?style=for-the-badge&logo=redis&logoColor=white)](https://upstash.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge)](LICENSE)

> ⚠️ **STATUS: DEMO / DEVELOPMENT** — Proyek ini masih dalam tahap pengembangan aktif dan belum siap untuk produksi.

[🚀 Quick Start](#-quick-start) · [📖 Dokumentasi](#-arsitektur) · [🤝 Kontribusi](CONTRIBUTING.md) · [🔒 Keamanan](SECURITY.md)

</div>

---

## 📋 Daftar Isi

- [Tentang Project](#-tentang-project)
- [Fitur Utama](#-fitur-utama)
- [Tech Stack](#-tech-stack)
- [Arsitektur](#-arsitektur)
- [Alur Kerja Sistem](#-alur-kerja-sistem)
- [Database Schema](#-database-schema)
- [API Endpoints](#-api-endpoints)
- [Quick Start](#-quick-start)
- [Konfigurasi](#-konfigurasi)
- [Akun Demo](#-akun-demo)
- [Sprint Roadmap](#-sprint-roadmap)
- [Kontribusi](#-kontribusi)
- [Lisensi](#-lisensi)

---

## 🎯 Tentang Project

**AkuBelajar** adalah platform Learning Management System (LMS) yang dirancang khusus untuk sekolah-sekolah di Indonesia. Dibangun dengan arsitektur modern, AkuBelajar menyediakan solusi lengkap untuk manajemen akademik digital — mulai dari pengelolaan penugasan, ujian berbasis komputer (CBT), presensi, hingga penilaian otomatis.

### Mengapa AkuBelajar?

| Masalah                               | Solusi AkuBelajar                                |
| :------------------------------------ | :----------------------------------------------- |
| Administrasi sekolah masih manual     | ✅ Digitalisasi penuh dengan dashboard analytics |
| Ujian kertas boros & lambat diperiksa | ✅ CBT dengan auto-grading & randomisasi soal    |
| Presensi manual sulit direkap         | ✅ Presensi digital batch + rekap otomatis       |
| Komunikasi guru-siswa tidak terpusat  | ✅ Notifikasi in-app real-time                   |
| Sulit membuat soal ujian bervariasi   | ✅ AI Quiz Generator (Google Gemini)             |

---

## ✨ Fitur Utama

### 🔐 Autentikasi & Keamanan

- Login dengan **Paseto v4** (lebih aman dari JWT)
- Password hashing **Argon2id** (state-of-the-art)
- **Role-Based Access Control** (Super Admin, Guru, Ketua Kelas, Siswa)
- Rate limiting & brute-force protection
- Refresh token rotation

### 📚 Manajemen Akademik

- CRUD Tahun Ajaran, Kelas, Mata Pelajaran
- Assign siswa ke kelas & guru ke mata pelajaran
- Multi-tenant support (per sekolah)

### 📝 Penugasan (Assignment)

- Guru membuat tugas dengan deadline & lampiran
- Siswa submit tugas + file upload ke Supabase Storage
- Fitur keterlambatan (late penalty otomatis)
- Grading dengan feedback

### 🎯 Kuis CBT (Computer-Based Testing)

- Ujian berbasis komputer dengan timer
- Randomisasi soal & opsi jawaban (anti-mencontek)
- Auto-grading untuk pilihan ganda
- Multiple attempts dengan batas percobaan
- **AI Quiz Generator** — buat soal dari prompt via Google Gemini

### ✅ Presensi Digital

- Input presensi batch per kelas
- Status: Hadir, Izin, Sakit, Absent, Terlambat
- Rekap presensi per siswa

### 📊 Penilaian & Rapor

- Input nilai per kategori (tugas, kuis, UTS, UAS)
- Perhitungan nilai tertimbang otomatis
- Rekap nilai per siswa & per kelas

### 🔔 Notifikasi

- Notifikasi in-app untuk semua role
- Broadcast notifikasi dari guru/admin
- Badge unread counter

### 📈 Dashboard Analytics

- **Admin**: Total users, guru, siswa, kelas, mapel
- **Guru**: Tugas dibuat, kuis aktif, submission pending
- **Siswa**: Tugas selesai, kuis diambil, rata-rata nilai

### 👥 Admin Panel

- CRUD user management
- Filter user by role
- Soft delete & aktivasi/deaktivasi akun

---

## 🛠 Tech Stack

### Backend

| Teknologi         | Kegunaan                                    |
| :---------------- | :------------------------------------------ |
| **Go 1.23**       | Language utama backend                      |
| **Gin**           | HTTP web framework                          |
| **pgx/v5**        | PostgreSQL driver (native, zero-dependency) |
| **Paseto v4**     | Token-based authentication                  |
| **Argon2id**      | Password hashing                            |
| **validator/v10** | Request validation                          |

### Frontend

| Teknologi         | Kegunaan                     |
| :---------------- | :--------------------------- |
| **Next.js 16**    | React framework (App Router) |
| **TypeScript**    | Type-safe JavaScript         |
| **Zustand**       | State management             |
| **CSS Variables** | Dark theme design system     |

### Infrastructure

| Teknologi         | Kegunaan                                         |
| :---------------- | :----------------------------------------------- |
| **Supabase**      | PostgreSQL 16 (managed) + Storage                |
| **Upstash**       | Redis 7 (serverless) untuk rate limiting & cache |
| **Google Gemini** | AI quiz question generation                      |

---

## 🏗 Arsitektur

### Arsitektur Sistem

```
┌────────────────────────────────────────────────────────────────────┐
│                        CLIENT (Browser)                            │
│                     Next.js 16 + TypeScript                        │
│              Zustand Store │ API Client (fetch)                     │
└──────────────────────┬─────────────────────────────────────────────┘
                       │ HTTPS (REST API)
                       ▼
┌────────────────────────────────────────────────────────────────────┐
│                      GIN HTTP SERVER (:8080)                       │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐   │
│  │   CORS   │→ │Rate Limit│→ │  Auth MW │→ │   RBAC (Role)    │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────────────┘   │
│                                                                    │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │                    HANDLER LAYER                             │  │
│  │  Auth │ Academic │ Assignment │ Quiz │ Attendance │ Grade    │  │
│  │  Notification │ Upload │ Dashboard │ Admin                   │  │
│  └────────────────────────┬─────────────────────────────────────┘  │
│                           │                                        │
│  ┌────────────────────────▼─────────────────────────────────────┐  │
│  │                   SERVICE LAYER                              │  │
│  │  Business Logic │ Validation │ Auto-grading │ AI Generation  │  │
│  └────────────────────────┬─────────────────────────────────────┘  │
│                           │                                        │
│  ┌────────────────────────▼─────────────────────────────────────┐  │
│  │                  REPOSITORY LAYER                            │  │
│  │              SQL Queries │ Data Access                        │  │
│  └──────┬──────────────────┬──────────────────┬─────────────────┘  │
└─────────┼──────────────────┼──────────────────┼────────────────────┘
          ▼                  ▼                  ▼
   ┌─────────────┐   ┌─────────────┐   ┌──────────────┐
   │ PostgreSQL  │   │   Redis     │   │  Supabase    │
   │ (Supabase)  │   │  (Upstash)  │   │  Storage     │
   │             │   │             │   │              │
   │ 15+ tables  │   │ Rate limit  │   │ File upload  │
   │ ENUM types  │   │ Session     │   │ CDN URL      │
   │ Indexes     │   │ Cache       │   │              │
   └─────────────┘   └─────────────┘   └──────────────┘
                                               │
                                               ▼
                                     ┌──────────────────┐
                                     │  Google Gemini    │
                                     │  API (AI Quiz)    │
                                     └──────────────────┘
```

### Struktur Folder

```
AkuBelajar/
├── backend/
│   ├── cmd/
│   │   ├── api/main.go              # Entry point server
│   │   └── seed/main.go             # Database seeder
│   ├── config/config.go             # Environment configuration
│   ├── internal/                    # Business modules
│   │   ├── auth/                    # 🔐 Login, register, token refresh
│   │   │   ├── handler.go           #    HTTP handlers
│   │   │   ├── service.go           #    Business logic
│   │   │   ├── repository.go        #    Database queries
│   │   │   ├── model.go             #    DTOs & domain models
│   │   │   └── routes.go            #    Route registration
│   │   ├── academic/                # 📚 Tahun ajaran, kelas, mapel
│   │   ├── assignment/              # 📝 Tugas + submissions
│   │   ├── quiz/                    # 🎯 CBT + AI generation
│   │   ├── attendance/              # ✅ Presensi digital
│   │   ├── grade/                   # 📊 Penilaian
│   │   ├── notification/            # 🔔 Notifikasi
│   │   ├── upload/                  # 📁 File upload
│   │   ├── dashboard/               # 📈 Analytics
│   │   ├── admin/                   # 👥 User management
│   │   └── middleware/              # 🛡 Auth, RBAC, rate limit, logger
│   ├── pkg/                         # Shared packages
│   │   ├── ai/gemini.go             #    Google Gemini client
│   │   ├── cache/redis.go           #    Redis connection
│   │   ├── database/postgres.go     #    PostgreSQL pool
│   │   ├── security/                #    Paseto + Argon2id
│   │   ├── storage/supabase.go      #    File storage client
│   │   ├── response/response.go     #    Standard API response
│   │   └── validator/validator.go   #    Request validation
│   ├── migrations/                  # SQL migration files
│   │   ├── 000001-000012            #    Individual migrations
│   │   ├── 000013_reconcile...      #    Schema reconciliation
│   │   ├── combined_all.sql         #    ⭐ Single file for fresh DB
│   │   └── seed.sql                 #    Reference seed data
│   ├── .env                         # Environment variables
│   └── go.mod                       # Go dependencies
│
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   │   ├── (auth)/login/        # 🔐 Halaman login
│   │   │   ├── (dashboard)/         # 📊 Dashboard layout + sidebar
│   │   │   │   ├── layout.tsx       #    Sidebar + topbar + role nav
│   │   │   │   ├── dashboard/       #    Analytics per role
│   │   │   │   ├── admin/           #    Admin pages (users, years, etc.)
│   │   │   │   ├── teacher/         #    Teacher pages (tugas, kuis, absen)
│   │   │   │   ├── student/         #    Student pages (tugas, kuis, nilai)
│   │   │   │   └── notifications/   #    Notifikasi
│   │   │   ├── layout.tsx           #    Root layout
│   │   │   └── page.tsx             #    Redirect → /login
│   │   ├── components/
│   │   │   └── Sidebar.tsx          #    Role-based sidebar component
│   │   └── lib/
│   │       ├── api/client.ts        #    API client (fetch + auto-refresh)
│   │       ├── hooks/useAuth.ts     #    Auth hook (login/logout/refresh)
│   │       └── store/
│   │           ├── authStore.ts     #    Zustand auth state
│   │           └── uiStore.ts      #    UI state (sidebar toggle)
│   ├── .env.local                   # Frontend env vars
│   └── package.json
│
├── README.md                        # 📖 Dokumentasi ini
├── LICENSE                          # MIT License
├── CODE_OF_CONDUCT.md               # Kode etik kontributor
├── CONTRIBUTING.md                  # Panduan kontribusi
└── SECURITY.md                      # Kebijakan keamanan
```

---

## 🔄 Alur Kerja Sistem

### 1. Alur Autentikasi

```
┌──────────┐     POST /auth/login       ┌──────────────┐
│  Client  │ ────────────────────────► │   Auth       │
│ (Login)  │     {email, password}      │  Handler     │
└──────────┘                            └──────┬───────┘
                                               │
                      ┌────────────────────────┤
                      ▼                        ▼
              ┌──────────────┐        ┌──────────────┐
              │  Find User   │        │ Verify Hash  │
              │  by Email    │        │ (Argon2id)   │
              └──────┬───────┘        └──────┬───────┘
                     │                       │
                     ▼                       ▼
              ┌──────────────┐        ┌──────────────┐
              │ Check Lock   │        │ Generate     │
              │ Status       │        │ Paseto Token │
              └──────┬───────┘        └──────┬───────┘
                     │                       │
                     ▼                       ▼
              ┌──────────────────────────────────────┐
              │  Response: {access_token, refresh_    │
              │  token, user: {id, email, role, ...}} │
              └──────────────────────────────────────┘
```

### 2. Alur CBT (Computer-Based Testing)

```
  GURU                                    SISWA
   │                                        │
   │  1. Buat Kuis                          │
   │  POST /quizzes/teacher                 │
   │  {title, time_limit, class_id}         │
   │──────────────────────────►             │
   │                                        │
   │  2. Tambah Soal                        │
   │  POST /quizzes/:id/questions           │
   │  {question_text, options, answer}      │
   │──────────────────────────►             │
   │                                        │
   │  3. Publish                            │
   │  POST /quizzes/:id/publish             │
   │──────────────────────────►             │
   │                                        │
   │                                        │  4. Start Exam
   │                                        │  POST /quizzes/student/:id/start
   │                                        │◄────────────────────────
   │                                        │
   │                                        │  → returns: session_id + questions
   │                                        │  → timer starts (time_limit mins)
   │                                        │  → questions randomized
   │                                        │
   │                                        │  5. Answer Questions
   │                                        │  POST /sessions/:id/answer
   │                                        │  {question_id, selected_key}
   │                                        │◄────────────────────────
   │                                        │  (repeat for each question)
   │                                        │
   │                                        │  6. Submit
   │                                        │  POST /sessions/:id/submit
   │                                        │◄────────────────────────
   │                                        │
   │  ┌─────────────────────────────────┐   │  → Auto-grade (count correct)
   │  │       AUTO-GRADING ENGINE       │   │  → Return score
   │  │                                 │   │
   │  │  for each answer:               │   │
   │  │    hash(selected_key) ==        │   │
   │  │    question.answer_hash ?       │   │
   │  │    → correct++ : skip           │   │
   │  │                                 │   │
   │  │  score = (correct/total) * 100  │   │
   │  └─────────────────────────────────┘   │
   │                                        │
   │  7. Review Results                     │
   │  GET /quizzes/:id/sessions             │
   │  → List all student scores             │
   │──────────────────────────►             │
```

### 3. Alur Penugasan

```
  GURU                                    SISWA
   │                                        │
   │  1. Buat Tugas                         │
   │  POST /assignments/teacher             │
   │  {title, description, deadline}        │
   │──────────────────────────►             │
   │                                        │
   │  2. Publish                            │
   │  POST /:id/publish                     │
   │──────────────────────────►             │
   │                                        │
   │                                        │  3. Lihat Tugas
   │                                        │  GET /assignments/student
   │                                        │◄────────────────────────
   │                                        │
   │                                        │  4. Submit + Upload File
   │                                        │  POST /:id/submit
   │                                        │  + POST /upload (multipart)
   │                                        │◄────────────────────────
   │                                        │
   │  5. Grade Submission                   │
   │  POST /submissions/:id/grade           │
   │  {grade: 85, feedback: "Bagus!"}       │
   │──────────────────────────►             │
   │                                        │
   │  ┌─────────────────────────────────┐   │
   │  │     LATE PENALTY CALCULATOR     │   │
   │  │                                 │   │
   │  │  if submitted > deadline:       │   │
   │  │    late_days = diff(days)        │   │
   │  │    penalty = late_days * 10%    │   │
   │  │    final = grade * (1-penalty)  │   │
   │  └─────────────────────────────────┘   │
```

### 4. Alur AI Quiz Generation

```
┌──────────┐                          ┌──────────────┐
│  Teacher  │  POST /:id/ai-generate  │  Quiz        │
│           │ ─────────────────────►  │  Handler     │
│  prompt:  │  {prompt, count: 5,     │              │
│  "Buat    │   question_type: "mc"}  └──────┬───────┘
│  soal     │                                │
│  MTK      │                                ▼
│  pecahan" │                         ┌──────────────┐
└──────────┘                          │  AI Service  │
                                      │              │
                                      │  1. Build    │
                                      │     prompt   │
                                      │  2. Call     │
                                      │     Gemini   │
                                      │  3. Parse    │
                                      │     JSON     │
                                      │  4. Hash     │
                                      │     answers  │
                                      │  5. Save to  │
                                      │     DB       │
                                      └──────┬───────┘
                                             │
                                             ▼
                                      ┌──────────────┐
                                      │ Google       │
                                      │ Gemini API   │
                                      │              │
                                      │ Returns:     │
                                      │ [{question,  │
                                      │   options,   │
                                      │   answer}]   │
                                      └──────────────┘
```

---

## 💾 Database Schema

### Entity Relationship Diagram

```
┌─────────────┐     ┌─────────────┐     ┌──────────────┐
│   schools    │────<│    users     │────<│ user_profiles│
│             │     │             │     │              │
│ id (PK)     │     │ id (PK)     │     │ user_id (FK) │
│ name        │     │ school_id   │     │ nisn         │
│ code        │     │ email       │     │ nip          │
│ config      │     │ password_hash│    │ phone_wa     │
└─────────────┘     │ role (enum) │     └──────────────┘
                    │ full_name   │
                    └──────┬──────┘
                           │
          ┌────────────────┼────────────────┐
          ▼                ▼                ▼
   ┌─────────────┐  ┌──────────┐   ┌─────────────┐
   │student_class│  │class_subj│   │ attendances  │
   │             │  │          │   │              │
   │ student_id  │  │ class_id │   │ student_id   │
   │ class_id    │  │subject_id│   │ class_id     │
   │ acad_year_id│  │teacher_id│   │ date, status │
   └──────┬──────┘  └────┬─────┘   └──────────────┘
          │              │
          ▼              ▼
   ┌─────────────┐  ┌──────────────┐
   │   classes    │  │   subjects   │
   │             │  │              │
   │ school_id   │  │ school_id    │
   │ acad_yr_id  │  │ name, code   │
   │ name, grade │  └──────────────┘
   └─────────────┘
          │
     ┌────┴─────────────┐
     ▼                  ▼
┌──────────┐     ┌──────────┐     ┌──────────────┐
│assignments│    │ quizzes   │    │   grades      │
│          │     │          │     │              │
│ class_id │     │ class_id │     │ student_id   │
│ title    │     │ title    │     │ subject_id   │
│ deadline │     │timelimit │     │ score        │
└────┬─────┘     └────┬─────┘     └──────────────┘
     │                │
     ▼                ▼
┌──────────┐     ┌──────────┐     ┌──────────────┐
│assignment│     │quiz_quest│     │quiz_sessions │
│submissions│    │          │     │              │
│          │     │quiz_id   │     │ quiz_id      │
│ grade    │     │question  │     │ student_id   │
│ feedback │     │options   │     │ score        │
└──────────┘     │answer    │     │ cheat_count  │
                 └──────────┘     └──────┬───────┘
                                        │
                                        ▼
                                  ┌──────────────┐
                                  │ quiz_answers  │
                                  │              │
                                  │ session_id   │
                                  │ question_id  │
                                  │ selected_key │
                                  │ is_correct   │
                                  └──────────────┘
```

### Tabel Lengkap (15 tabel)

| #   | Tabel                    | Deskripsi                    | Relasi                          |
| :-- | :----------------------- | :--------------------------- | :------------------------------ |
| 1   | `schools`                | Data sekolah                 | —                               |
| 2   | `users`                  | Akun pengguna                | → schools                       |
| 3   | `user_profiles`          | Profil detail                | → users                         |
| 4   | `academic_years`         | Tahun ajaran                 | → schools                       |
| 5   | `classes`                | Kelas                        | → schools, academic_years       |
| 6   | `subjects`               | Mata pelajaran               | → schools                       |
| 7   | `class_subjects`         | Guru mengajar mapel di kelas | → classes, subjects, users      |
| 8   | `student_classes`        | Siswa terdaftar di kelas     | → users, classes                |
| 9   | `assignments`            | Tugas dari guru              | → classes, subjects, users      |
| 10  | `assignment_submissions` | Submit tugas siswa           | → assignments, users            |
| 11  | `quizzes`                | Kuis/ujian                   | → classes, subjects, users      |
| 12  | `quiz_questions`         | Soal kuis                    | → quizzes                       |
| 13  | `quiz_sessions`          | Sesi ujian siswa             | → quizzes, users                |
| 14  | `quiz_answers`           | Jawaban siswa                | → quiz_sessions, quiz_questions |
| 15  | `attendances`            | Presensi                     | → users, classes                |
| 16  | `grades`                 | Nilai                        | → users, subjects, classes      |
| 17  | `notifications`          | Notifikasi                   | → users                         |
| 18  | `audit_logs`             | Catatan audit                | —                               |
| 19  | `active_sessions`        | Sesi login aktif             | → users                         |

---

## 🔌 API Endpoints

### Auth (`/api/v1/auth`)

| Method | Endpoint           | Deskripsi               | Auth |
| :----- | :----------------- | :---------------------- | :--- |
| POST   | `/login`           | Login                   | ❌   |
| POST   | `/register`        | Register (invite-based) | ❌   |
| POST   | `/refresh`         | Refresh token           | ❌   |
| POST   | `/logout`          | Logout                  | ✅   |
| GET    | `/me`              | Current user info       | ✅   |
| PUT    | `/change-password` | Ubah password           | ✅   |

### Academic (`/api/v1/academic`) — Admin Only

| Method | Endpoint                | Deskripsi             |
| :----- | :---------------------- | :-------------------- |
| CRUD   | `/years`                | Tahun ajaran          |
| CRUD   | `/classes`              | Kelas                 |
| CRUD   | `/subjects`             | Mata pelajaran        |
| POST   | `/classes/:id/students` | Assign siswa ke kelas |
| POST   | `/classes/:id/subjects` | Assign guru ke mapel  |

### Assignments (`/api/v1/assignments`)

| Method   | Endpoint                 | Role  | Deskripsi         |
| :------- | :----------------------- | :---- | :---------------- |
| GET/POST | `/teacher`               | Guru  | List/Create tugas |
| POST     | `/:id/publish`           | Guru  | Publish tugas     |
| GET      | `/student`               | Siswa | List tugas        |
| POST     | `/:id/submit`            | Siswa | Submit tugas      |
| POST     | `/submissions/:id/grade` | Guru  | Beri nilai        |

### Quiz CBT (`/api/v1/quizzes`)

| Method   | Endpoint               | Role  | Deskripsi            |
| :------- | :--------------------- | :---- | :------------------- |
| GET/POST | `/teacher`             | Guru  | List/Create kuis     |
| POST     | `/:id/questions`       | Guru  | Tambah soal          |
| POST     | `/:id/publish`         | Guru  | Publish kuis         |
| POST     | `/:id/ai-generate`     | Guru  | Generate soal via AI |
| POST     | `/student/:id/start`   | Siswa | Mulai ujian          |
| POST     | `/sessions/:id/answer` | Siswa | Jawab soal           |
| POST     | `/sessions/:id/submit` | Siswa | Submit ujian         |

### Other Endpoints

| Module       | Endpoints                                              | Deskripsi          |
| :----------- | :----------------------------------------------------- | :----------------- |
| Attendance   | `POST/GET /attendance/teacher`, `GET /student/history` | Presensi           |
| Grade        | `POST/GET /grades/teacher`, `GET /student`             | Nilai              |
| Notification | `GET/POST /notifications`, `POST /:id/read`            | Notifikasi         |
| Upload       | `POST /upload` (multipart)                             | File upload        |
| Dashboard    | `GET /dashboard/stats`                                 | Analytics per role |
| Admin        | `CRUD /admin/users`                                    | User management    |

---

## 🚀 Quick Start

### Prasyarat

- **Go** 1.23+
- **Node.js** 20+ & **pnpm** 9+
- **Supabase** account (untuk database PostgreSQL)
- **Upstash** account (untuk Redis) — opsional

### 1. Clone Repository

```bash
git clone https://github.com/kazanaruishere-max/AkuBelajar.git
cd AkuBelajar
```

### 2. Setup Database

1. Buat project di [Supabase Dashboard](https://supabase.com/dashboard)
2. Buka **SQL Editor → New Query**
3. Copy-paste isi `backend/migrations/combined_all.sql`
4. Klik **Run**

### 3. Setup Backend

```bash
cd backend

# Copy dan edit environment variables
cp .env.example .env
# Edit .env → isi DATABASE_URL, PASETO_KEY, dll.

# Install dependencies
go mod tidy

# Jalankan seed data
go run cmd/seed/main.go

# Jalankan server
go run cmd/api/main.go
# → Server berjalan di http://localhost:8080
```

### 4. Setup Frontend

```bash
cd frontend

# Install dependencies
pnpm install

# Buat file .env.local
echo "NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1" > .env.local

# Jalankan dev server
pnpm dev
# → Frontend berjalan di http://localhost:3000
```

### 5. Buka Browser

```
http://localhost:3000
```

---

## ⚙️ Konfigurasi

### Backend `.env`

```env
# App
PORT=8080
APP_ENV=development
FRONTEND_URL=http://localhost:3000

# Database (Supabase)
DATABASE_URL=postgresql://postgres.[project-id]:[password]@aws-0-[region].pooler.supabase.com:6543/postgres?sslmode=require
DB_MAX_CONNS=10

# Redis (Upstash) — opsional
REDIS_URL=rediss://default:[password]@[region].upstash.io:6379

# Auth
PASETO_KEY=<generate-32-byte-hex>     # openssl rand -hex 32
ACCESS_TOKEN_EXPIRY_MIN=15
REFRESH_TOKEN_EXPIRY_DAY=7

# AI (opsional)
GEMINI_API_KEY=<your-gemini-api-key>
GEMINI_MODEL=gemini-2.0-flash

# Storage (opsional)
SUPABASE_URL=https://[project-id].supabase.co
SUPABASE_ANON_KEY=<your-anon-key>
SUPABASE_BUCKET=uploads

# CORS
CORS_ORIGIN=http://localhost:3000
```

### Frontend `.env.local`

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

---

## 👤 Akun Demo

Setelah menjalankan `go run cmd/seed/main.go`:

| Role           | Email                      | Password     |
| :------------- | :------------------------- | :----------- |
| 🔴 Super Admin | `admin@akubelajar.id`      | `Admin@123!` |
| 🟢 Guru 1      | `guru@akubelajar.id`       | `Guru@123!`  |
| 🟢 Guru 2      | `guru2@akubelajar.id`      | `Guru@123!`  |
| 🔵 Siswa 1     | `siswa@akubelajar.id`      | `Siswa@123!` |
| 🔵 Siswa 2-5   | `siswa{2-5}@akubelajar.id` | `Siswa@123!` |

### Data Seeded

- 🏫 **Sekolah**: SMP Nusantara Demo
- 📅 **Tahun Ajaran**: 2025/2026 (aktif)
- 📚 **Mata Pelajaran**: Matematika, Bahasa Indonesia, IPA, Bahasa Inggris
- 🏛️ **Kelas**: 7A, 8A, 9A
- 👨‍🎓 5 siswa → Kelas 7A
- 👨‍🏫 guru@akubelajar.id → Mengajar MTK & BIN di 7A

---

## 📅 Sprint Roadmap

| Sprint       | Status     | Fitur                                                           |
| :----------- | :--------- | :-------------------------------------------------------------- |
| **Sprint 1** | ✅ Done    | Auth, login, register, RBAC, Paseto, Argon2id                   |
| **Sprint 2** | ✅ Done    | Academic (tahun ajaran, kelas, mapel), Assignment CRUD          |
| **Sprint 3** | ✅ Done    | Quiz CBT + auto-grading, Attendance, Grade module               |
| **Sprint 4** | ✅ Done    | Notifications, File Upload (Supabase Storage), AI Quiz (Gemini) |
| **Sprint 5** | ✅ Done    | Dashboard analytics, Admin user management, Sidebar navigation  |
| **Sprint 6** | 🔲 Planned | Report cards, PDF generation, email notifications               |
| **Sprint 7** | 🔲 Planned | Mobile responsive, PWA, offline support                         |
| **Sprint 8** | 🔲 Planned | Multi-school, school admin panel, white-labeling                |

---

## 🔒 Keamanan

- **Paseto v4** — lebih aman dari JWT (tanpa `alg: none` vulnerability)
- **Argon2id** — password hashing terbaik (OWASP recommended)
- **Rate Limiting** — 120 req/min per IP via Redis
- **Brute-force Protection** — lock akun setelah 5 failed login (15 menit)
- **RBAC Middleware** — role-based access per endpoint
- **Input Validation** — validator/v10 pada semua request body
- **SQL Injection** — parameterized queries via pgx
- **CORS** — konfigurasi strict per origin
- **Soft Delete** — data tidak benar-benar dihapus

Untuk melaporkan vulnerability, lihat [SECURITY.md](SECURITY.md).

---

## 🤝 Kontribusi

Kami menerima kontribusi! Silakan baca [CONTRIBUTING.md](CONTRIBUTING.md) untuk panduan kontribusi dan [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) untuk kode etik komunitas.

---

## 📄 Lisensi

Proyek ini dilisensikan di bawah [MIT License](LICENSE).

---

<div align="center">

**AkuBelajar** — Dibuat dengan ❤️ untuk pendidikan Indonesia

_© 2026 AkuBelajar. All rights reserved._

</div>
