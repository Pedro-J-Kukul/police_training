CREATE TABLE "enrollment_statuses" (
  "id" bigserial PRIMARY KEY,
  "status" text UNIQUE NOT NULL,
  "created_at" timestamp DEFAULT NOW()
);