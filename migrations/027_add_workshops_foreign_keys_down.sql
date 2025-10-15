-- Remove foreign key constraints for workshops table
ALTER TABLE "workshops" DROP CONSTRAINT IF EXISTS "workshop_training_type";
ALTER TABLE "workshops" DROP CONSTRAINT IF EXISTS "workshop_category";