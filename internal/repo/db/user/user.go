package user

import (
	"context"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	xormlib "github.com/faisalhardin/employee-payroll-system/pkg/xorm"
	"github.com/pkg/errors"
)

const (
	MstUserTable = "mst_user"

	WrapMsgGetUser = "conn.GetUser"
)

type Conn struct {
	DB *xormlib.DBConnect
}

func New(conn *Conn) *Conn {
	return conn
}

func (c *Conn) GetUser(ctx context.Context, params model.SignInRequest) (res model.MstUser, err error) {
	session := c.DB.MasterDB.Table(MstUserTable)
	_, err = session.
		Where("username = ?", params.Username).
		Where("password_hash = ?", params.Password).
		Get(&res)
	if err != nil {
		err = errors.Wrap(err, WrapMsgGetUser)
		return
	}

	return
}

func (c *Conn) ListUser(ctx context.Context) (res []model.MstUser, err error) {
	session := c.DB.MasterDB.Table(MstUserTable)
	err = session.Find(&res)
	if err != nil {
		err = errors.Wrap(err, "conn.ListUser")
		return
	}
	return
}
