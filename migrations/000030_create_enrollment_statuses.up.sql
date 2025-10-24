CREATE TABLE "enrollment_statuses" (
  "id" bigserial PRIMARY KEY,
  "status" text UNIQUE NOT NULL
);