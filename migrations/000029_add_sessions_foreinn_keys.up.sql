
ALTER TABLE "training_sessions"
ADD CONSTRAINT training_sessions_facilitator_id_fkey FOREIGN KEY (facilitator_id) REFERENCES "users"(id) ON DELETE CASCADE;   
ALTER TABLE "training_sessions"
ADD CONSTRAINT training_sessions_workshop_id_fkey FOREIGN KEY (workshop_id) REFERENCES "workshops"(id) ON DELETE CASCADE;   
ALTER TABLE "training_sessions"
ADD CONSTRAINT training_sessions_formation_id_fkey FOREIGN KEY (formation_id) REFERENCES "formations"(id) ON DELETE CASCADE;   
ALTER TABLE "training_sessions"
ADD CONSTRAINT training_sessions_region_id_fkey FOREIGN KEY (region_id) REFERENCES "regions"(id) ON DELETE CASCADE;   
ALTER TABLE "training_sessions"
ADD CONSTRAINT training_sessions_training_status_id_fkey FOREIGN KEY (training_status_id) REFERENCES "training_status"(id) ON DELETE CASCADE;