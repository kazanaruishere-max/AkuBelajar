-- Seed data for local development
-- Credentials: password hashes are Argon2id of the passwords listed below

-- School: SMP Nusantara Demo
INSERT INTO schools (id, name, code, address, config) VALUES
    ('a0000000-0000-0000-0000-000000000001', 'SMP Nusantara Demo', 'SMND01', 'Jl. Pendidikan No. 1, Jakarta',
     '{"grading_weights": {"assignment": 60, "quiz": 40}, "kkm_default": 70, "attendance_minimum_pct": 75, "late_penalty_pct_per_day": 10, "max_late_days": 5}'
    );

-- Users (passwords need to be hashed at runtime via seed script)
-- Super Admin: admin@akubelajar.id / Admin@123!
-- Guru 1: guru@akubelajar.id / Guru@123!
-- Guru 2: guru2@akubelajar.id / Guru@123!
-- Siswa 1: siswa@akubelajar.id / Siswa@123!
-- Siswa 2-5: siswa2-5@akubelajar.id / Siswa@123!

-- Note: Run the Go seed command instead of raw SQL for proper Argon2id hashing:
--   go run cmd/seed/main.go
-- This SQL serves as a reference for the seed data structure.

-- Academic Year
INSERT INTO academic_years (id, school_id, name, start_date, end_date, is_active) VALUES
    ('b0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001',
     '2025/2026', '2025-07-14', '2026-06-20', TRUE);

-- Subjects
INSERT INTO subjects (id, school_id, name, code) VALUES
    ('c0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', 'Matematika', 'MTK'),
    ('c0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', 'Bahasa Indonesia', 'BIN'),
    ('c0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', 'IPA', 'IPA'),
    ('c0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000001', 'Bahasa Inggris', 'BIG');

-- Classes
INSERT INTO classes (id, school_id, academic_year_id, name, grade_level) VALUES
    ('d0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001',
     'b0000000-0000-0000-0000-000000000001', '7A', 7),
    ('d0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001',
     'b0000000-0000-0000-0000-000000000001', '8A', 8),
    ('d0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001',
     'b0000000-0000-0000-0000-000000000001', '9A', 9);
