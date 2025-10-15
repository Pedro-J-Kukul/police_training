-- Remove foreign key constraints for formations table
ALTER TABLE "formations" DROP CONSTRAINT IF EXISTS "formation_region";