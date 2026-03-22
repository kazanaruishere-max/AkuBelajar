-- ============================================================
-- AkuBelajar - Combined Migration (ALL tables)
-- Run this in Supabase SQL Editor: Dashboard → SQL Editor → New Query
-- ============================================================

-- ==============================
-- 000001: Schools
-- ==============================
CREATE TABLE IF NOT EXISTS schools (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(200) NOT NULL,
    code            VARCHAR(6) UNIQUE NOT NULL,
    address         TEXT,
    logo_url        VARCHAR(500),
    theme_color     VARCHAR(7) DEFAULT '#3B82F6',
    config          JSONB DEFAULT '{}',
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_schools_code ON schools(code) WHERE deleted_at IS NULL;

-- ==============================
-- 000002: Users
-- ==============================
DO $$ BEGIN
    CREATE TYPE user_role AS ENUM ('super_admin', 'teacher', 'class_leader', 'student');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS users (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id           UUID REFERENCES schools(id),
    email               VARCHAR(255) UNIQUE NOT NULL,
    password_hash       VARCHAR(255) NOT NULL,
    full_name           VARCHAR(200) DEFAULT '',
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
CREATE INDEX IF NOT EXISTS idx_users_school ON users(school_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_role ON users(school_id, role) WHERE deleted_at IS NULL;

-- ==============================
-- 000003: User Profiles
-- ==============================
CREATE TABLE IF NOT EXISTS user_profiles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nisn            VARCHAR(10),
    nip             VARCHAR(18),
    birth_date      DATE,
    phone_wa        VARCHAR(15),
    parent_name     VARCHAR(200),
    parent_phone    VARCHAR(15),
    photo_url       VARCHAR(500),
    bio             TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_profiles_nisn ON user_profiles(nisn) WHERE nisn IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_profiles_user ON user_profiles(user_id);

-- ==============================
-- 000004: Academic Tables
-- ==============================
CREATE TABLE IF NOT EXISTS academic_years (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id       UUID NOT NULL REFERENCES schools(id),
    name            VARCHAR(20) NOT NULL,
    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,
    is_active       BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS classes (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id           UUID NOT NULL REFERENCES schools(id),
    academic_year_id    UUID NOT NULL REFERENCES academic_years(id),
    name                VARCHAR(20) NOT NULL,
    grade_level         INTEGER NOT NULL,
    homeroom_teacher_id UUID REFERENCES users(id),
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_classes_school_year ON classes(school_id, academic_year_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS subjects (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id       UUID NOT NULL REFERENCES schools(id),
    name            VARCHAR(100) NOT NULL,
    code            VARCHAR(10),
    description     TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_subjects_school ON subjects(school_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS class_subjects (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id        UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    subject_id      UUID NOT NULL REFERENCES subjects(id),
    teacher_id      UUID NOT NULL REFERENCES users(id),
    schedule_json   JSONB DEFAULT '[]',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(class_id, subject_id)
);

CREATE TABLE IF NOT EXISTS student_classes (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id          UUID NOT NULL REFERENCES users(id),
    class_id            UUID NOT NULL REFERENCES classes(id),
    academic_year_id    UUID NOT NULL REFERENCES academic_years(id),
    joined_at           TIMESTAMPTZ DEFAULT NOW(),
    left_at             TIMESTAMPTZ,
    UNIQUE(student_id, class_id, academic_year_id)
);
CREATE INDEX IF NOT EXISTS idx_sc_class ON student_classes(class_id);
CREATE INDEX IF NOT EXISTS idx_sc_student ON student_classes(student_id);

-- ==============================
-- 000005: Attendances
-- ==============================
DO $$ BEGIN
    CREATE TYPE attendance_status AS ENUM ('present', 'permission', 'sick', 'absent', 'late', 'excused');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS attendances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id      UUID NOT NULL REFERENCES users(id),
    class_id        UUID NOT NULL REFERENCES classes(id),
    subject_id      UUID REFERENCES subjects(id),
    date            DATE NOT NULL,
    status          attendance_status NOT NULL DEFAULT 'present',
    notes           TEXT,
    proof_url       VARCHAR(500),
    recorded_by     UUID NOT NULL REFERENCES users(id),
    is_draft        BOOLEAN DEFAULT FALSE,
    approved_by     UUID REFERENCES users(id),
    approved_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(student_id, class_id, subject_id, date)
);
CREATE INDEX IF NOT EXISTS idx_att_class_date ON attendances(class_id, date);
CREATE INDEX IF NOT EXISTS idx_att_student ON attendances(student_id, date);

-- ==============================
-- 000006: Assignments
-- ==============================
DO $$ BEGIN
    CREATE TYPE assignment_status AS ENUM ('draft', 'published', 'closed');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;
DO $$ BEGIN
    CREATE TYPE submission_status AS ENUM ('draft', 'submitted', 'graded', 'revision_requested');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS assignments (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id            UUID NOT NULL REFERENCES classes(id),
    subject_id          UUID NOT NULL REFERENCES subjects(id),
    teacher_id          UUID NOT NULL REFERENCES users(id),
    title               VARCHAR(200) NOT NULL,
    description         TEXT,
    deadline_at         TIMESTAMPTZ NOT NULL,
    allow_late          BOOLEAN DEFAULT FALSE,
    late_penalty_pct    INTEGER DEFAULT 10,
    max_late_days       INTEGER DEFAULT 5,
    max_file_count      INTEGER DEFAULT 5,
    max_file_size_mb    INTEGER DEFAULT 20,
    allowed_extensions  TEXT[] DEFAULT '{pdf,docx,pptx,xlsx,jpg,png,zip}',
    weight_pct          INTEGER DEFAULT 100,
    status              assignment_status DEFAULT 'draft',
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS assignment_attachments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_id   UUID NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
    file_url        VARCHAR(500) NOT NULL,
    file_name       VARCHAR(255) NOT NULL,
    file_size       BIGINT NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS assignment_submissions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_id   UUID NOT NULL REFERENCES assignments(id),
    student_id      UUID NOT NULL REFERENCES users(id),
    submitted_at    TIMESTAMPTZ,
    is_late         BOOLEAN DEFAULT FALSE,
    late_days       INTEGER DEFAULT 0,
    status          submission_status DEFAULT 'draft',
    grade           INTEGER CHECK (grade >= 0 AND grade <= 100),
    grade_after_penalty INTEGER,
    feedback        TEXT,
    graded_by       UUID REFERENCES users(id),
    graded_at       TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(assignment_id, student_id)
);

CREATE TABLE IF NOT EXISTS submission_files (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id   UUID NOT NULL REFERENCES assignment_submissions(id) ON DELETE CASCADE,
    file_url        VARCHAR(500) NOT NULL,
    file_name       VARCHAR(255) NOT NULL,
    file_size       BIGINT NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- ==============================
-- 000007: Quizzes
-- ==============================
DO $$ BEGIN
    CREATE TYPE quiz_status AS ENUM ('draft', 'published', 'active', 'ended', 'graded', 'locked');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;
DO $$ BEGIN
    CREATE TYPE question_type AS ENUM ('multiple_choice', 'essay');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;
DO $$ BEGIN
    CREATE TYPE review_mode AS ENUM ('immediately', 'after_all_submit', 'manual');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;
DO $$ BEGIN
    CREATE TYPE session_status AS ENUM ('in_progress', 'submitted', 'auto_submitted', 'locked');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS quizzes (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id            UUID NOT NULL REFERENCES classes(id),
    subject_id          UUID NOT NULL REFERENCES subjects(id),
    teacher_id          UUID NOT NULL REFERENCES users(id),
    title               VARCHAR(200) NOT NULL,
    time_limit          INTEGER NOT NULL,
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

CREATE TABLE IF NOT EXISTS quiz_questions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quiz_id         UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    question_text   TEXT NOT NULL,
    question_type   question_type NOT NULL DEFAULT 'multiple_choice',
    options         JSONB,
    answer_hash     VARCHAR(255),
    explanation     TEXT,
    order_num       INTEGER NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS quiz_sessions (
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
    cheat_log       JSONB DEFAULT '[]',
    question_order  INTEGER[] NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS quiz_answers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id      UUID NOT NULL REFERENCES quiz_sessions(id) ON DELETE CASCADE,
    question_id     UUID NOT NULL REFERENCES quiz_questions(id),
    selected_key    VARCHAR(1),
    essay_text      TEXT,
    is_correct      BOOLEAN,
    answered_at     TIMESTAMPTZ,
    UNIQUE(session_id, question_id)
);

-- ==============================
-- 000008: Grades
-- ==============================
CREATE TABLE IF NOT EXISTS grades (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id          UUID NOT NULL REFERENCES users(id),
    subject_id          UUID NOT NULL REFERENCES subjects(id),
    class_id            UUID REFERENCES classes(id),
    academic_year_id    UUID NOT NULL REFERENCES academic_years(id),
    semester            INTEGER NOT NULL CHECK (semester IN (1, 2)),
    category            VARCHAR(50),
    description         VARCHAR(200),
    score               NUMERIC(5,2),
    weight_pct          NUMERIC(5,2) DEFAULT 100,
    weighted_score      NUMERIC(5,2),
    assignment_avg      NUMERIC(5,2),
    quiz_avg            NUMERIC(5,2),
    final_score         NUMERIC(5,2),
    grade_letter        VARCHAR(2),
    teacher_id          UUID REFERENCES users(id),
    is_locked           BOOLEAN DEFAULT FALSE,
    locked_at           TIMESTAMPTZ,
    locked_by           UUID REFERENCES users(id),
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(student_id, subject_id, academic_year_id, semester)
);
CREATE INDEX IF NOT EXISTS idx_grades_student ON grades(student_id);
CREATE INDEX IF NOT EXISTS idx_grades_class ON grades(class_id) WHERE class_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS report_cards (
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

-- ==============================
-- 000009: Notifications
-- ==============================
CREATE TABLE IF NOT EXISTS notifications (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                 UUID NOT NULL REFERENCES users(id),
    type                    VARCHAR(50) NOT NULL,
    title                   VARCHAR(100) NOT NULL,
    message                 TEXT NOT NULL,
    link                    VARCHAR(500) DEFAULT '',
    is_read                 BOOLEAN DEFAULT FALSE,
    read_at                 TIMESTAMPTZ,
    related_entity_type     VARCHAR(50),
    related_entity_id       UUID,
    created_at              TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_notif_user_unread ON notifications(user_id, is_read) WHERE is_read = FALSE;
CREATE INDEX IF NOT EXISTS idx_notif_user ON notifications(user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS notification_preferences (
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

-- ==============================
-- 000010: Auth Tables
-- ==============================
CREATE TABLE IF NOT EXISTS invite_tokens (
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

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    token_hash      VARCHAR(255) NOT NULL,
    expires_at      TIMESTAMPTZ NOT NULL,
    used_at         TIMESTAMPTZ,
    ip_address      INET,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS password_histories (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    password_hash   VARCHAR(255) NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_pw_history_user ON password_histories(user_id);

CREATE TABLE IF NOT EXISTS active_sessions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id),
    refresh_token_hash  VARCHAR(255) NOT NULL,
    device_info         JSONB,
    ip_address          INET,
    is_remember_me      BOOLEAN DEFAULT FALSE,
    expires_at          TIMESTAMPTZ NOT NULL,
    created_at          TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_sessions_user ON active_sessions(user_id);

-- ==============================
-- 000011: Audit Logs
-- ==============================
CREATE TABLE IF NOT EXISTS audit_logs (
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
CREATE INDEX IF NOT EXISTS idx_audit_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_actor ON audit_logs(actor_id, created_at);

-- ============================================================
-- DONE! All tables created.
-- Next step: Run the seed: go run cmd/seed/main.go
-- ============================================================
