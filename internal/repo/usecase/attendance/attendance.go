package attendance

import (
	"context"
	"time"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	attendaceDB "github.com/faisalhardin/employee-payroll-system/internal/repo/db/attendance"
	"github.com/faisalhardin/employee-payroll-system/pkg/middlewares/auth"
	"github.com/pkg/errors"
)

type Usecase struct {
	AttendanceDB *attendaceDB.Conn
}

func New(u Usecase) *Usecase {
	return &u
}

func (u *Usecase) TapIn(ctx context.Context) (resp model.TapInResponse, err error) {

	user, found := auth.GetUserDetailFromCtx(ctx)
	if !found {
		err = errors.Wrap(errors.New("user not found"), "Usecase.TapIn")
		return
	}

	mstAttendace := &model.MstAttendance{
		IDMstUser:      user.ID,
		AttendanceDate: time.Now(),
	}

	existingAttendance, err := u.AttendanceDB.GetAttendance(ctx, *mstAttendace)
	if err != nil {
		err = errors.Wrap(err, "Usecase.TapIn")
		return
	}

	if !existingAttendance.CreatedAt.IsZero() {
		resp = model.TapInResponse{
			AttendanceDate: existingAttendance.CreatedAt,
		}
		return
	}
	err = u.AttendanceDB.Create(ctx, mstAttendace)
	if err != nil {
		err = errors.Wrap(err, "Usecase.TapIn")
		return
	}

	resp = model.TapInResponse{
		AttendanceDate: mstAttendace.AttendanceDate,
	}
	return
}
