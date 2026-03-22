DROP POLICY IF EXISTS att_school_isolation ON attendances;
ALTER TABLE attendances DISABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS users_school_isolation ON users;
ALTER TABLE users DISABLE ROW LEVEL SECURITY;
