-- Remove foreign key constraint linking users and officers tables
ALTER TABLE "users" DROP CONSTRAINT IF EXISTS "officer_user";