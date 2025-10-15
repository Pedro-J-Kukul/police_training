-- Add foreign key constraints for training_enrollments table
ALTER TABLE "training_enrollments" ADD CONSTRAINT "enrollment_officer" FOREIGN KEY ("officer_id") REFERENCES "officers" ("id");
ALTER TABLE "training_enrollments" ADD CONSTRAINT "enrollment_session" FOREIGN KEY ("session_id") REFERENCES "training_sessions" ("id");
ALTER TABLE "training_enrollments" ADD CONSTRAINT "enrollment_enrollment_status" FOREIGN KEY ("enrollment_status_id") REFERENCES "enrollment_statuses" ("id");
ALTER TABLE "training_enrollments" ADD CONSTRAINT "enrollment_attendance_status" FOREIGN KEY ("attendance_status_id") REFERENCES "attendance_statuses" ("id");
ALTER TABLE "training_enrollments" ADD CONSTRAINT "enrollment_progress_status" FOREIGN KEY ("progress_status_id") REFERENCES "progress_statuses" ("id");