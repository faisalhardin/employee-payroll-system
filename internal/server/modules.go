package server

import (
	attendancehandler "github.com/faisalhardin/employee-payroll-system/internal/repo/handler/attendance"
	userhandler "github.com/faisalhardin/employee-payroll-system/internal/repo/handler/user"
	"github.com/faisalhardin/employee-payroll-system/pkg/middlewares/auth"
)

type Handlers struct {
	UserHandler       *userhandler.UserHandler
	AttendanceHandler *attendancehandler.AttendanceHandler
}

type Modules struct {
	Handlers       *Handlers
	AuthMiddleware *auth.Options
}
