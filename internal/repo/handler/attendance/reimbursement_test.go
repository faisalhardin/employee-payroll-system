package attendance

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	"github.com/faisalhardin/employee-payroll-system/pkg/common/binding"
	"github.com/golang/mock/gomock"
)

func Test_SubmitReimbursement(t *testing.T) {
	ctrl := initMocks(t)
	defer ctrl.Finish()

	mockReqBody := `{
		"amount": 100000,
		"description": "Transportation cost"
	}`

	mockSuccessfulResponse := model.SubmitReimbursementResponse{
		ID:          1,
		Amount:      100000,
		Description: "Transportation cost",
		Status:      "pending",
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
					req := httptest.NewRequest(http.MethodPost, "/attendance/reimbursement", bytes.NewBufferString(mockReqBody))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				mockAttendanceUC.EXPECT().SubmitReimbursement(gomock.Any(), gomock.Any()).
					Return(mockSuccessfulResponse, nil).Times(1)
			},
		},
		{
			name:       "Failed",
			statusCode: http.StatusInternalServerError,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/attendance/reimbursement", bytes.NewBufferString(mockReqBody))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				mockAttendanceUC.EXPECT().SubmitReimbursement(gomock.Any(), gomock.Any()).
					Return(model.SubmitReimbursementResponse{}, errFoo).Times(1)
			},
		},
		{
			name:       "Failed at binding",
			statusCode: http.StatusInternalServerError,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/attendance/reimbursement", bytes.NewBufferString(`{invalid json`))
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
			h.SubmitReimbursement(tt.args.w, tt.args.r)
			resp := tt.args.w.Result()
			if resp.StatusCode != tt.statusCode {
				t.Errorf("handler.SubmitReimbursement expected status %v, got %d", tt.statusCode, resp.StatusCode)
			}
		})
	}
}
