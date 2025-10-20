CREATE VIEW "public"."PERMISSIONS" AS 
SELECT 
    r.role AS role_name,
    p.code AS permission_code
FROM 
    roles_permissions rp
    INNER JOIN roles r ON rp.role_id = r.id
    INNER JOIN permissions p ON rp.permission_id = p.id
ORDER BY r.role, p.code;