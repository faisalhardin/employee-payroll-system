package attendance

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/constant"
	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	"github.com/faisalhardin/employee-payroll-system/pkg/middlewares/auth"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_calculatePayslipSummaryTotalSalary(t *testing.T) {
	type args struct {
		payslipSummary      map[int64]model.TrxUserPayslip
		numberOfWorkingDays int
	}
	testCases := []struct {
		name                       string
		args                       args
		patch                      func()
		unpatch                    func()
		wantModifiedPayslipSummary map[int64]model.TrxUserPayslip
		wantTotalTakeHomePay       int64
	}{
		{
			name: "success - single employee with full attendance",
			args: args{
				payslipSummary: map[int64]model.TrxUserPayslip{
					1: {
						UserID:              1,
						BaseSalary:          1000000,
						AttendedDays:        20,
						OvertimeHours:       4,
						TotalReimbursements: 50000,
					},
				},
				numberOfWorkingDays: 20,
			},
			wantModifiedPayslipSummary: map[int64]model.TrxUserPayslip{
				1: {
					UserID:              1,
					BaseSalary:          1000000,
					AttendedDays:        20,
					OvertimeHours:       4,
					TotalReimbursements: 50000,
					ProratedSalary:      1000000, // (20/20) * 1000000
					OvertimePay:         50000,   // (1000000/20/8) * 4 * 2
					TotalTakeHome:       1100000, // 1000000 + 50000 + 50000
				},
			},
			wantTotalTakeHomePay: 1100000,
			patch:                func() {},
			unpatch:              func() {},
		},
		{
			name: "success - single employee with partial attendance",
			args: args{
				payslipSummary: map[int64]model.TrxUserPayslip{
					1: {
						UserID:              1,
						BaseSalary:          2000000,
						AttendedDays:        15,
						OvertimeHours:       2,
						TotalReimbursements: 100000,
					},
				},
				numberOfWorkingDays: 20,
			},
			wantModifiedPayslipSummary: map[int64]model.TrxUserPayslip{
				1: {
					UserID:              1,
					BaseSalary:          2000000,
					AttendedDays:        15,
					OvertimeHours:       2,
					TotalReimbursements: 100000,
					ProratedSalary:      1500000, // (15/20) * 2000000
					OvertimePay:         50000,   // (2000000/20/8) * 2 * 2
					TotalTakeHome:       1650000, // 1500000 + 62500 + 100000
				},
			},
			wantTotalTakeHomePay: 1650000,
			patch:                func() {},
			unpatch:              func() {},
		},
		{
			name: "success - multiple employees",
			args: args{
				payslipSummary: map[int64]model.TrxUserPayslip{
					1: {
						UserID:              1,
						BaseSalary:          1000000,
						AttendedDays:        20,
						OvertimeHours:       4,
						TotalReimbursements: 50000,
					},
					2: {
						UserID:              2,
						BaseSalary:          1500000,
						AttendedDays:        18,
						OvertimeHours:       2,
						TotalReimbursements: 75000,
					},
				},
				numberOfWorkingDays: 20,
			},
			wantModifiedPayslipSummary: map[int64]model.TrxUserPayslip{
				1: {
					UserID:              1,
					BaseSalary:          1000000,
					AttendedDays:        20,
					OvertimeHours:       4,
					TotalReimbursements: 50000,
					ProratedSalary:      1000000, // (20/20) * 1000000
					OvertimePay:         50000,   // (1000000/20/8) * 4 * 2
					TotalTakeHome:       1100000, // 1000000 + 50000 + 50000
				},
				2: {
					UserID:              2,
					BaseSalary:          1500000,
					AttendedDays:        18,
					OvertimeHours:       2,
					TotalReimbursements: 75000,
					ProratedSalary:      1350000, // (18/20) * 1500000
					OvertimePay:         37500,   // (1500000/20/8) * 2 * 2
					TotalTakeHome:       1462500, // 1350000 + 37500 + 75000
				},
			},
			wantTotalTakeHomePay: 2562500, // 1100000 + 1462500
			patch:                func() {},
			unpatch:              func() {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := &Usecase{}
			tc.patch()
			defer tc.unpatch()

			gotModifiedPayslipSummary, gotTotalTakeHomePay := u.calculatePayslipSummaryTotalSalary(tc.args.payslipSummary, tc.args.numberOfWorkingDays)

			assert.Equal(t, tc.wantTotalTakeHomePay, gotTotalTakeHomePay)
			assert.Equal(t, tc.wantModifiedPayslipSummary, gotModifiedPayslipSummary)
		})
	}
}

