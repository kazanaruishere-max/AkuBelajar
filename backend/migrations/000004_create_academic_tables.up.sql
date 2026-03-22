-- 000004: Create academic tables
CREATE TABLE academic_years (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id       UUID NOT NULL REFERENCES schools(id),
    name            VARCHAR(20) NOT NULL,
    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,
    is_active       BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_ay_active ON academic_years(school_id) WHERE is_active = TRUE;

CREATE TABLE classes (
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

CREATE INDEX idx_classes_school_year ON classes(school_id, academic_year_id) WHERE deleted_at IS NULL;

CREATE TABLE subjects (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id       UUID NOT NULL REFERENCES schools(id),
    name            VARCHAR(100) NOT NULL,
    code            VARCHAR(10),
    description     TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_subjects_school ON subjects(school_id) WHERE deleted_at IS NULL;

CREATE TABLE class_subjects (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id        UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    subject_id      UUID NOT NULL REFERENCES subjects(id),
    teacher_id      UUID NOT NULL REFERENCES users(id),
    schedule_json   JSONB DEFAULT '[]',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(class_id, subject_id)
);

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
