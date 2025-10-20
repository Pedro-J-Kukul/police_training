-- Add foreign key constraints for roles_users table
ALTER TABLE "roles_users" ADD CONSTRAINT "roles_users_role" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON DELETE CASCADE;
ALTER TABLE "roles_users" ADD CONSTRAINT "roles_users_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;