CREATE TABLE mst_payroll_period (
    id SERIAL PRIMARY KEY,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    payroll_processed_date TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    created_by BIGINT NULL,
    updated_by BIGINT NULL
);

create INDEX idx_payroll_dates on mst_payroll_period (start_date, end_date);
CREATE INDEX idx_payroll_processed_status ON mst_payroll_period (payroll_processed_date);
CREATE INDEX idx_payroll_created_at on mst_payroll_period (created_at);