func Test_calculatePayslipSummaryTotalSalary_DecimalPrecision(t *testing.T) {
	// Test to verify decimal calculations are handled correctly
	u := &Usecase{}

	payslipSummary := map[int64]model.TrxUserPayslip{
		1: {
			UserID:              1,
			BaseSalary:          1000001, // Odd number to test precision
			AttendedDays:        13,      // Partial attendance
			OvertimeHours:       3,       // Odd overtime hours
			TotalReimbursements: 33333,   // Odd reimbursement
		},
	}
	numberOfWorkingDays := 21 // Odd working days

	gotModifiedPayslipSummary, gotTotalTakeHomePay := u.calculatePayslipSummaryTotalSalary(payslipSummary, numberOfWorkingDays)

	employee := gotModifiedPayslipSummary[1]

	// Verify calculations manually
	expectedProratedSalary := decimal.NewFromInt(13).
		Div(decimal.NewFromInt(21)).
		Mul(decimal.NewFromInt(1000001)).IntPart()

	expectedHourlyPay := decimal.NewFromInt(1000001).
		Div(decimal.NewFromInt(21)).
		Div(decimal.NewFromInt(8))

	expectedOvertimePay := decimal.NewFromInt(3).
		Mul(expectedHourlyPay).
		Mul(decimal.NewFromInt(2)).IntPart()

	expectedTotalTakeHome := expectedProratedSalary + expectedOvertimePay + 33333

	assert.Equal(t, expectedProratedSalary, employee.ProratedSalary)
	assert.Equal(t, expectedOvertimePay, employee.OvertimePay)
	assert.Equal(t, expectedTotalTakeHome, employee.TotalTakeHome)
	assert.Equal(t, expectedTotalTakeHome, gotTotalTakeHomePay)
}

