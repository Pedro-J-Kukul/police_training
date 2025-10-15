CREATE TABLE "attendance_statuses" (
  "id" bigserial PRIMARY KEY,
  "status" text UNIQUE NOT NULL,
  "counts_as_present" boolean NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL
);