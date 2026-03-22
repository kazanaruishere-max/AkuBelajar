-- 000011: Create audit_logs table (IMMUTABLE — no UPDATE/DELETE allowed)
CREATE TABLE audit_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id        UUID NOT NULL,
    action          VARCHAR(50) NOT NULL,
    entity_type     VARCHAR(50) NOT NULL,
    entity_id       UUID NOT NULL,
    old_value       JSONB,
    new_value       JSONB,
    ip_address      INET,
    user_agent      TEXT,
    request_id      UUID,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_audit_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_actor ON audit_logs(actor_id, created_at);

-- IMPORTANT: Revoke mutation permissions — audit logs are IMMUTABLE
-- This will be applied after the app_user role is created
-- REVOKE UPDATE, DELETE, TRUNCATE ON audit_logs FROM app_user;
