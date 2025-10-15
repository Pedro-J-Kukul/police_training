-- Remove foreign key constraints for training_sessions table
ALTER TABLE "training_sessions" DROP CONSTRAINT IF EXISTS "session_status";
ALTER TABLE "training_sessions" DROP CONSTRAINT IF EXISTS "session_region";
ALTER TABLE "training_sessions" DROP CONSTRAINT IF EXISTS "session_formation";
ALTER TABLE "training_sessions" DROP CONSTRAINT IF EXISTS "session_facilitator";
ALTER TABLE "training_sessions" DROP CONSTRAINT IF EXISTS "session_workshop";