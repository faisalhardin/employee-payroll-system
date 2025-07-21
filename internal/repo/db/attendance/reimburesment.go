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

func (c *Conn) ListReimbursementByParams(ctx context.Context, params model.ListReimbursementParams) (resp []model.TrxReimbursement, err error) {
	session := c.DB.MasterDB.Table(TrxReimbursementTable)

	if !params.StartDate.IsZero() && !params.EndDate.IsZero() {
		session.Where("created_at BETWEEN ? and ?", params.StartDate.Format("2006-01-02"), params.EndDate.Format("2006-01-02"))
	}

	if params.Status != "" {
		session.Where("status = ?", params.Status)
	}

	err = session.Find(&resp)
	if err != nil {
		return nil, errors.Wrap(err, "conn.ListReimbursementByParams")
	}
	return resp, nil
}

func (c *Conn) UpdateReimbursement(ctx context.Context, reimbursement *model.TrxReimbursement) (err error) {
	session := c.DB.MasterDB.Table(TrxReimbursementTable)
	_, err = session.Where("id = ?", reimbursement.ID).Update(reimbursement)
	if err != nil {
		return errors.Wrap(err, "conn.UpdateReimbursementStatus")
	}
	return nil
}
