package attendance

import (
	"context"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	xormlib "github.com/faisalhardin/employee-payroll-system/pkg/xorm"
	"github.com/pkg/errors"
)

const (
	MstAttendanceTable = "mst_attendance"

	WrapMsgCreate = "conn.Create"
)

type Conn struct {
	DB *xormlib.DBConnect
}

func New(conn *Conn) *Conn {
	return conn
}

func (c *Conn) Create(ctx context.Context, attendance *model.MstAttendance) error {
	session := c.DB.MasterDB.Table(MstAttendanceTable)
	_, err := session.InsertOne(attendance)
	if err != nil {
		return errors.Wrap(err, WrapMsgCreate)
	}
	return nil
}

func (c *Conn) GetAttendance(ctx context.Context, params model.MstAttendance) (res model.MstAttendance, err error) {
	session := c.DB.MasterDB.Table(MstAttendanceTable)
	_, err = session.Where("id_mst_user = ? AND attendance_date = ?", params.IDMstUser, params.AttendanceDate.Format("2006-01-02")).Get(&res)
	if err != nil {
		return res, errors.Wrap(err, "conn.GetAttendance")
	}
	return res, nil
}
