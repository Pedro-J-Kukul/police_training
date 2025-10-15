-- Add foreign key constraints for roles_permissions table
ALTER TABLE "roles_permissions" ADD CONSTRAINT "roles_permissions_role" FOREIGN KEY ("role_id") REFERENCES "roles" ("id");
ALTER TABLE "roles_permissions" ADD CONSTRAINT "roles_permissions_permission" FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id");