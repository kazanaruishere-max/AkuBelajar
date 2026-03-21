# 🗄️ Database Schema Full — AkuBelajar

> DDL SQL lengkap SEMUA tabel di sistem AkuBelajar. Sumber kebenaran tunggal untuk struktur database.

---

## Konvensi Global

| Konvensi | Detail |
|:---|:---|
| Primary Key | UUID v7 (`gen_random_uuid()`) |
| Timestamps | `created_at`, `updated_at` di semua tabel |
| Soft Delete | `deleted_at TIMESTAMPTZ` |
| Multi-tenant | `school_id` di semua tabel data sekolah |
| Password | Argon2id hash |
| Kunci Jawaban | Argon2id hash (tidak boleh plain text) |

---

## 1. CORE TABLES

### schools

```sql
CREATE TABLE schools (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(200) NOT NULL,
    code            VARCHAR(6) UNIQUE NOT NULL,        -- 6 char uppercase alphanumeric
    address         TEXT,
    logo_url        VARCHAR(500),
    theme_color     VARCHAR(7) DEFAULT '#3B82F6',      -- hex color
    config          JSONB DEFAULT '{}',                 -- grading weights, KKM, etc.
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

-- Config JSONB structure:
-- {
--   "grading_weights": { "assignment": 60, "quiz": 40 },
--   "kkm_default": 70,
--   "attendance_minimum_pct": 75,
--   "late_penalty_pct_per_day": 10,
--   "max_late_days": 5,
--   "grade_scale": [
--     { "min": 90, "max": 100, "letter": "A", "label": "Sangat Baik" },
--     { "min": 80, "max": 89,  "letter": "B", "label": "Baik" },
--     { "min": 70, "max": 79,  "letter": "C", "label": "Cukup" },
--     { "min": 60, "max": 69,  "letter": "D", "label": "Kurang" },
--     { "min": 0,  "max": 59,  "letter": "E", "label": "Sangat Kurang" }
--   ]
-- }

CREATE INDEX idx_schools_code ON schools(code) WHERE deleted_at IS NULL;
```

### users

```sql
CREATE TYPE user_role AS ENUM ('super_admin', 'teacher', 'class_leader', 'student');

CREATE TABLE users (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id           UUID REFERENCES schools(id),
    email               VARCHAR(255) UNIQUE NOT NULL,
    password            VARCHAR(255) NOT NULL,           -- Argon2id hash
    role                user_role NOT NULL,
    is_active           BOOLEAN DEFAULT TRUE,
    is_first_login      BOOLEAN DEFAULT TRUE,
    failed_login_count  INTEGER DEFAULT 0,
    locked_until        TIMESTAMPTZ,
    last_login_at       TIMESTAMPTZ,
    last_login_ip       INET,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);

CREATE INDEX idx_users_school ON users(school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(school_id, role) WHERE deleted_at IS NULL;

-- RLS
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
CREATE POLICY users_school_isolation ON users
    USING (school_id = current_setting('app.current_school_id')::UUID);
```

### user_profiles

```sql
CREATE TABLE user_profiles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nisn            VARCHAR(10),                        -- 10 digit, siswa only
    nip             VARCHAR(18),                        -- 18 digit ASN, guru only
    birth_date      DATE,
    phone_wa        VARCHAR(15),                        -- format E.164
    parent_name     VARCHAR(200),
    parent_phone    VARCHAR(15),
    photo_url       VARCHAR(500),
    bio             TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_profiles_nisn ON user_profiles(nisn)
    WHERE nisn IS NOT NULL;
CREATE INDEX idx_profiles_user ON user_profiles(user_id);
```

---

## 2. ACADEMIC TABLES

### academic_years

```sql
CREATE TABLE academic_years (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id       UUID NOT NULL REFERENCES schools(id),
    name            VARCHAR(20) NOT NULL,               -- "2025/2026"
    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,
    is_active       BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_ay_active ON academic_years(school_id)
    WHERE is_active = TRUE;                             -- Hanya 1 aktif per sekolah
```

### classes

```sql
CREATE TABLE classes (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id           UUID NOT NULL REFERENCES schools(id),
    academic_year_id    UUID NOT NULL REFERENCES academic_years(id),
    name                VARCHAR(20) NOT NULL,            -- "8A", "XII-IPA-1"
    grade_level         INTEGER NOT NULL,                -- 7, 8, 9, 10, 11, 12
    homeroom_teacher_id UUID REFERENCES users(id),
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);

CREATE INDEX idx_classes_school_year ON classes(school_id, academic_year_id)
    WHERE deleted_at IS NULL;
```

