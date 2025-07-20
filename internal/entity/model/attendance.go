package model

import (
	"database/sql"
	"time"
)

// MstAttendance represents employee attendance records
type MstAttendance struct {
	ID             int64         `json:"id" xorm:"'id' pk autoincr"`
	IDMstUser      int64         `json:"user_id" xorm:"id_mst_user"`
	AttendanceDate time.Time     `json:"attendance_date" xorm:"attendance_date"`
	CreatedAt      time.Time     `json:"created_at" xorm:"'created_at' created"`
	UpdatedAt      time.Time     `json:"updated_at" xorm:"'updated_at' updated"`
	CreatedBy      sql.NullInt64 `json:"created_by,omitempty" xorm:"created_by"`
	UpdatedBy      sql.NullInt64 `json:"updated_by,omitempty" xorm:"updated_by"`
}

type TapInResponse struct {
	AttendanceDate time.Time `json:"attendance_date" xorm:"attendance_date"`
}

// MstPayrollPeriod represents payroll periods
type MstPayrollPeriod struct {
	ID                 int64         `json:"id" xorm:"'id' pk autoincr"`
	StartDate          time.Time     `json:"start_date" xorm:"start_date"`
	EndDate            time.Time     `json:"end_date" xorm:"end_date"`
	IsPayrollProcessed bool          `json:"is_payroll_processed" xorm:"is_payroll_processed"`
	CreatedAt          time.Time     `json:"created_at" xorm:"'created_at' created"`
	UpdatedAt          time.Time     `json:"updated_at" xorm:"'updated_at' updated"`
	CreatedBy          sql.NullInt64 `json:"created_by,omitempty" xorm:"created_by"`
	UpdatedBy          sql.NullInt64 `json:"updated_by,omitempty" xorm:"updated_by"`
}

type PayrollPeriodResponse struct {
	ID                 int64     `json:"id" xorm:"id"`
	StartDate          time.Time `json:"start_date" xorm:"start_date"`
	EndDate            time.Time `json:"end_date" xorm:"end_date"`
	IsPayrollProcessed bool      `json:"is_payroll_processed" xorm:"is_payroll_processed"`
}

type PayrollPeriodRequest struct {
	StartDate time.Time `json:"start_date" xorm:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" xorm:"end_date" validate:"required,gtfield=StartDate"`
}

// Overtime represents overtime submissions
type TrxOvertime struct {
	ID                 int64         `json:"id" xorm:"id"`
	UserID             int64         `json:"user_id" xorm:"user_id"`
	IDMstPayrollPeriod int64         `json:"id_mst_payroll_period" xorm:"id_mst_payroll_period"`
	OvertimeDate       time.Time     `json:"overtime_date" xorm:"overtime_date"`
	Hours              float64       `json:"hours" xorm:"hours"`
	CreatedAt          time.Time     `json:"created_at" xorm:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at" xorm:"updated_at"`
	CreatedBy          sql.NullInt64 `json:"created_by,omitempty" xorm:"created_by"`
	UpdatedBy          sql.NullInt64 `json:"updated_by,omitempty" xorm:"updated_by"`
}

// Reimbursement represents reimbursement requests
type TrxReimbursement struct {
	ID                 int64         `json:"id" xorm:"id"`
	UserID             int64         `json:"user_id" xorm:"user_id"`
	IDMstPayrollPeriod int64         `json:"id_mst_payroll_period" xorm:"id_mst_payroll_period"`
	Amount             float64       `json:"amount" xorm:"amount"`
	Description        string        `json:"description" xorm:"description"`
	CreatedAt          time.Time     `json:"created_at" xorm:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at" xorm:"updated_at"`
	CreatedBy          sql.NullInt64 `json:"created_by,omitempty" xorm:"created_by"`
	UpdatedBy          sql.NullInt64 `json:"updated_by,omitempty" xorm:"updated_by"`
}

// Payroll represents processed payroll records
type Payroll struct {
	ID                  int64         `json:"id" xorm:"id"`
	IDMstPayrollPeriod  int64         `json:"id_mst_payroll_period" xorm:"id_mst_payroll_period"`
	UserID              int64         `json:"user_id" xorm:"user_id"`
	BaseSalary          int64         `json:"base_salary" xorm:"base_salary"`
	WorkingDays         int           `json:"working_days" xorm:"working_days"`
	AttendedDays        int           `json:"attended_days" xorm:"attended_days"`
	ProratedSalary      float64       `json:"prorated_salary" xorm:"prorated_salary"`
	OvertimeHours       float64       `json:"overtime_hours" xorm:"overtime_hours"`
	OvertimePay         float64       `json:"overtime_pay" xorm:"overtime_pay"`
	TotalReimbursements float64       `json:"total_reimbursements" xorm:"total_reimbursements"`
	TotalTakeHome       float64       `json:"total_take_home" xorm:"total_take_home"`
	CreatedAt           time.Time     `json:"created_at" xorm:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at" xorm:"updated_at"`
	CreatedBy           sql.NullInt64 `json:"created_by,omitempty" xorm:"created_by"`
	UpdatedBy           sql.NullInt64 `json:"updated_by,omitempty" xorm:"updated_by"`
}
