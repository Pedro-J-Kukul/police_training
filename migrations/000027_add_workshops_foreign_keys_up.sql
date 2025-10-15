-- Add foreign key constraints for workshops table
ALTER TABLE "workshops" ADD CONSTRAINT "workshop_category" FOREIGN KEY ("category_id") REFERENCES "training_categories" ("id");
ALTER TABLE "workshops" ADD CONSTRAINT "workshop_training_type" FOREIGN KEY ("training_type_id") REFERENCES "training_types" ("id");