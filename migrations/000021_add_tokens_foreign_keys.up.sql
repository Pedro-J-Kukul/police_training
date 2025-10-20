-- Add foreign key constraints for tokens table
ALTER TABLE "tokens" ADD CONSTRAINT "user_tokens" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;