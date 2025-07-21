CREATE TABLE dtl_payroll (
    id SERIAL,
    id_mst_payroll_period BIGINT NOT NULL,
    total_take_home BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    created_by BIGINT NOT NULL,
    updated_by BIGINT NULL
);