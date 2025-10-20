DELETE FROM "roles_permissions"
WHERE ("permission_id", "role_id") IN (
    ---- Admin-level permissions
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_CREATE'),
        (SELECT id FROM "roles" WHERE "role" = 'ADMIN')
    ),
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_READ'),
        (SELECT id FROM "roles" WHERE "role" = 'ADMIN')
    ),
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_MODIFY'),
        (SELECT id FROM "roles" WHERE "role" = 'ADMIN')
    ),
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_DELETE'),
        (SELECT id FROM "roles" WHERE "role" = 'ADMIN')
    )

    ---- User-level permissions
    ,
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_READ'),
        (SELECT id FROM "roles" WHERE "role" = 'USER')
    ),
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_READ_SELF'),
        (SELECT id FROM "roles" WHERE "role" = 'USER')
    ),
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_MODIFY_SELF'),
        (SELECT id FROM "roles" WHERE "role" = 'USER')
    ),
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_DELETE_SELF'),
        (SELECT id FROM "roles" WHERE "role" = 'USER')
    ),

    --- Trainer-level permissions
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_READ_TRAININGS'),
        (SELECT id FROM "roles" WHERE "role" = 'TRAINER')
    ),
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_ASSIGN_TRAININGS'),
        (SELECT id FROM "roles" WHERE "role" = 'TRAINER')
    ),
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_READ_WORKSHOPS'),
        (SELECT id FROM "roles" WHERE "role" = 'TRAINER')
    ),
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_CREATE_WORKSHOPS'),
        (SELECT id FROM "roles" WHERE "role" = 'TRAINER')
    ),
    --- Content-creator-level permissions
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_CREATE_WORKSHOPS'),
        (SELECT id FROM "roles" WHERE "role" = 'CONTENT_CREATOR')
    ),
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_READ_WORKSHOPS'),
        (SELECT id FROM "roles" WHERE "role" = 'CONTENT_CREATOR')
    ),
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_MODIFY_WORKSHOPS'),
        (SELECT id FROM "roles" WHERE "role" = 'CONTENT_CREATOR')
    ),
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_DELETE_WORKSHOPS'),
        (SELECT id FROM "roles" WHERE "role" = 'CONTENT_CREATOR')
    ),
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_CREATE_TRAININGS'),
        (SELECT id FROM "roles" WHERE "role" = 'CONTENT_CREATOR')
    ),
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_READ_TRAININGS'),
        (SELECT id FROM "roles" WHERE "role" = 'CONTENT_CREATOR')
    ),
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_MODIFY_TRAININGS'),
        (SELECT id FROM "roles" WHERE "role" = 'CONTENT_CREATOR')
    ),
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_DELETE_TRAININGS'),
        (SELECT id FROM "roles" WHERE "role" = 'CONTENT_CREATOR')
    ),
    --- Officer-level permissions
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_READ_WORKSHOPS'),
        (SELECT id FROM "roles" WHERE "role" = 'OFFICER')
    ),
    (
        (SELECT id FROM "permissions" WHERE "code" = 'CAN_READ_TRAININGS'),
        (SELECT id FROM "roles" WHERE "role" = 'OFFICER')
    ),
    --- Anonymous-level permissions
    --- Nothing
);
