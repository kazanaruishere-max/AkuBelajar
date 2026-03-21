# 📡 API Specification — AkuBelajar

> Kontrak API yang konkret: path, HTTP method, request body, response shape, status code, dan contoh nyata. Ini adalah **sumber kebenaran** agar Go backend dan Next.js frontend tetap sinkron.

---

## Base URL

| Environment | Base URL |
|:---|:---|
| Development | `http://localhost:8080/api/v1` |
| Staging | `https://staging-api.akubelajar.id/api/v1` |
| Production | `https://api.akubelajar.id/api/v1` |

---

## Konvensi Umum

### Headers (Semua Request)

```
Content-Type: application/json
Authorization: Bearer <access_token>    // Wajib kecuali login/register
X-Request-ID: <uuid>                    // Auto-generated oleh client
```

### Response Envelope

```json
// Success
{
  "success": true,
  "data": { ... },
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 150,
    "total_pages": 8
  }
}

// Error
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Email is required",
    "details": [
      { "field": "email", "message": "must not be empty" }
    ]
  }
}
```

### Status Codes

| Code | Arti | Kapan Digunakan |
|:---|:---|:---|
| `200` | OK | GET berhasil, UPDATE berhasil |
| `201` | Created | POST berhasil membuat resource baru |
| `204` | No Content | DELETE berhasil |
| `400` | Bad Request | Validasi input gagal |
| `401` | Unauthorized | Token tidak ada atau expired |
| `403` | Forbidden | Role tidak memiliki izin |
| `404` | Not Found | Resource tidak ditemukan |
| `409` | Conflict | Duplikat data (email, NISN) |
| `422` | Unprocessable Entity | Request valid tapi logic gagal |
| `429` | Too Many Requests | Rate limit terlampaui |
| `500` | Internal Server Error | Bug di server |

---

## Auth Endpoints

### `POST /auth/login`

Login dan dapatkan access + refresh token.

**Request:**
```json
{
  "email": "guru@akubelajar.id",
  "password": "Guru@123!"
}
```

**Response `200`:**
```json
{
  "success": true,
  "data": {
    "access_token": "v4.public.eyJ...",
    "refresh_token": "v4.public.eyJ...",
    "token_type": "Bearer",
    "expires_in": 900,
    "user": {
      "id": "019516a2-7c1e-7a3b-8d2f-1a2b3c4d5e6f",
      "email": "guru@akubelajar.id",
      "role": "teacher",
      "name": "Budi Santoso",
      "school_id": "019516a2-1111-7a3b-8d2f-aabbccddeeff"
    }
  }
}
```

