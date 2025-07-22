package attendance

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	mocksusecase "github.com/faisalhardin/employee-payroll-system/internal/entity/repo/_mocks"
	"github.com/faisalhardin/employee-payroll-system/pkg/common/binding"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

var (
	mockAttendanceUC *mocksusecase.MockAttendanceUsecaseRepository

	errFoo = errors.New("err")
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		opt  *AttendanceHandler
		want *AttendanceHandler
	}{
		{
			name: "should return the same pointer that was passed in",
			opt:  &AttendanceHandler{},
		},
		{
			name: "should handle nil input",
			opt:  nil,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.opt != nil {
				tt.want = tt.opt
			}

			got := New(tt.opt)
			if got != tt.want {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func initMocks(t *testing.T) *gomock.Controller {
	ctrl := gomock.NewController(t)
	mockAttendanceUC = mocksusecase.NewMockAttendanceUsecaseRepository(ctrl)

	return ctrl
}

func Test_CreatePayrollPeriod(t *testing.T) {
	ctrl := initMocks(t)
	defer ctrl.Finish()

	mockReqBody := `{
		"start_date": "2024-01-01T00:00:00Z",
		"end_date": "2024-01-31T00:00:00Z"
	}`
	startDate, _ := time.Parse("2006-01-02T15:04:05Z", "2024-01-01T00:00:00Z")
	endDate, _ := time.Parse("2006-01-02T15:04:05Z", "2024-01-31T00:00:00Z")
	mockSuccessfulResponse := model.PayrollPeriodResponse{
		ID:        1,
		StartDate: startDate,
		EndDate:   endDate,
	}

	type args struct {
		r *http.Request
		w *httptest.ResponseRecorder
	}

	tests := []struct {
		name       string
		statusCode int
		args       args
		patch      func()
	}{
		{
			name:       "Successful",
			statusCode: http.StatusOK,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/payroll-period", bytes.NewBufferString(mockReqBody))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				mockAttendanceUC.EXPECT().CreatePayrollPeriod(gomock.Any(), gomock.Any()).
					Return(mockSuccessfulResponse, nil).Times(1)
			},
		},
		{
			name:       "Failed",
			statusCode: http.StatusInternalServerError,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/payroll-period", bytes.NewBufferString(mockReqBody))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				mockAttendanceUC.EXPECT().CreatePayrollPeriod(gomock.Any(), gomock.Any()).
					Return(model.PayrollPeriodResponse{}, errFoo).Times(1)
			},
		},
		{
			name:       "Failed at binding",
			statusCode: http.StatusInternalServerError,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/payroll-period", bytes.NewBufferString(mockReqBody))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				bindingBind = func(r *http.Request, targetDecode interface{}) error {
					defer func() {
						bindingBind = binding.Bind
					}()
					return errFoo
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := AttendanceHandler{
				AttendanceUsecase: mockAttendanceUC,
			}
			tt.patch()
			h.CreatePayrollPeriod(tt.args.w, tt.args.r)
			resp := tt.args.w.Result()
			if resp.StatusCode != tt.statusCode {
				t.Errorf("handler.CreatePayrollPeriod expected status %v, got %d", tt.statusCode, resp.StatusCode)
			}
		})
	}
}

func Test_TapIn(t *testing.T) {
	ctrl := initMocks(t)
	defer ctrl.Finish()

	mockSuccessfulResponse := model.TapInResponse{
		AttendanceDate: time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
	}

	type args struct {
		r *http.Request
		w *httptest.ResponseRecorder
	}

	tests := []struct {
		name       string
		statusCode int
		args       args
		patch      func()
	}{
		{
			name:       "Successful",
			statusCode: http.StatusOK,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/attendance/tap-in", nil)
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				mockAttendanceUC.EXPECT().TapIn(gomock.Any(), model.MstAttendance{}).
					Return(mockSuccessfulResponse, nil).Times(1)
			},
		},
		{
			name:       "Failed",
			statusCode: http.StatusInternalServerError,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/attendance/tap-in", nil)
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				mockAttendanceUC.EXPECT().TapIn(gomock.Any(), model.MstAttendance{}).
					Return(model.TapInResponse{}, errFoo).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := AttendanceHandler{
				AttendanceUsecase: mockAttendanceUC,
			}
			tt.patch()
			h.TapIn(tt.args.w, tt.args.r)
			resp := tt.args.w.Result()
			if resp.StatusCode != tt.statusCode {
				t.Errorf("handler.TapIn expected status %v, got %d", tt.statusCode, resp.StatusCode)
			}
		})
	}
}

