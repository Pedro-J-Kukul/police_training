-- Foreign Keys for roles_users and tokens
ALTER TABLE "roles_users" ADD CONSTRAINT "user_roles_users" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "roles_users" ADD CONSTRAINT "role_roles_users" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON DELETE CASCADE;
ALTER TABLE "tokens" ADD CONSTRAINT "user_tokens" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;