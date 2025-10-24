
-- Insert Enrollment Status
INSERT INTO enrollment_statuses (status) VALUES 
('Enrolled'),
('Waitlisted'),
('Confirmed'),
('Cancelled'),
('No Show'),
('Completed')
ON CONFLICT (status) DO NOTHING;

-- Insert Progress Status
INSERT INTO progress_statuses (status) VALUES 
('Not Started'),
('In Progress'),
('Completed'),
('Failed'),
('Withdrawn')
ON CONFLICT (status) DO NOTHING;

-- Insert Attendance Status
INSERT INTO attendance_statuses (status, counts_as_present) VALUES 
('Present', true),
('Absent', false),
('Late', true),
('Excused', false),
('Sick Leave', false),
('Emergency Leave', false)
ON CONFLICT (status) DO NOTHING;