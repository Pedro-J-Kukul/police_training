CREATE TABLE "roles_users" (
  "role_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("role_id", "user_id")
);