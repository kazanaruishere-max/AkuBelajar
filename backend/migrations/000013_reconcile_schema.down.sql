-- 000013 down: Revert reconciliation
ALTER TABLE users RENAME COLUMN password_hash TO password;
ALTER TABLE users DROP COLUMN IF EXISTS full_name;

ALTER TABLE notifications RENAME COLUMN message TO body;
ALTER TABLE notifications DROP COLUMN IF EXISTS link;

ALTER TABLE assignment_submissions RENAME TO submissions;

ALTER TABLE grades DROP COLUMN IF EXISTS class_id;
ALTER TABLE grades DROP COLUMN IF EXISTS category;
ALTER TABLE grades DROP COLUMN IF EXISTS score;
ALTER TABLE grades DROP COLUMN IF EXISTS weight_pct;
ALTER TABLE grades DROP COLUMN IF EXISTS weighted_score;
ALTER TABLE grades DROP COLUMN IF EXISTS description;
ALTER TABLE grades DROP COLUMN IF EXISTS teacher_id;

ALTER TABLE attendances RENAME COLUMN notes TO reason;

DROP INDEX IF EXISTS idx_grades_student;
DROP INDEX IF EXISTS idx_grades_class;
DROP INDEX IF EXISTS idx_notif_user;
