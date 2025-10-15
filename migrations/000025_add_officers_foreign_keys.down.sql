-- Remove foreign key constraints for officers table
ALTER TABLE "officers" DROP CONSTRAINT IF EXISTS "officer_region";
ALTER TABLE "officers" DROP CONSTRAINT IF EXISTS "officer_posting";
ALTER TABLE "officers" DROP CONSTRAINT IF EXISTS "officer_rank";
ALTER TABLE "officers" DROP CONSTRAINT IF EXISTS "officer_formation";