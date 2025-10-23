CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "first_name" varchar(100) NOT NULL,
  "last_name" varchar(100) NOT NULL,
  "gender" char NOT NULL,
  "email" text UNIQUE NOT NULL,
  "password_hash" bytea NOT NULL,
  "is_facilitator" boolean NOT NULL DEFAULT false,
  "is_officer" boolean NOT NULL DEFAULT false,
  "is_activated" boolean NOT NULL DEFAULT false,
  "is_deleted" boolean NOT NULL DEFAULT false,
  "version" int NOT NULL DEFAULT 1,
  "created_at" timestamp DEFAULT NOW(),
  "updated_at" timestamp DEFAULT NOW()
);