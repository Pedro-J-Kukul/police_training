-- Ranks
CREATE TABLE "ranks" (
    "id" BIGSERIAL PRIMARY KEY,
    "rank" text NOT NULL UNIQUE,
    "code" text NOT NULL UNIQUE,
    "annual_training_hours" integer NOT NULL
);