func Test_reimbursementCalculation(t *testing.T) {
	ctrl := initMock(t)
	defer ctrl.Finish()

	type args struct {
		ctx             context.Context
		startDate       time.Time
		endDate         time.Time
		payslipSummary  map[int64]model.TrxUserPayslip
		payrollPeriodID int64
		userID          int64
	}
	testCases := []struct {
		name                    string
		args                    args
		patch                   func()
		unpatch                 func()
		wantErr                 bool
		wantPayslipSummary      map[int64]model.TrxUserPayslip
		wantListOfReimbursement []model.TrxReimbursement
	}{
		{
			name: "success - multiple users with reimbursements",
			args: args{
				ctx:       context.Background(),
				startDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				endDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
				payslipSummary: map[int64]model.TrxUserPayslip{
					123: {
						UserID:              123,
						TotalReimbursements: 0,
					},
					456: {
						UserID:              456,
						TotalReimbursements: 0,
					},
				},
				payrollPeriodID: 1,
				userID:          789,
			},
			wantPayslipSummary: map[int64]model.TrxUserPayslip{
				123: {
					UserID:              123,
					TotalReimbursements: 50000,
				},
				456: {
					UserID:              456,
					TotalReimbursements: 85000, // 10000 + 75000
				},
			},
			wantListOfReimbursement: []model.TrxReimbursement{
				{
					ID:     1,
					UserID: 123,
					Amount: 50000,
					Status: ReimbursementStatusPaid,
					IDMstPayrollPeriod: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
					UpdatedBy: sql.NullInt64{
						Int64: 789,
						Valid: true,
					},
				},
				{
					ID:     2,
					UserID: 456,
					Amount: 75000,
					Status: ReimbursementStatusPaid,
					IDMstPayrollPeriod: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
					UpdatedBy: sql.NullInt64{
						Int64: 789,
						Valid: true,
					},
				},
				{
					ID:     3,
					UserID: 456,
					Amount: 10000,
					Status: ReimbursementStatusPaid,
					IDMstPayrollPeriod: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
					UpdatedBy: sql.NullInt64{
						Int64: 789,
						Valid: true,
					},
				},
			},
			patch: func() {
				mockAttendanceRepo.
					EXPECT().ListReimbursementByParams(gomock.Any(), gomock.Any()).
					Return([]model.TrxReimbursement{
						{
							ID:     1,
							UserID: 123,
							Amount: 50000,
							Status: ReimbursementStatusPending,
						},
						{
							ID:     2,
							UserID: 456,
							Amount: 75000,
							Status: ReimbursementStatusPending,
						},
						{
							ID:     3,
							UserID: 456,
							Amount: 10000,
							Status: ReimbursementStatusPending,
						},
					}, nil).
					Times(1)
			},
			unpatch: func() {},
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := Usecase{
				AttendanceDB: mockAttendanceRepo,
			}
			tc.patch()
			defer tc.unpatch()

			gotPayslipSummary, gotListOfReimbursement, err := u.reimbursementCalculation(
				tc.args.ctx,
				tc.args.startDate,
				tc.args.endDate,
				tc.args.payslipSummary,
				tc.args.payrollPeriodID,
				tc.args.userID,
			)

			if assert.Equal(t, tc.wantErr, err != nil) {
				if !tc.wantErr {
					assert.Equal(t, tc.wantPayslipSummary, gotPayslipSummary)
					assert.Equal(t, tc.wantListOfReimbursement, gotListOfReimbursement)
				}
			}
		})
	}
}

func Test_overtimeCalculation(t *testing.T) {
	ctrl := initMock(t)
	defer ctrl.Finish()

	type args struct {
		ctx             context.Context
		startDate       time.Time
		endDate         time.Time
		payslipSummary  map[int64]model.TrxUserPayslip
		payrollPeriodID int64
		userID          int64
	}
	testCases := []struct {
		name               string
		args               args
		patch              func()
		unpatch            func()
		wantErr            bool
		wantPayslipSummary map[int64]model.TrxUserPayslip
		wantListOfOvertime []model.TrxOvertime
	}{
		{
			name: "success - multiple users with overtime",
			args: args{
				ctx:       context.Background(),
				startDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				endDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
				payslipSummary: map[int64]model.TrxUserPayslip{
					123: {
						UserID:        123,
						OvertimeHours: 0,
					},
					456: {
						UserID:        456,
						OvertimeHours: 1,
					},
				},
				payrollPeriodID: 1,
				userID:          789,
			},
			wantPayslipSummary: map[int64]model.TrxUserPayslip{
				123: {
					UserID:        123,
					OvertimeHours: 5,
				},
				456: {
					UserID:        456,
					OvertimeHours: 4, // 1 + 3
				},
			},
			wantListOfOvertime: []model.TrxOvertime{
				{
					ID:     1,
					UserID: 123,
					Hours:  5,
					IDMstPayrollPeriod: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
					UpdatedBy: sql.NullInt64{
						Int64: 789,
						Valid: true,
					},
				},
				{
					ID:     2,
					UserID: 456,
					Hours:  3,
					IDMstPayrollPeriod: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
					UpdatedBy: sql.NullInt64{
						Int64: 789,
						Valid: true,
					},
				},
			},
			patch: func() {
				mockAttendanceRepo.
					EXPECT().ListOvertimeByParams(gomock.Any(), gomock.Any()).
					Return([]model.TrxOvertime{
						{
							ID:     1,
							UserID: 123,
							Hours:  5,
						},
						{
							ID:     2,
							UserID: 456,
							Hours:  3,
						},
					}, nil).
					Times(1)
			},
			unpatch: func() {},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := Usecase{
				AttendanceDB: mockAttendanceRepo,
			}
			tc.patch()
			defer tc.unpatch()

			gotPayslipSummary, gotListOfOvertime, err := u.overtimeCalculation(
				tc.args.ctx,
				tc.args.startDate,
				tc.args.endDate,
				tc.args.payslipSummary,
				tc.args.payrollPeriodID,
				tc.args.userID,
			)

			if assert.Equal(t, tc.wantErr, err != nil) {
				if !tc.wantErr {
					assert.Equal(t, tc.wantPayslipSummary, gotPayslipSummary)
					assert.Equal(t, tc.wantListOfOvertime, gotListOfOvertime)
				}
			}
		})
	}
}

