-- Create Regions table
CREATE TABLE regions (
    "id" BIGSERIAL PRIMARY KEY,
    "region" TEXT NOT NULL UNIQUE
);