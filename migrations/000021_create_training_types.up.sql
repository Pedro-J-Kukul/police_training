-- Training Types
CREATE TABLE "training_types" (
    "id" BIGSERIAL PRIMARY KEY,
    "name" TEXT NOT NULL UNIQUE,
    "is_active" BOOLEAN NOT NULL DEFAULT TRUE
);