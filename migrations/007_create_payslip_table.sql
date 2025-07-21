CREATE TABLE trx_user_payslip (
    id SERIAL,
    id_mst_user BIGINT NOT NULL,
    username VARCHAR(100) NOT NULL,
    id_mst_payroll_period BIGINT NOT NULL,
    base_salary BIGINT NOT NULL,
    working_days INTEGER NOT NULL,
    attended_days INTEGER NOT NULL,
    prorated_salary BIGINT NOT NULL,
    overtime_hours INTEGER NOT NULL,
    overtime_pay BIGINT NOT NULL,
    total_reimbursements BIGINT NOT NULL,
    total_take_home BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    created_by BIGINT NULL,
    updated_by BIGINT NULL
);

CREATE INDEX idx_payslip_user 
ON trx_user_payslip (id_mst_user);

CREATE INDEX idx_payslip_payroll_period 
ON trx_user_payslip (id_mst_payroll_period);

CREATE INDEX idx_payslip_user_period 
ON trx_user_payslip (id_mst_user, id_mst_payroll_period);