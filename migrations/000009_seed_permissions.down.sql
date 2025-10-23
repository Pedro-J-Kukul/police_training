DELETE FROM permissions
WHERE code IN (
-- Self management permissions
    'self:view',
    'self:update',
-- User management permissions
    'user:create',
    'user:view',
    'user:update',
    'user:delete',
    'user:reset_password',
-- Officer management permissions
    'officer:create',
    'officer:view',
    'officer:update',
    'officer:delete',
-- Role management permissions
    'role:create',
    'role:view',
    'role:update',
    'role:delete',
-- Permission management permissions
    'permission:create',
    'permission:view',
    'permission:update',
    'permission:delete',
-- Ranks management permissions
    'rank:create',
    'rank:view',
    'rank:update',
    'rank:delete',
-- Formations management permissions
    'formation:create',
    'formation:view',
    'formation:update',
    'formation:delete',
-- Postings management permissions
    'posting:create',
    'posting:view',
    'posting:update',
    'posting:delete',
-- Regions management permissions
    'region:create',
    'region:view',
    'region:update',
    'region:delete',
-- Workshops management permissions
    'workshop:create',
    'workshop:view',
    'workshop:update',
    'workshop:delete',
-- Trainings management permissions
    'training:create',
    'training:view',
    'training:update',
    'training:delete',
-- Enrollment management permissions
    'enrollment:create',
    'enrollment:view',
    'enrollment:update',
    'enrollment:delete'
-- Sessions management permissions
    ,'session:create'
    ,'session:view'
    ,'session:update'
    ,'session:delete'
);