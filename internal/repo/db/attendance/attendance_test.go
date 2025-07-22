package attendance

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	xormlib "github.com/faisalhardin/employee-payroll-system/pkg/xorm"
	"github.com/pkg/errors"
)

func Test_RecordAttendance(t *testing.T) {
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
		ctx        context.Context
		attendance *model.MstAttendance
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
				attendance: &model.MstAttendance{
					IDMstUser:      1,
					AttendanceDate: time.Now(),
				},
			},
			wantErr: false,
			patch: func() {
				mockDB.
					ExpectQuery("^INSERT INTO \"mst_attendance\"").
					WillReturnRows(
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
				attendance: &model.MstAttendance{
					IDMstUser:      1,
					AttendanceDate: time.Now(),
				},
			},
			wantErr: true,
			patch: func() {
				mockDB.
					ExpectQuery("^INSERT INTO \"mst_attendance\"").
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
			err := c.RecordAttendance(tt.args.ctx, tt.args.attendance)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn.RecordAttendance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_GetAttendance(t *testing.T) {
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
		params model.MstAttendance
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.MstAttendance
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
				params: model.MstAttendance{
					IDMstUser:      1,
					AttendanceDate: time.Now(),
				},
			},
			want: model.MstAttendance{
				ID:             1,
				IDMstUser:      1,
				AttendanceDate: time.Now(),
			},
			wantErr: false,
			patch: func() {
				mockDB.ExpectQuery("^SELECT .*").
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "id_mst_user", "attendance_date"}).
							AddRow(1, 1, time.Now()),
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
			got, err := c.GetAttendance(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn.GetAttendance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.ID, tt.want.ID) {
				t.Errorf("Conn.GetAttendance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ListAttendanceByParams(t *testing.T) {
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
		params model.ListAttendanceParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.MstAttendance
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
				params: model.ListAttendanceParams{
					IDsMstUser: []int64{1, 2},
					StartDate:  time.Now(),
					EndDate:    time.Now().AddDate(0, 0, 7),
				},
			},
			want: []model.MstAttendance{
				{
					ID:             1,
					IDMstUser:      1,
					AttendanceDate: time.Now(),
				},
			},
			wantErr: false,
			patch: func() {
				mockDB.ExpectQuery("^SELECT .*").
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "id_mst_user", "attendance_date"}).
							AddRow(1, 1, time.Now()),
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
			got, err := c.ListAttendanceByParams(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn.ListAttendanceByParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("Conn.ListAttendanceByParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_UpdateAttendance(t *testing.T) {
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
		ctx        context.Context
		attendance *model.MstAttendance
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
				attendance: &model.MstAttendance{
					ID:             1,
					IDMstUser:      1,
					AttendanceDate: time.Now(),
				},
			},
			wantErr: false,
			patch: func() {
				mockDB.ExpectExec("^UPDATE \"mst_attendance\"").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conn{
				DB: tt.fields.DB,
			}
			tt.patch()
			err := c.UpdateAttendance(tt.args.ctx, tt.args.attendance)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn.UpdateAttendance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_SubmitOvertime(t *testing.T) {
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
		overtime *model.TrxOvertime
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
				overtime: &model.TrxOvertime{
					UserID:       1,
					OvertimeDate: time.Now(),
					Hours:        2,
				},
			},
			wantErr: false,
			patch: func() {
				mockDB.
					ExpectQuery("^INSERT INTO \"trx_overtime\"").
					WillReturnRows(
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
				overtime: &model.TrxOvertime{
					UserID:       1,
					OvertimeDate: time.Now(),
					Hours:        2,
				},
			},
			wantErr: true,
			patch: func() {
				mockDB.
					ExpectQuery("^INSERT INTO \"trx_overtime\"").
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
			err := c.SubmitOvertime(tt.args.ctx, tt.args.overtime)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn.SubmitOvertime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_CreatePayrollPeriod(t *testing.T) {
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
		ctx          context.Context
		payrolPeriod *model.MstPayrollPeriod
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
				payrolPeriod: &model.MstPayrollPeriod{
					StartDate: time.Now(),
					EndDate:   time.Now().AddDate(0, 1, 0),
				},
			},
			wantErr: false,
			patch: func() {
				mockDB.
					ExpectQuery("^INSERT INTO \"mst_payroll_period\"").
					WillReturnRows(
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
				payrolPeriod: &model.MstPayrollPeriod{
					StartDate: time.Now(),
					EndDate:   time.Now().AddDate(0, 1, 0),
				},
			},
			wantErr: true,
			patch: func() {
				mockDB.
					ExpectQuery("^INSERT INTO \"mst_payroll_period\"").
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
			err := c.CreatePayrollPeriod(tt.args.ctx, tt.args.payrolPeriod)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn.CreatePayrollPeriod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_GetPayrollPeriod(t *testing.T) {
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
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.MstPayrollPeriod
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
				id:  1,
			},
			want: model.MstPayrollPeriod{
				ID:        1,
				StartDate: time.Now(),
				EndDate:   time.Now().AddDate(0, 1, 0),
			},
			wantErr: false,
			patch: func() {
				mockDB.ExpectQuery("^SELECT .*").
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "start_date", "end_date"}).
							AddRow(1, time.Now(), time.Now().AddDate(0, 1, 0)),
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
				id:  1,
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
			got, err := c.GetPayrollPeriod(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn.GetPayrollPeriod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.ID, tt.want.ID) {
				t.Errorf("Conn.GetPayrollPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_UpdatePayrollPeriod(t *testing.T) {
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
		ctx          context.Context
		payrolPeriod *model.MstPayrollPeriod
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
				payrolPeriod: &model.MstPayrollPeriod{
					ID:        1,
					StartDate: time.Now(),
					EndDate:   time.Now().AddDate(0, 1, 0),
				},
			},
			wantErr: false,
			patch: func() {
				mockDB.ExpectExec("^UPDATE \"mst_payroll_period\"").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Failed at Update",
			fields: fields{
				DB: &xormlib.DBConnect{
					MasterDB: mockConn,
				},
			},
			args: args{
				ctx: context.Background(),
				payrolPeriod: &model.MstPayrollPeriod{
					ID:        1,
					StartDate: time.Now(),
					EndDate:   time.Now().AddDate(0, 1, 0),
				},
			},
			wantErr: true,
			patch: func() {
				mockDB.ExpectExec("^UPDATE \"mst_payroll_period\"").
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
			err := c.UpdatePayrollPeriod(tt.args.ctx, tt.args.payrolPeriod)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn.UpdatePayrollPeriod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_GetOvertime(t *testing.T) {
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
		params model.TrxOvertime
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.TrxOvertime
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
				params: model.TrxOvertime{
					UserID:       1,
					OvertimeDate: time.Now(),
				},
			},
			want: model.TrxOvertime{
				ID:           1,
				UserID:       1,
				OvertimeDate: time.Now(),
				Hours:        2,
			},
			wantErr: false,
			patch: func() {
				mockDB.ExpectQuery("^SELECT .*").
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "id_mst_user", "overtime_date", "hours"}).
							AddRow(1, 1, time.Now(), 2),
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
				params: model.TrxOvertime{
					UserID:       1,
					OvertimeDate: time.Now(),
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
			got, err := c.GetOvertime(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn.GetOvertime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.ID, tt.want.ID) {
				t.Errorf("Conn.GetOvertime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_UpdateOvertime(t *testing.T) {
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
		overtime *model.TrxOvertime
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
				overtime: &model.TrxOvertime{
					ID:           1,
					UserID:       1,
					OvertimeDate: time.Now(),
					Hours:        3,
				},
			},
			wantErr: false,
			patch: func() {
				mockDB.ExpectExec("^UPDATE \"trx_overtime\"").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Failed at Update",
			fields: fields{
				DB: &xormlib.DBConnect{
					MasterDB: mockConn,
				},
			},
			args: args{
				ctx: context.Background(),
				overtime: &model.TrxOvertime{
					ID:           1,
					UserID:       1,
					OvertimeDate: time.Now(),
					Hours:        3,
				},
			},
			wantErr: true,
			patch: func() {
				mockDB.ExpectExec("^UPDATE \"trx_overtime\"").
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
			err := c.UpdateOvertime(tt.args.ctx, tt.args.overtime)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn.UpdateOvertime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ListOvertimeByParams(t *testing.T) {
	mockConn, mockDB := xormlib.NewMockDB()
	defer func() {
		mockConn.Close()
		err := mockDB.ExpectationsWereMet()
		if err != nil {
			t.Error(err)
		}
	}()

	jakartaTZ, _ := time.LoadLocation("Asia/Jakarta")
	fixedTime1 := time.Date(2025, 7, 22, 10, 0, 0, 0, jakartaTZ)
	fixedTime2 := time.Date(2025, 7, 23, 10, 0, 0, 0, jakartaTZ)

	type fields struct {
		DB *xormlib.DBConnect
	}
	type args struct {
		ctx    context.Context
		params model.ListOvertimeParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.TrxOvertime
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
				params: model.ListOvertimeParams{
					StartDate:              fixedTime1,
					EndDate:                fixedTime2,
					UserIDs:                []int64{1, 2},
					IsForGeneratingPayroll: true,
					IDMstPayrollPeriod:     0,
				},
			},
			want: []model.TrxOvertime{
				{
					ID:           1,
					UserID:       1,
					OvertimeDate: fixedTime1,
					Hours:        2,
				},
				{
					ID:           2,
					UserID:       2,
					OvertimeDate: fixedTime2,
					Hours:        3,
				},
			},
			wantErr: false,
			patch: func() {
				mockDB.ExpectQuery("^SELECT .*").
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "id_mst_user", "overtime_date", "hours"}).
							AddRow(1, 1, fixedTime1, 2).
							AddRow(2, 2, fixedTime2, 3),
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
				params: model.ListOvertimeParams{
					StartDate:              fixedTime1,
					EndDate:                fixedTime2,
					UserIDs:                []int64{999},
					IsForGeneratingPayroll: true,
				},
			},
			want:    ([]model.TrxOvertime)(nil),
			wantErr: false,
			patch: func() {
				mockDB.ExpectQuery("^SELECT .*").
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "id_mst_user", "overtime_date", "hours"}),
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
				params: model.ListOvertimeParams{
					StartDate: fixedTime1,
					EndDate:   fixedTime2,
					UserIDs:   []int64{1, 2},
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
			_, err := c.ListOvertimeByParams(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn.ListOvertimeByParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
