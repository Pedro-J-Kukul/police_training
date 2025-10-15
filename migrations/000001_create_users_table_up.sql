CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "first_name" varchar(100) NOT NULL,
  "last_name" varchar(100) NOT NULL,
  "gender" char NOT NULL,
  "email" text UNIQUE NOT NULL,
  "password_hash" bytea NOT NULL,
  "is_activated" bool NOT NULL DEFAULT false,
  "is_facilitator" bool NOT NULL DEFAULT false,
  "version" int NOT NULL DEFAULT 1,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL
);