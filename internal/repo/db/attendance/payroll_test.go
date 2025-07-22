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

func Test_SubmitPayslips(t *testing.T) {
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
		ctx      context.Context
		payslips []model.TrxUserPayslip
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
				payslips: []model.TrxUserPayslip{
					{
						UserID:             1,
						IDMstPayrollPeriod: 1,
						BaseSalary:         15000000,
						WorkingDays:        22,
						AttendedDays:       20,
						TotalTakeHome:      14500000,
					},
				},
			},
			wantErr: false,
			patch: func() {
				mockDB.
					ExpectExec("^INSERT .*").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Failed at Insert",
			fields: fields{
				DB: &xormlib.DBConnect{
					MasterDB: mockConn,
				},
			},
			args: args{
				ctx: context.Background(),
				payslips: []model.TrxUserPayslip{
					{
						UserID:             1,
						IDMstPayrollPeriod: 1,
						BaseSalary:         15000000,
						WorkingDays:        22,
						AttendedDays:       20,
						TotalTakeHome:      14500000,
					},
				},
			},
			wantErr: true,
			patch: func() {
				mockDB.
					ExpectExec("^INSERT INTO \"trx_user_payslip\"").
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
			err := c.SubmitPayslips(tt.args.ctx, tt.args.payslips)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn.SubmitPayslips() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_GetPayslips(t *testing.T) {
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
		params model.GetPayslipRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantPayslips []model.TrxUserPayslip
		wantErr      bool
		patch        func()
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
				params: model.GetPayslipRequest{
					IDMstPayrollPeriod: 1,
					UserID:             1,
				},
			},
			wantPayslips: []model.TrxUserPayslip{
				{
					ID:                 1,
					UserID:             1,
					IDMstPayrollPeriod: 1,
					BaseSalary:         15000000,
					WorkingDays:        22,
					AttendedDays:       20,
					TotalTakeHome:      14500000,
				},
			},
			wantErr: false,
			patch: func() {
				mockDB.ExpectQuery("^SELECT .*").
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "id_mst_user", "id_mst_payroll_period", "base_salary", "working_days", "attended_days", "total_take_home"}).
							AddRow(1, 1, 1, 15000000, 22, 20, 14500000),
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
				params: model.GetPayslipRequest{
					IDMstPayrollPeriod: 1,
					UserID:             1,
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
			gotPayslips, err := c.GetPayslips(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn.GetPayslips() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPayslips, tt.wantPayslips) {
				t.Errorf("Conn.GetPayslips() = %v, want %v", gotPayslips, tt.wantPayslips)
			}
		})
	}
}

func Test_SubmitPayroll(t *testing.T) {
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
		ctx     context.Context
		payroll model.DtlPayroll
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
				payroll: model.DtlPayroll{
					IDMstPayrollPeriod: 1,
					TotalTakeHome:      157000000,
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
			name: "Failed at Insert",
			fields: fields{
				DB: &xormlib.DBConnect{
					MasterDB: mockConn,
				},
			},
			args: args{
				ctx: context.Background(),
				payroll: model.DtlPayroll{
					IDMstPayrollPeriod: 1,
					TotalTakeHome:      157000000,
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
			err := c.SubmitPayroll(tt.args.ctx, tt.args.payroll)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn.SubmitPayroll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_GetPayrollDetail(t *testing.T) {
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
		params model.GetDtlPayrollRequest
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantPayrollDetail model.DtlPayroll
		wantErr           bool
		patch             func()
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
				params: model.GetDtlPayrollRequest{
					IDMstPayrollPeriod: 1,
				},
			},
			wantPayrollDetail: model.DtlPayroll{
				ID:                 1,
				IDMstPayrollPeriod: 1,
				TotalTakeHome:      157000000,
			},
			wantErr: false,
			patch: func() {
				mockDB.ExpectQuery("^SELECT .*").
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "id_mst_payroll_period", "total_employees", "total_base_salary", "total_overtime_pay", "total_reimbursement", "total_take_home"}).
							AddRow(1, 1, 10, 150000000, 5000000, 2000000, 157000000),
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
				params: model.GetDtlPayrollRequest{
					IDMstPayrollPeriod: 1,
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
			gotPayrollDetail, err := c.GetPayrollDetail(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn.GetPayrollDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPayrollDetail, tt.wantPayrollDetail) {
				t.Errorf("Conn.GetPayrollDetail() = %v, want %v", gotPayrollDetail, tt.wantPayrollDetail)
			}
		})
	}
}
