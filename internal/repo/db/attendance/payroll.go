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

func (c *Conn) SubmitPayroll(ctx context.Context, payroll model.DtlPayroll) (err error) {
	session := c.DB.MasterDB.Table(DtlPayrollTable)
	_, err = session.Insert(payroll)
	if err != nil {
		return errors.Wrap(err, "SubmitPayroll")
	}
	return
}
