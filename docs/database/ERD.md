# 📊 Entity Relationship Diagram — AkuBelajar

> Visualisasi relasi antar semua tabel di database AkuBelajar menggunakan Mermaid erDiagram.

---

## Bagian 1 — Core & User Management

```mermaid
erDiagram
    schools ||--o{ users : "has"
    schools ||--o{ academic_years : "has"
    schools ||--o{ classes : "has"
    schools ||--o{ subjects : "has"

    users ||--o| user_profiles : "has"
    users ||--o{ student_classes : "enrolled_in"

    academic_years ||--o{ classes : "contains"
    academic_years ||--o{ student_classes : "for"

    classes ||--o{ class_subjects : "teaches"
    classes ||--o{ student_classes : "has"
    classes }o--|| users : "homeroom_teacher"

    subjects ||--o{ class_subjects : "taught_in"
    class_subjects }o--|| users : "teacher"

    schools {
        UUID id PK
        VARCHAR name
        VARCHAR code UK
        JSONB config
        BOOLEAN is_active
    }

    users {
        UUID id PK
        UUID school_id FK
        VARCHAR email UK
        VARCHAR password
        user_role role
        BOOLEAN is_first_login
        BOOLEAN is_active
    }

    user_profiles {
        UUID id PK
        UUID user_id FK
        VARCHAR nisn
        VARCHAR nip
        VARCHAR phone_wa
        VARCHAR photo_url
    }

    academic_years {
        UUID id PK
        UUID school_id FK
        VARCHAR name
        DATE start_date
        DATE end_date
        BOOLEAN is_active
    }

    classes {
        UUID id PK
        UUID school_id FK
        UUID academic_year_id FK
        VARCHAR name
        INTEGER grade_level
        UUID homeroom_teacher_id FK
    }

    subjects {
        UUID id PK
        UUID school_id FK
        VARCHAR name
        VARCHAR code
    }

    class_subjects {
        UUID id PK
        UUID class_id FK
        UUID subject_id FK
        UUID teacher_id FK
        JSONB schedule_json
    }

    student_classes {
        UUID id PK
        UUID student_id FK
        UUID class_id FK
        UUID academic_year_id FK
    }
```

### Relasi yang Perlu Penjelasan

| Relasi | Penjelasan |
|:---|:---|
| `schools → users` | Semua user terkait dengan satu sekolah (multi-tenant) |
| `users → user_profiles` | 1:1 — profil terpisah agar tabel users tetap ringan |
| `classes → class_subjects` | Satu kelas bisa diajar beberapa mata pelajaran oleh guru berbeda |
| `student_classes` | Junction table — siswa bisa pindah kelas antar tahun ajaran |

---

## Bagian 2 — Attendance & Assignment

```mermaid
erDiagram
    users ||--o{ attendances : "has"
    classes ||--o{ attendances : "for"
    subjects ||--o{ attendances : "for"
    users ||--o{ attendances : "recorded_by"

    classes ||--o{ assignments : "for"
    subjects ||--o{ assignments : "for"
    users ||--o{ assignments : "created_by"

    assignments ||--o{ assignment_attachments : "has"
    assignments ||--o{ submissions : "receives"
    users ||--o{ submissions : "submits"
    submissions ||--o{ submission_files : "has"

    attendances {
        UUID id PK
        UUID student_id FK
        UUID class_id FK
        UUID subject_id FK
        DATE date
        attendance_status status
        TEXT reason
        UUID recorded_by FK
        BOOLEAN is_draft
    }

    assignments {
        UUID id PK
        UUID class_id FK
        UUID subject_id FK
        UUID teacher_id FK
        VARCHAR title
        TIMESTAMPTZ deadline_at
        BOOLEAN allow_late
        INTEGER late_penalty_pct
        assignment_status status
    }

    assignment_attachments {
        UUID id PK
        UUID assignment_id FK
        VARCHAR file_url
        VARCHAR file_name
        BIGINT file_size
    }

    submissions {
        UUID id PK
        UUID assignment_id FK
        UUID student_id FK
        BOOLEAN is_late
        INTEGER grade
        INTEGER grade_after_penalty
        TEXT feedback
        submission_status status
    }

    submission_files {
        UUID id PK
        UUID submission_id FK
        VARCHAR file_url
        BIGINT file_size
    }
```

### Relasi yang Perlu Penjelasan

| Relasi | Penjelasan |
|:---|:---|
| `attendances.recorded_by` | Bisa guru atau ketua kelas — tergantung role |
| `attendances.is_draft` | TRUE jika diinput ketua kelas (belum diapprove guru) |
| `submissions` unique constraint | Satu siswa hanya bisa submit sekali per tugas |
| `grade_after_penalty` | Dihitung otomatis: `grade × (1 - late_penalty_pct × late_days / 100)` |

---

## Bagian 3 — Quiz & Grades

