-- 000013: Reconcile schema with Go backend code
-- This migration aligns the existing database schema with the column names
-- and types used in the Sprint 2-5 Go backend code.

-- ============================================================
-- 1. USERS TABLE
-- Go code expects: password_hash, full_name
-- Current schema has: password (no full_name)
-- ============================================================
ALTER TABLE users RENAME COLUMN password TO password_hash;
ALTER TABLE users ADD COLUMN IF NOT EXISTS full_name VARCHAR(200) DEFAULT '';

-- ============================================================
-- 2. NOTIFICATIONS TABLE
-- Go code expects: message, link
-- Current schema has: body (no link, no message)
-- ============================================================
ALTER TABLE notifications RENAME COLUMN body TO message;
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS link VARCHAR(500) DEFAULT '';

-- ============================================================
-- 3. ASSIGNMENT SUBMISSIONS
-- Go code references table: assignment_submissions
-- Current schema has table: submissions
-- Also Go dashboard references: assignment_submissions.status
-- ============================================================
ALTER TABLE submissions RENAME TO assignment_submissions;

-- ============================================================
-- 4. GRADES TABLE
-- Go code expects: class_id, category, score, weight_pct, weighted_score, description, teacher_id
-- Current schema has: assignment_avg, quiz_avg, final_score, grade_letter, etc.
-- Solution: Add missing columns to support the Go grade handler
-- ============================================================
ALTER TABLE grades ADD COLUMN IF NOT EXISTS class_id UUID REFERENCES classes(id);
ALTER TABLE grades ADD COLUMN IF NOT EXISTS category VARCHAR(50);
ALTER TABLE grades ADD COLUMN IF NOT EXISTS score NUMERIC(5,2);
ALTER TABLE grades ADD COLUMN IF NOT EXISTS weight_pct NUMERIC(5,2) DEFAULT 100;
ALTER TABLE grades ADD COLUMN IF NOT EXISTS weighted_score NUMERIC(5,2);
ALTER TABLE grades ADD COLUMN IF NOT EXISTS description VARCHAR(200);
ALTER TABLE grades ADD COLUMN IF NOT EXISTS teacher_id UUID REFERENCES users(id);

-- ============================================================
-- 5. ATTENDANCES TABLE
-- Go code expects: notes (not reason), status 'excused' (not 'permission')
-- ============================================================
ALTER TABLE attendances RENAME COLUMN reason TO notes;

-- Add 'excused' to attendance_status enum
ALTER TYPE attendance_status ADD VALUE IF NOT EXISTS 'excused';

-- ============================================================
-- 6. INDEXES
-- ============================================================
CREATE INDEX IF NOT EXISTS idx_grades_student ON grades(student_id);
CREATE INDEX IF NOT EXISTS idx_grades_class ON grades(class_id) WHERE class_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_notif_user ON notifications(user_id, created_at DESC);
