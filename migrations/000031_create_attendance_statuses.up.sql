CREATE TABLE attendance_statuses (
    id SERIAL PRIMARY KEY,
    status TEXT NOT NULL UNIQUE,
    counts_as_present BOOLEAN NOT NULL DEFAULT FALSE
);