func Test_attendanceCalculation(t *testing.T) {
	ctrl := initMock(t)
	defer ctrl.Finish()

	type args struct {
		ctx             context.Context
		startDate       time.Time
		endDate         time.Time
		payslipSummary  map[int64]model.TrxUserPayslip
		payrollPeriodID int64
		userID          int64
	}
	testCases := []struct {
		name                 string
		args                 args
		patch                func()
		unpatch              func()
		wantErr              bool
		wantPayslipSummary   map[int64]model.TrxUserPayslip
		wantListOfAttendance []model.MstAttendance
	}{
		{
			name: "success - multiple users with attendance",
			args: args{
				ctx:       context.Background(),
				startDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				endDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
				payslipSummary: map[int64]model.TrxUserPayslip{
					123: {
						UserID:       123,
						AttendedDays: 0,
					},
					456: {
						UserID:       456,
						AttendedDays: 2,
					},
				},
				payrollPeriodID: 1,
				userID:          789,
			},
			wantPayslipSummary: map[int64]model.TrxUserPayslip{
				123: {
					UserID:       123,
					AttendedDays: 2, // 0 + 2 attendances
				},
				456: {
					UserID:       456,
					AttendedDays: 3, // 2 + 1 attendance
				},
			},
			wantListOfAttendance: []model.MstAttendance{
				{
					ID:        1,
					IDMstUser: 123,
					IDMstPayrollPeriod: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
					UpdatedBy: sql.NullInt64{
						Int64: 789,
						Valid: true,
					},
				},
				{
					ID:        2,
					IDMstUser: 123,
					IDMstPayrollPeriod: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
					UpdatedBy: sql.NullInt64{
						Int64: 789,
						Valid: true,
					},
				},
				{
					ID:        3,
					IDMstUser: 456,
					IDMstPayrollPeriod: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
					UpdatedBy: sql.NullInt64{
						Int64: 789,
						Valid: true,
					},
				},
			},
			patch: func() {
				mockAttendanceRepo.
					EXPECT().ListAttendanceByParams(gomock.Any(), gomock.Any()).
					Return([]model.MstAttendance{
						{
							ID:        1,
							IDMstUser: 123,
						},
						{
							ID:        2,
							IDMstUser: 123,
						},
						{
							ID:        3,
							IDMstUser: 456,
						},
					}, nil).
					Times(1)
			},
			unpatch: func() {},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := Usecase{
				AttendanceDB: mockAttendanceRepo,
			}
			tc.patch()
			defer tc.unpatch()

			gotPayslipSummary, gotListOfAttendance, err := u.attendanceCalculation(
				tc.args.ctx,
				tc.args.startDate,
				tc.args.endDate,
				tc.args.payslipSummary,
				tc.args.payrollPeriodID,
				tc.args.userID,
			)

			if assert.Equal(t, tc.wantErr, err != nil) {
				if !tc.wantErr {
					assert.Equal(t, tc.wantPayslipSummary, gotPayslipSummary)
					assert.Equal(t, tc.wantListOfAttendance, gotListOfAttendance)
				}
			}
		})
	}
}

