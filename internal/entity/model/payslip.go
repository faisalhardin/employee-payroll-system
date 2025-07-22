package model

import (
	"database/sql"
	"time"
)

// TrxUserPayslip represents payslip data for employees
type TrxUserPayslip struct {
	ID                  int64         `xorm:"'id' pk autoincr" json:"-"`
	UserID              int64         `xorm:"id_mst_user" json:"-"`
	Username            string        `xorm:"username" json:"username"`
	IDMstPayrollPeriod  int64         `xorm:"id_mst_payroll_period" json:"payrol_period_id"`
	BaseSalary          int64         `xorm:"base_salary" json:"base_salary"`
	WorkingDays         int           `xorm:"working_days" json:"working_days"`
	AttendedDays        int           `xorm:"attended_days" json:"attended_days"`
	ProratedSalary      int64         `xorm:"prorated_salary" json:"prorated_salary"`
	OvertimeHours       int           `xorm:"overtime_hours" json:"overtime_hours"`
	OvertimePay         int64         `xorm:"overtime_pay" json:"overtime_pay"`
	TotalReimbursements int64         `xorm:"total_reimbursements" json:"total_reimbursements"`
	TotalTakeHome       int64         `xorm:"total_take_home" json:"total_take_home_pay"`
	CreatedAt           time.Time     `xorm:"'created_at' created" json:"-"`
	UpdatedAt           time.Time     `xorm:"'updated_at' updated" json:"-"`
	CreatedBy           sql.NullInt64 `xorm:"created_by" json:"-"`
	UpdatedBy           sql.NullInt64 `xorm:"updated_by" json:"-"`
}

type GeneratePayrollRequest struct {
	IDMstPayrollPeriod int `json:"payroll_period_id"`
}

type GetPayrollRequest struct {
	IDMstPayrollPeriod int64 `schema:"payroll_period_id"`
}

type GetPayrollResponse struct {
	StartDate        time.Time        `json:"start_date"`
	EndDate          time.Time        `json:"end_date"`
	EmployeesPayslip []TrxUserPayslip `json:"employees_payslip"`
	TotalTakeHomePay int64            `json:"total_take_home_pay"`
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

type GetDtlPayrollRequest struct {
	IDMstPayrollPeriod int64
}

type GetPayslipResponse struct {
	StartDate           time.Time                     `json:"start_date"`
	EndDate             time.Time                     `json:"end_date"`
	TotalTakeHomePay    int64                         `json:"total_take_home_pay"`
	AttendanceDate      []string                      `json:"attendance_date"`
	WorkingDays         int                           `json:"working_days"`
	AttendedDays        int                           `json:"attended_days"`
	ProratedSalary      int64                         `json:"prorated_salary"`
	OvertimeHours       int                           `json:"overtime_hours"`
	OvertimePay         int64                         `json:"overtime_pay"`
	OvertimeDetails     []GetOvertimeResponse         `json:"overtime_details"`
	ReimbursementList   []SubmitReimbursementResponse `json:"reimbursement_list"`
	TotalReimbursements int64                         `json:"total_reimbursements"`
}

type GetPayslipRequest struct {
	IDMstPayrollPeriod int64 `schema:"payroll_period_id"`
	UserID             int64 `json:"-"`
}
