-- Remove foreign key constraints for roles_permissions table
ALTER TABLE "roles_permissions" DROP CONSTRAINT IF EXISTS "roles_permissions_permission";
ALTER TABLE "roles_permissions" DROP CONSTRAINT IF EXISTS "roles_permissions_role";