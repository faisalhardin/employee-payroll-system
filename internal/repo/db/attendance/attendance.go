package attendance

import (
	"context"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	xormlib "github.com/faisalhardin/employee-payroll-system/pkg/xorm"
	"github.com/lib/pq"
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

func (c *Conn) ListAttendanceByParams(ctx context.Context, params model.ListAttendanceParams) (res []model.MstAttendance, err error) {
	session := c.DB.MasterDB.Table(MstAttendanceTable)

	if len(params.IDsMstUser) > 0 {
		session.Where("id_mst_user = ANY(?)", pq.Array(params.IDsMstUser))
	}

	if !params.EndDate.IsZero() && !params.StartDate.IsZero() {
		session.Where("attendance_date BETWEEN ? AND  ?", params.StartDate.Format("2006-01-02"), params.EndDate.Format("2006-01-02"))
	}

	if params.IDMstPayrollPeriod != 0 {
		session.Where("id_mst_payroll_period = ?", params.IDMstPayrollPeriod)
	}

	if params.IsForGeneratingPayroll {
		session.Where("id_mst_payroll_period is null")
	}

	err = session.
		OrderBy("id ASC").
		Find(&res)
	if err != nil {
		return res, errors.Wrap(err, "conn.ListAttenanceByParams")
	}
	return res, nil
}

func (c *Conn) UpdateAttendance(ctx context.Context, attendance *model.MstAttendance) (err error) {
	session := c.DB.MasterDB.Table(MstAttendanceTable)
	_, err = session.Where("id = ?", attendance.ID).Update(attendance)
	if err != nil {
		return errors.Wrap(err, "conn.UpdateAttendance")
	}
	return nil
}

func (c *Conn) CreatePayrollPeriod(ctx context.Context, payrolPeriod *model.MstPayrollPeriod) (err error) {
	session := c.DB.MasterDB.Table(MstPayrollPeriodTable)
	_, err = session.InsertOne(payrolPeriod)
	if err != nil {
		return errors.Wrap(err, "conn.CreatePayrollPeriod")
	}
	return nil
}

func (c *Conn) GetPayrollPeriod(ctx context.Context, id int64) (res model.MstPayrollPeriod, err error) {
	session := c.DB.MasterDB.Table(MstPayrollPeriodTable)
	_, err = session.Where("id = ?", id).Get(&res)
	if err != nil {
		return res, errors.Wrap(err, "conn.GetPayrollPeriod")
	}
	return res, nil
}

func (c *Conn) UpdatePayrollPeriod(ctx context.Context, payrolPeriod *model.MstPayrollPeriod) (err error) {
	session := c.DB.MasterDB.Table(MstPayrollPeriodTable)
	_, err = session.Where("id = ?", payrolPeriod.ID).Update(payrolPeriod)
	if err != nil {
		return errors.Wrap(err, "conn.UpdatePayrollPeriod")
	}
	return nil
}

func (c *Conn) SubmitOvertime(ctx context.Context, overtime *model.TrxOvertime) (err error) {
	session := c.DB.MasterDB.Table(TrxOvertime)
	_, err = session.InsertOne(overtime)
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

func (c *Conn) UpdateOvertime(ctx context.Context, overtime *model.TrxOvertime) (err error) {
	session := c.DB.MasterDB.Table(TrxOvertime)
	_, err = session.Where("id = ?", overtime.ID).Update(overtime)
	if err != nil {
		return errors.Wrap(err, "conn.UpdateOvertime")
	}
	return nil
}

func (c *Conn) ListOvertimeByParams(ctx context.Context, params model.ListOvertimeParams) (res []model.TrxOvertime, err error) {
	session := c.DB.MasterDB.Table(TrxOvertime)

	if !params.StartDate.IsZero() && !params.EndDate.IsZero() {
		session.Where("overtime_date BETWEEN ? and ?", params.StartDate.Format("2006-01-02"), params.EndDate.Format("2006-01-02"))
	}

	if params.IsForGeneratingPayroll {
		session.Where("id_mst_payroll_period is null")
	}

	if params.IDMstPayrollPeriod > 0 {
		session.Where("id_mst_payroll_period = ?", params.IDMstPayrollPeriod)
	}
	if len(params.UserIDs) > 0 {
		session.Where("id_mst_user = any(?)", pq.Array(params.UserIDs))
	}

	err = session.Find(&res)
	if err != nil {
		return res, errors.Wrap(err, "conn.ListOvertime")
	}
	return res, nil
}
