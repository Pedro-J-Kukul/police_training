-- Removing foreign keys from formations table
ALTER TABLE formations DROP CONSTRAINT IF EXISTS fk_region;