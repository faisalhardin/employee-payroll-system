CREATE TABLE trx_overtime (
    id BIGSERIAL PRIMARY KEY,
    id_mst_user BIGINT NOT NULL,
    id_mst_payroll_period BIGINT NULL,
    overtime_date DATE NOT NULL,
    hours INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT NULL,
    updated_by BIGINT NULL
);

CREATE INDEX idx_overtime_user ON trx_overtime(id_mst_user);
CREATE INDEX idx_overtime_date ON trx_overtime(overtime_date);
CREATE INDEX idx_overtime_payroll_period ON trx_overtime(id_mst_payroll_period);
CREATE INDEX idx_overtime_user_date ON trx_overtime(id_mst_user, overtime_date);

-- Insert 50 overtime records with random hours between 1-4 for different users and dates
INSERT INTO trx_overtime (id_mst_user, id_mst_payroll_period, overtime_date, hours, created_by) VALUES
(2, NULL, '2025-06-20', 2, 1),
(5, NULL, '2025-06-20', 3, 1),
(8, NULL, '2025-06-21', 1, 1),
(12, NULL, '2025-06-21', 2, 1),
(15, NULL, '2025-06-22', 4, 1),
(18, NULL, '2025-06-22', 2, 1),
(22, NULL, '2025-06-23', 3, 1),
(25, NULL, '2025-06-23', 1, 1),
(28, NULL, '2025-06-24', 2, 1),
(32, NULL, '2025-06-24', 4, 1),
(35, NULL, '2025-06-25', 2, 1),
(38, NULL, '2025-06-25', 3, 1),
(42, NULL, '2025-06-26', 1, 1),
(45, NULL, '2025-06-26', 2, 1),
(48, NULL, '2025-06-27', 4, 1),
(52, NULL, '2025-06-27', 2, 1),
(55, NULL, '2025-06-28', 3, 1),
(58, NULL, '2025-06-28', 1, 1),
(62, NULL, '2025-06-29', 2, 1),
(65, NULL, '2025-06-29', 4, 1),
(68, NULL, '2025-06-30', 2, 1),
(72, NULL, '2025-06-30', 3, 1),
(75, NULL, '2025-07-01', 1, 1),
(78, NULL, '2025-07-01', 2, 1),
(82, NULL, '2025-07-02', 4, 1),
(85, NULL, '2025-07-02', 2, 1),
(88, NULL, '2025-07-03', 3, 1),
(92, NULL, '2025-07-03', 1, 1),
(95, NULL, '2025-07-04', 2, 1),
(98, NULL, '2025-07-04', 4, 1),
(3, NULL, '2025-07-05', 2, 1),
(6, NULL, '2025-07-05', 3, 1),
(9, NULL, '2025-07-06', 1, 1),
(13, NULL, '2025-07-06', 2, 1),
(16, NULL, '2025-07-07', 4, 1),
(19, NULL, '2025-07-07', 2, 1),
(23, NULL, '2025-07-08', 3, 1),
(26, NULL, '2025-07-08', 1, 1),
(29, NULL, '2025-07-09', 2, 1),
(33, NULL, '2025-07-09', 4, 1),
(36, NULL, '2025-07-10', 2, 1),
(39, NULL, '2025-07-10', 3, 1),
(43, NULL, '2025-07-11', 1, 1),
(46, NULL, '2025-07-11', 2, 1),
(49, NULL, '2025-07-12', 4, 1),
(53, NULL, '2025-07-12', 2, 1),
(56, NULL, '2025-07-13', 3, 1),
(59, NULL, '2025-07-13', 1, 1),
(63, NULL, '2025-07-14', 2, 1),
(66, NULL, '2025-07-14', 4, 1);
