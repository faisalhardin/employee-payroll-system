package config

import (
	"fmt"
	"os"

	auth "github.com/faisalhardin/employee-payroll-system/pkg/middlewares/auth"
	xormlib "github.com/faisalhardin/employee-payroll-system/pkg/xorm"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    Server         `yaml:"server"`
	JWTConfig auth.JWTConfig `yaml:"jwt_config"`
	DBConfig  DBConfig       `yaml:"db_config"`
}

type DBConfig struct {
	DBMaster xormlib.Config `yaml:"db_master"`
}

type Server struct {
	Name    string `yaml:"name"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	BaseURL string `yaml:"base_url"`
}

func New(repoName string) (*Config, error) {
	dir, _ := os.Getwd()
	filename := "files/env/" + repoName + ".yaml"

	f, err := os.Open(fmt.Sprintf("%s/%s", dir, filename))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Failed to close file: %s\n", err)
		}
	}()

	var cfg Config
	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