func Test_submitPayslips(t *testing.T) {
	ctrl := initMock(t)
	defer ctrl.Finish()

	type args struct {
		ctx           context.Context
		mapOfPayslips map[int64]model.TrxUserPayslip
	}
	testCases := []struct {
		name    string
		args    args
		patch   func()
		unpatch func()
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				mapOfPayslips: map[int64]model.TrxUserPayslip{
					123: {
						UserID:              123,
						Username:            "john.doe",
						BaseSalary:          1000000,
						WorkingDays:         20,
						AttendedDays:        20,
						ProratedSalary:      1000000,
						OvertimeHours:       2,
						OvertimePay:         25000,
						TotalReimbursements: 0,
						TotalTakeHome:       1025000,
						IDMstPayrollPeriod:  1,
					},
					456: {
						UserID:              456,
						Username:            "jane.smith",
						BaseSalary:          1500000,
						WorkingDays:         20,
						AttendedDays:        19,
						ProratedSalary:      1425000,
						OvertimeHours:       3,
						OvertimePay:         56250,
						TotalReimbursements: 50000,
						TotalTakeHome:       1531250,
						IDMstPayrollPeriod:  1,
					},
				},
			},
			patch: func() {
				mockAttendanceRepo.
					EXPECT().SubmitPayslips(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			unpatch: func() {},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := Usecase{
				AttendanceDB: mockAttendanceRepo,
			}
			tc.patch()
			defer tc.unpatch()

			err := u.submitPayslips(tc.args.ctx, tc.args.mapOfPayslips)

			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func Test_updateReimbursementInBulk(t *testing.T) {
	ctrl := initMock(t)
	defer ctrl.Finish()

	type args struct {
		ctx            context.Context
		reimbursements []model.TrxReimbursement
	}
	testCases := []struct {
		name    string
		args    args
		patch   func()
		unpatch func()
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				reimbursements: []model.TrxReimbursement{
					{
						ID:          1,
						UserID:      123,
						Amount:      30000,
						Description: "Meal allowance",
						Status:      ReimbursementStatusPaid,
						IDMstPayrollPeriod: sql.NullInt64{
							Int64: 1,
							Valid: true,
						},
						UpdatedBy: sql.NullInt64{
							Int64: 456,
							Valid: true,
						},
					},
				},
			},
			patch: func() {
				mockAttendanceRepo.
					EXPECT().UpdateReimbursement(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			unpatch: func() {},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := Usecase{
				AttendanceDB: mockAttendanceRepo,
			}
			tc.patch()
			defer tc.unpatch()

			err := u.updateReimbursementInBulk(tc.args.ctx, tc.args.reimbursements)

			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func Test_updateOvertimeInBulk(t *testing.T) {
	ctrl := initMock(t)
	defer ctrl.Finish()

	type args struct {
		ctx       context.Context
		overtimes []model.TrxOvertime
	}
	testCases := []struct {
		name    string
		args    args
		patch   func()
		unpatch func()
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				overtimes: []model.TrxOvertime{
					{
						ID:           1,
						UserID:       123,
						OvertimeDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
						Hours:        3,
						IDMstPayrollPeriod: sql.NullInt64{
							Int64: 1,
							Valid: true,
						},
						UpdatedBy: sql.NullInt64{
							Int64: 456,
							Valid: true,
						},
					},
				},
			},
			patch: func() {
				mockAttendanceRepo.
					EXPECT().UpdateOvertime(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			unpatch: func() {},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := Usecase{
				AttendanceDB: mockAttendanceRepo,
			}
			tc.patch()
			defer tc.unpatch()

			err := u.updateOvertimeInBulk(tc.args.ctx, tc.args.overtimes)

			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func Test_updateAttendanceInBulk(t *testing.T) {
	ctrl := initMock(t)
	defer ctrl.Finish()

	type args struct {
		ctx         context.Context
		attendances []model.MstAttendance
	}
	testCases := []struct {
		name    string
		args    args
		patch   func()
		unpatch func()
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				attendances: []model.MstAttendance{
					{
						ID:             1,
						IDMstUser:      123,
						AttendanceDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
						CreatedAt:      time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
						IDMstPayrollPeriod: sql.NullInt64{
							Int64: 1,
							Valid: true,
						},
						UpdatedBy: sql.NullInt64{
							Int64: 456,
							Valid: true,
						},
					},
				},
			},
			patch: func() {
				mockAttendanceRepo.
					EXPECT().UpdateAttendance(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			unpatch: func() {},
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := Usecase{
				AttendanceDB: mockAttendanceRepo,
			}
			tc.patch()
			defer tc.unpatch()

			err := u.updateAttendanceInBulk(tc.args.ctx, tc.args.attendances)

			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func Test_updatePayrollPeriod(t *testing.T) {
	ctrl := initMock(t)
	defer ctrl.Finish()

	type args struct {
		ctx           context.Context
		payrollPeriod *model.MstPayrollPeriod
		userID        int64
	}
	testCases := []struct {
		name    string
		args    args
		patch   func()
		unpatch func()
		wantErr bool
	}{
		{
			name: "success - update payroll period",
			args: args{
				ctx: context.Background(),
				payrollPeriod: &model.MstPayrollPeriod{
					ID:        1,
					StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
					CreatedBy: sql.NullInt64{
						Int64: 100,
						Valid: true,
					},
					PayrollProcessedDate: sql.NullTime{
						Time:  time.Time{},
						Valid: false,
					},
					UpdatedBy: sql.NullInt64{
						Int64: 0,
						Valid: false,
					},
				},
				userID: 456,
			},
			patch: func() {
				mockAttendanceRepo.
					EXPECT().UpdatePayrollPeriod(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			unpatch: func() {},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := Usecase{
				AttendanceDB: mockAttendanceRepo,
			}
			tc.patch()
			defer tc.unpatch()

			err := u.updatePayrollPeriod(tc.args.ctx, tc.args.payrollPeriod, tc.args.userID)

			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func Test_getMapOfPayslipSummary(t *testing.T) {
	type args struct {
		employees       []model.MstUser
		workingDays     int
		payrollPeriodID int64
	}
	testCases := []struct {
		name    string
		args    args
		patch   func()
		unpatch func()
		want    map[int64]model.TrxUserPayslip
	}{
		{
			name: "success - multiple employees",
			args: args{
				employees: []model.MstUser{
					{
						ID:       123,
						Username: "john.doe",
						Salary:   1000000,
					},
					{
						ID:       456,
						Username: "jane.smith",
						Salary:   1500000,
					},
					{
						ID:       789,
						Username: "bob.wilson",
						Salary:   800000,
					},
				},
				workingDays:     22,
				payrollPeriodID: 5,
			},
			want: map[int64]model.TrxUserPayslip{
				123: {
					UserID:             123,
					Username:           "john.doe",
					BaseSalary:         1000000,
					WorkingDays:        22,
					IDMstPayrollPeriod: 5,
				},
				456: {
					UserID:             456,
					Username:           "jane.smith",
					BaseSalary:         1500000,
					WorkingDays:        22,
					IDMstPayrollPeriod: 5,
				},
				789: {
					UserID:             789,
					Username:           "bob.wilson",
					BaseSalary:         800000,
					WorkingDays:        22,
					IDMstPayrollPeriod: 5,
				},
			},
			patch:   func() {},
			unpatch: func() {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := &Usecase{}
			tc.patch()
			defer tc.unpatch()

			got := u.getMapOfPayslipSummary(tc.args.employees, tc.args.workingDays, tc.args.payrollPeriodID)

			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_getNumberOfWorkingDays(t *testing.T) {
	type args struct {
		startDate time.Time
		endDate   time.Time
	}
	testCases := []struct {
		name    string
		args    args
		patch   func()
		unpatch func()
		want    int
	}{
		{
			name: "success - full month with weekends",
			args: args{
				startDate: time.Date(2025, 7, 13, 0, 0, 0, 0, time.UTC), // Monday
				endDate:   time.Date(2025, 7, 19, 0, 0, 0, 0, time.UTC), // Sunday
			},
			want:    5, // 7 days - 2 weekend days (Saturdays + Sundays)
			patch:   func() {},
			unpatch: func() {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := &Usecase{}
			tc.patch()
			defer tc.unpatch()

			got := u.getNumberOfWorkingDays(tc.args.startDate, tc.args.endDate)

			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_GetEmployeePayslip(t *testing.T) {
	ctrl := initMock(t)
	defer ctrl.Finish()

	type args struct {
		ctx     context.Context
		request model.GetPayslipRequest
	}
	testCases := []struct {
		name    string
		args    args
		patch   func()
		unpatch func()
		wantErr bool
		want    model.GetPayslipResponse
	}{
		{
			name: "success - complete payslip data",
			args: args{
				ctx: context.Background(),
				request: model.GetPayslipRequest{
					IDMstPayrollPeriod: 1,
				},
			},
			want: model.GetPayslipResponse{
				StartDate:           time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:             time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
				TotalTakeHomePay:    1175000,
				AttendanceDate:      []string{"2024-01-15", "2024-01-16", "2024-01-17"},
				WorkingDays:         22,
				AttendedDays:        20,
				ProratedSalary:      909090,
				OvertimeHours:       4,
				OvertimePay:         50000,
				TotalReimbursements: 75000,
				OvertimeDetails: []model.GetOvertimeResponse{
					{
						OvertimeDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
						Hours:        2,
					},
					{
						OvertimeDate: time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
						Hours:        2,
					},
				},
				ReimbursementList: []model.SubmitReimbursementResponse{
					{
						Description: "Transportation",
						Amount:      50000,
						Status:      ReimbursementStatusPaid,
					},
					{
						Description: "Meal allowance",
						Amount:      25000,
						Status:      ReimbursementStatusPaid,
					},
				},
			},
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						ID: 123,
					}, true
				}

				mockAttendanceRepo.
					EXPECT().GetPayrollPeriod(gomock.Any(), int64(1)).
					Return(model.MstPayrollPeriod{
						ID:        1,
						StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
					}, nil).
					Times(1)

				mockAttendanceRepo.
					EXPECT().GetPayslips(gomock.Any(), gomock.Any()).
					Return([]model.TrxUserPayslip{
						{
							UserID:              123,
							WorkingDays:         22,
							AttendedDays:        20,
							ProratedSalary:      909090,
							OvertimeHours:       4,
							OvertimePay:         50000,
							TotalReimbursements: 75000,
							TotalTakeHome:       1175000,
						},
					}, nil).
					Times(1)

				mockAttendanceRepo.
					EXPECT().ListAttendanceByParams(gomock.Any(), gomock.Any()).
					Return([]model.MstAttendance{
						{AttendanceDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
						{AttendanceDate: time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC)},
						{AttendanceDate: time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC)},
					}, nil).
					Times(1)

				mockAttendanceRepo.
					EXPECT().ListOvertimeByParams(gomock.Any(), gomock.Any()).
					Return([]model.TrxOvertime{
						{
							OvertimeDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
							Hours:        2,
						},
						{
							OvertimeDate: time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
							Hours:        2,
						},
					}, nil).
					Times(1)

				mockAttendanceRepo.
					EXPECT().ListReimbursementByParams(gomock.Any(), gomock.Any()).
					Return([]model.TrxReimbursement{
						{
							Description: "Transportation",
							Amount:      50000,
							Status:      ReimbursementStatusPaid,
						},
						{
							Description: "Meal allowance",
							Amount:      25000,
							Status:      ReimbursementStatusPaid,
						},
					}, nil).
					Times(1)
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := Usecase{
				AttendanceDB: mockAttendanceRepo,
			}
			tc.patch()
			defer tc.unpatch()

			got, err := u.GetEmployeePayslip(tc.args.ctx, tc.args.request)

			if assert.Equal(t, tc.wantErr, err != nil) {
				if !tc.wantErr {
					assert.Equal(t, tc.want, got)
				}
			}
		})
	}
}

func Test_GetPayroll(t *testing.T) {
	ctrl := initMock(t)
	defer ctrl.Finish()

	type args struct {
		ctx     context.Context
		request model.GetPayrollRequest
	}
	testCases := []struct {
		name    string
		args    args
		patch   func()
		unpatch func()
		wantErr bool
		want    model.GetPayrollResponse
	}{
		{
			name: "success - complete payroll data",
			args: args{
				ctx: context.Background(),
				request: model.GetPayrollRequest{
					IDMstPayrollPeriod: 1,
				},
			},
			want: model.GetPayrollResponse{
				StartDate:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
				TotalTakeHomePay: 5500000,
				EmployeesPayslip: []model.TrxUserPayslip{
					{
						UserID:              123,
						Username:            "john.doe",
						BaseSalary:          2000000,
						WorkingDays:         22,
						AttendedDays:        20,
						ProratedSalary:      1818181,
						OvertimeHours:       4,
						OvertimePay:         100000,
						TotalReimbursements: 50000,
						TotalTakeHome:       1968181,
						IDMstPayrollPeriod:  1,
					},
					{
						UserID:              456,
						Username:            "jane.smith",
						BaseSalary:          2500000,
						WorkingDays:         22,
						AttendedDays:        22,
						ProratedSalary:      2500000,
						OvertimeHours:       6,
						OvertimePay:         187500,
						TotalReimbursements: 75000,
						TotalTakeHome:       2762500,
						IDMstPayrollPeriod:  1,
					},
					{
						UserID:              789,
						Username:            "bob.wilson",
						BaseSalary:          1500000,
						WorkingDays:         22,
						AttendedDays:        18,
						ProratedSalary:      1227272,
						OvertimeHours:       2,
						OvertimePay:         37500,
						TotalReimbursements: 25000,
						TotalTakeHome:       1289772,
						IDMstPayrollPeriod:  1,
					},
				},
			},
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						ID:   999,
						Role: constant.UserRoleAdmin,
					}, true
				}

				mockAttendanceRepo.
					EXPECT().GetPayrollPeriod(gomock.Any(), int64(1)).
					Return(model.MstPayrollPeriod{
						ID:        1,
						StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
						CreatedAt: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
					}, nil).
					Times(1)

				mockAttendanceRepo.
					EXPECT().GetPayslips(gomock.Any(), gomock.Any()).
					Return([]model.TrxUserPayslip{
						{
							UserID:              123,
							Username:            "john.doe",
							BaseSalary:          2000000,
							WorkingDays:         22,
							AttendedDays:        20,
							ProratedSalary:      1818181,
							OvertimeHours:       4,
							OvertimePay:         100000,
							TotalReimbursements: 50000,
							TotalTakeHome:       1968181,
							IDMstPayrollPeriod:  1,
						},
						{
							UserID:              456,
							Username:            "jane.smith",
							BaseSalary:          2500000,
							WorkingDays:         22,
							AttendedDays:        22,
							ProratedSalary:      2500000,
							OvertimeHours:       6,
							OvertimePay:         187500,
							TotalReimbursements: 75000,
							TotalTakeHome:       2762500,
							IDMstPayrollPeriod:  1,
						},
						{
							UserID:              789,
							Username:            "bob.wilson",
							BaseSalary:          1500000,
							WorkingDays:         22,
							AttendedDays:        18,
							ProratedSalary:      1227272,
							OvertimeHours:       2,
							OvertimePay:         37500,
							TotalReimbursements: 25000,
							TotalTakeHome:       1289772,
							IDMstPayrollPeriod:  1,
						},
					}, nil).
					Times(1)

				mockAttendanceRepo.
					EXPECT().GetPayrollDetail(gomock.Any(), gomock.Any()).
					Return(model.DtlPayroll{
						IDMstPayrollPeriod: 1,
						TotalTakeHome:      5500000,
						CreatedBy:          999,
					}, nil).
					Times(1)
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := Usecase{
				AttendanceDB: mockAttendanceRepo,
			}
			tc.patch()
			defer tc.unpatch()

			got, err := u.GetPayroll(tc.args.ctx, tc.args.request)

			if assert.Equal(t, tc.wantErr, err != nil) {
				if !tc.wantErr {
					assert.Equal(t, tc.want, got)
				}
			}
		})
	}
}
