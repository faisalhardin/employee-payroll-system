package user

import (
	"context"
	"time"

	"github.com/faisalhardin/employee-payroll-system/internal/config"
	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	"github.com/faisalhardin/employee-payroll-system/internal/repo/db/user"
	"github.com/faisalhardin/employee-payroll-system/pkg/middlewares/auth"
	"github.com/pkg/errors"
)

type Usecase struct {
	Cfg      *config.Config
	UserDB   *user.Conn
	AuthRepo *auth.Options
}

func New(opt *Usecase) *Usecase {
	return opt
}

func (u *Usecase) SignIn(ctx context.Context, params model.SignInRequest) (jwt string, err error) {
	resp, err := u.UserDB.GetUser(ctx, params)
	if err != nil {
		err = errors.Wrap(err, "Usecase.SignIn")
		return
	}
	if resp.ID == 0 {
		err = errors.New("user not found")
		return
	}

	currTime := time.Now()
	expireDuration := time.Duration(u.Cfg.JWTConfig.DurationInHours) * time.Hour
	expiredTime := currTime.Add(expireDuration)
	token, err := u.AuthRepo.CreateJWTToken(ctx, auth.UserJWTPayload{
		ID:       resp.ID,
		Username: resp.Username,
		Role:     resp.Role,
	}, currTime, expiredTime)
	if err != nil {
		return
	}
	return token, nil
}