```mermaid
erDiagram
    classes ||--o{ quizzes : "for"
    subjects ||--o{ quizzes : "for"
    users ||--o{ quizzes : "created_by"

    quizzes ||--o{ quiz_questions : "contains"
    quizzes ||--o{ quiz_sessions : "has"
    users ||--o{ quiz_sessions : "takes"

    quiz_sessions ||--o{ quiz_answers : "contains"
    quiz_questions ||--o{ quiz_answers : "answered_by"

    users ||--o{ grades : "has"
    subjects ||--o{ grades : "for"
    academic_years ||--o{ grades : "for"

    users ||--o{ report_cards : "has"
    academic_years ||--o{ report_cards : "for"

    quizzes {
        UUID id PK
        UUID class_id FK
        UUID subject_id FK
        UUID teacher_id FK
        VARCHAR title
        INTEGER time_limit
        BOOLEAN randomize_questions
        INTEGER max_attempts
        review_mode review_mode
        quiz_status status
    }

    quiz_questions {
        UUID id PK
        UUID quiz_id FK
        TEXT question_text
        question_type question_type
        JSONB options
        VARCHAR answer_hash
        INTEGER order_num
    }

    quiz_sessions {
        UUID id PK
        UUID quiz_id FK
        UUID student_id FK
        TIMESTAMPTZ expires_at
        session_status status
        INTEGER score
        INTEGER cheat_count
    }

    quiz_answers {
        UUID id PK
        UUID session_id FK
        UUID question_id FK
        VARCHAR selected_key
        BOOLEAN is_correct
    }

    grades {
        UUID id PK
        UUID student_id FK
        UUID subject_id FK
        UUID academic_year_id FK
        INTEGER semester
        NUMERIC final_score
        VARCHAR grade_letter
        BOOLEAN is_locked
    }

    report_cards {
        UUID id PK
        UUID student_id FK
        UUID academic_year_id FK
        INTEGER semester
        VARCHAR pdf_url
        VARCHAR qr_verification_code
        BOOLEAN is_published
    }
```

### Relasi yang Perlu Penjelasan

| Relasi | Penjelasan |
|:---|:---|
| `quiz_questions.answer_hash` | Jawaban benar di-hash Argon2id — tidak pernah plain text |
| `quiz_sessions.question_order` | Array integer untuk urutan soal yang sudah diacak per siswa |
| `grades` unique constraint | Satu siswa × satu mapel × satu tahun ajaran × satu semester |
| `report_cards.qr_verification_code` | QR code unik untuk verifikasi keaslian rapor online |

---

## Bagian 4 — Notifications & System

```mermaid
erDiagram
    users ||--o{ notifications : "receives"
    users ||--o| notification_preferences : "has"
    users ||--o{ active_sessions : "has"
    users ||--o{ password_reset_tokens : "requests"
    users ||--o{ password_histories : "has"

    schools ||--o{ invite_tokens : "generates"
    users ||--o{ invite_tokens : "created_by"
    users ||--o{ audit_logs : "performs"

    notifications {
        UUID id PK
        UUID user_id FK
        VARCHAR type
        VARCHAR title
        TEXT body
        BOOLEAN is_read
    }

    notification_preferences {
        UUID id PK
        UUID user_id FK
        BOOLEAN email_enabled
        BOOLEAN wa_enabled
        TIME quiet_start
        TIME quiet_end
    }

    invite_tokens {
        UUID id PK
        UUID school_id FK
        UUID created_by FK
        VARCHAR token_hash
        user_role role
        INTEGER max_uses
        INTEGER uses_count
        TIMESTAMPTZ expires_at
    }

    password_reset_tokens {
        UUID id PK
        UUID user_id FK
        VARCHAR token_hash
        TIMESTAMPTZ expires_at
        TIMESTAMPTZ used_at
    }

    active_sessions {
        UUID id PK
        UUID user_id FK
        VARCHAR refresh_token_hash
        JSONB device_info
        BOOLEAN is_remember_me
        TIMESTAMPTZ expires_at
    }

    audit_logs {
        UUID id PK
        UUID actor_id
        VARCHAR action
        VARCHAR entity_type
        UUID entity_id
        JSONB old_value
        JSONB new_value
    }
```

---

## Ringkasan Relasi

| Tabel | Relasi Ke |
|:---|:---|
| `schools` | users, academic_years, classes, subjects, invite_tokens |
| `users` | user_profiles, student_classes, attendances, assignments, submissions, quiz_sessions, grades, report_cards, notifications, active_sessions |
| `classes` | class_subjects, student_classes, attendances, assignments, quizzes |
| `quizzes` | quiz_questions, quiz_sessions |
| `assignments` | assignment_attachments, submissions |
| `audit_logs` | _(standalone — referensi by entity_type + entity_id)_ |

---

*Terakhir diperbarui: 21 Maret 2026*
