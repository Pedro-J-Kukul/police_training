CREATE TABLE "formations" (
  "id" bigserial PRIMARY KEY,
  "formation" text UNIQUE NOT NULL,
  "region_id" bigint NOT NULL,
  "created_at" timestamp NOT NULL
);