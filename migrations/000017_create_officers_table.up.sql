CREATE TABLE "officers" (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "user_id" BIGINT NOT NULL UNIQUE, -- One-to-one constraint
    "regulation_number" VARCHAR(50) UNIQUE NOT NULL,
    "rank_id" BIGINT NOT NULL,
    "posting_id" BIGINT NOT NULL,
    "formation_id" BIGINT NOT NULL,
    "region_id" BIGINT NOT NULL,
    "created_at" timestamp DEFAULT NOW(),
    "updated_at" timestamp DEFAULT NOW()
);