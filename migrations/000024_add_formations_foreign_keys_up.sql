-- Add foreign key constraints for formations table
ALTER TABLE "formations" ADD CONSTRAINT "formation_region" FOREIGN KEY ("region_id") REFERENCES "regions" ("id");