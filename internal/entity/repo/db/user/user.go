package user

import (
	"context"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -self_package=github.com/faisalhardin/employee-payroll-system/internal/entity/repo/db/user -destination=../_mocks/user/mock_user.go -package=user github.com/faisalhardin/employee-payroll-system/internal/entity/repo/db/user UserRepository
type UserRepository interface {
	GetUser(ctx context.Context, params model.SignInRequest) (res model.MstUser, err error)
	ListUser(ctx context.Context) (res []model.MstUser, err error)
}
