package attendance

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/constant"
	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	mockattendancedb "github.com/faisalhardin/employee-payroll-system/internal/entity/repo/db/_mocks/attendance"
	mockuserdb "github.com/faisalhardin/employee-payroll-system/internal/entity/repo/db/_mocks/user"
	"github.com/faisalhardin/employee-payroll-system/pkg/middlewares/auth"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	mockAttendanceRepo *mockattendancedb.MockAttendanceRepository
	mockUserRepo       *mockuserdb.MockUserRepository

	errFoo = errors.New("errFoo")
)

func initMock(t *testing.T) *gomock.Controller {
	ctrl := gomock.NewController(t)
	mockAttendanceRepo = mockattendancedb.NewMockAttendanceRepository(ctrl)
	mockUserRepo = mockuserdb.NewMockUserRepository(ctrl)

	return ctrl
}

func Test_TapIn(t *testing.T) {
	ctrl := initMock(t)
	defer ctrl.Finish()

	mockNow, _ := time.Parse("2006-01-02", "2025-07-21")    //monday
	mockSunday, _ := time.Parse("2006-01-02", "2025-07-20") //sunday
	type args struct {
		ctx context.Context
		req model.MstAttendance
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    model.TapInResponse
		patch   func()
		unpatch func()
	}{
		{
			name: "Successful first tap in",
			args: args{
				ctx: context.Background(),
				req: model.MstAttendance{
					AttendanceDate: mockNow,
				},
			},
			want: model.TapInResponse{
				AttendanceDate: mockNow,
			},
			wantErr: false,
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						Username: "username",
					}, true
				}

				mockAttendanceRepo.EXPECT().
					GetAttendance(gomock.Any(), gomock.Any()).
					Return(model.MstAttendance{}, nil).Times(1)

				mockAttendanceRepo.EXPECT().
					RecordAttendance(gomock.Any(), gomock.Any()).
					Return(nil).Times(1)
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
		},
		{
			name: "Successful second tap in",
			args: args{
				ctx: context.Background(),
				req: model.MstAttendance{
					AttendanceDate: mockSunday,
				},
			},
			want: model.TapInResponse{
				AttendanceDate: time.Time{},
			},
			wantErr: true,
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						Username: "username",
					}, true
				}
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
		},
		{
			name: "Tap in during weekend",
			args: args{
				ctx: context.Background(),
				req: model.MstAttendance{
					AttendanceDate: mockNow,
				},
			},
			want: model.TapInResponse{
				AttendanceDate: mockNow,
			},
			wantErr: false,
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						Username: "username",
					}, true
				}

				mockAttendanceRepo.EXPECT().
					GetAttendance(gomock.Any(), gomock.Any()).
					Return(model.MstAttendance{}, nil).Times(1)

				mockAttendanceRepo.EXPECT().
					RecordAttendance(gomock.Any(), gomock.Any()).
					Return(nil).Times(1)
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
		},
		{
			name: "error during submission",
			args: args{
				ctx: context.Background(),
				req: model.MstAttendance{
					AttendanceDate: mockNow,
				},
			},
			want:    model.TapInResponse{},
			wantErr: true,
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						Username: "username",
					}, true
				}

				mockAttendanceRepo.EXPECT().
					GetAttendance(gomock.Any(), gomock.Any()).
					Return(model.MstAttendance{}, nil).Times(1)

				mockAttendanceRepo.EXPECT().
					RecordAttendance(gomock.Any(), gomock.Any()).
					Return(errFoo).Times(1)
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
		},
		{
			name: "error when GetAttendance",
			args: args{
				ctx: context.Background(),
				req: model.MstAttendance{
					AttendanceDate: mockNow,
				},
			},
			want:    model.TapInResponse{},
			wantErr: true,
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						Username: "username",
					}, true
				}

				mockAttendanceRepo.EXPECT().
					GetAttendance(gomock.Any(), gomock.Any()).
					Return(model.MstAttendance{}, errFoo).Times(1)
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Usecase{
				AttendanceDB: mockAttendanceRepo,
			}

			tt.patch()
			defer tt.unpatch()
			got, err := u.TapIn(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("TapIn() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_CreatePayrollPeriod(t *testing.T) {
	ctrl := initMock(t)
	defer ctrl.Finish()

	type args struct {
		ctx                  context.Context
		payrollPeriodRequest model.PayrollPeriodRequest
	}
	testCases := []struct {
		name    string
		args    args
		patch   func()
		unpatch func()
		wantErr bool
		want    model.PayrollPeriodResponse
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				payrollPeriodRequest: model.PayrollPeriodRequest{
					StartDate: time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC),
				},
			},
			want: model.PayrollPeriodResponse{
				ID:        1,
				StartDate: time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC),
			},
			patch: func() {
				mockAttendanceRepo.
					EXPECT().CreatePayrollPeriod(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, period *model.MstPayrollPeriod) error {
						period.ID = 1 // Simulate DB setting ID
						return nil
					}).Times(1)

				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						ID:   1,
						Role: constant.UserRoleAdmin,
					}, true
				}

			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
			wantErr: false,
		},
		{
			name: "error - user not found in context",
			args: args{
				ctx: context.Background(),
				payrollPeriodRequest: model.PayrollPeriodRequest{
					StartDate: time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC),
				},
			},
			want: model.PayrollPeriodResponse{},
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{}, false
				}
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
			wantErr: true,
		},
		{
			name: "error - user is not admin",
			args: args{
				ctx: context.Background(),
				payrollPeriodRequest: model.PayrollPeriodRequest{
					StartDate: time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC),
				},
			},
			want: model.PayrollPeriodResponse{},
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						ID:   1,
						Role: constant.UserRoleEmployee,
					}, true
				}
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
			wantErr: true,
		},
		{
			name: "error during create payroll period",
			args: args{
				ctx: context.Background(),
				payrollPeriodRequest: model.PayrollPeriodRequest{
					StartDate: time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC),
				},
			},
			want: model.PayrollPeriodResponse{},
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						ID:   1,
						Role: constant.UserRoleAdmin,
					}, true
				}

				mockAttendanceRepo.EXPECT().
					CreatePayrollPeriod(gomock.Any(), gomock.Any()).
					Return(errFoo)
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := Usecase{
				AttendanceDB: mockAttendanceRepo,
			}
			tc.patch()
			defer tc.unpatch()
			got, err := u.CreatePayrollPeriod(tc.args.ctx, tc.args.payrollPeriodRequest)
			if assert.Equal(t, tc.wantErr, err != nil) {
				if !tc.wantErr {
					assert.Equal(t, tc.want, got)
				}
			}
		})
	}
}

