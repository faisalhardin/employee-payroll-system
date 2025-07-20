package user

import (
	"net/http"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	userUsecase "github.com/faisalhardin/employee-payroll-system/internal/repo/usecase/user"
	"github.com/faisalhardin/employee-payroll-system/pkg/common/binding"
	commonwriter "github.com/faisalhardin/employee-payroll-system/pkg/common/writer"
)

var (
	bindingBind = binding.Bind
)

type UserHandler struct {
	UserUsecase *userUsecase.Usecase
}

func New(userHandler *UserHandler) *UserHandler {
	return userHandler
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request := model.SignInRequest{}
	err := bindingBind(r, &request)
	if err != nil {
		commonwriter.SetError(ctx, w, err)
		return
	}
	token, err := h.UserUsecase.SignIn(ctx, request)
	if err != nil {
		commonwriter.SetError(ctx, w, err)
		return
	}

	commonwriter.SetOKWithData(ctx, w, token)
}
