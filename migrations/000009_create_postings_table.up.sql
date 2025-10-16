CREATE TABLE "postings" (
  "id" bigserial PRIMARY KEY,
  "posting" text UNIQUE NOT NULL,
  "code" text UNIQUE,
  "created_at" timestamp DEFAULT NOW()
);