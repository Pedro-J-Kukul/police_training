CREATE TABLE "workshops" (
  "id" bigserial PRIMARY KEY,
  "workshop_name" text NOT NULL,
  "category_id" bigint NOT NULL,
  "training_type_id" bigint NOT NULL,
  "credit_hours" int NOT NULL,
  "description" text,
  "objectives" text,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamp DEFAULT NOW(),
  "updated_at" timestamp DEFAULT NOW()
);