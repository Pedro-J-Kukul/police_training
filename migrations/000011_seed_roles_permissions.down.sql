-- Remove all permissions from the Admin role
DELETE FROM roles_permissions
WHERE role_id = (SELECT id FROM roles WHERE role = 'Admin');


-- Remove specific permissions from the Content-Contributor role
DO $$
DECLARE
    perm_code TEXT;
    cc_role_id INT;
    perm_codes TEXT[] := ARRAY[
        'workshop:create', 'workshop:view', 'workshop:update',
        'training:create', 'training:view', 'training:update',
        'enrollment:create', 'enrollment:view', 'enrollment:update',
        'session:create', 'session:view', 'session:update',
        'user:create', 'user:view', 'user:update',
        'officer:create', 'officer:view', 'officer:update',
        'self:view', 'self:update'
    ];
BEGIN
    -- Get the role_id for the Content-Contributor role
    SELECT id INTO cc_role_id FROM roles WHERE role = 'Content-Contributor';

    -- Loop through the permission codes and remove them
    FOREACH perm_code IN ARRAY perm_codes LOOP
        DELETE FROM roles_permissions
        WHERE role_id = cc_role_id
        AND permission_id = (SELECT id FROM permissions WHERE code = perm_code);
    END LOOP;
END $$;


-- Remove specific permissions from the Officer role
DO $$
DECLARE
    perm_code TEXT;
    officer_role_id INT;
    perm_codes TEXT[] := ARRAY[
        'self:view', 'self:update',
        'training:view',
        'enrollment:view',
        'session:view',
        'workshop:view'
    ];
BEGIN
    -- Get the role_id for the Officer role
    SELECT id INTO officer_role_id FROM roles WHERE role = 'Officer';

    -- Loop through the permission codes and remove them
    FOREACH perm_code IN ARRAY perm_codes LOOP
        DELETE FROM roles_permissions
        WHERE role_id = officer_role_id
        AND permission_id = (SELECT id FROM permissions WHERE code = perm_code);
    END LOOP;
END $$;
