package server

import (
	userhandler "github.com/faisalhardin/employee-payroll-system/internal/repo/handler/user"
)

type Handlers struct {
	UserHandler *userhandler.UserHandler
}

type Modules struct {
	Handlers *Handlers
}
