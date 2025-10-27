CREATE TABLE workshops (
    "id" BIGSERIAL PRIMARY KEY,
    "workshop_name" TEXT NOT NULL UNIQUE,
    "category_id" BIGINT NOT NULL,
    "type_id" BIGINT NOT NULL,
    "credit_hours" INT NOT NULL,
    "description" TEXT,
    "is_active" BOOLEAN NOT NULL DEFAULT TRUE,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);