func Test_SubmitOvertime(t *testing.T) {
	ctrl := initMocks(t)
	defer ctrl.Finish()

	mockReqBody := `{
		"overtime_date": "2024-01-15T00:00:00Z",
		"hours": 4
	}`

	overtimeDate, _ := time.Parse(time.RFC3339, "2024-01-15T00:00:00Z")
	mockSuccessfulResponse := model.SubmitOvertimeResponse{
		OvertimeDate: overtimeDate,
		Hours:        4,
	}

	type args struct {
		r *http.Request
		w *httptest.ResponseRecorder
	}

	tests := []struct {
		name       string
		statusCode int
		args       args
		patch      func()
	}{
		{
			name:       "Successful",
			statusCode: http.StatusOK,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/attendance/overtime", bytes.NewBufferString(mockReqBody))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				mockAttendanceUC.EXPECT().SubmitOvertime(gomock.Any(), gomock.Any()).
					Return(mockSuccessfulResponse, nil).Times(1)
			},
		},
		{
			name:       "Failed",
			statusCode: http.StatusInternalServerError,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/attendance/overtime", bytes.NewBufferString(mockReqBody))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				mockAttendanceUC.EXPECT().SubmitOvertime(gomock.Any(), gomock.Any()).
					Return(model.SubmitOvertimeResponse{}, errFoo).Times(1)
			},
		},
		{
			name:       "Failed at binding",
			statusCode: http.StatusInternalServerError,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/attendance/overtime", bytes.NewBufferString(`{invalid json`))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				bindingBind = func(r *http.Request, targetDecode interface{}) error {
					defer func() {
						bindingBind = binding.Bind
					}()
					return errFoo
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := AttendanceHandler{
				AttendanceUsecase: mockAttendanceUC,
			}
			tt.patch()
			h.SubmitOvertime(tt.args.w, tt.args.r)
			resp := tt.args.w.Result()
			if resp.StatusCode != tt.statusCode {
				t.Errorf("handler.SubmitOvertime expected status %v, got %d", tt.statusCode, resp.StatusCode)
			}
		})
	}
}

func Test_GeneratePayroll(t *testing.T) {
	ctrl := initMocks(t)
	defer ctrl.Finish()

	mockReqBody := `{
		"id_mst_payroll_period": 1
	}`

	type args struct {
		r *http.Request
		w *httptest.ResponseRecorder
	}

	tests := []struct {
		name       string
		statusCode int
		args       args
		patch      func()
	}{
		{
			name:       "Successful",
			statusCode: http.StatusOK,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/payroll/generate", bytes.NewBufferString(mockReqBody))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				mockAttendanceUC.EXPECT().GeneratePayroll(gomock.Any(), gomock.Any()).
					Return(nil).Times(1)
			},
		},
		{
			name:       "Failed",
			statusCode: http.StatusInternalServerError,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/payroll/generate", bytes.NewBufferString(mockReqBody))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				mockAttendanceUC.EXPECT().GeneratePayroll(gomock.Any(), gomock.Any()).
					Return(errFoo).Times(1)
			},
		},
		{
			name:       "Failed at binding",
			statusCode: http.StatusInternalServerError,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/payroll/generate", bytes.NewBufferString(`{invalid json`))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				bindingBind = func(r *http.Request, targetDecode interface{}) error {
					defer func() {
						bindingBind = binding.Bind
					}()
					return errFoo
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := AttendanceHandler{
				AttendanceUsecase: mockAttendanceUC,
			}
			tt.patch()
			h.GeneratePayroll(tt.args.w, tt.args.r)
			resp := tt.args.w.Result()
			if resp.StatusCode != tt.statusCode {
				t.Errorf("handler.GeneratePayroll expected status %v, got %d", tt.statusCode, resp.StatusCode)
			}
		})
	}
}

