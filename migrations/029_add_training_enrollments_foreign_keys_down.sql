-- Remove foreign key constraints for training_enrollments table
ALTER TABLE "training_enrollments" DROP CONSTRAINT IF EXISTS "enrollment_progress_status";
ALTER TABLE "training_enrollments" DROP CONSTRAINT IF EXISTS "enrollment_attendance_status";
ALTER TABLE "training_enrollments" DROP CONSTRAINT IF EXISTS "enrollment_enrollment_status";
ALTER TABLE "training_enrollments" DROP CONSTRAINT IF EXISTS "enrollment_session";
ALTER TABLE "training_enrollments" DROP CONSTRAINT IF EXISTS "enrollment_officer";