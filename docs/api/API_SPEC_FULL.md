# 📡 API Specification Full — AkuBelajar

> Kontrak API lengkap semua modul. Sumber kebenaran tunggal antara Go backend dan Next.js frontend.

---

## Konvensi Global

- **Base URL:** `/api/v1/`
- **Auth:** `Authorization: Bearer <paseto_token>`
- **Content-Type:** `application/json`
- **Success:** `{ "data": {...}, "meta": {...} }`
- **Error:** `{ "error": { "code": "string", "message": "string", "details": [...] } }`

---

## AUTH MODULE

### `POST /auth/login`
**Auth:** Tidak  
**Request:** `{ "email": "guru@akubelajar.id", "password": "Guru@123!" }`  
**Response 200:** `{ "data": { "access_token", "refresh_token", "expires_in": 900, "user": { id, email, role, school_id, is_first_login } } }`  
**Errors:** `AUTH_001` (401), `AUTH_002` (429), `AUTH_006` (403 suspended)

### `POST /auth/refresh`
**Auth:** Tidak (refresh token di body)  
**Request:** `{ "refresh_token": "v4.public.eyJ..." }`  
**Response 200:** `{ "data": { "access_token", "refresh_token" (rotated), "expires_in" } }`  
**Errors:** `AUTH_003` (401 expired), `AUTH_004` (401 invalid)

### `POST /auth/logout`
**Auth:** Ya  
**Request:** `{ "refresh_token": "v4.public.eyJ..." }`  
**Response 204:** No Content

### `POST /auth/password-reset/request`
**Auth:** Tidak  
**Request:** `{ "identifier": "guru@akubelajar.id" }` (email atau nomor WA)  
**Response 200:** `{ "data": { "message": "Jika akun terdaftar, OTP telah dikirim" } }` (pesan generik anti-enumeration)  
**Errors:** `AUTH_011` (429 rate limited)

### `POST /auth/password-reset/verify`
**Auth:** Tidak  
**Request:** `{ "identifier": "guru@akubelajar.id", "otp": "123456", "new_password": "NewPass@123" }`  
**Response 200:** `{ "data": { "message": "Password berhasil diubah" } }`  
**Errors:** `AUTH_009` (400 OTP invalid), `AUTH_010` (400 OTP expired)

### `POST /auth/change-password`
**Auth:** Ya  
**Request:** `{ "current_password": "Old@123", "new_password": "New@456" }`  
**Response 200:** `{ "data": { "message": "Password berhasil diubah" } }`  
**Errors:** `AUTH_001` (401 current password salah), `AUTH_008` (400 password change required)

---

## USER MODULE

### `GET /users/me`
**Auth:** Ya | **Role:** Semua  
**Response 200:** `{ "data": { id, email, role, school_id, is_first_login, profile: { nisn, nip, phone_wa, photo_url, ... }, last_login_at } }`

### `PUT /users/me`
**Auth:** Ya | **Role:** Semua  
**Request:** `{ "name": "...", "phone_wa": "+6281234567890", "birth_date": "2010-05-15" }`  
**Response 200:** Updated user data

### `PUT /users/me/avatar`
**Auth:** Ya | **Role:** Semua  
**Request:** `multipart/form-data` — field `avatar` (max 2MB, JPG/PNG/WebP)  
**Response 200:** `{ "data": { "photo_url": "https://..." } }`

### `GET /users` (admin only)
**Auth:** Ya | **Role:** `super_admin`  
**Query:** `?page=1&per_page=20&role=student&search=budi`  
**Response 200:** Array of users + meta pagination

### `POST /users` (admin create single)
**Auth:** Ya | **Role:** `super_admin`  
**Request:** `{ "email", "name", "role", "class_id" (optional) }`  
**Response 201:** Created user (password auto-generated, returned once)

### `PUT /users/:id/role`
**Auth:** Ya | **Role:** `super_admin`, `teacher`  
**Request:** `{ "role": "class_leader" }`  
**Response 200:** Updated user  
**Errors:** `USER_004` (400 invalid role assignment)

