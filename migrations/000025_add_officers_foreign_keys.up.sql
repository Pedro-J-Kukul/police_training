-- Add foreign key constraints for officers table
ALTER TABLE "officers" ADD CONSTRAINT "officer_formation" FOREIGN KEY ("formation_id") REFERENCES "formations" ("id");
ALTER TABLE "officers" ADD CONSTRAINT "officer_rank" FOREIGN KEY ("rank_id") REFERENCES "ranks" ("id");
ALTER TABLE "officers" ADD CONSTRAINT "officer_posting" FOREIGN KEY ("posting_id") REFERENCES "postings" ("id");
ALTER TABLE "officers" ADD CONSTRAINT "officer_region" FOREIGN KEY ("region_id") REFERENCES "regions" ("id");
ALTER TABLE "officers" ADD CONSTRAINT "officer_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE SET NULL;