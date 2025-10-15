-- Add foreign key constraint linking users and officers tables
ALTER TABLE "users" ADD CONSTRAINT "officer_user" FOREIGN KEY ("id") REFERENCES "officers" ("id");