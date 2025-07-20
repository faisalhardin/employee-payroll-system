CREATE TABLE mst_attendance (
    id SERIAL PRIMARY KEY,
    id_mst_user BIGINT NOT NULL,
    attendance_date DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    created_by BIGINT NULL,
    updated_by BIGINT NULL
);

CREATE INDEX idx_user_id ON mst_attendance (id_mst_user);
CREATE INDEX idx_attendance_date ON mst_attendance (attendance_date);
CREATE INDEX idx_user_attendance ON mst_attendance (id_mst_user, attendance_date);
