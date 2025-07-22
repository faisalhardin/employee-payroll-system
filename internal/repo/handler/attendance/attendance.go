package attendance

import (
	"net/http"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	"github.com/faisalhardin/employee-payroll-system/internal/repo/usecase/attendance"
	"github.com/faisalhardin/employee-payroll-system/pkg/common/binding"
	commonwriter "github.com/faisalhardin/employee-payroll-system/pkg/common/writer"
)

var (
	bindingBind = binding.Bind
)

type AttendanceHandler struct {
	AttendanceUsecase *attendance.Usecase
}

func New(attendanceHandler *AttendanceHandler) *AttendanceHandler {
	return attendanceHandler
}

func (h *AttendanceHandler) TapIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := h.AttendanceUsecase.TapIn(ctx, model.MstAttendance{})
	if err != nil {
		commonwriter.SetError(ctx, w, err)
		return
	}

	commonwriter.SetOKWithData(ctx, w, resp)
}

func (h *AttendanceHandler) CreatePayrollPeriod(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := model.PayrollPeriodRequest{}
	err := bindingBind(r, &req)
	if err != nil {
		commonwriter.SetError(ctx, w, err)
		return
	}

	resp, err := h.AttendanceUsecase.CreatePayrollPeriod(ctx, req)
	if err != nil {
		commonwriter.SetError(ctx, w, err)
		return
	}

	commonwriter.SetOKWithData(ctx, w, resp)
}

func (h *AttendanceHandler) SubmitOvertime(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := model.SubmitOvertimeRequest{}
	err := bindingBind(r, &req)
	if err != nil {
		commonwriter.SetError(ctx, w, err)
		return
	}

	resp, err := h.AttendanceUsecase.SubmitOvertime(ctx, req)
	if err != nil {
		commonwriter.SetError(ctx, w, err)
		return
	}

	commonwriter.SetOKWithData(ctx, w, resp)
}

func (h *AttendanceHandler) GeneratePayroll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := model.GeneratePayrollRequest{}
	err := bindingBind(r, &req)
	if err != nil {
		commonwriter.SetError(ctx, w, err)
		return
	}

	err = h.AttendanceUsecase.GeneratePayroll(ctx, req)
	if err != nil {
		commonwriter.SetError(ctx, w, err)
		return
	}

	commonwriter.SetOKWithData(ctx, w, "OK")
}

func (h *AttendanceHandler) GetEmployeePayslip(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := model.GetPayslipRequest{}
	err := bindingBind(r, &req)
	if err != nil {
		commonwriter.SetError(ctx, w, err)
		return
	}

	resp, err := h.AttendanceUsecase.GetEmployeePayslip(ctx, req)
	if err != nil {
		commonwriter.SetError(ctx, w, err)
		return
	}

	commonwriter.SetOKWithData(ctx, w, resp)
}
