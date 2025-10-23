-- Remove the constraint referencing users
ALTER TABLE "training_enrollments" DROP CONSTRAINT IF EXISTS "enrollment_officer";

-- Add the constraint referencing officers
ALTER TABLE "training_enrollments" ADD CONSTRAINT "enrollment_officer" FOREIGN KEY ("officer_id") REFERENCES "officers" ("id");