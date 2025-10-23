-- Create Formations table
CREATE TABLE formations (
    "id" BIGSERIAL PRIMARY KEY,
    "formation" TEXT NOT NULL UNIQUE,
    "region_id" BIGINT NOT NULL 
);