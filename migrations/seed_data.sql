-- Seed Patients
INSERT INTO Patients (name, age, gender) VALUES
('John Doe', 30, 'male'),
('Jane Smith', 25, 'female'),
('Michael Johnson', 40, 'male'),
('Emily Davis', 35, 'female');

-- Seed MasterSignas
INSERT INTO MasterSignas (signa, description) VALUES
('3x1', 'Take one tablet three times a day after meals'),
('2x1', 'Take one tablet twice a day'),
('1x1', 'Take one tablet once a day before bedtime');

-- Seed MasterMedicines
INSERT INTO MasterMedicines (name, quantity, price) VALUES
('Paracetamol', 100, 500),
('Amoxicillin', 50, 1000),
('Vitamin C', 200, 300),
('Ibuprofen', 75, 800);
