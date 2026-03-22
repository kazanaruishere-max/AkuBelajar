-- 000006: Create assignment tables
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
    grade_after_penalty INTEGER,
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
