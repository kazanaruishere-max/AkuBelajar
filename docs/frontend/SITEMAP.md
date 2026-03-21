# 🗺️ Sitemap — AkuBelajar

> Semua halaman di aplikasi Next.js 15 App Router: path, role access, rendering strategy, dan data fetching.

---

## Public Routes (Tidak Perlu Login)

| Path | Halaman | Layout | Rendering | Data Fetching |
|:---|:---|:---|:---|:---|
| `/` | Landing Page | `MarketingLayout` | SSG | Static |
| `/login` | Login | `AuthLayout` | CSR | — |
| `/forgot-password` | Lupa Password | `AuthLayout` | CSR | — |
| `/reset-password/[token]` | Reset Password | `AuthLayout` | CSR | Token validation |
| `/invite/[token]` | Claim Invite | `AuthLayout` | CSR | Token validation |

---

## Onboarding Routes (Login, Belum Lengkap)

| Path | Halaman | Role | Rendering |
|:---|:---|:---|:---|
| `/onboarding/change-password` | Ganti Password | Semua (first login) | CSR |
| `/onboarding/complete-profile` | Lengkapi Profil | Semua (first login) | CSR |
| `/onboarding/select-class` | Pilih Kelas | Siswa only | CSR |

**Middleware:** `onboardingMiddleware` — block semua route lain jika `is_first_login = true`

---

## Shared Routes (Semua Role yang Sudah Login)

| Path | Halaman | Rendering | Data Fetching |
|:---|:---|:---|:---|
| `/dashboard` | Dashboard (redirect by role) | SSR | Server Component |
| `/notifications` | Daftar Notifikasi | CSR | react-query |
| `/settings` | Settings Overview | CSR | — |
| `/settings/profile` | Edit Profil | CSR | react-query |
| `/settings/security` | Keamanan (password, sessions) | CSR | react-query |
| `/settings/notifications` | Preferensi Notifikasi | CSR | react-query |

---

## Admin Routes (`super_admin` only)

| Path | Halaman | Rendering | Data Fetching |
|:---|:---|:---|:---|
| `/admin/users` | Daftar User | SSR | Server Component + pagination |
| `/admin/users/[id]` | Detail User | SSR | Server Component |
| `/admin/users/import` | Bulk Import Excel | CSR | File upload |
| `/admin/schools` | Pengaturan Sekolah | SSR | Server Component |
| `/admin/academic-years` | Tahun Ajaran | SSR | Server Component |
| `/admin/classes` | Daftar Kelas | SSR | Server Component |
| `/admin/subjects` | Daftar Mapel | SSR | Server Component |
| `/admin/invite-tokens` | Kelola Invite Token | CSR | react-query + realtime |
| `/admin/audit-logs` | Audit Log | SSR | Server Component + pagination |
| `/admin/system` | System Settings | CSR | react-query |

---

## Teacher Routes (`teacher` only)

| Path | Halaman | Rendering | Data Fetching |
|:---|:---|:---|:---|
| `/teacher/dashboard` | Dashboard Guru | SSR | Server Component |
| `/teacher/attendance/[class_id]` | Input Absensi | CSR | react-query + WebSocket |
| `/teacher/assignments` | Daftar Tugas | SSR | Server Component |
| `/teacher/assignments/new` | Buat Tugas Baru | CSR | react-hook-form |
| `/teacher/assignments/[id]` | Detail Tugas | SSR | Server Component |
| `/teacher/assignments/[id]/submissions` | Daftar Submission | SSR | Server Component |
| `/teacher/quizzes` | Daftar Kuis | SSR | Server Component |
| `/teacher/quizzes/new` | Buat Kuis | CSR | react-hook-form + AI |
| `/teacher/quizzes/[id]` | Detail Kuis | SSR | Server Component |
| `/teacher/quizzes/[id]/results` | Hasil Kuis | SSR | Server Component |
| `/teacher/grades/[class_id]` | Buku Nilai | CSR | react-query |
| `/teacher/report-cards/[class_id]` | Rapor | CSR | react-query |

---

## Student Routes (`student`, `class_leader`)

| Path | Halaman | Rendering | Data Fetching |
|:---|:---|:---|:---|
| `/student/dashboard` | Dashboard Siswa | SSR | Server Component |
| `/student/assignments` | Daftar Tugas | SSR | Server Component |
| `/student/assignments/[id]` | Detail + Submit | CSR | react-query + file upload |
| `/student/quizzes` | Daftar Kuis | SSR | Server Component |
| `/student/quizzes/[id]/take` | **CBT Interface** (fullscreen) | CSR | WebSocket + Zustand |
| `/student/quizzes/[id]/review` | Review Jawaban | SSR | Server Component |
| `/student/grades` | Nilai Saya | SSR | Server Component (RLS) |
| `/student/attendance` | Absensi Saya | SSR | Server Component (RLS) |
| `/student/report-cards` | Rapor Saya | SSR | Server Component (RLS) |

---

## Middleware Stack

```typescript
// middleware.ts
export function middleware(request: NextRequest) {
  const path = request.nextUrl.pathname;
  
  // 1. Public routes — skip auth check
  if (isPublicRoute(path)) return NextResponse.next();
  
  // 2. Auth check — redirect to /login if no token
  const token = getAccessToken(request);
  if (!token) return redirect('/login');
  
  // 3. Onboarding check — block all unless onboarding routes
  if (user.is_first_login && !isOnboardingRoute(path))
    return redirect('/onboarding/change-password');
  
  // 4. Role check — verify user.role has access to this route
  if (!hasAccess(user.role, path))
    return redirect('/dashboard');
}
```

---

## Layout Hierarchy

```
RootLayout
├── MarketingLayout (/)
├── AuthLayout (/login, /forgot-password, ...)
├── OnboardingLayout (/onboarding/*)
└── DashboardLayout (authenticated routes)
    ├── AdminLayout (/admin/*)
    ├── TeacherLayout (/teacher/*)
    └── StudentLayout (/student/*)
```

---

*Terakhir diperbarui: 21 Maret 2026*
