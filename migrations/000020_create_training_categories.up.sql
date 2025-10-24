-- Training Categories
CREATE TABLE training_categories (
    "id" BIGSERIAL PRIMARY KEY,
    "name" TEXT NOT NULL UNIQUE,
    "is_active" BOOLEAN NOT NULL DEFAULT TRUE
);