### subjects

```sql
CREATE TABLE subjects (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id       UUID NOT NULL REFERENCES schools(id),
    name            VARCHAR(100) NOT NULL,
    code            VARCHAR(10),                        -- "MTK", "BIN", "IPA"
    description     TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_subjects_school ON subjects(school_id) WHERE deleted_at IS NULL;
```

### class_subjects

```sql
CREATE TABLE class_subjects (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id        UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    subject_id      UUID NOT NULL REFERENCES subjects(id),
    teacher_id      UUID NOT NULL REFERENCES users(id),
    schedule_json   JSONB DEFAULT '[]',                 -- [{"day":"monday","start":"08:00","end":"09:30"}]
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(class_id, subject_id)
);
```

### student_classes

```sql
CREATE TABLE student_classes (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id          UUID NOT NULL REFERENCES users(id),
    class_id            UUID NOT NULL REFERENCES classes(id),
    academic_year_id    UUID NOT NULL REFERENCES academic_years(id),
    joined_at           TIMESTAMPTZ DEFAULT NOW(),
    left_at             TIMESTAMPTZ,
    UNIQUE(student_id, class_id, academic_year_id)
);

CREATE INDEX idx_sc_class ON student_classes(class_id);
CREATE INDEX idx_sc_student ON student_classes(student_id);
```

---

## 3. ATTENDANCE TABLE

```sql
CREATE TYPE attendance_status AS ENUM ('present', 'permission', 'sick', 'absent', 'late');

CREATE TABLE attendances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id      UUID NOT NULL REFERENCES users(id),
    class_id        UUID NOT NULL REFERENCES classes(id),
    subject_id      UUID REFERENCES subjects(id),
    date            DATE NOT NULL,
    status          attendance_status NOT NULL DEFAULT 'present',
    reason          TEXT,
    proof_url       VARCHAR(500),
    recorded_by     UUID NOT NULL REFERENCES users(id),
    is_draft        BOOLEAN DEFAULT FALSE,
    approved_by     UUID REFERENCES users(id),
    approved_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(student_id, class_id, subject_id, date)
);

CREATE INDEX idx_att_class_date ON attendances(class_id, date);
CREATE INDEX idx_att_student ON attendances(student_id, date);

ALTER TABLE attendances ENABLE ROW LEVEL SECURITY;
CREATE POLICY att_school_isolation ON attendances
    USING (student_id IN (SELECT id FROM users WHERE school_id = current_setting('app.current_school_id')::UUID));
```

---

## 4. ASSIGNMENT TABLES

```sql
CREATE TYPE assignment_status AS ENUM ('draft', 'published', 'closed');
CREATE TYPE submission_status AS ENUM ('draft', 'submitted', 'graded', 'revision_requested');

CREATE TABLE assignments (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id            UUID NOT NULL REFERENCES classes(id),
    subject_id          UUID NOT NULL REFERENCES subjects(id),
    teacher_id          UUID NOT NULL REFERENCES users(id),
    title               VARCHAR(200) NOT NULL,
    description         TEXT,
    deadline_at         TIMESTAMPTZ NOT NULL,
    allow_late          BOOLEAN DEFAULT FALSE,
    late_penalty_pct    INTEGER DEFAULT 10,              -- % per day
    max_late_days       INTEGER DEFAULT 5,
    max_file_count      INTEGER DEFAULT 3,
    max_file_size_mb    INTEGER DEFAULT 20,
    allowed_extensions  TEXT[] DEFAULT '{pdf,docx,pptx,xlsx,jpg,png,zip}',
    weight_pct          INTEGER DEFAULT 100,             -- bobot nilai %
    status              assignment_status DEFAULT 'draft',
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);

CREATE TABLE assignment_attachments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_id   UUID NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
    file_url        VARCHAR(500) NOT NULL,
    file_name       VARCHAR(255) NOT NULL,
    file_size       BIGINT NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE submissions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_id   UUID NOT NULL REFERENCES assignments(id),
    student_id      UUID NOT NULL REFERENCES users(id),
    submitted_at    TIMESTAMPTZ,
    is_late         BOOLEAN DEFAULT FALSE,
    late_days       INTEGER DEFAULT 0,
    status          submission_status DEFAULT 'draft',
    grade           INTEGER CHECK (grade >= 0 AND grade <= 100),
    grade_after_penalty INTEGER,                         -- grade × (1 - penalty%)
    feedback        TEXT,
    graded_by       UUID REFERENCES users(id),
    graded_at       TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(assignment_id, student_id)
);

CREATE TABLE submission_files (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id   UUID NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
    file_url        VARCHAR(500) NOT NULL,
    file_name       VARCHAR(255) NOT NULL,
    file_size       BIGINT NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 5. QUIZ TABLES

```sql
CREATE TYPE quiz_status AS ENUM ('draft', 'published', 'active', 'ended', 'graded', 'locked');
CREATE TYPE question_type AS ENUM ('multiple_choice', 'essay');
CREATE TYPE review_mode AS ENUM ('immediately', 'after_all_submit', 'manual');
CREATE TYPE session_status AS ENUM ('in_progress', 'submitted', 'auto_submitted', 'locked');

