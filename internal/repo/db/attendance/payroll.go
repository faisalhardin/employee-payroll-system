package attendance

import (
	"context"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	"github.com/pkg/errors"
)

const (
	TrxUserPayslipTable = "trx_user_payslip"
	DtlPayrollTable     = "dtl_payroll"
)

func (c *Conn) SubmitPayslips(ctx context.Context, payslips []model.TrxUserPayslip) (err error) {
	session := c.DB.MasterDB.Table(TrxUserPayslipTable)
	_, err = session.Insert(payslips)
	if err != nil {
		return errors.Wrap(err, "SubmitPayslip")
	}
	return
}

func (c *Conn) GetPayslips(ctx context.Context, params model.GetPayslipRequest) (payslips []model.TrxUserPayslip, err error) {
	session := c.DB.MasterDB.Table(TrxUserPayslipTable)

	if params.IDMstPayrollPeriod > 0 {
		session.Where("id_mst_payroll_period = ?", params.IDMstPayrollPeriod)
	}
	if params.UserID > 0 {
		session.Where("id_mst_user = ?", params.UserID)
	}
	err = session.Find(&payslips)
	if err != nil {
		return nil, errors.Wrap(err, "GetPayslips")
	}
	return
}

func (c *Conn) SubmitPayroll(ctx context.Context, payroll model.DtlPayroll) (err error) {
	session := c.DB.MasterDB.Table(DtlPayrollTable)
	_, err = session.Insert(payroll)
	if err != nil {
		return errors.Wrap(err, "SubmitPayroll")
	}
	return
}
