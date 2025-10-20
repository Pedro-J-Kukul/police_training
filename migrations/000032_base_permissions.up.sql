--- Create base permissions ---
INSERT INTO "permissions" ("code") VALUES
    ('CAN_CREATE'), -- Admin permissions
    ('CAN_READ'),
    ('CAN_MODIFY'),
    ('CAN_DELETE'), -- Self permissions

    ('CAN_READ_SELF'), --- Domain-specific user permissions
    ('CAN_MODIFY_SELF'),
    ('CAN_DELETE_SELF'), --- Domain-specific user management permissions

    ('CAN_CREATE_USER'),
    ('CAN_READ_USER'),
    ('CAN_MODIFY_USER'), 
    ('CAN_DELETE_USER'), --- Domain-specific formation permissions

    ('CAN_CREATE_FORMATIONS'),
    ('CAN_READ_FORMATIONS'),
    ('CAN_MODIFY_FORMATIONS'),
    ('CAN_DELETE_FORMATIONS'), --- Domain-specific officer permissions

    ('CAN_CREATE_OFFICERS'),
    ('CAN_READ_OFFICERS'),
    ('CAN_MODIFY_OFFICERS'),
    ('CAN_DELETE_OFFICERS'), --- Domain-specific postings permissions

    ('CAN_CREATE_POSTINGS'),
    ('CAN_READ_POSTINGS'),
    ('CAN_MODIFY_POSTINGS'),
    ('CAN_DELETE_POSTINGS'), --- Domain-specific permissions permissions

    ('CAN_CREATE_PERMISSIONS'),
    ('CAN_READ_PERMISSIONS'),
    ('CAN_MODIFY_PERMISSIONS'),
    ('CAN_DELETE_PERMISSIONS'), --- Domain-specific ranks permissions

    ('CAN_CREATE_RANKS'),
    ('CAN_READ_RANKS'),
    ('CAN_MODIFY_RANKS'),
    ('CAN_DELETE_RANKS'), --- Domain-specific regions permissions

    ('CAN_CREATE_REGIONS'),
    ('CAN_READ_REGIONS'),
    ('CAN_MODIFY_REGIONS'),
    ('CAN_DELETE_REGIONS'), --- Domain-specific roles permissions

    ('CAN_CREATE_ROLES'),
    ('CAN_READ_ROLES'),
    ('CAN_MODIFY_ROLES'),
    ('CAN_DELETE_ROLES'), --- Domain-specific trainings permissions

    ('CAN_CREATE_TRAININGS'),
    ('CAN_READ_TRAININGS'),
    ('CAN_MODIFY_TRAININGS'),
    ('CAN_DELETE_TRAININGS'),
    ('CAN_ASSIGN_TRAININGS'), --- Domain-specific training_categories permissions

    ('CAN_CREATE_TRAINING_CATEGORIES'),
    ('CAN_READ_TRAINING_CATEGORIES'),
    ('CAN_MODIFY_TRAINING_CATEGORIES'),
    ('CAN_DELETE_TRAINING_CATEGORIES'), --- Domain-specific workshops permissions

    ('CAN_CREATE_WORKSHOPS'),
    ('CAN_READ_WORKSHOPS'),
    ('CAN_MODIFY_WORKSHOPS'),
    ('CAN_DELETE_WORKSHOPS'), --- Domain-specific enrollments permissions
    
    ('CAN_CREATE_ENROLLMENTS'),
    ('CAN_READ_ENROLLMENTS'),
    ('CAN_MODIFY_ENROLLMENTS'),
    ('CAN_DELETE_ENROLLMENTS'),  --- Domain-specific enrollment_statuses permissions

    ('CAN_CREATE_ENROLLMENT_STATUSES'),
    ('CAN_READ_ENROLLMENT_STATUSES'),
    ('CAN_MODIFY_ENROLLMENT_STATUSES'),
    ('CAN_DELETE_ENROLLMENT_STATUSES'), --- Domain-specific training_sessions permissions

    ('CAN_CREATE_TRAINING_SESSIONS'),
    ('CAN_READ_TRAINING_SESSIONS'),
    ('CAN_MODIFY_TRAINING_SESSIONS'),
    ('CAN_DELETE_TRAINING_SESSIONS'), --- Domain-specific training_enrollments permissions

    ('CAN_CREATE_TRAINING_ENROLLMENTS'),
    ('CAN_READ_TRAINING_ENROLLMENTS'),
    ('CAN_MODIFY_TRAINING_ENROLLMENTS'),
    ('CAN_DELETE_TRAINING_ENROLLMENTS')
    ;