**Response `401`:**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Email or password is incorrect"
  }
}
```

**Response `429`:**
```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMITED",
    "message": "Too many login attempts. Try again in 15 minutes.",
    "retry_after": 900
  }
}
```

---

### `POST /auth/refresh`

Refresh access token menggunakan refresh token.

**Request:**
```json
{
  "refresh_token": "v4.public.eyJ..."
}
```

**Response `200`:**
```json
{
  "success": true,
  "data": {
    "access_token": "v4.public.eyJ...(new)",
    "refresh_token": "v4.public.eyJ...(new, rotated)",
    "expires_in": 900
  }
}
```

---

### `POST /auth/logout`

Invalidasi refresh token.

**Request:**
```json
{
  "refresh_token": "v4.public.eyJ..."
}
```

**Response `204`:** _(No Content)_

---

## User Endpoints

### `GET /users` — List Users

> **Role:** `super_admin` only

**Query Params:**

| Param | Type | Default | Deskripsi |
|:---|:---|:---|:---|
| `page` | int | 1 | Halaman |
| `per_page` | int | 20 | Item per halaman (max 100) |
| `role` | string | — | Filter by role |
| `search` | string | — | Search by name/email |

**Response `200`:**
```json
{
  "success": true,
  "data": [
    {
      "id": "019516a2-...",
      "email": "guru@akubelajar.id",
      "name": "Budi Santoso",
      "role": "teacher",
      "is_active": true,
      "created_at": "2026-03-01T10:00:00Z"
    }
  ],
  "meta": { "page": 1, "per_page": 20, "total": 45, "total_pages": 3 }
}
```

---

### `GET /users/:id` — Get User Detail

> **Role:** `super_admin`, atau user sendiri

**Response `200`:**
```json
{
  "success": true,
  "data": {
    "id": "019516a2-...",
    "email": "guru@akubelajar.id",
    "name": "Budi Santoso",
    "role": "teacher",
    "school_id": "019516a2-1111-...",
    "is_active": true,
    "profile": {
      "nip": "198501012010011001",
      "phone": "+6281234567890",
      "subjects": ["Matematika", "Fisika"]
    },
    "last_login_at": "2026-03-21T08:30:00Z",
    "created_at": "2026-03-01T10:00:00Z"
  }
}
```

---

## Quiz Endpoints

### `POST /quizzes` — Create Quiz

> **Role:** `teacher`

**Request:**
```json
{
  "title": "Ujian Tengah Semester - Biologi Kelas 10",
  "subject_id": "019516a2-2222-...",
  "class_id": "019516a2-3333-...",
  "time_limit": 60,
  "start_at": "2026-04-01T08:00:00+07:00",
  "end_at": "2026-04-01T10:00:00+07:00",
  "ai_generate": true,
  "ai_config": {
    "topic": "Struktur dan Fungsi Sel",
    "question_count": 20,
    "difficulty": "mixed",
    "curriculum": "merdeka"
  }
}
```

**Response `201`:**
```json
{
  "success": true,
  "data": {
    "id": "019516a2-4444-...",
    "title": "Ujian Tengah Semester - Biologi Kelas 10",
    "status": "draft",
    "questions_count": 20,
    "ai_generated": true,
    "created_at": "2026-03-21T10:00:00Z"
  }
}
```

---

### `POST /quizzes/:id/start` — Start Quiz Session

> **Role:** `student`

**Response `200`:**
```json
{
  "success": true,
  "data": {
    "session_id": "019516a2-5555-...",
    "quiz_id": "019516a2-4444-...",
    "questions": [
      {
        "id": "019516a2-6666-...",
        "number": 1,
        "question": "Organel yang berfungsi sebagai pusat kontrol sel adalah...",
        "options": [
          { "key": "A", "text": "Ribosom" },
          { "key": "B", "text": "Nukleus" },
          { "key": "C", "text": "Mitokondria" },
          { "key": "D", "text": "Vakuola" }
        ]
      }
    ],
    "time_limit": 60,
    "server_time": "2026-04-01T08:00:00Z",
    "expires_at": "2026-04-01T09:00:00Z"
  }
}
```

---

### `POST /quizzes/:id/submit` — Submit Quiz Answers

> **Role:** `student`

**Request:**
```json
{
  "session_id": "019516a2-5555-...",
  "answers": [
    { "question_id": "019516a2-6666-...", "selected": "B" },
    { "question_id": "019516a2-7777-...", "selected": "A" }
  ]
}
```

**Response `200`:**
```json
{
  "success": true,
  "data": {
    "submission_id": "019516a2-8888-...",
    "score": 85,
    "correct": 17,
    "incorrect": 3,
    "total": 20,
    "time_taken": 2847,
    "submitted_at": "2026-04-01T08:47:27Z"
  }
}
```

---

## Attendance Endpoints

### `POST /attendances` — Record Attendance

> **Role:** `teacher`, `class_leader`

**Request:**
```json
{
  "class_id": "019516a2-3333-...",
  "date": "2026-03-21",
  "records": [
    { "student_id": "019516a2-aaaa-...", "status": "present" },
    { "student_id": "019516a2-bbbb-...", "status": "sick", "note": "Demam" },
    { "student_id": "019516a2-cccc-...", "status": "absent" }
  ]
}
```

**Response `201`:**
```json
{
  "success": true,
  "data": {
    "class_id": "019516a2-3333-...",
    "date": "2026-03-21",
    "summary": { "present": 28, "permission": 1, "sick": 1, "absent": 2 },
    "recorded_by": "019516a2-...",
    "recorded_at": "2026-03-21T07:30:00Z"
  }
}
```

---

## Grade Endpoints

### `GET /grades` — Get Student Grades

> **Role:** `student` (own only via RLS), `teacher` (class only)

**Query Params:**

| Param | Type | Deskripsi |
|:---|:---|:---|
| `student_id` | uuid | Filter by student (teacher only) |
| `subject_id` | uuid | Filter by subject |
| `academic_year_id` | uuid | Filter by academic year |

**Response `200`:**
```json
{
  "success": true,
  "data": [
    {
      "id": "019516a2-dddd-...",
      "subject": { "id": "...", "name": "Biologi" },
      "assignment_score": 82,
      "quiz_score": 88,
      "final_score": 84.4,
      "grade_letter": "A-",
      "is_locked": false
    }
  ]
}
```

---

## Pagination Convention

Semua list endpoint menggunakan **offset-based pagination**:

```
GET /users?page=2&per_page=20
```

Untuk endpoint yang membutuhkan performa tinggi (audit log, notifications), gunakan **cursor-based pagination**:

```
GET /audit-logs?cursor=019516a2-eeee-...&limit=50
```

---

## Rate Limits

| Endpoint | Limit | Window |
|:---|:---|:---|
| `POST /auth/login` | 5 requests | 15 menit |
| `POST /auth/refresh` | 10 requests | 15 menit |
| `POST /quizzes` (AI generate) | 10 requests | 1 jam |
| Semua GET endpoints | 100 requests | 1 menit |
| Semua POST/PUT/DELETE | 30 requests | 1 menit |

---

## Referensi Terkait

- [Data Models](DATA_MODELS.md) — Mapping struct Go ↔ TypeScript types
- [System Overview](SYSTEM_OVERVIEW.md) — Arsitektur keseluruhan
- [Backend Guide](../fase-2/BACKEND_GUIDE.md) — Pattern handler/service/repo

---

*Terakhir diperbarui: 21 Maret 2026*