### `PUT /users/:id/status`
**Auth:** Ya | **Role:** `super_admin`  
**Request:** `{ "is_active": false }`  
**Response 200:** Updated user

### `DELETE /users/:id`
**Auth:** Ya | **Role:** `super_admin`  
**Response 204:** Soft delete  
**Errors:** `USER_005` (400 cannot delete own account)

---

## INVITE TOKEN MODULE

### `POST /invite-tokens`
**Auth:** Ya | **Role:** `super_admin`, `teacher`  
**Request:** `{ "role": "student", "class_id": "uuid", "max_uses": 30, "expires_in_hours": 72 }`  
**Response 201:** `{ "data": { "id", "token": "ABcd1234", "url": "https://app.akubelajar.id/invite/ABcd1234", "expires_at" } }`

### `POST /invite-tokens/claim`
**Auth:** Tidak  
**Request:** `{ "token": "ABcd1234", "email": "siswa@...", "password": "...", "name": "..." }`  
**Response 201:** Created user

### `GET /invite-tokens`
**Auth:** Ya | **Role:** `super_admin`  
**Response 200:** Array invite tokens with usage stats

### `DELETE /invite-tokens/:id`
**Auth:** Ya | **Role:** `super_admin`  
**Response 204:** Token revoked

---

## ATTENDANCE MODULE

### `POST /attendances`
**Auth:** Ya | **Role:** `teacher`, `class_leader`  
**Request:**
```json
{
  "class_id": "uuid", "subject_id": "uuid", "date": "2026-03-21",
  "records": [
    { "student_id": "uuid", "status": "present" },
    { "student_id": "uuid", "status": "sick", "reason": "Demam" }
  ]
}
```
**Response 201:** Summary (counts per status)

### `GET /attendances?class_id=&date=`
**Auth:** Ya | **Role:** `teacher`  
**Response 200:** Array attendance records for the class

### `PUT /attendances/:id`
**Auth:** Ya | **Role:** `teacher`, `super_admin`  
**Request:** `{ "status": "permission", "reason": "Koreksi: surat izin terlambat masuk" }`  
**Response 200:** Updated record  
**Errors:** `ATT_003` (403 edit window expired > 24h without reason)

### `GET /attendances/my`
**Auth:** Ya | **Role:** `student`  
**Response 200:** Own attendance records (RLS enforced)

### `GET /attendances/summary?student_id=&month=`
**Auth:** Ya | **Role:** `teacher`, `student` (own only)  
**Response 200:** `{ "data": { "total_days", "present", "permission", "sick", "absent", "late", "percentage" } }`

---

## ASSIGNMENT MODULE

### `POST /assignments`
**Auth:** Ya | **Role:** `teacher`  
**Request:** `{ "class_id", "subject_id", "title", "description", "deadline_at", "allow_late", "late_penalty_pct", "weight_pct" }`  
**Response 201:** Created assignment

### `GET /assignments?class_id=`
**Auth:** Ya | **Role:** `teacher`, `student`  
**Response 200:** Array assignments (student sees only published)

### `PUT /assignments/:id`
**Auth:** Ya | **Role:** `teacher` (owner only)  
**Response 200:** Updated assignment

### `DELETE /assignments/:id`
**Auth:** Ya | **Role:** `teacher` (owner only)  
**Response 204:** Soft delete  
**Errors:** `ASSIGN_005` (400 has submissions — archive instead)

### `POST /assignments/:id/submissions`
**Auth:** Ya | **Role:** `student`  
**Request:** `multipart/form-data` — files (max 3, max 20MB each)  
**Response 201:** `{ "data": { "id", "is_late", "late_days", "status": "submitted" } }`  
**Errors:** `ASSIGN_002` (400 deadline passed), `ASSIGN_003` (400 not allowed)

