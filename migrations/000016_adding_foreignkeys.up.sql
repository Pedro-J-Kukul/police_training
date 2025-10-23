-- Foreign keys for formations table
ALTER TABLE formations ADD CONSTRAINT fk_region FOREIGN KEY (region_id) REFERENCES regions(id) ON DELETE CASCADE;