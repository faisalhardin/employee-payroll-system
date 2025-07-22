package attendance

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	xormlib "github.com/faisalhardin/employee-payroll-system/pkg/xorm"
	"github.com/pkg/errors"
)

func Test_UpdateReimburse(t *testing.T) {
	mockConn, mockDB := xormlib.NewMockDB()
	defer func() {
		mockConn.Close()
		err := mockDB.ExpectationsWereMet()
		if err != nil {
			t.Error(err)
		}
	}()

	type args struct {
		ctx           context.Context
		reimbursement *model.TrxReimbursement
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		patch   func()
	}{
		{
			name: "success update reimbursement",
			args: args{
				ctx: context.Background(),
				reimbursement: &model.TrxReimbursement{
					ID:          1,
					UserID:      1,
					Status:      "approved",
					Amount:      50000,
					Description: "Updated description",
				},
			},
			wantErr: false,
			patch: func() {
				mockDB.ExpectExec("^UPDATE .*").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "error update reimbursement - database error",
			args: args{
				ctx: context.Background(),
				reimbursement: &model.TrxReimbursement{
					ID:          2,
					UserID:      2,
					Status:      "rejected",
					Amount:      25000,
					Description: "Test description",
				},
			},
			wantErr: true,
			patch: func() {
				mockDB.ExpectExec("^UPDATE .*").
					WillReturnError(errors.New("database error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.patch()

			conn := &Conn{
				DB: &xormlib.DBConnect{
					MasterDB: mockConn,
				}}
			err := conn.UpdateReimbursement(tt.args.ctx, tt.args.reimbursement)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateReimbursement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_SubmitReimbursement(t *testing.T) {
	mockConn, mockDB := xormlib.NewMockDB()
	defer func() {
		mockConn.Close()
		err := mockDB.ExpectationsWereMet()
		if err != nil {
			t.Error(err)
		}
	}()

	type fields struct {
		DB *xormlib.DBConnect
	}
	type args struct {
		ctx           context.Context
		reimbursement *model.TrxReimbursement
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		patch   func()
	}{
		{
			name: "Successful",
			fields: fields{
				DB: &xormlib.DBConnect{
					MasterDB: mockConn,
				},
			},
			args: args{
				ctx: context.Background(),
				reimbursement: &model.TrxReimbursement{
					UserID:      1,
					Status:      "pending",
					Amount:      50000,
					Description: "Transportation allowance",
				},
			},
			wantErr: false,
			patch: func() {
				mockDB.
					ExpectQuery("^INSERT .*").WillReturnRows(
					sqlmock.NewRows([]string{"id"}).AddRow(1),
				)
			},
		},
		{
			name: "Failed at InsertOne",
			fields: fields{
				DB: &xormlib.DBConnect{
					MasterDB: mockConn,
				},
			},
			args: args{
				ctx: context.Background(),
				reimbursement: &model.TrxReimbursement{
					UserID:      1,
					Status:      "pending",
					Amount:      50000,
					Description: "Transportation allowance",
				},
			},
			wantErr: true,
			patch: func() {
				mockDB.
					ExpectQuery("^INSERT .*").
					WillReturnError(
						errors.New("database error"),
					)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conn{
				DB: tt.fields.DB,
			}
			tt.patch()
			err := c.SubmitReimbursement(tt.args.ctx, tt.args.reimbursement)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn.SubmitReimbursement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_ListReimbursementByParams(t *testing.T) {
	mockConn, mockDB := xormlib.NewMockDB()
	defer func() {
		mockConn.Close()
		err := mockDB.ExpectationsWereMet()
		if err != nil {
			t.Error(err)
		}
	}()

	type fields struct {
		DB *xormlib.DBConnect
	}
	type args struct {
		ctx    context.Context
		params model.ListReimbursementParams
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp []model.TrxReimbursement
		wantErr  bool
		patch    func()
	}{
		{
			name: "Successful",
			fields: fields{
				DB: &xormlib.DBConnect{
					MasterDB: mockConn,
				},
			},
			args: args{
				ctx: context.Background(),
				params: model.ListReimbursementParams{
					UserID: 1,
					Status: "pending",
				},
			},
			wantResp: []model.TrxReimbursement{
				{
					ID:          1,
					UserID:      1,
					Status:      "pending",
					Amount:      50000,
					Description: "Transportation allowance",
				},
			},
			wantErr: false,
			patch: func() {
				mockDB.ExpectQuery("^SELECT .*").
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "id_mst_user", "status", "amount", "description"}).
							AddRow(1, 1, "pending", 50000, "Transportation allowance"),
					)
			},
		},
		{
			name: "Failed because no entry found",
			fields: fields{
				DB: &xormlib.DBConnect{
					MasterDB: mockConn,
				},
			},
			args: args{
				ctx: context.Background(),
				params: model.ListReimbursementParams{
					UserID: 999,
					Status: "pending",
				},
			},
			wantResp: ([]model.TrxReimbursement)(nil),
			wantErr:  false,
			patch: func() {
				mockDB.ExpectQuery("^SELECT .*").
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "id_mst_user", "status", "amount", "description"}),
					)
			},
		},
		{
			name: "Failed because get method",
			fields: fields{
				DB: &xormlib.DBConnect{
					MasterDB: mockConn,
				},
			},
			args: args{
				ctx: context.Background(),
				params: model.ListReimbursementParams{
					UserID: 1,
					Status: "pending",
				},
			},
			wantErr: true,
			patch: func() {
				mockDB.ExpectQuery("^SELECT .*").
					WillReturnError(errors.New("database error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conn{
				DB: tt.fields.DB,
			}
			tt.patch()
			gotResp, err := c.ListReimbursementByParams(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn.ListReimbursementByParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("Conn.ListReimbursementByParams() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
