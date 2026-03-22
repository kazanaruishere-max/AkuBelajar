-- 000010: Create auth system tables
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
    token_hash      VARCHAR(255) NOT NULL,
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
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id),
    refresh_token_hash  VARCHAR(255) NOT NULL,
    device_info         JSONB,
    ip_address          INET,
    is_remember_me      BOOLEAN DEFAULT FALSE,
    expires_at          TIMESTAMPTZ NOT NULL,
    created_at          TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_sessions_user ON active_sessions(user_id);
