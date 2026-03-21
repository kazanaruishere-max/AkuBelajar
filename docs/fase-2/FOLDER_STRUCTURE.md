# 📁 Folder Structure — AkuBelajar

> Panduan definitif di mana setiap file harus diletakkan. Menghilangkan tebakan agar codebase tetap konsisten.

---

## Aturan Utama

1. **Group by domain/feature** — `quiz/handler.go` bukan `handlers/quiz.go`
2. **Flat is better than nested** — maksimal 4 level kedalaman
3. **Nama file = isinya** — `quiz_service.go` berisi quiz service

---

## Backend (Go)

```
akubelajar-backend/
├── cmd/
│   ├── api/main.go                 # Entry point API server
│   └── worker/main.go              # Entry point background worker
├── internal/                        # Kode private
│   ├── auth/                        # 🔐 Autentikasi & Otorisasi
│   │   ├── handler.go              #    HTTP handlers
│   │   ├── service.go              #    Business logic
│   │   ├── repository.go           #    Database queries
│   │   ├── model.go                #    Structs & types
│   │   └── auth_test.go            #    Tests
│   ├── academic/                    # 🎓 Kelas, siswa, tahun ajaran
│   ├── quiz/                        # 📝 Kuis & CBT
│   ├── assignment/                  # 📋 Tugas
│   ├── attendance/                  # ✅ Absensi
│   ├── grade/                       # 📊 Penilaian & Rapor
│   ├── ai/                          # 🤖 Gemini AI integration
│   │   ├── service.go, prompt.go, sanitizer.go
│   ├── notification/                # 🔔 Notifikasi
│   └── middleware/                  # 🛡️ Auth, RBAC, rate limit, logger
├── pkg/                             # Reusable packages
│   ├── database/postgres.go
│   ├── cache/redis.go
│   ├── security/hash.go, token.go
│   ├── validator/validator.go
│   └── response/response.go
├── migrations/                      # SQL migration files
├── config/config.go
├── Makefile, Dockerfile, .env.example
```

### Dimana Meletakkan File Baru (Backend)?

| Saya mau membuat... | Letakkan di... |
|:---|:---|
| Endpoint HTTP baru | `internal/{domain}/handler.go` |
| Business logic baru | `internal/{domain}/service.go` |
| Query database baru | `internal/{domain}/repository.go` |
| Struct request/response | `internal/{domain}/model.go` |
| Middleware baru | `internal/middleware/{nama}.go` |
| Utility reusable | `pkg/{kategori}/{nama}.go` |
| Migration database | `migrations/{timestamp}_{deskripsi}.up.sql` |

---

## Frontend (Next.js)

```
akubelajar-frontend/
├── app/                              # Next.js App Router
│   ├── (auth)/                       # Login/Register (no sidebar layout)
│   │   ├── login/page.tsx
│   │   └── register/page.tsx
│   ├── (dashboard)/                  # Dashboard (sidebar layout)
│   │   ├── admin/                    # Super Admin pages
│   │   │   ├── users/page.tsx
│   │   │   └── settings/page.tsx
│   │   ├── teacher/                  # Guru pages
│   │   │   ├── quizzes/page.tsx
│   │   │   ├── quizzes/create/page.tsx
│   │   │   └── attendance/page.tsx
│   │   ├── student/                  # Siswa pages
│   │   │   ├── quizzes/page.tsx
│   │   │   ├── grades/page.tsx
│   │   │   └── assignments/page.tsx
│   │   └── layout.tsx
│   ├── api/                          # BFF routes
│   ├── layout.tsx, page.tsx, globals.css
├── components/
│   ├── ui/                           # 🧩 Shadcn UI primitives
│   ├── shared/                       # 🔄 Cross-module (Navbar, Sidebar, etc.)
│   └── features/                     # 🎯 Feature-specific
│       ├── quiz/QuizCard.tsx, QuizForm.tsx, CBTInterface.tsx
│       ├── attendance/AttendanceForm.tsx
│       └── grade/GradeTable.tsx, GradeChart.tsx
├── lib/
│   ├── api/client.ts                 # Type-safe API client
│   ├── hooks/useAuth.ts, useQuiz.ts
│   ├── store/authStore.ts, uiStore.ts
│   └── utils/formatDate.ts, cn.ts
├── types/
│   ├── user.ts, quiz.ts, grade.ts, attendance.ts, api.ts
├── public/sw.js, manifest.json, icons/
```

### Dimana Meletakkan File Baru (Frontend)?

| Saya mau membuat... | Letakkan di... |
|:---|:---|
| Halaman baru | `app/(dashboard)/{role}/{fitur}/page.tsx` |
| Komponen UI primitif | `components/ui/{nama}.tsx` (via Shadcn CLI) |
| Komponen shared | `components/shared/{Nama}.tsx` |
| Komponen fitur spesifik | `components/features/{domain}/{Nama}.tsx` |
| Custom hook | `lib/hooks/use{Nama}.ts` |
| Zustand store | `lib/store/{nama}Store.ts` |
| TypeScript interface | `types/{domain}.ts` |

---

## Anti-Patterns

| ❌ Salah | ✅ Benar | Alasan |
|:---|:---|:---|
| `handlers/quiz_handler.go` | `internal/quiz/handler.go` | Group by domain |
| `components/QuizCard.tsx` (root) | `components/features/quiz/QuizCard.tsx` | Organize by feature |
| `utils/helpers.ts` (God file) | `lib/utils/formatDate.ts` (focused) | Single responsibility |

---

*Terakhir diperbarui: 21 Maret 2026*
