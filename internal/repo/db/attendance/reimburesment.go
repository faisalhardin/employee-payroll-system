package attendance

import (
	"context"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	"github.com/pkg/errors"
)

const (
	TrxReimbursementTable = "trx_reimbursement"
)

func (c *Conn) SubmitReimbursement(ctx context.Context, reimbursement *model.TrxReimbursement) (err error) {
	session := c.DB.MasterDB.Table(TrxReimbursementTable)
	_, err = session.InsertOne(reimbursement)
	if err != nil {
		return errors.Wrap(err, "conn.SubmitReimbursement")
	}
	return nil
}
