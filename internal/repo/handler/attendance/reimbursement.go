package attendance

import (
	"net/http"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	commonwriter "github.com/faisalhardin/employee-payroll-system/pkg/common/writer"
)

func (h *AttendanceHandler) SubmitReimbursement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := model.SubmitReimbursementRequest{}
	err := bindingBind(r, &req)
	if err != nil {
		commonwriter.SetError(ctx, w, err)
		return
	}

	resp, err := h.AttendanceUsecase.SubmitReimbursement(ctx, req)
	if err != nil {
		commonwriter.SetError(ctx, w, err)
		return
	}

	commonwriter.SetOKWithData(ctx, w, resp)
}
