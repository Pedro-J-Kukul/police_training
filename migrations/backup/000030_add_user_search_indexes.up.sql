CREATE INDEX IF NOT EXISTS users_first_name_search_idx
    ON users USING GIN (to_tsvector('simple', first_name));

CREATE INDEX IF NOT EXISTS users_last_name_search_idx
    ON users USING GIN (to_tsvector('simple', last_name));

CREATE INDEX IF NOT EXISTS users_email_search_idx
    ON users USING GIN (to_tsvector('simple', email));

CREATE INDEX IF NOT EXISTS tokens_hash_scope_idx
    ON tokens (hash, scope);
