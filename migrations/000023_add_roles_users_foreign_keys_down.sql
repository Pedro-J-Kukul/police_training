-- Remove foreign key constraints for roles_users table
ALTER TABLE "roles_users" DROP CONSTRAINT IF EXISTS "roles_users_user";
ALTER TABLE "roles_users" DROP CONSTRAINT IF EXISTS "roles_users_role";