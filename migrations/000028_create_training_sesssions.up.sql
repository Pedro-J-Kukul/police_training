CREATE TABLE "training_sessions" (
    "id" bigserial PRIMARY KEY,
    "facilitator_id" bigint NOT NULL,
    "workshop_id" bigint NOT NULL,
    "formation_id" bigint NOT NULL,
    "region_id" bigint NOT NULL,
    "training_status_id" bigint NOT NULL,
    "session_date" date NOT NULL,
    "start_time" time NOT NULL,
    "end_time" time NOT NULL,
    "location" text,
    "max_capacity" int,
    "notes" text,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);