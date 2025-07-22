package user

import (
	"context"
	"testing"

	"github.com/faisalhardin/employee-payroll-system/internal/config"
	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	mockrepo "github.com/faisalhardin/employee-payroll-system/internal/entity/repo/_mocks"
	mockuserdb "github.com/faisalhardin/employee-payroll-system/internal/entity/repo/db/_mocks/user"
	"github.com/faisalhardin/employee-payroll-system/pkg/middlewares/auth"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	mockUserDB   *mockuserdb.MockUserRepository
	mockAuthRepo *mockrepo.MockAuthenticator

	errFoo = errors.New("errFoo")
)

func initMock(t *testing.T) *gomock.Controller {
	ctrl := gomock.NewController(t)
	mockUserDB = mockuserdb.NewMockUserRepository(ctrl)
	mockAuthRepo = mockrepo.NewMockAuthenticator(ctrl)
	return ctrl
}

func Test_SignIn(t *testing.T) {
	ctrl := initMock(t)
	defer ctrl.Finish()

	type args struct {
		ctx    context.Context
		params model.SignInRequest
	}
	testCases := []struct {
		name    string
		args    args
		patch   func()
		wantErr bool
		want    string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				params: model.SignInRequest{
					Username: "fooname",
					Password: "foopass",
				},
			},
			want: "token",
			patch: func() {
				mockUserDB.
					EXPECT().GetUser(gomock.Any(), gomock.Any()).
					Return(model.MstUser{
						ID:       1,
						Username: "fooname",
					}, nil).
					Times(1)
				mockAuthRepo.
					EXPECT().
					CreateJWTToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("token", nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "error during jwt creation",
			args: args{
				ctx: context.Background(),
				params: model.SignInRequest{
					Username: "fooname",
					Password: "foopass",
				},
			},
			want: "",
			patch: func() {
				mockUserDB.
					EXPECT().GetUser(gomock.Any(), gomock.Any()).
					Return(model.MstUser{
						ID:       1,
						Username: "fooname",
					}, nil).
					Times(1)
				mockAuthRepo.
					EXPECT().
					CreateJWTToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("token", errFoo).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "error user not found",
			args: args{
				ctx: context.Background(),
				params: model.SignInRequest{
					Username: "fooname",
					Password: "foopass",
				},
			},
			want: "",
			patch: func() {
				mockUserDB.
					EXPECT().GetUser(gomock.Any(), gomock.Any()).
					Return(model.MstUser{}, nil).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "error during user retrieval",
			args: args{
				ctx: context.Background(),
				params: model.SignInRequest{
					Username: "fooname",
					Password: "foopass",
				},
			},
			want: "",
			patch: func() {
				mockUserDB.
					EXPECT().GetUser(gomock.Any(), gomock.Any()).
					Return(model.MstUser{}, errFoo).
					Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			u := Usecase{
				Cfg: &config.Config{
					JWTConfig: auth.JWTConfig{
						DurationInHours: 1,
					},
				},
				UserDB:   mockUserDB,
				AuthRepo: mockAuthRepo,
			}
			tc.patch()
			got, err := u.SignIn(tc.args.ctx, tc.args.params)

			if assert.Equal(t, tc.wantErr, err != nil) {
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestNew(t *testing.T) {
	mockUsecase := &Usecase{
		UserDB:   mockUserDB,
		AuthRepo: mockAuthRepo,
	}
	tests := []struct {
		name     string
		input    *Usecase
		expected *Usecase
	}{
		{
			name:     "returns the same usecase instance",
			input:    mockUsecase,
			expected: mockUsecase,
		},
		{
			name:     "handles nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "returns empty usecase",
			input:    &Usecase{},
			expected: &Usecase{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := New(tt.input)
			assert.Equal(t, tt.input, result)
			if tt.input != nil {
				assert.Same(t, tt.input, result)
			}
		})
	}
}
