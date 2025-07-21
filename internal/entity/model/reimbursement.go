package model

import (
	"database/sql"
	"time"
)

// TrxReimbursement represents reimbursement requests
type TrxReimbursement struct {
	ID                 int64         `xorm:"'id' pk autoincr"`
	UserID             int64         `xorm:"'id_mst_user'"`
	IDMstPayrollPeriod sql.NullInt64 `xorm:"id_mst_payroll_period"`
	Status             string        `xorm:"'status'"`
	Amount             int64         `xorm:"amount"`
	Description        string        `xorm:"description"`
	CreatedAt          time.Time     `xorm:"'created_at' created"`
	UpdatedAt          time.Time     `xorm:"'updated_at' updated"`
	CreatedBy          sql.NullInt64 `xorm:"created_by"`
	UpdatedBy          sql.NullInt64 `xorm:"updated_by"`
}

type SubmitReimbursementRequest struct {
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

type SubmitReimbursementResponse struct {
	ID          int64  `json:"id"`
	Amount      int64  `json:"amount"`
	Status      string `json:"status"`
	Description string `json:"description"`
}
