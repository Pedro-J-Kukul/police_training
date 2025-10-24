
-- DROPPING FOREIGN KEYS
ALTER TABLE "training_enrollments"
DROP CONSTRAINT IF EXISTS fk_training_enrollments_officer;

ALTER TABLE "training_enrollments"
DROP CONSTRAINT IF EXISTS fk_training_enrollments_session;

ALTER TABLE "training_enrollments"
DROP CONSTRAINT IF EXISTS fk_training_enrollments_enrollment_status;

ALTER TABLE "training_enrollments"
DROP CONSTRAINT IF EXISTS fk_training_enrollments_attendance_status;

ALTER TABLE "training_enrollments"
DROP CONSTRAINT IF EXISTS fk_training_enrollments_progress_status;