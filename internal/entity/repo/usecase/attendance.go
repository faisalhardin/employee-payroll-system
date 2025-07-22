package usecase

import (
	"context"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -self_package=github.com/faisalhardin/employee-payroll-system/internal/entity/repo/usecase -destination=../_mocks/mock_attendance_usecase.go -package=mock github.com/faisalhardin/employee-payroll-system/internal/entity/repo/usecase AttendanceUsecaseRepository
type AttendanceUsecaseRepository interface {
	TapIn(ctx context.Context, tapInRequest model.MstAttendance) (resp model.TapInResponse, err error)
	CreatePayrollPeriod(ctx context.Context, payrollPeriodRequest model.PayrollPeriodRequest) (resp model.PayrollPeriodResponse, err error)
	SubmitOvertime(ctx context.Context, overtimeRequest model.SubmitOvertimeRequest) (resp model.SubmitOvertimeResponse, err error)

	GeneratePayroll(ctx context.Context, request model.GeneratePayrollRequest) (err error)
	GetPayroll(ctx context.Context, request model.GetPayrollRequest) (payrollSummary model.GetPayrollResponse, err error)
	GetEmployeePayslip(ctx context.Context, request model.GetPayslipRequest) (payslip model.GetPayslipResponse, err error)

	SubmitReimbursement(ctx context.Context, submitReimbursementRequest model.SubmitReimbursementRequest) (resp model.SubmitReimbursementResponse, err error)
}
