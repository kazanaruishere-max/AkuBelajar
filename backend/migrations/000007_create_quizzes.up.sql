-- 000007: Create quiz tables
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

CREATE TABLE quiz_questions (
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
    cheat_log       JSONB DEFAULT '[]',
    question_order  INTEGER[] NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE quiz_answers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id      UUID NOT NULL REFERENCES quiz_sessions(id) ON DELETE CASCADE,
    question_id     UUID NOT NULL REFERENCES quiz_questions(id),
    selected_key    VARCHAR(1),
    essay_text      TEXT,
    is_correct      BOOLEAN,
    answered_at     TIMESTAMPTZ,
    UNIQUE(session_id, question_id)
);
