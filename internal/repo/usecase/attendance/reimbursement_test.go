package attendance

import (
	"context"
	"testing"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	"github.com/faisalhardin/employee-payroll-system/pkg/middlewares/auth"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_SubmitReimbursement(t *testing.T) {
	ctrl := initMock(t)
	defer ctrl.Finish()

	type args struct {
		ctx                        context.Context
		submitReimbursementRequest model.SubmitReimbursementRequest
	}
	testCases := []struct {
		name    string
		args    args
		patch   func()
		unpatch func()
		wantErr bool
		want    model.SubmitReimbursementResponse
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				submitReimbursementRequest: model.SubmitReimbursementRequest{
					Amount:      100000,
					Description: "Transportation cost",
				},
			},
			want: model.SubmitReimbursementResponse{
				ID:          1,
				Amount:      100000,
				Description: "Transportation cost",
				Status:      ReimbursementStatusPending,
			},
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						ID: 123,
					}, true
				}

				mockAttendanceRepo.
					EXPECT().SubmitReimbursement(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, reimbursement *model.TrxReimbursement) error {
						reimbursement.ID = 1
						return nil
					}).
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
				submitReimbursementRequest: model.SubmitReimbursementRequest{
					Amount:      50000,
					Description: "Meal allowance",
				},
			},
			want: model.SubmitReimbursementResponse{},
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
			name: "error during submit reimbursement",
			args: args{
				ctx: context.Background(),
				submitReimbursementRequest: model.SubmitReimbursementRequest{
					Amount:      75000,
					Description: "Office supplies",
				},
			},
			want: model.SubmitReimbursementResponse{},
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						ID: 456,
					}, true
				}

				mockAttendanceRepo.
					EXPECT().SubmitReimbursement(gomock.Any(), gomock.Any()).
					Return(errFoo).
					Times(1)
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
			wantErr: true,
		},
		{
			name: "success - verify created by field",
			args: args{
				ctx: context.Background(),
				submitReimbursementRequest: model.SubmitReimbursementRequest{
					Amount:      200000,
					Description: "Business trip expenses",
				},
			},
			want: model.SubmitReimbursementResponse{
				ID:          2,
				Amount:      200000,
				Description: "Business trip expenses",
				Status:      ReimbursementStatusPending,
			},
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						ID: 789,
					}, true
				}

				mockAttendanceRepo.
					EXPECT().SubmitReimbursement(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, reimbursement *model.TrxReimbursement) error {
						reimbursement.ID = 2
						return nil
					}).
					Times(1)
			},
			unpatch: func() {
				authGetUserDetailFromCtx = auth.GetUserDetailFromCtx
			},
			wantErr: false,
		},
		{
			name: "success - empty description",
			args: args{
				ctx: context.Background(),
				submitReimbursementRequest: model.SubmitReimbursementRequest{
					Amount:      25000,
					Description: "",
				},
			},
			want: model.SubmitReimbursementResponse{
				ID:          3,
				Amount:      25000,
				Description: "",
				Status:      ReimbursementStatusPending,
			},
			patch: func() {
				authGetUserDetailFromCtx = func(ctx context.Context) (auth.UserJWTPayload, bool) {
					return auth.UserJWTPayload{
						ID: 111,
					}, true
				}

				mockAttendanceRepo.
					EXPECT().SubmitReimbursement(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, reimbursement *model.TrxReimbursement) error {
						reimbursement.ID = 3
						return nil
					}).
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
			got, err := u.SubmitReimbursement(tc.args.ctx, tc.args.submitReimbursementRequest)

			if assert.Equal(t, tc.wantErr, err != nil) {
				if !tc.wantErr {
					assert.Equal(t, tc.want, got)
				}
			}
		})
	}
}
