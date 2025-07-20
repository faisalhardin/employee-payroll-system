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
