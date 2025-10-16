CREATE TABLE "training_sessions" (
  "id" bigserial PRIMARY KEY,
  "formation_id" bigint NOT NULL,
  "region_id" bigint NOT NULL,
  "facilitator_id" bigint NOT NULL,
  "workshop_id" bigint NOT NULL,
  "session_date" date NOT NULL,
  "start_time" time NOT NULL,
  "end_time" time NOT NULL,
  "location" text,
  "max_capacity" int,
  "training_status_id" bigint NOT NULL,
  "notes" text,
  "created_at" timestamp DEFAULT NOW(),
  "updated_at" timestamp DEFAULT NOW()
);