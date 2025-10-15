CREATE TABLE "roles_permissions" (
  "permission_id" bigint NOT NULL,
  "role_id" bigint NOT NULL,
  PRIMARY KEY ("permission_id", "role_id")
);