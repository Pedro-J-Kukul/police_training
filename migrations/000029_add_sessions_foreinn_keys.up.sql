-- CREATE TABLE "training_sessions" (
--     "id" bigserial PRIMARY KEY,
--     "facilitator_id" bigint NOT NULL,
--     "workshop_id" bigint NOT NULL,
--     "formation_id" bigint NOT NULL,
--     "region_id" bigint NOT NULL,
--     "session_date" date NOT NULL,
--     "start_time" time NOT NULL,
--     "end_time" time NOT NULL,
--     "location" text,
--     "max_capacity" int,
--     "training_status_id" bigint NOT NULL,
--     "notes" text,
--     "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
-- -- );
ALTER TABLE "training_sessions"
ADD CONSTRAINT training_sessions_facilitator_id_fkey FOREIGN KEY (facilitator_id) REFERENCES "instructors"(id) ON DELETE CASCADE;   
ALTER TABLE "training_sessions"
ADD CONSTRAINT training_sessions_workshop_id_fkey FOREIGN KEY (workshop_id) REFERENCES "workshops"(id) ON DELETE CASCADE;   
ALTER TABLE "training_sessions"
ADD CONSTRAINT training_sessions_formation_id_fkey FOREIGN KEY (formation_id) REFERENCES "formations"(id) ON DELETE CASCADE;   
ALTER TABLE "training_sessions"
ADD CONSTRAINT training_sessions_region_id_fkey FOREIGN KEY (region_id) REFERENCES "regions"(id) ON DELETE CASCADE;   
ALTER TABLE "training_sessions"
ADD CONSTRAINT training_sessions_training_status_id_fkey FOREIGN KEY (training_status_id) REFERENCES "