package attendance

import (
	"context"
	"database/sql"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	"github.com/faisalhardin/employee-payroll-system/pkg/middlewares/auth"
	"github.com/pkg/errors"
)

const (
	ReimbursementStatusPending = "pending"
	ReimbursementStatusPaid    = "paid"
)

func (u *Usecase) SubmitReimbursement(ctx context.Context, submitReimbursementRequest model.SubmitReimbursementRequest) (resp model.SubmitReimbursementResponse, err error) {

	user, found := auth.GetUserDetailFromCtx(ctx)
	if !found {
		err = errors.Wrap(errors.New("user not found"), "Usecase.TapIn")
		return
	}

	trxReimbursement := &model.TrxReimbursement{
		UserID:      user.ID,
		Amount:      submitReimbursementRequest.Amount,
		Description: submitReimbursementRequest.Description,
		Status:      ReimbursementStatusPending,
		CreatedBy: sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		},
	}

	err = u.AttendanceDB.SubmitReimbursement(ctx, trxReimbursement)
	if err != nil {
		err = errors.Wrap(err, "Usecase.SubmitReimbursement")
		return
	}

	resp = model.SubmitReimbursementResponse{
		ID:          trxReimbursement.ID,
		Amount:      trxReimbursement.Amount,
		Description: trxReimbursement.Description,
		Status:      trxReimbursement.Status,
	}

	return resp, nil
}
