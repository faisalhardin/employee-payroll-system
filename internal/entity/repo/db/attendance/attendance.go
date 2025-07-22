package attendance

import (
	"context"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -self_package=github.com/faisalhardin/employee-payroll-system/internal/entity/repo/db/attendance -destination=../_mocks/attendance/mock_attendance.go -package=attendance github.com/faisalhardin/employee-payroll-system/internal/entity/repo/db/attendance AttendanceRepository
type AttendanceRepository interface {
	RecordAttendance(ctx context.Context, attendance *model.MstAttendance) error
	GetAttendance(ctx context.Context, params model.MstAttendance) (res model.MstAttendance, err error)
	ListAttendanceByParams(ctx context.Context, params model.ListAttendanceParams) (res []model.MstAttendance, err error)
	UpdateAttendance(ctx context.Context, attendance *model.MstAttendance) (err error)
	CreatePayrollPeriod(ctx context.Context, payrolPeriod *model.MstPayrollPeriod) (err error)
	GetPayrollPeriod(ctx context.Context, id int64) (res model.MstPayrollPeriod, err error)
	UpdatePayrollPeriod(ctx context.Context, payrolPeriod *model.MstPayrollPeriod) (err error)
	SubmitOvertime(ctx context.Context, overtime *model.TrxOvertime) (err error)
	GetOvertime(ctx context.Context, params model.TrxOvertime) (res model.TrxOvertime, err error)
	UpdateOvertime(ctx context.Context, overtime *model.TrxOvertime) (err error)
	ListOvertimeByParams(ctx context.Context, params model.ListOvertimeParams) (res []model.TrxOvertime, err error)

	SubmitReimbursement(ctx context.Context, reimbursement *model.TrxReimbursement) (err error)
	ListReimbursementByParams(ctx context.Context, params model.ListReimbursementParams) (resp []model.TrxReimbursement, err error)
	UpdateReimbursement(ctx context.Context, reimbursement *model.TrxReimbursement) (err error)

	SubmitPayslips(ctx context.Context, payslips []model.TrxUserPayslip) (err error)
	GetPayslips(ctx context.Context, params model.GetPayslipRequest) (payslips []model.TrxUserPayslip, err error)
	SubmitPayroll(ctx context.Context, payroll model.DtlPayroll) (err error)
	GetPayrollDetail(ctx context.Context, params model.GetDtlPayrollRequest) (payrollDetail model.DtlPayroll, err error)
}