func Test_GetEmployeePayslip(t *testing.T) {
	ctrl := initMocks(t)
	defer ctrl.Finish()

	mockReqBody := `{
		"id_mst_payroll_period": 1
	}`

	mockSuccessfulResponse := model.GetPayslipResponse{
		StartDate:           time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:             time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		TotalTakeHomePay:    1500000,
		AttendanceDate:      []string{"2024-01-15", "2024-01-16"},
		WorkingDays:         22,
		AttendedDays:        20,
		ProratedSalary:      1363636,
		OvertimeHours:       4,
		OvertimePay:         100000,
		TotalReimbursements: 50000,
	}

	type args struct {
		r *http.Request
		w *httptest.ResponseRecorder
	}

	tests := []struct {
		name       string
		statusCode int
		args       args
		patch      func()
	}{
		{
			name:       "Successful",
			statusCode: http.StatusOK,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/payslip", bytes.NewBufferString(mockReqBody))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				mockAttendanceUC.EXPECT().GetEmployeePayslip(gomock.Any(), gomock.Any()).
					Return(mockSuccessfulResponse, nil).Times(1)
			},
		},
		{
			name:       "Failed",
			statusCode: http.StatusInternalServerError,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/payslip", bytes.NewBufferString(mockReqBody))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				mockAttendanceUC.EXPECT().GetEmployeePayslip(gomock.Any(), gomock.Any()).
					Return(model.GetPayslipResponse{}, errFoo).Times(1)
			},
		},
		{
			name:       "Failed at binding",
			statusCode: http.StatusInternalServerError,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/payslip", bytes.NewBufferString(`{invalid json`))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				bindingBind = func(r *http.Request, targetDecode interface{}) error {
					defer func() {
						bindingBind = binding.Bind
					}()
					return errFoo
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := AttendanceHandler{
				AttendanceUsecase: mockAttendanceUC,
			}
			tt.patch()
			h.GetEmployeePayslip(tt.args.w, tt.args.r)
			resp := tt.args.w.Result()
			if resp.StatusCode != tt.statusCode {
				t.Errorf("handler.GetEmployeePayslip expected status %v, got %d", tt.statusCode, resp.StatusCode)
			}
		})
	}
}

func Test_GetPayroll(t *testing.T) {
	ctrl := initMocks(t)
	defer ctrl.Finish()

	mockReqBody := `{
		"id_mst_payroll_period": 1
	}`

	mockSuccessfulResponse := model.GetPayrollResponse{
		StartDate:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:          time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		TotalTakeHomePay: 5500000,
		EmployeesPayslip: []model.TrxUserPayslip{
			{
				UserID:        123,
				Username:      "john.doe",
				TotalTakeHome: 1500000,
			},
			{
				UserID:        456,
				Username:      "jane.smith",
				TotalTakeHome: 2000000,
			},
		},
	}

	type args struct {
		r *http.Request
		w *httptest.ResponseRecorder
	}

	tests := []struct {
		name       string
		statusCode int
		args       args
		patch      func()
	}{
		{
			name:       "Successful",
			statusCode: http.StatusOK,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/payroll", bytes.NewBufferString(mockReqBody))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				mockAttendanceUC.EXPECT().GetPayroll(gomock.Any(), gomock.Any()).
					Return(mockSuccessfulResponse, nil).Times(1)
			},
		},
		{
			name:       "Failed",
			statusCode: http.StatusInternalServerError,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/payroll", bytes.NewBufferString(mockReqBody))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				mockAttendanceUC.EXPECT().GetPayroll(gomock.Any(), gomock.Any()).
					Return(model.GetPayrollResponse{}, errFoo).Times(1)
			},
		},
		{
			name:       "Failed at binding",
			statusCode: http.StatusInternalServerError,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/payroll", bytes.NewBufferString(`{invalid json`))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				bindingBind = func(r *http.Request, targetDecode interface{}) error {
					defer func() {
						bindingBind = binding.Bind
					}()
					return errFoo
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := AttendanceHandler{
				AttendanceUsecase: mockAttendanceUC,
			}
			tt.patch()
			h.GetPayroll(tt.args.w, tt.args.r)
			resp := tt.args.w.Result()
			if resp.StatusCode != tt.statusCode {
				t.Errorf("handler.GetPayroll expected status %v, got %d", tt.statusCode, resp.StatusCode)
			}
		})
	}
}
