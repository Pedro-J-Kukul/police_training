-- Add foreign key constraints for training_sessions table
ALTER TABLE "training_sessions" ADD CONSTRAINT "session_workshop" FOREIGN KEY ("workshop_id") REFERENCES "workshops" ("id");
ALTER TABLE "training_sessions" ADD CONSTRAINT "session_facilitator" FOREIGN KEY ("facilitator_id") REFERENCES "users" ("id");
ALTER TABLE "training_sessions" ADD CONSTRAINT "session_formation" FOREIGN KEY ("formation_id") REFERENCES "formations" ("id");
ALTER TABLE "training_sessions" ADD CONSTRAINT "session_region" FOREIGN KEY ("region_id") REFERENCES "regions" ("id");
ALTER TABLE "training_sessions" ADD CONSTRAINT "session_status" FOREIGN KEY ("training_status_id") REFERENCES "training_status" ("id");