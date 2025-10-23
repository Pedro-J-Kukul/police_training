-- Foreign Keys
ALTER TABLE officers ADD CONSTRAINT fk_rank FOREIGN KEY (rank_id) REFERENCES ranks(id) ON DELETE CASCADE;
ALTER TABLE officers ADD CONSTRAINT fk_posting FOREIGN KEY (posting_id) REFERENCES postings(id) ON DELETE CASCADE;
ALTER TABLE officers ADD CONSTRAINT fk_formation FOREIGN KEY (formation_id) REFERENCES formations(id) ON DELETE CASCADE;
ALTER TABLE officers ADD CONSTRAINT fk_region FOREIGN KEY (region_id) REFERENCES regions(id) ON DELETE CASCADE;
ALTER TABLE officers ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;


CREATE INDEX idx_officers_user_id ON officers(user_id);