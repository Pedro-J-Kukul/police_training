-- POstings
CREATE TABLE postings (
    "id" BIGSERIAL PRIMARY KEY,
    "posting" text NOT NULL,
    "code" text
);