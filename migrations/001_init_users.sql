-- Create table for MstUser struct
CREATE TABLE IF NOT EXISTS mst_user (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    salary DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_by VARCHAR(255)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_mst_user_username ON mst_user(username);
CREATE INDEX IF NOT EXISTS idx_mst_user_role ON mst_user(role);

-- Insert sample data
INSERT INTO mst_user (username, password_hash, role, salary, created_by) 
VALUES 
    ('admin', '$2a$10$example_hash_here', 'admin', 75000.00, 'system'),
ON CONFLICT (username) DO NOTHING;
