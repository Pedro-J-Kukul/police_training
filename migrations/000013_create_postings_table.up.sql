-- POstings
CREATE TABLE postings (
    "id" BIGSERIAL PRIMARY KEY,
    "posting" text NOT NULL UNIQUE,
    "code" text NOT NULL UNIQUE
);