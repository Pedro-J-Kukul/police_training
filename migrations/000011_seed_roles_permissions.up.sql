-- Assign all permissions to Admin role
DO $$
DECLARE
    perm RECORD;
    admin_role_id INT;
BEGIN
    -- Get the role_id for the Admin role
    SELECT id INTO admin_role_id FROM roles WHERE role = 'Admin';

    -- Loop through all permissions and assign them to the Admin role
    FOR perm IN SELECT id FROM permissions LOOP
        INSERT INTO roles_permissions (role_id, permission_id)
        VALUES (admin_role_id, perm.id)
        ON CONFLICT DO NOTHING; -- Avoid duplicate inserts
    END LOOP;
END $$;

-- Assign specific permissions to Content-Contributor role
DO $$
DECLARE
    perm_code TEXT;
    cc_role_id INT;
    perm_codes TEXT[] := ARRAY[
        'workshops:create', 'workshops:view', 'workshops:edit',
        'training:categories:create', 'training:categories:view', 'training:categories:edit',
        'training:types:create', 'training:types:view', 'training:types:edit',
        'training:status:create', 'training:status:view', 'training:status:edit',
        'training:sessions:create', 'training:sessions:view', 'training:sessions:edit',
        'training:enrollments:create', 'training:enrollments:view', 'training:enrollments:edit',
        'users:view', 'users:edit',
        'officers:view', 'officers:create', 'officers:edit',
        'self:view', 'self:update'
    ];
BEGIN
    -- Get the role_id for the Content-Contributor role
    SELECT id INTO cc_role_id FROM roles WHERE role = 'Content-Contributor';

    -- Loop through permission codes and assign them
    FOREACH perm_code IN ARRAY perm_codes LOOP
        INSERT INTO roles_permissions (role_id, permission_id)
        VALUES (
            cc_role_id,
            (SELECT id FROM permissions WHERE code = perm_code)
        )
        ON CONFLICT DO NOTHING;
    END LOOP;
END $$;


-- Assign specific permissions to Officer role
DO $$
DECLARE
    perm_code TEXT;
    officer_role_id INT;
    perm_codes TEXT[] := ARRAY[
        'self:view', 'self:update',
        'training:sessions:view',
        'training:enrollments:view',
        'workshops:view'
    ];
BEGIN
    -- Get the role_id for the Officer role
    SELECT id INTO officer_role_id FROM roles WHERE role = 'Officer';

    -- Loop through permission codes and assign them
    FOREACH perm_code IN ARRAY perm_codes LOOP
        INSERT INTO roles_permissions (role_id, permission_id)
        VALUES (
            officer_role_id,
            (SELECT id FROM permissions WHERE code = perm_code)
        )
        ON CONFLICT DO NOTHING;
    END LOOP;
END $$;