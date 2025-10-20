-- Drop base roles --
DELETE FROM "roles" WHERE "role" IN (
    'ADMIN',
    'USER',
    'TRAINER',
    'OFFICER',
    'CONTENT_CREATOR',
    'ANONYMOUS'
);