### `GET /assignments/:id/submissions`
**Auth:** Ya | **Role:** `teacher`  
**Response 200:** Array submissions with student info

### `PUT /assignments/:id/submissions/:sub_id/grade`
**Auth:** Ya | **Role:** `teacher`  
**Request:** `{ "grade": 85, "feedback": "Bagus, perhatikan tata bahasa." }`  
**Response 200:** Graded submission (auto-calculates penalty if late)

---

## QUIZ MODULE

### `POST /quizzes`
**Auth:** Ya | **Role:** `teacher`  
**Request:** `{ "class_id", "subject_id", "title", "time_limit", "start_at", "end_at", "randomize_questions": true }`  
**Response 201:** Created quiz (status: draft)

### `POST /quizzes/generate-ai`
**Auth:** Ya | **Role:** `teacher`  
**Request:** `{ "quiz_id": "uuid", "topic": "Struktur Sel", "question_count": 20, "difficulty": "mixed" }`  
**Response 200:** `{ "data": { "questions_generated": 20, "quiz_id" } }`  
**Errors:** `QUIZ_007` (503 AI generation failed)

### `PUT /quizzes/:id/publish`
**Auth:** Ya | **Role:** `teacher`  
**Response 200:** Status → published, notifications sent to students

### `POST /quizzes/:id/sessions`
**Auth:** Ya | **Role:** `student`  
**Response 200:** `{ "data": { "session_id", "questions" (shuffled), "expires_at", "server_time" } }`  
**Errors:** `QUIZ_002` (400 not in time window), `QUIZ_004` (409 already submitted)

### `POST /quizzes/:id/sessions/submit`
**Auth:** Ya | **Role:** `student`  
**Request:** `{ "session_id": "uuid", "answers": [{ "question_id", "selected_key": "B" }] }`  
**Response 200:** `{ "data": { "score", "correct", "incorrect", "total", "time_taken" } }`

### `GET /quizzes/:id/results`
**Auth:** Ya | **Role:** `teacher`  
**Response 200:** Array per-student results + statistics

---

## GRADE MODULE

### `GET /grades?student_id=&subject_id=&semester=`
**Auth:** Ya | **Role:** `teacher`, `student` (own only via RLS)  
**Response 200:** Array grades with calculated final_score

### `GET /grades/report-card/:student_id`
**Auth:** Ya | **Role:** `teacher`, `student` (own only), `super_admin`  
**Response 200:** Full report card data

### `POST /grades/report-card/:student_id/generate`
**Auth:** Ya | **Role:** `super_admin`  
**Response 200:** `{ "data": { "pdf_url" (signed, 15 min TTL), "qr_code" } }`

---

## NOTIFICATION MODULE

### `GET /notifications?page=&per_page=`
**Auth:** Ya | **Role:** Semua  
**Response 200:** Array notifications (newest first)

### `PUT /notifications/:id/read`
**Auth:** Ya  
**Response 200:** Marked as read

### `PUT /notifications/read-all`
**Auth:** Ya  
**Response 200:** All marked as read

### `GET /notifications/preferences`
**Auth:** Ya  
**Response 200:** `{ "data": { "email_enabled", "wa_enabled", "in_app_enabled", "quiet_start", "quiet_end" } }`

### `PUT /notifications/preferences`
**Auth:** Ya  
**Request:** `{ "wa_enabled": false, "quiet_start": "22:00", "quiet_end": "06:00" }`  
**Response 200:** Updated preferences

---

## Rate Limits Summary

| Endpoint Group | Limit |
|:---|:---|
| `POST /auth/login` | 5/min per IP |
| `POST /auth/password-reset/*` | 3/hour per IP |
| `POST /quizzes/generate-ai` | 10/hour per user |
| `POST */submissions` (file upload) | 20/hour per user |
| `GET /*` (general) | 120/min per user |
| `POST/PUT/DELETE /*` (general) | 30/min per user |

---

*Terakhir diperbarui: 21 Maret 2026*
