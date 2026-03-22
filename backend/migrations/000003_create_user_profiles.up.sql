-- 000003: Create user_profiles table (Core)
CREATE TABLE user_profiles (
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

CREATE UNIQUE INDEX idx_profiles_nisn ON user_profiles(nisn) WHERE nisn IS NOT NULL;
CREATE INDEX idx_profiles_user ON user_profiles(user_id);
