package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/faisalhardin/employee-payroll-system/internal/config"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
}

func NewServer(cfg *config.Config, m *Modules) *http.Server {
	port, _ := strconv.Atoi(cfg.Server.Port)
	NewServer := &Server{
		port: port,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(m),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
