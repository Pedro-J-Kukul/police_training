CREATE TABLE  IF NOT EXISTS "officers" (
  "id" bigserial PRIMARY KEY,
  "regulation_number" text NOT NULL,
  "posting_id" bigint NOT NULL,
  "rank_id" bigint NOT NULL,
  "formation_id" bigint NOT NULL,
  "region_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "created_at" timestamp DEFAULT NOW(),
  "updated_at" timestamp DEFAULT NOW()
);