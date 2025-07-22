package attendance

import (
	"context"
	"database/sql"
	"time"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/constant"
	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	"github.com/faisalhardin/employee-payroll-system/pkg/common/commonerr"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

const (
	OvertimeMultiplier = int64(2)
	WorkingHours       = int64(8)
)

var (
	usecaseGetNumberOfWorkingDays             = (*Usecase).getNumberOfWorkingDays
	usecaseCalculatePayslipSummaryTotalSalary = (*Usecase).calculatePayslipSummaryTotalSalary
	usecaseReimbursementCalculation           = (*Usecase).reimbursementCalculation
	usecaseOvertimeCalculation                = (*Usecase).overtimeCalculation
	usecaseAttendanceCalculation              = (*Usecase).attendanceCalculation
	usecaseSubmitPayslips                     = (*Usecase).submitPayslips
	usecaseUpdateReimbursementInBulk          = (*Usecase).updateReimbursementInBulk
	usecaseUpdateOvertimeInBulk               = (*Usecase).updateOvertimeInBulk
	usecaseUpdateAttendanceInBulk             = (*Usecase).updateAttendanceInBulk
	usecaseUpdatePayrollPeriod                = (*Usecase).updatePayrollPeriod
	usecaseGetMapOfPayslipSummary             = (*Usecase).getMapOfPayslipSummary
)

func (u *Usecase) GeneratePayroll(ctx context.Context, request model.GeneratePayrollRequest) (err error) {
	user, found := authGetUserDetailFromCtx(ctx)
	if !found {
		err = errors.Wrap(errors.New("user not found"), "Usecase.GeneratePayroll")
		return
	}

	payrollPeriod, err := u.AttendanceDB.GetPayrollPeriod(ctx, int64(request.IDMstPayrollPeriod))
	if err != nil {
		return errors.Wrap(err, "Usecase.GeneratePayroll")
	}

	if payrollPeriod.ID == 0 {
		return commonerr.SetNewBadRequest("invalid", "payroll period not found")
	}

	if !payrollPeriod.PayrollProcessedDate.Time.IsZero() {
		return commonerr.SetNewBadRequest("invalid", "payroll has been processed")
	}

	numberOfWorkingDays := u.getNumberOfWorkingDays(payrollPeriod.StartDate, payrollPeriod.EndDate)

	employees, err := u.UserDB.ListUser(ctx)
	if err != nil {
		return errors.Wrap(err, "Usecase.GeneratePayroll")
	}

	payslipSummary := u.getMapOfPayslipSummary(employees, numberOfWorkingDays, payrollPeriod.ID)

	// START attendance calculation
	payslipSummary, listOfAttendance, err := u.attendanceCalculation(
		ctx,
		payrollPeriod.StartDate, payrollPeriod.EndDate,
		payslipSummary,
		payrollPeriod.ID,
		user.ID,
	)
	if err != nil {
		return errors.Wrap(err, "Usecase.GeneratePayroll")
	}
	// END attendance calculation

	// START overtime calculation
	payslipSummary, listOfOvertime, err := u.overtimeCalculation(
		ctx,
		payrollPeriod.StartDate, payrollPeriod.EndDate,
		payslipSummary,
		payrollPeriod.ID,
		user.ID,
	)
	if err != nil {
		return errors.Wrap(err, "Usecase.GeneratePayroll")
	}
	// END overtime calculation

	//  START reimbursement calculation
	payslipSummary, listOfReimbursement, err := u.reimbursementCalculation(
		ctx,
		payrollPeriod.StartDate, payrollPeriod.EndDate,
		payslipSummary,
		payrollPeriod.ID,
		user.ID,
	)
	if err != nil {
		return errors.Wrap(err, "Usecase.GeneratePayroll")
	}
	// END reimbursement calculation

	payslipSummary, totalTakeHomePay := u.calculatePayslipSummaryTotalSalary(payslipSummary, numberOfWorkingDays)
	payrollSummary := model.DtlPayroll{
		IDMstPayrollPeriod: payrollPeriod.ID,
		CreatedBy:          user.ID,
		TotalTakeHome:      totalTakeHomePay,
	}

	// store payroll summary
	err = u.AttendanceDB.SubmitPayroll(ctx, payrollSummary)
	if err != nil {
		return errors.Wrap(err, "Usecase.GeneratePayroll")
	}

	// store payslip summary
	err = u.submitPayslips(ctx, payslipSummary)
	if err != nil {
		return errors.Wrap(err, "Usecase.GeneratePayroll")
	}

	// update attendance
	err = u.updateAttendanceInBulk(ctx, listOfAttendance)
	if err != nil {
		return errors.Wrap(err, "Usecase.GeneratePayroll")
	}

	// update reimbursement
	err = u.updateReimbursementInBulk(ctx, listOfReimbursement)
	if err != nil {
		return errors.Wrap(err, "Usecase.GeneratePayroll")
	}

	// update overtime
	err = u.updateOvertimeInBulk(ctx, listOfOvertime)
	if err != nil {
		return errors.Wrap(err, "Usecase.GeneratePayroll")
	}

	// update payroll period
	err = u.updatePayrollPeriod(ctx, &payrollPeriod, user.ID)
	if err != nil {
		return errors.Wrap(err, "Usecase.GeneratePayroll")
	}

	return
}

func (u *Usecase) getNumberOfWorkingDays(startDate, endDate time.Time) int {
	workingDays := 0
	attendanceDate := startDate
	for attendanceDate.Before(endDate) || attendanceDate.Equal(endDate) {
		if attendanceDate.Weekday() != time.Saturday && attendanceDate.Weekday() != time.Sunday {
			workingDays++
		}
		attendanceDate = attendanceDate.AddDate(0, 0, 1)
	}
	return workingDays
}

func (u *Usecase) getMapOfPayslipSummary(employees []model.MstUser, workingDays int, payrollPeriodID int64) map[int64]model.TrxUserPayslip {
	payslipSummary := map[int64]model.TrxUserPayslip{}

	for _, employee := range employees {
		payslipSummary[employee.ID] = model.TrxUserPayslip{
			UserID:             employee.ID,
			Username:           employee.Username,
			BaseSalary:         employee.Salary,
			WorkingDays:        workingDays,
			IDMstPayrollPeriod: payrollPeriodID,
		}
	}

	return payslipSummary
}

func (u *Usecase) updatePayrollPeriod(ctx context.Context, payrollPeriod *model.MstPayrollPeriod, userID int64) error {
	payrollPeriod.PayrollProcessedDate = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	payrollPeriod.UpdatedBy = sql.NullInt64{
		Int64: userID,
		Valid: true,
	}
	err := u.AttendanceDB.UpdatePayrollPeriod(ctx, payrollPeriod)
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) updateAttendanceInBulk(ctx context.Context, attendances []model.MstAttendance) (err error) {
	for _, attendance := range attendances {
		err = u.AttendanceDB.UpdateAttendance(ctx, &attendance)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Usecase) updateOvertimeInBulk(ctx context.Context, overtimes []model.TrxOvertime) (err error) {
	for _, overtime := range overtimes {
		err = u.AttendanceDB.UpdateOvertime(ctx, &overtime)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Usecase) updateReimbursementInBulk(ctx context.Context, reimbursements []model.TrxReimbursement) (err error) {
	for _, reimbursement := range reimbursements {
		err = u.AttendanceDB.UpdateReimbursement(ctx, &reimbursement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Usecase) submitPayslips(ctx context.Context, mapOfPayslips map[int64]model.TrxUserPayslip) (err error) {
	payslips := []model.TrxUserPayslip{}

	for _, payslip := range mapOfPayslips {
		payslips = append(payslips, payslip)
	}
	err = u.AttendanceDB.SubmitPayslips(ctx, payslips)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) attendanceCalculation(
	ctx context.Context,
	startDate, endDate time.Time,
	payslipSummary map[int64]model.TrxUserPayslip,
	payrollPeriodID int64,
	userID int64,
) (map[int64]model.TrxUserPayslip, []model.MstAttendance, error) {

	listAttendanceParams := model.ListAttendanceParams{
		StartDate:              startDate,
		EndDate:                endDate,
		IsForGeneratingPayroll: true,
	}

	listOfAttendance, err := u.AttendanceDB.ListAttendanceByParams(ctx, listAttendanceParams)
	if err != nil {
		return payslipSummary, []model.MstAttendance{}, err
	}

	for i, attendance := range listOfAttendance {
		userAttendance := payslipSummary[attendance.IDMstUser]
		userAttendance.AttendedDays += 1
		payslipSummary[attendance.IDMstUser] = userAttendance

		attendance.IDMstPayrollPeriod = sql.NullInt64{
			Int64: payrollPeriodID,
			Valid: true,
		}
		attendance.UpdatedBy = sql.NullInt64{
			Int64: userID,
			Valid: true,
		}
		listOfAttendance[i] = attendance
	}

	return payslipSummary, listOfAttendance, nil
}

func (u *Usecase) overtimeCalculation(
	ctx context.Context,
	startDate time.Time, endDate time.Time,
	payslipSummary map[int64]model.TrxUserPayslip,
	payrollPeriodID int64,
	userID int64,
) (map[int64]model.TrxUserPayslip, []model.TrxOvertime, error) {
	listOvertimeParams := model.ListOvertimeParams{
		StartDate: startDate,
		EndDate:   endDate,
	}

	listOfOvertime, err := u.AttendanceDB.ListOvertimeByParams(ctx, listOvertimeParams)
	if err != nil {
		return payslipSummary, []model.TrxOvertime{}, err
	}

	for i, overtime := range listOfOvertime {
		userOvertime := payslipSummary[overtime.UserID]
		userOvertime.OvertimeHours += overtime.Hours
		payslipSummary[overtime.UserID] = userOvertime

		overtime.IDMstPayrollPeriod = sql.NullInt64{
			Int64: payrollPeriodID,
			Valid: true,
		}
		overtime.UpdatedBy = sql.NullInt64{
			Int64: userID,
			Valid: true,
		}
		listOfOvertime[i] = overtime
	}

	return payslipSummary, listOfOvertime, nil
}

func (u *Usecase) reimbursementCalculation(ctx context.Context,
	startDate time.Time, endDate time.Time,
	payslipSummary map[int64]model.TrxUserPayslip,
	payrollPeriodID int64,
	userID int64,
) (map[int64]model.TrxUserPayslip, []model.TrxReimbursement, error) {
	listReimbursementParams := model.ListReimbursementParams{
		StartDate: startDate,
		EndDate:   endDate,
		Status:    ReimbursementStatusPending,
	}

	listOfReimbursement, err := u.AttendanceDB.ListReimbursementByParams(ctx, listReimbursementParams)
	if err != nil {
		return payslipSummary, []model.TrxReimbursement{}, err
	}

	for i, reimbursement := range listOfReimbursement {
		userReimbursement := payslipSummary[reimbursement.UserID]
		userReimbursement.TotalReimbursements += reimbursement.Amount
		payslipSummary[reimbursement.UserID] = userReimbursement

		reimbursement.IDMstPayrollPeriod = sql.NullInt64{
			Int64: payrollPeriodID,
			Valid: true,
		}
		reimbursement.UpdatedBy = sql.NullInt64{
			Int64: userID,
			Valid: true,
		}
		reimbursement.Status = ReimbursementStatusPaid
		listOfReimbursement[i] = reimbursement
	}

	return payslipSummary, listOfReimbursement, err
}

func (*Usecase) calculatePayslipSummaryTotalSalary(
	payslipSummary map[int64]model.TrxUserPayslip,
	numberOfWorkingDays int,
) (
	modifiedPayslipSummary map[int64]model.TrxUserPayslip,
	totalTakeHomePay int64,
) {
	totalTakeHomePay = 0
	for i, employeeSummary := range payslipSummary {
		// attendance calculation
		proratedSalary := decimal.NewFromInt(int64(employeeSummary.AttendedDays)).
			Div(decimal.NewFromInt(int64(numberOfWorkingDays))).
			Mul(decimal.NewFromInt(int64(employeeSummary.BaseSalary)))
		employeeSummary.ProratedSalary = proratedSalary.IntPart()

		// overtime calculation
		hourlyPay := decimal.NewFromInt(int64(employeeSummary.BaseSalary)).
			Div(decimal.NewFromInt(int64(numberOfWorkingDays))).
			Div(decimal.NewFromInt(WorkingHours))

		overtimePay := decimal.NewFromInt(int64(employeeSummary.OvertimeHours)).
			Mul(hourlyPay).
			Mul(decimal.NewFromInt(OvertimeMultiplier))

		employeeSummary.OvertimePay = overtimePay.IntPart()

		reimbursement := decimal.NewFromInt(int64(employeeSummary.TotalReimbursements))

		totalTakeHome := proratedSalary.Add(overtimePay).Add(reimbursement)
		employeeSummary.TotalTakeHome = totalTakeHome.IntPart()

		totalTakeHomePay += employeeSummary.TotalTakeHome
		payslipSummary[i] = employeeSummary
	}

	return payslipSummary, totalTakeHomePay
}

func (u *Usecase) GetEmployeePayslip(ctx context.Context, request model.GetPayslipRequest) (payslip model.GetPayslipResponse, err error) {
	user, found := authGetUserDetailFromCtx(ctx)
	if !found {
		err = errors.Wrap(errors.New("user not found"), "Usecase.GetEmployeePayslip")
		return
	}
	request.UserID = user.ID

	payrollPeriod, err := u.AttendanceDB.GetPayrollPeriod(ctx, int64(request.IDMstPayrollPeriod))
	if err != nil {
		return
	}

	payslips, err := u.AttendanceDB.GetPayslips(ctx, request)
	if err != nil {
		return
	}
	if len(payslips) == 0 {
		err = errors.Wrap(commonerr.SetNewBadRequest("not found", "payslip not found"), "Usecase.GetEmployeePayslip")
		return
	}
	employeePayslip := payslips[0]

	listOfAttendance, err := u.AttendanceDB.ListAttendanceByParams(ctx, model.ListAttendanceParams{
		IDsMstUser:         []int64{user.ID},
		IDMstPayrollPeriod: int64(request.IDMstPayrollPeriod),
	})
	if err != nil {
		return
	}

	attendanceDate := []string{}
	for _, attendance := range listOfAttendance {
		attendanceDate = append(attendanceDate, attendance.AttendanceDate.Format("2006-01-02"))
	}

	listOfOvertime, err := u.AttendanceDB.ListOvertimeByParams(ctx, model.ListOvertimeParams{
		UserIDs:            []int64{user.ID},
		IDMstPayrollPeriod: int64(request.IDMstPayrollPeriod),
	})
	if err != nil {
		return
	}
	listOfOvertimeResponse := []model.GetOvertimeResponse{}
	for _, overtime := range listOfOvertime {
		listOfOvertimeResponse = append(listOfOvertimeResponse, model.GetOvertimeResponse{
			OvertimeDate: overtime.OvertimeDate,
			Hours:        overtime.Hours,
		})
	}

	listOfReimbursement, err := u.AttendanceDB.ListReimbursementByParams(ctx, model.ListReimbursementParams{
		UserID:             user.ID,
		IDMstPayrollPeriod: int64(request.IDMstPayrollPeriod),
	})
	if err != nil {
		return
	}
	reimbursementListResponse := []model.SubmitReimbursementResponse{}
	for _, reimbursement := range listOfReimbursement {
		reimbursementListResponse = append(reimbursementListResponse, model.SubmitReimbursementResponse{
			Description: reimbursement.Description,
			Amount:      reimbursement.Amount,
			Status:      reimbursement.Status,
		})
	}

	payslip = model.GetPayslipResponse{
		StartDate:           payrollPeriod.StartDate,
		EndDate:             payrollPeriod.EndDate,
		TotalTakeHomePay:    employeePayslip.TotalTakeHome,
		AttendanceDate:      attendanceDate,
		WorkingDays:         employeePayslip.WorkingDays,
		AttendedDays:        employeePayslip.AttendedDays,
		ProratedSalary:      employeePayslip.ProratedSalary,
		OvertimeHours:       employeePayslip.OvertimeHours,
		OvertimePay:         employeePayslip.OvertimePay,
		TotalReimbursements: employeePayslip.TotalReimbursements,
		OvertimeDetails:     listOfOvertimeResponse,
		ReimbursementList:   reimbursementListResponse,
	}

	return
}

func (u *Usecase) GetPayroll(ctx context.Context, request model.GetPayrollRequest) (payrollSummary model.GetPayrollResponse, err error) {
	user, found := authGetUserDetailFromCtx(ctx)
	if !found {
		err = errors.Wrap(errors.New("user not found"), "Usecase.GetEmployeePayslip")
		return
	}
	if user.Role != constant.UserRoleAdmin {
		err = errors.Wrap(errors.New("unauthorized"), "Usecase.GetEmployeePayslip")
		return
	}

	payrollPeriod, err := u.AttendanceDB.GetPayrollPeriod(ctx, request.IDMstPayrollPeriod)
	if err != nil {
		return
	}
	if payrollPeriod.CreatedAt.IsZero() {
		err = errors.Wrap(commonerr.SetNewBadRequest("not found", "payroll period not found"), "Usecase.GetEmployeePayslip")
		return
	}

	payslips, err := u.AttendanceDB.GetPayslips(ctx, model.GetPayslipRequest{
		IDMstPayrollPeriod: request.IDMstPayrollPeriod,
	})
	if err != nil {
		return
	}
	if len(payslips) == 0 {
		err = errors.Wrap(commonerr.SetNewBadRequest("not found", "payslip not found"), "Usecase.GetEmployeePayslip")
		return
	}

	payrollDetail, err := u.AttendanceDB.GetPayrollDetail(ctx, model.GetDtlPayrollRequest{
		IDMstPayrollPeriod: request.IDMstPayrollPeriod,
	})
	if err != nil {
		return
	}

	payrollSummary = model.GetPayrollResponse{
		StartDate:        payrollPeriod.StartDate,
		EndDate:          payrollPeriod.EndDate,
		EmployeesPayslip: payslips,
		TotalTakeHomePay: payrollDetail.TotalTakeHome,
	}

	return
}
