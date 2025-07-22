package user

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mocksusecase "github.com/faisalhardin/employee-payroll-system/internal/entity/repo/_mocks"
	"github.com/faisalhardin/employee-payroll-system/pkg/common/binding"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

var (
	mockUserUC *mocksusecase.MockUserUsecaseRepository

	errFoo = errors.New("err")
)

func initMocks(t *testing.T) *gomock.Controller {
	ctrl := gomock.NewController(t)
	mockUserUC = mocksusecase.NewMockUserUsecaseRepository(ctrl)

	return ctrl
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		opt  *UserHandler
		want *UserHandler
	}{
		{
			name: "should return the same pointer that was passed in",
			opt:  &UserHandler{},
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

func Test_SignIn(t *testing.T) {
	ctrl := initMocks(t)
	defer ctrl.Finish()

	mockReqBody := `{
		"username": "john.doe",
		"password": "password123"
	}`

	mockSuccessfulResponse := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

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
					req := httptest.NewRequest(http.MethodPost, "/auth/signin", bytes.NewBufferString(mockReqBody))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				mockUserUC.EXPECT().SignIn(gomock.Any(), gomock.Any()).
					Return(mockSuccessfulResponse, nil).Times(1)
			},
		},
		{
			name:       "Failed",
			statusCode: http.StatusInternalServerError,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/auth/signin", bytes.NewBufferString(mockReqBody))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
				w: httptest.NewRecorder(),
			},
			patch: func() {
				mockUserUC.EXPECT().SignIn(gomock.Any(), gomock.Any()).
					Return("", errFoo).Times(1)
			},
		},
		{
			name:       "Failed at binding",
			statusCode: http.StatusInternalServerError,
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/auth/signin", bytes.NewBufferString(`{invalid json`))
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
			h := UserHandler{
				UserUsecase: mockUserUC,
			}
			tt.patch()
			h.SignIn(tt.args.w, tt.args.r)
			resp := tt.args.w.Result()
			if resp.StatusCode != tt.statusCode {
				t.Errorf("handler.SignIn expected status %v, got %d", tt.statusCode, resp.StatusCode)
			}
		})
	}
}
