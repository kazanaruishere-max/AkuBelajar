-- 000008: Create grade tables
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
