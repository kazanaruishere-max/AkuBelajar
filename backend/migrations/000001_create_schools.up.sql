-- 000001: Create schools table (Core)
CREATE TABLE schools (
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

CREATE INDEX idx_schools_code ON schools(code) WHERE deleted_at IS NULL;
