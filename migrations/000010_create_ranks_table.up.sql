CREATE TABLE "ranks" (
  "id" bigserial PRIMARY KEY,
  "rank" text UNIQUE NOT NULL,
  "code" text UNIQUE NOT NULL,
  "annual_training_hours_required" int NOT NULL DEFAULT 40
);