CREATE TABLE quizzes (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id            UUID NOT NULL REFERENCES classes(id),
    subject_id          UUID NOT NULL REFERENCES subjects(id),
    teacher_id          UUID NOT NULL REFERENCES users(id),
    title               VARCHAR(200) NOT NULL,
    time_limit          INTEGER NOT NULL,                -- menit
    is_published        BOOLEAN DEFAULT FALSE,
    ai_generated        BOOLEAN DEFAULT FALSE,
    randomize_questions BOOLEAN DEFAULT TRUE,
    randomize_options   BOOLEAN DEFAULT TRUE,
    max_attempts        INTEGER DEFAULT 1,
    allow_review        BOOLEAN DEFAULT TRUE,
    review_mode         review_mode DEFAULT 'immediately',
    start_at            TIMESTAMPTZ,
    end_at              TIMESTAMPTZ,
    status              quiz_status DEFAULT 'draft',
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);

CREATE TABLE quiz_questions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quiz_id         UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    question_text   TEXT NOT NULL,
    question_type   question_type NOT NULL DEFAULT 'multiple_choice',
    options         JSONB,                               -- [{"key":"A","text":"..."},...]
    answer_hash     VARCHAR(255),                        -- Argon2id hash
    explanation     TEXT,
    order_num       INTEGER NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE quiz_sessions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quiz_id         UUID NOT NULL REFERENCES quizzes(id),
    student_id      UUID NOT NULL REFERENCES users(id),
    started_at      TIMESTAMPTZ DEFAULT NOW(),
    submitted_at    TIMESTAMPTZ,
    expires_at      TIMESTAMPTZ NOT NULL,
    status          session_status DEFAULT 'in_progress',
    score           INTEGER,
    ip_address      INET,
    user_agent      TEXT,
    cheat_count     INTEGER DEFAULT 0,
    cheat_log       JSONB DEFAULT '[]',                  -- [{"type":"tab_switch","at":"..."}]
    question_order  INTEGER[] NOT NULL,                   -- shuffled question IDs order
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE quiz_answers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id      UUID NOT NULL REFERENCES quiz_sessions(id) ON DELETE CASCADE,
    question_id     UUID NOT NULL REFERENCES quiz_questions(id),
    selected_key    VARCHAR(1),                          -- "A","B","C","D" or NULL
    essay_text      TEXT,
    is_correct      BOOLEAN,
    answered_at     TIMESTAMPTZ,
    UNIQUE(session_id, question_id)
);
```

---

## 6. GRADE TABLES

```sql
CREATE TABLE grades (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id          UUID NOT NULL REFERENCES users(id),
    subject_id          UUID NOT NULL REFERENCES subjects(id),
    academic_year_id    UUID NOT NULL REFERENCES academic_years(id),
    semester            INTEGER NOT NULL CHECK (semester IN (1, 2)),
    assignment_avg      NUMERIC(5,2),
    quiz_avg            NUMERIC(5,2),
    final_score         NUMERIC(5,2),
    grade_letter        VARCHAR(2),
    is_locked           BOOLEAN DEFAULT FALSE,
    locked_at           TIMESTAMPTZ,
    locked_by           UUID REFERENCES users(id),
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(student_id, subject_id, academic_year_id, semester)
);

