package model

import (
	"database/sql"
	"time"
)

// MstAttendance represents employee attendance records
type MstAttendance struct {
	ID                 int64         `json:"id" xorm:"'id' pk autoincr"`
	IDMstUser          int64         `json:"user_id" xorm:"id_mst_user"`
	AttendanceDate     time.Time     `json:"attendance_date" xorm:"attendance_date"`
	IDMstPayrollPeriod sql.NullInt64 `json:"payroll_period_id" xorm:"id_mst_payroll_period"`
	CreatedAt          time.Time     `json:"created_at" xorm:"'created_at' created"`
	UpdatedAt          time.Time     `json:"updated_at" xorm:"'updated_at' updated"`
	CreatedBy          sql.NullInt64 `json:"created_by,omitempty" xorm:"created_by"`
	UpdatedBy          sql.NullInt64 `json:"updated_by,omitempty" xorm:"updated_by"`
}

type ListAttendanceParams struct {
	IDsMstUser             []int64   `xorm:"id_mst_user"`
	IDMstPayrollPeriod     int64     `xorm:"id_mst_payroll_period"`
	StartDate              time.Time `xorm:"start_date"`
	EndDate                time.Time `xorm:"end_date"`
	IsForGeneratingPayroll bool      `xorm:"is_for_generating_payroll"`
}

type TapInResponse struct {
	AttendanceDate time.Time `json:"attendance_date" xorm:"attendance_date"`
}

// MstPayrollPeriod represents payroll periods
type MstPayrollPeriod struct {
	ID                   int64         `json:"id" xorm:"'id' pk autoincr"`
	StartDate            time.Time     `json:"start_date" xorm:"start_date"`
	EndDate              time.Time     `json:"end_date" xorm:"end_date"`
	PayrollProcessedDate sql.NullTime  `json:"payroll_processed_date" xorm:"payroll_processed_date"`
	CreatedAt            time.Time     `json:"created_at" xorm:"'created_at' created"`
	UpdatedAt            time.Time     `json:"updated_at" xorm:"'updated_at' updated"`
	CreatedBy            sql.NullInt64 `json:"created_by,omitempty" xorm:"created_by"`
	UpdatedBy            sql.NullInt64 `json:"updated_by,omitempty" xorm:"updated_by"`
}

type PayrollPeriodResponse struct {
	ID        int64     `json:"id" xorm:"id"`
	StartDate time.Time `json:"start_date" xorm:"start_date"`
	EndDate   time.Time `json:"end_date" xorm:"end_date"`
}

type PayrollPeriodRequest struct {
	StartDate time.Time `json:"start_date" xorm:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" xorm:"end_date" validate:"required,gtfield=StartDate"`
}

// TrxOvertime represents overtime submissions
type TrxOvertime struct {
	ID                 int64         `xorm:"'id' pk autoincr"`
	UserID             int64         `xorm:"id_mst_user"`
	IDMstPayrollPeriod sql.NullInt64 `xorm:"id_mst_payroll_period"`
	OvertimeDate       time.Time     `xorm:"overtime_date"`
	Hours              int           `xorm:"hours"`
	CreatedAt          time.Time     `xorm:"'created_at' created"`
	UpdatedAt          time.Time     `xorm:"'updated_at' updated"`
	CreatedBy          sql.NullInt64 `xorm:"created_by"`
	UpdatedBy          sql.NullInt64 `xorm:"updated_by"`
}

type GetOvertimeResponse struct {
	OvertimeDate time.Time `json:"overtime_date"`
	Hours        int       `json:"hours"`
}

type SubmitOvertimeRequest struct {
	OvertimeDate time.Time `json:"overtime_date" validate:"required"`
	Hours        int       `json:"hours" validate:"required,gt=0"`
}

type ListOvertimeParams struct {
	StartDate              time.Time
	EndDate                time.Time
	UserIDs                []int64
	IDMstPayrollPeriod     int64
	IsForGeneratingPayroll bool
}

type SubmitOvertimeResponse struct {
	OvertimeDate time.Time `json:"overtime_date"`
	Hours        int       `json:"hours"`
}
