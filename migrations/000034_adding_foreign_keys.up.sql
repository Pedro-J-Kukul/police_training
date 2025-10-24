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


-- ADDING FOREIGN KEYS
ALTER TABLE "training_enrollments"
ADD CONSTRAINT fk_training_enrollments_officer
FOREIGN KEY ("officer_id") REFERENCES "officers" ("id")
ON DELETE CASCADE;

ALTER TABLE "training_enrollments"
ADD CONSTRAINT fk_training_enrollments_session
FOREIGN KEY ("session_id") REFERENCES "training_sessions" ("id")
ON DELETE CASCADE;

ALTER TABLE "training_enrollments"
ADD CONSTRAINT fk_training_enrollments_enrollment_status
FOREIGN KEY ("enrollment_status_id") REFERENCES "enrollment_statuses" ("id")
ON DELETE CASCADE;

ALTER TABLE "training_enrollments"
ADD CONSTRAINT fk_training_enrollments_attendance_status
FOREIGN KEY ("attendance_status_id") REFERENCES "attendance_statuses" ("id")
ON DELETE CASCADE;

ALTER TABLE "training_enrollments"
ADD CONSTRAINT fk_training_enrollments_progress_status
FOREIGN KEY ("progress_status_id") REFERENCES "progress_statuses" ("id")
ON DELETE CASCADE;