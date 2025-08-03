-- Seed Patients
INSERT INTO Patients (name, age, gender) VALUES
('John Doe', 30, 'male'),
('Jane Smith', 25, 'female'),
('Michael Johnson', 40, 'male'),
('Emily Davis', 35, 'female'),
('Robert Brown', 50, 'male');

-- Seed MasterSignas
INSERT INTO MasterSignas (signa, description) VALUES
('3x1', 'Take one tablet three times a day after meals'),
('2x1', 'Take one tablet twice a day'),
('1x1', 'Take one tablet once a day before bedtime'),
('1x2', 'Take one tablet once in the morning and two tablets at night'),
('2x2', 'Take two tablets twice a day after meals');


-- Seed MasterMedicines
INSERT INTO MasterMedicines (name, quantity, price) VALUES
('Paracetamol', 1000, 15000),
('Amoxicillin', 5000, 10000),
('Cetirizine', 3000, 12000),
('Omeprazole', 2000, 18000),
('Metformin', 4000, 25000),
('Vitamin C', 2000, 30000),
('Ibuprofen', 7500, 20000),
('Lisinopril', 6000, 22000);