func Test_SubmitOvertime(t *testing.T) {
	ctrl := initMock(t)
	defer ctrl.Finish()

	type args struct {
		ctx             context.Context
		overtimeRequest model.SubmitOvertimeRequest
	}
	testCases := []struct {
		name    string
		args    args
		patch   func()
		unpatch func()
		wantErr bool
		want    model.SubmitOvertimeResponse
	}{
		{
			name: "success - weekday with attendance",
			args: args{
				ctx: context.Background(),
				overtimeRequest: model.SubmitOvertimeRequest{
					OvertimeDate: time.Date(2025, 7, 21, 0, 0, 0, 0, time.UTC), // Monday
					Hours:        2,
				},
			},
			want: model.SubmitOvertimeResponse{
				OvertimeDate: time.Date(2025, 7, 21, 0, 0, 0, 0, time.UTC), // Monday
				Hours:        2,
			},
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						ID: 123,
					}, true
				}

				mockAttendanceRepo.
					EXPECT().GetAttendance(gomock.Any(), gomock.Any()).
					Return(model.MstAttendance{ID: 1}, nil).
					Times(1)

				mockAttendanceRepo.
					EXPECT().GetOvertime(gomock.Any(), gomock.Any()).
					Return(model.TrxOvertime{}, nil).
					Times(1)

				mockAttendanceRepo.
					EXPECT().SubmitOvertime(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
			wantErr: false,
		},
		{
			name: "success - weekend (no attendance check needed)",
			args: args{
				ctx: context.Background(),
				overtimeRequest: model.SubmitOvertimeRequest{
					OvertimeDate: time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC), // Sunday
					Hours:        3,
				},
			},
			want: model.SubmitOvertimeResponse{
				OvertimeDate: time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC), // Monday
				Hours:        3,
			},
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						ID: 123,
					}, true
				}

				mockAttendanceRepo.
					EXPECT().GetOvertime(gomock.Any(), gomock.Any()).
					Return(model.TrxOvertime{}, nil).
					Times(1)

				mockAttendanceRepo.
					EXPECT().SubmitOvertime(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
			wantErr: false,
		},
		{
			name: "error - user not found in context",
			args: args{
				ctx: context.Background(),
				overtimeRequest: model.SubmitOvertimeRequest{
					OvertimeDate: time.Date(2025, 7, 21, 0, 0, 0, 0, time.UTC), // Monday
					Hours:        2,
				},
			},
			want: model.SubmitOvertimeResponse{},
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{}, false
				}
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
			wantErr: true,
		},
		{
			name: "error - weekday with no attendance",
			args: args{
				ctx: context.Background(),
				overtimeRequest: model.SubmitOvertimeRequest{
					OvertimeDate: time.Date(2025, 7, 21, 0, 0, 0, 0, time.UTC), // Monday
					Hours:        2,
				},
			},
			want: model.SubmitOvertimeResponse{},
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						ID: 123,
					}, true
				}

				mockAttendanceRepo.
					EXPECT().GetAttendance(gomock.Any(), gomock.Any()).
					Return(model.MstAttendance{ID: 0}, nil).
					Times(1)
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
			wantErr: true,
		},
		{
			name: "error - existing overtime for date",
			args: args{
				ctx: context.Background(),
				overtimeRequest: model.SubmitOvertimeRequest{
					OvertimeDate: time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC), // Sunday
					Hours:        2,
				},
			},
			want: model.SubmitOvertimeResponse{},
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						ID: 123,
					}, true
				}

				mockAttendanceRepo.
					EXPECT().GetOvertime(gomock.Any(), gomock.Any()).
					Return(model.TrxOvertime{ID: 1}, nil).
					Times(1)
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
			wantErr: true,
		},
		{
			name: "error during get attendance",
			args: args{
				ctx: context.Background(),
				overtimeRequest: model.SubmitOvertimeRequest{
					OvertimeDate: time.Date(2025, 7, 21, 0, 0, 0, 0, time.UTC), // Monday
					Hours:        2,
				},
			},
			want: model.SubmitOvertimeResponse{},
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						ID: 123,
					}, true
				}

				mockAttendanceRepo.
					EXPECT().GetAttendance(gomock.Any(), gomock.Any()).
					Return(model.MstAttendance{}, errFoo).
					Times(1)
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
			wantErr: true,
		},
		{
			name: "error during get overtime",
			args: args{
				ctx: context.Background(),
				overtimeRequest: model.SubmitOvertimeRequest{
					OvertimeDate: time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC), // Sunday
					Hours:        2,
				},
			},
			want: model.SubmitOvertimeResponse{},
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						ID: 123,
					}, true
				}

				mockAttendanceRepo.
					EXPECT().GetOvertime(gomock.Any(), gomock.Any()).
					Return(model.TrxOvertime{}, errFoo).
					Times(1)
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
			wantErr: true,
		},
		{
			name: "error during submit overtime",
			args: args{
				ctx: context.Background(),
				overtimeRequest: model.SubmitOvertimeRequest{
					OvertimeDate: time.Date(2024, 1, 13, 0, 0, 0, 0, time.UTC), // Saturday
					Hours:        2,
				},
			},
			want: model.SubmitOvertimeResponse{},
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						ID: 123,
					}, true
				}

				mockAttendanceRepo.
					EXPECT().GetOvertime(gomock.Any(), gomock.Any()).
					Return(model.TrxOvertime{}, nil).
					Times(1)

				mockAttendanceRepo.
					EXPECT().SubmitOvertime(gomock.Any(), gomock.Any()).
					Return(errFoo).
					Times(1)
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
			wantErr: true,
		},
		{
			name: "success - hours exceeding max are capped",
			args: args{
				ctx: context.Background(),
				overtimeRequest: model.SubmitOvertimeRequest{
					OvertimeDate: time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC), // Sunday
					Hours:        MaxOvertimeHours + 2,
				},
			},
			want: model.SubmitOvertimeResponse{
				OvertimeDate: time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC), // Sunday
				Hours:        MaxOvertimeHours,
			},
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						ID: 123,
					}, true
				}

				mockAttendanceRepo.
					EXPECT().GetOvertime(gomock.Any(), gomock.Any()).
					Return(model.TrxOvertime{}, nil).
					Times(1)

				mockAttendanceRepo.
					EXPECT().SubmitOvertime(gomock.Any(), &model.TrxOvertime{
					UserID:       123,
					OvertimeDate: time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC), // Sunday
					Hours:        MaxOvertimeHours,
					CreatedBy: sql.NullInt64{
						Int64: 123,
						Valid: true,
					},
				}).
					Return(nil).
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
			got, err := u.SubmitOvertime(tc.args.ctx, tc.args.overtimeRequest)

			if assert.Equal(t, tc.wantErr, err != nil) {
				if !tc.wantErr {
					assert.Equal(t, tc.want, got)
				}
			}
		})
	}
}
