package attendance

import (
	"context"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	xormlib "github.com/faisalhardin/employee-payroll-system/pkg/xorm"
	"github.com/pkg/errors"
)

const (
	MstAttendanceTable    = "mst_attendance"
	MstPayrollPeriodTable = "mst_payroll_period"
	TrxOvertime           = "trx_overtime"

	WrapMsgCreate = "conn.Create"
)

type Conn struct {
	DB *xormlib.DBConnect
}

func New(conn *Conn) *Conn {
	return conn
}

func (c *Conn) RecordAttendance(ctx context.Context, attendance *model.MstAttendance) error {
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

func (c *Conn) CreatePayrollPeriod(ctx context.Context, payrolPeriod *model.MstPayrollPeriod) (err error) {
	session := c.DB.MasterDB.Table(MstPayrollPeriodTable)
	session.InsertOne(payrolPeriod)
	if err != nil {
		return errors.Wrap(err, "conn.CreatePayrollPeriod")
	}
	return nil
}

func (c *Conn) SubmitOvertime(ctx context.Context, overtime *model.TrxOvertime) (err error) {
	session := c.DB.MasterDB.Table(TrxOvertime)
	session.InsertOne(overtime)
	if err != nil {
		return errors.Wrap(err, "conn.SubmitOvertime")
	}
	return nil
}

func (c *Conn) GetOvertime(ctx context.Context, params model.TrxOvertime) (res model.TrxOvertime, err error) {
	session := c.DB.MasterDB.Table(TrxOvertime)
	_, err = session.Where("id_mst_user = ? AND overtime_date = ?", params.UserID, params.OvertimeDate.Format("2006-01-02")).Get(&res)
	if err != nil {
		return res, errors.Wrap(err, "conn.GetOvertime")
	}
	return res, nil
}
