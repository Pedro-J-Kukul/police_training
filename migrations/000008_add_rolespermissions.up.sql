-- For Roles and permissions table

ALTER TABLE "roles_permissions" ADD CONSTRAINT "role_roles_permissions" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON DELETE CASCADE;
ALTER TABLE "roles_permissions" ADD CONSTRAINT "permission_roles_permissions" FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id") ON DELETE CASCADE;