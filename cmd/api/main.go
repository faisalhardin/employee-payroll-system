package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/faisalhardin/employee-payroll-system/internal/config"
	"github.com/faisalhardin/employee-payroll-system/internal/server"

	attendancedb "github.com/faisalhardin/employee-payroll-system/internal/repo/db/attendance"
	userdb "github.com/faisalhardin/employee-payroll-system/internal/repo/db/user"
	"github.com/faisalhardin/employee-payroll-system/pkg/middlewares/auth"

	attendanceusecase "github.com/faisalhardin/employee-payroll-system/internal/repo/usecase/attendance"
	userusecase "github.com/faisalhardin/employee-payroll-system/internal/repo/usecase/user"

	attendancehandler "github.com/faisalhardin/employee-payroll-system/internal/repo/handler/attendance"
	userhandler "github.com/faisalhardin/employee-payroll-system/internal/repo/handler/user"

	xormlib "github.com/faisalhardin/employee-payroll-system/pkg/xorm"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {

	loc, _ := time.LoadLocation("Asia/Jakarta")
	time.Local = loc

	// init config
	cfg, err := config.New("envconfig")
	if err != nil {
		log.Fatalf("failed to init the config: %v", err)
	}

	db, err := xormlib.NewDBConnection(cfg.DBConfig.DBMaster)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
		return
	}

	authRepo, err := auth.New(&cfg.JWTConfig)
	if err != nil {
		log.Fatalf("failed to init auth repo: %v", err)
		return
	}
	attendanceRepo := attendancedb.New(&attendancedb.Conn{
		DB: db,
	})

	userDB := userdb.New(&userdb.Conn{
		DB: db,
	})

	userUC := userusecase.New(&userusecase.Usecase{
		Cfg:      cfg,
		UserDB:   userDB,
		AuthRepo: authRepo,
	})
	attendanceUC := attendanceusecase.New(attendanceusecase.Usecase{
		AttendanceDB: attendanceRepo,
	})

	userHandler := userhandler.New(&userhandler.UserHandler{
		UserUsecase: userUC,
	})

	attendanceHandler := attendancehandler.New(&attendancehandler.AttendanceHandler{
		AttendanceUsecase: attendanceUC,
	})

	handlers := &server.Handlers{
		UserHandler:       userHandler,
		AttendanceHandler: attendanceHandler,
	}

	server := server.NewServer(cfg, &server.Modules{
		Handlers:       handlers,
		AuthMiddleware: authRepo,
	})

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