CREATE TABLE report_cards (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id              UUID NOT NULL REFERENCES users(id),
    academic_year_id        UUID NOT NULL REFERENCES academic_years(id),
    semester                INTEGER NOT NULL CHECK (semester IN (1, 2)),
    pdf_url                 VARCHAR(500),
    generated_at            TIMESTAMPTZ,
    digital_signature_hash  VARCHAR(255),
    qr_verification_code    VARCHAR(50) UNIQUE,
    is_published            BOOLEAN DEFAULT FALSE,
    published_at            TIMESTAMPTZ,
    created_at              TIMESTAMPTZ DEFAULT NOW(),
    updated_at              TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(student_id, academic_year_id, semester)
);
```

---

## 7. NOTIFICATION & SYSTEM TABLES

```sql
CREATE TABLE notifications (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                 UUID NOT NULL REFERENCES users(id),
    type                    VARCHAR(50) NOT NULL,
    title                   VARCHAR(100) NOT NULL,
    body                    TEXT NOT NULL,
    is_read                 BOOLEAN DEFAULT FALSE,
    read_at                 TIMESTAMPTZ,
    related_entity_type     VARCHAR(50),
    related_entity_id       UUID,
    created_at              TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_notif_user_unread ON notifications(user_id, is_read)
    WHERE is_read = FALSE;

CREATE TABLE notification_preferences (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID UNIQUE NOT NULL REFERENCES users(id),
    email_enabled   BOOLEAN DEFAULT TRUE,
    wa_enabled      BOOLEAN DEFAULT TRUE,
    in_app_enabled  BOOLEAN DEFAULT TRUE,
    quiet_start     TIME DEFAULT '22:00',
    quiet_end       TIME DEFAULT '06:00',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE invite_tokens (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id       UUID NOT NULL REFERENCES schools(id),
    created_by      UUID NOT NULL REFERENCES users(id),
    token_hash      VARCHAR(255) NOT NULL,
    role            user_role NOT NULL,
    class_id        UUID REFERENCES classes(id),
    max_uses        INTEGER DEFAULT 1,
    uses_count      INTEGER DEFAULT 0,
    expires_at      TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE password_reset_tokens (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    token_hash      VARCHAR(255) NOT NULL,               -- SHA-256 hash of OTP
    expires_at      TIMESTAMPTZ NOT NULL,
    used_at         TIMESTAMPTZ,
    ip_address      INET,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE password_histories (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    password_hash   VARCHAR(255) NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_pw_history_user ON password_histories(user_id);

CREATE TABLE active_sessions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    refresh_token_hash VARCHAR(255) NOT NULL,
    device_info     JSONB,
    ip_address      INET,
    is_remember_me  BOOLEAN DEFAULT FALSE,
    expires_at      TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_sessions_user ON active_sessions(user_id);

-- Audit log — IMMUTABLE
CREATE TABLE audit_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id        UUID NOT NULL,
    action          VARCHAR(50) NOT NULL,
    entity_type     VARCHAR(50) NOT NULL,
    entity_id       UUID NOT NULL,
    old_value       JSONB,
    new_value       JSONB,
    ip_address      INET,
    user_agent      TEXT,
    request_id      UUID,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

REVOKE UPDATE, DELETE, TRUNCATE ON audit_logs FROM app_user;
CREATE INDEX idx_audit_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_actor ON audit_logs(actor_id, created_at);
```

---

## Ringkasan Tabel

| # | Tabel | Kolom | Kategori |
|:---|:---|:---|:---|
| 1 | `schools` | 11 | Core |
| 2 | `users` | 14 | Core |
| 3 | `user_profiles` | 11 | Core |
| 4 | `academic_years` | 7 | Academic |
| 5 | `classes` | 8 | Academic |
| 6 | `subjects` | 7 | Academic |
| 7 | `class_subjects` | 6 | Academic |
| 8 | `student_classes` | 6 | Academic |
| 9 | `attendances` | 13 | Attendance |
| 10 | `assignments` | 17 | Assignment |
| 11 | `assignment_attachments` | 6 | Assignment |
| 12 | `submissions` | 14 | Assignment |
| 13 | `submission_files` | 6 | Assignment |
| 14 | `quizzes` | 17 | Quiz |
| 15 | `quiz_questions` | 9 | Quiz |
| 16 | `quiz_sessions` | 14 | Quiz |
| 17 | `quiz_answers` | 7 | Quiz |
| 18 | `grades` | 13 | Grade |
| 19 | `report_cards` | 10 | Grade |
| 20 | `notifications` | 9 | System |
| 21 | `notification_preferences` | 8 | System |
| 22 | `invite_tokens` | 10 | System |
| 23 | `password_reset_tokens` | 7 | System |
| 24 | `password_histories` | 4 | System |
| 25 | `active_sessions` | 8 | System |
| 26 | `audit_logs` | 10 | System |

**Total: 26 tabel**

---

*Terakhir diperbarui: 21 Maret 2026*
