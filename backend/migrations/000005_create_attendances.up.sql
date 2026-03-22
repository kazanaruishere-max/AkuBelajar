-- 000005: Create attendances table
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
