-- 000012: Enable Row Level Security on multi-tenant tables
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
CREATE POLICY users_school_isolation ON users
    USING (school_id = current_setting('app.current_school_id', true)::UUID);

ALTER TABLE attendances ENABLE ROW LEVEL SECURITY;
CREATE POLICY att_school_isolation ON attendances
    USING (student_id IN (
        SELECT id FROM users
        WHERE school_id = current_setting('app.current_school_id', true)::UUID
    ));

-- Note: RLS is enforced via SET config in database.SetSchoolContext()
-- The 'true' flag in current_setting makes it return NULL instead of error if not set
