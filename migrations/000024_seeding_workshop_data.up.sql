-- Insert Training Types
INSERT INTO training_types ("name", is_active) VALUES 
('Mandatory', true),
('Optional', true),
('Specialized', true),
('Refresher', true),
('Certification', true),
('Leadership' , true);

-- Insert Training Categories
INSERT INTO training_categories ("name", is_active) VALUES 
('Firearms Training', true),
('First Aid & Medical', true),
('Legal & Constitutional', true),
('Physical Fitness', true),
('Communication Skills', true),
('Technology & Cyber', true),
('Investigation Techniques', true),
('Community Relations', true),
('Leadership Development', false),
('Specialized Operations', true);

