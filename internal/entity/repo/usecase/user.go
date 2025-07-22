package usecase

import (
	"context"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -self_package=github.com/faisalhardin/employee-payroll-system/internal/entity/repo/usecase -destination=../_mocks/mock_user_usecase.go -package=mock github.com/faisalhardin/employee-payroll-system/internal/entity/repo/usecase UserUsecaseRepository
type UserUsecaseRepository interface {
	SignIn(ctx context.Context, params model.SignInRequest) (jwt string, err error)
}
