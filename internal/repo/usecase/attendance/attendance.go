package attendance

import (
	"context"
	"database/sql"
	"time"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/constant"
	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	attendancerepo "github.com/faisalhardin/employee-payroll-system/internal/entity/repo/db/attendance"
	userdbrepo "github.com/faisalhardin/employee-payroll-system/internal/entity/repo/db/user"
	"github.com/faisalhardin/employee-payroll-system/pkg/common/commonerr"
	"github.com/faisalhardin/employee-payroll-system/pkg/middlewares/auth"
	"github.com/pkg/errors"
)

const (
	MaxOvertimeHours = 3
)

var (
	authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
)

type Usecase struct {
	AttendanceDB attendancerepo.AttendanceRepository
	UserDB       userdbrepo.UserRepository
}

func New(u Usecase) *Usecase {
	return &u
}

func (u *Usecase) TapIn(ctx context.Context, tapInRequest model.MstAttendance) (resp model.TapInResponse, err error) {

	user, found := authGetUserDetailFromCtx(ctx)
	if !found {
		err = errors.Wrap(errors.New("user not found"), "Usecase.TapIn")
		return
	}

	if tapInRequest.AttendanceDate.IsZero() {
		tapInRequest.AttendanceDate = time.Now()
	}
	mstAttendace := &model.MstAttendance{
		IDMstUser:      user.ID,
		AttendanceDate: tapInRequest.AttendanceDate,
		CreatedBy: sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		},
	}

	if isWeekend(mstAttendace.AttendanceDate) {
		err = commonerr.SetNewBadRequest("invalid", "cannot tap in on weekend")
		return
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
	err = u.AttendanceDB.RecordAttendance(ctx, mstAttendace)
	if err != nil {
		err = errors.Wrap(err, "Usecase.TapIn")
		return
	}

	resp = model.TapInResponse{
		AttendanceDate: mstAttendace.AttendanceDate,
	}
	return
}

func isWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

func (u *Usecase) CreatePayrollPeriod(ctx context.Context, payrollPeriodRequest model.PayrollPeriodRequest) (resp model.PayrollPeriodResponse, err error) {

	user, found := authGetUserDetailFromCtx(ctx)
	if !found {
		err = errors.Wrap(errors.New("forbidden"), "Usecase.CreatePayrollPeriod")
		return
	}

	if user.Role != constant.UserRoleAdmin {
		err = errors.Wrap(errors.New("forbidden"), "Usecase.CreatePayrollPeriod")
		return
	}

	mstPayrollPeriod := &model.MstPayrollPeriod{
		StartDate: payrollPeriodRequest.StartDate,
		EndDate:   payrollPeriodRequest.EndDate,
		CreatedBy: sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		},
	}

	err = u.AttendanceDB.CreatePayrollPeriod(ctx, mstPayrollPeriod)
	if err != nil {
		err = errors.Wrap(err, "Usecase.CreatePayrollPeriod")
		return
	}

	resp = model.PayrollPeriodResponse{
		ID:        mstPayrollPeriod.ID,
		StartDate: mstPayrollPeriod.StartDate,
		EndDate:   mstPayrollPeriod.EndDate,
	}
	return
}

func (u *Usecase) SubmitOvertime(ctx context.Context, overtimeRequest model.SubmitOvertimeRequest) (resp model.SubmitOvertimeResponse, err error) {
	user, found := authGetUserDetailFromCtx(ctx)
	if !found {
		err = errors.Wrap(errors.New("forbidden"), "Usecase.CreatePayrollPeriod")
		return
	}

	if !isWeekend(overtimeRequest.OvertimeDate) {

		// check if user attended on the overtime date during business days
		mstAttendance, e := u.AttendanceDB.GetAttendance(ctx, model.MstAttendance{
			IDMstUser:      user.ID,
			AttendanceDate: overtimeRequest.OvertimeDate,
		})
		if e != nil {
			err = errors.Wrap(e, "Usecase.SubmitOvertime")
			return
		}

		if mstAttendance.ID == 0 {
			err = commonerr.SetNewBadRequest("invalid", "user did not attend on the overtime date")
			return
		}
	}

	if overtimeRequest.Hours > MaxOvertimeHours {
		overtimeRequest.Hours = MaxOvertimeHours
	}

	overtime := &model.TrxOvertime{
		UserID:       user.ID,
		OvertimeDate: overtimeRequest.OvertimeDate,
	}

	// check if overtime already submitted for this date
	existingOvertime, err := u.AttendanceDB.GetOvertime(ctx, *overtime)
	if err != nil {
		err = errors.Wrap(err, "Usecase.SubmitOvertime")
		return
	}

	if existingOvertime.ID != 0 {
		err = commonerr.SetNewBadRequest("invalid", "overtime already submitted for this date")
		return
	}

	// complete for insertion
	overtime.Hours = overtimeRequest.Hours
	overtime.CreatedBy = sql.NullInt64{
		Int64: user.ID,
		Valid: true,
	}

	err = u.AttendanceDB.SubmitOvertime(ctx, overtime)
	if err != nil {
		err = errors.Wrap(err, "Usecase.SubmitOvertime")
		return
	}

	resp = model.SubmitOvertimeResponse{
		OvertimeDate: overtimeRequest.OvertimeDate,
		Hours:        overtimeRequest.Hours,
	}

	return
}
