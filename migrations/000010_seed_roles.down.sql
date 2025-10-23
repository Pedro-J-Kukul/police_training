DELETE FROM roles
WHERE role IN (
    'Admin',
    'Content-Contributor',
    'Officer'
);