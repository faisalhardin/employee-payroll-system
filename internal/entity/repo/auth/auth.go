package auth

import (
	"context"
	"net/http"
	"time"

	authrepo "github.com/faisalhardin/employee-payroll-system/pkg/middlewares/auth"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -self_package=github.com/faisalhardin/employee-payroll-system/internal/entity/repo/auth -destination=../_mocks/mock_auth.go -package=mock github.com/faisalhardin/employee-payroll-system/internal/entity/repo/auth Authenticator
type Authenticator interface {
	CreateJWTToken(ctx context.Context, payload authrepo.UserJWTPayload, timeNow, timeExpired time.Time) (tokenStr string, err error)
	AuthHandler(next http.Handler) http.Handler
	VerifyJWT(jwtToken string, claims any) (err error)
	GetTokenClaims(token string) (claims *authrepo.Claims, err error)
	HandleAuthMiddleware(ctx context.Context, token string) (ret authrepo.UserJWTPayload, err error)
}
