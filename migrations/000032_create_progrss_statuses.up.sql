CREATE TABLE progress_statuses (
    id SERIAL PRIMARY KEY,
    status TEXT NOT NULL UNIQUE
);