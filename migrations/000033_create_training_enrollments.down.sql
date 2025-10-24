-- CREATE TABLE "training_enrollments" (
--   "id" bigserial PRIMARY KEY,
--   "officer_id" bigint NOT NULL,
--   "session_id" bigint NOT NULL,
--   "enrollment_status_id" bigint NOT NULL,
--   "attendance_status_id" bigint,
--   "progress_status_id" bigint NOT NULL,
--   "completion_date" date,
--   "certificate_issued" boolean DEFAULT false,
--   "certificate_number" text,
--   "created_at" timestamp DEFAULT NOW(),
--   "updated_at" timestamp DEFAULT NOW()
-- );

-- CREATE UNIQUE INDEX idx_training_enrollments_session_officer ON "training_enrollments" ("session_id", "officer_id");

-- CREATE INDEX idx_training_enrollments_session_id ON "training_enrollments" ("session_id");

-- CREATE INDEX idx_training_enrollments_officer_id ON "training_enrollments" ("officer_id");

-- CREATE INDEX idx_training_enrollments_completion_date ON "training_enrollments" ("completion_date");

-- DROPPING
DROP TABLE IF EXISTS "training_enrollments";
DROP INDEX IF EXISTS idx_training_enrollments_session_officer;
DROP INDEX IF EXISTS idx_training_enrollments_session_id;
DROP INDEX IF EXISTS idx_training_enrollments_officer_id;
DROP INDEX IF EXISTS idx_training_enrollments_completion_date;