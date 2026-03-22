-- 000002: Create users table (Core)
CREATE TYPE user_role AS ENUM ('super_admin', 'teacher', 'class_leader', 'student');

CREATE TABLE users (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id           UUID REFERENCES schools(id),
    email               VARCHAR(255) UNIQUE NOT NULL,
    password            VARCHAR(255) NOT NULL,
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
