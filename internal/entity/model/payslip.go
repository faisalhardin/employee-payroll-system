package model

import (
	"database/sql"
	"time"
)

// TrxUserPayslip represents payslip data for employees
type TrxUserPayslip struct {
	ID                  int64         `xorm:"'id' pk autoincr"`
	UserID              int64         `xorm:"id_mst_user"`
	Username            string        `xorm:"username"`
	IDMstPayrollPeriod  int64         `xorm:"id_mst_payroll_period"`
	BaseSalary          int64         `xorm:"base_salary"`
	WorkingDays         int           `xorm:"working_days"`
	AttendedDays        int           `xorm:"attended_days"`
	ProratedSalary      int64         `xorm:"prorated_salary"`
	OvertimeHours       int           `xorm:"overtime_hours"`
	OvertimePay         int64         `xorm:"overtime_pay"`
	TotalReimbursements int64         `xorm:"total_reimbursements"`
	TotalTakeHome       int64         `xorm:"total_take_home"`
	CreatedAt           time.Time     `xorm:"'created_at' created"`
	UpdatedAt           time.Time     `xorm:"'updated_at' updated"`
	CreatedBy           sql.NullInt64 `xorm:"created_by"`
	UpdatedBy           sql.NullInt64 `xorm:"updated_by"`
}

type GeneratePayrollRequest struct {
	IDMstPayrollPeriod int `json:"payroll_period_id"`
}

// DtlPayroll represents summary for companywide payroll
type DtlPayroll struct {
	ID                 int64         `xorm:"'id' pk autoincr"`
	IDMstPayrollPeriod int64         `xorm:"id_mst_payroll_period"`
	TotalTakeHome      int64         `xorm:"total_take_home"`
	CreatedAt          time.Time     `xorm:"created_at"`
	UpdatedAt          time.Time     `xorm:"updated_at"`
	CreatedBy          int64         `xorm:"created_by"`
	UpdatedBy          sql.NullInt64 `xorm:"updated_by"`
}
