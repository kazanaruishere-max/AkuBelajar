-- 000009: Create notification tables
CREATE TABLE notifications (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                 UUID NOT NULL REFERENCES users(id),
    type                    VARCHAR(50) NOT NULL,
    title                   VARCHAR(100) NOT NULL,
    body                    TEXT NOT NULL,
    is_read                 BOOLEAN DEFAULT FALSE,
    read_at                 TIMESTAMPTZ,
    related_entity_type     VARCHAR(50),
    related_entity_id       UUID,
    created_at              TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_notif_user_unread ON notifications(user_id, is_read) WHERE is_read = FALSE;

CREATE TABLE notification_preferences (
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
