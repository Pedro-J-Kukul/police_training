-- Remove foreign key constraints for tokens table
ALTER TABLE "tokens" DROP CONSTRAINT IF EXISTS "user_tokens";