package xorm

import (
	"errors"
	"fmt"

	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

type DBConnect struct {
	MasterDB *xorm.Engine
}

type dbContext string

// Context keys
var dbContextSession = dbContext("session")

func NewDBConnection(cfg Config) (dbConnection *DBConnect, err error) {

	masterDB, err := generateXormEngineInstance(cfg.DSN)
	if err != nil {
		return nil, errors.New("failed to make connection to master db")
	}

	return &DBConnect{
		MasterDB: masterDB,
	}, nil
}

func generateXormEngineInstance(dsn string) (*xorm.Engine, error) {

	engine, err := xorm.NewEngine("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create engine: %v", err)
	}

	engine.ShowSQL(true)

	// Ping the database to verify the connection
	if err := engine.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	engine.SetTableMapper(core.GonicMapper{})
	engine.SetColumnMapper(core.GonicMapper{})

	return engine, nil

}
