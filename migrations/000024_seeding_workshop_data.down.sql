

DELETE FROM training_categories WHERE "name" IN (
'Firearms Training',
'First Aid & Medical',
'Legal & Constitutional',
'Physical Fitness',
'Communication Skills',
'Technology & Cyber',
'Investigation Techniques',
'Community Relations',
'Leadership Development',
'Specialized Operations'
);

DELETE FROM training_types WHERE "name" IN (
'Mandatory',
'Optional',
'Specialized',
'Refresher',
'Certification',
'Leadership'
);