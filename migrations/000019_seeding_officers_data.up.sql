-- regions
INSERT INTO regions (region) VALUES
('Northern Region'),
('Eastern Division'),
('Western Region'),
('Southern Region')
ON CONFLICT (region) DO NOTHING;

-- INSERT Ranks
INSERT INTO ranks(rank, code, annual_training_hours) VALUES
('Special Constable', 'SC', 40),
('Constable', 'PC', 40),
('Corporal', 'CPL', 50),
('Sergeant', 'SGT', 60),
('Inspector of Police', 'INSP', 70),
('Assistant Superintendent of Police', 'ASP', 80),
('Superintendent of Police', 'Sr. Supt', 90),
('Assistant Commisioner of Police', 'ACP', 100),
('Deputy Commissioner of Police', 'DCP', 110),
('Commissioner of Police', 'COMPOL', 120),
('Not Applicable', 'NA', 0)
ON CONFLICT (rank) DO NOTHING;

-- Insert Formations
INSERT INTO formations (formation, region_id) VALUES
-- Northern Region
('Corozal Police Formation', (SELECT id FROM regions WHERE region = 'Northern Region')),
('Orange Walk Police Formation', (SELECT id FROM regions WHERE region = 'Northern Region')),
-- Western Region
('Police Headquarters - Belmopan', (SELECT id FROM regions WHERE region = 'Western Region')),
('San Ignacio Police Formation', (SELECT id FROM regions WHERE region = 'Western Region')),
('Belmopan Police Formation', (SELECT id FROM regions WHERE region = 'Western Region')),
('Benque Viejo Police Formation', (SELECT id FROM regions WHERE region = 'Western Region')),
('Roaring Creek Police Sub-Formation', (SELECT id FROM regions WHERE region = 'Western Region')),
-- Eastern Division
('Police Headquarters - Eastern Division', (SELECT id FROM regions WHERE region = 'Eastern Division')),
('Precint 1', (SELECT id FROM regions WHERE region = 'Eastern Division')),
('Precint 2', (SELECT id FROM regions WHERE region = 'Eastern Division')),
('Precint 3', (SELECT id FROM regions WHERE region = 'Eastern Division')),
('Precint 4', (SELECT id FROM regions WHERE region = 'Eastern Division')),
('Ladyville Police Sub-Formation', (SELECT id FROM regions WHERE region = 'Eastern Division')),
('Hattieville Police Sub-Formation', (SELECT id FROM regions WHERE region = 'Eastern Division')),
('Caye Caulker Police Formation', (SELECT id FROM regions WHERE region = 'Eastern Division')),
('San Pedro Police Formation', (SELECT id FROM regions WHERE region = 'Eastern Division')),
-- Southern Region
('Punta Gorda Police Formation', (SELECT id FROM regions WHERE region = 'Southern Region')),
('Intermediate Southern Formation', (SELECT id FROM regions WHERE region = 'Southern Region')),
('Placencia Police Sub-Formation', (SELECT id FROM regions WHERE region = 'Southern Region')),
('Seine Bright Police Sub-Formation', (SELECT id FROM regions WHERE region = 'Southern Region')),
('Hopkins Police Sub-Formation', (SELECT id FROM regions WHERE region = 'Southern Region')),
('Dangriga Police Formation', (SELECT id FROM regions WHERE region = 'Southern Region'))
ON CONFLICT (formation) DO NOTHING;

-- Insert Postings
INSERT INTO postings (posting, code) VALUES
('Relief', ''),
('Staff Duties', ''),
('Station Manager', ''),
('Crimes Investigation Branch', 'CIB'),
('Special Branch', 'SB'),
('Quick Response Team', 'QRT'),
('Prosecution Branch', ''),
('Gang Intelligence, Intediction & Investigation', 'GI3'),
('Anti-Narcotics Unit', 'ANU'),
('Special Patrol Unit', 'SPU'),
('Tourism Police Unit', 'TPU'),
('Major Crimes Unit', 'MCU'),
('Mobile Interdiction Unit', 'MIU'),
('K-9 Unit', 'K9'),
('Professional Standards Branch', 'PSB'),
('Deputy Officer Commanding/ Deputy Commander', ''),
('Officer Commanding/ Commander', ''),
('Regional Commander', ''),
('Special Assignment', ''),
('Other', '')
ON CONFLICT (code) DO NOTHING;
