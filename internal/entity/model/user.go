package model

import (
	"database/sql"
	"time"
)

// MstUser represents both employees and admin for mst_user table
type MstUser struct {
	ID           int64          `json:"id" xorm:"'id'"`
	Username     string         `json:"username" xorm:"'username'"`
	PasswordHash string         `json:"-" xorm:"'password_hash'"`
	Role         string         `json:"role" xorm:"'role'"`
	Salary       int64          `json:"salary,omitempty" xorm:"'salary'"`
	CreatedAt    time.Time      `json:"created_at" xorm:"'created_at' created"`
	UpdatedAt    time.Time      `json:"updated_at" xorm:"'updated_at' updated"`
	CreatedBy    sql.NullString `json:"created_by,omitempty" xorm:"created_by"`
	UpdatedBy    sql.NullString `json:"updated_by,omitempty" xorm:"updated_by"`
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
