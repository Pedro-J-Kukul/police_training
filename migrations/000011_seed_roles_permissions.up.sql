-- Give Admin all permissions
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
        ON CONFLICT DO NOTHING;  -- avoids duplicate inserts
    END LOOP;
END $$;


-- Give Content-Contributor specific permissions (view, edit, create)
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


-- Give Officer specific permissions (view/edit self, view trainings/enrollments/workshops)
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
