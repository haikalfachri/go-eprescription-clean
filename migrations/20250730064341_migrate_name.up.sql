CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Patients table
CREATE TABLE Patients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR,
    age INT,
    gender VARCHAR CHECK (gender IN ('male', 'female')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- MasterSignas table
CREATE TABLE MasterSignas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    signa VARCHAR UNIQUE, 
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- MasterMedicines table
CREATE TABLE MasterMedicines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR UNIQUE,  
    quantity INT,
    price INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Transactions table
CREATE TABLE Transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    patient_id UUID REFERENCES Patients(id),
    medicine_type TEXT CHECK (medicine_type IN ('compound', 'non_compound')),
    total_price INT,
    total_medicines INT,
    status TEXT DEFAULT 'pending' CHECK (status IN (
        'pending', 
        'paid', 
        'failed',
        'cancelled',
        'expired',
        'challenge'
    )),
    payment_provider TEXT DEFAULT 'other' CHECK (payment_provider IN ('midtrans', 'xendit', 'other')),
    payment_redirect_url TEXT,       
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- MedicineDetails table
CREATE TABLE MedicineDetails (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID REFERENCES Transactions(id),
    signa_id UUID REFERENCES MasterSignas(id),
    medicine_id UUID REFERENCES MasterMedicines(id),
    quantity INT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


