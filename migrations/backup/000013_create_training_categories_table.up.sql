CREATE TABLE "training_categories" (
  "id" bigserial PRIMARY KEY,
  "name" text UNIQUE NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamp DEFAULT NOW()
);