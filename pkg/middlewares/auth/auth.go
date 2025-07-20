package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/cristalhq/jwt/v5"
	"github.com/faisalhardin/employee-payroll-system/pkg/common/commonerr"
	commonwriter "github.com/faisalhardin/employee-payroll-system/pkg/common/writer"
)

const (
	ContentLength         = "Content-Length"
	ContentType           = "Content-Type"
	Authorization         = "Authorization"
	AccountsAuthorization = "accounts-authorization"
	Bearer                = "Bearer %s"
	Key                   = "key=%s"
	Basic                 = "Basic %s"
	XAppKey               = "X-App-Key"
)

var AllowedHeaders = []string{
	"Accept",
	ContentType,
	ContentLength,
	"Authorization",
	"Accept-Encoding",
	"accounts-authorization",
	"X-CSRF-Token",
	"API-KEY",
	"X-Device",
	"X-Element-ID",
	"x-requested-with",
	XAppKey,
}

var AllowedMethodRequest = []string{
	"OPTIONS",
	"GET",
	"POST",
	"PUT",
	"DELETE",
	"PATCH",
}

type userAuth struct{}

var (
	userContextKey = userAuth{}
)

func (opt *Options) AuthHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		bearerToken := r.Header.Get("Authorization")
		token, err := GetBearerToken(bearerToken)
		if err != nil {
			handleError(ctx, w, r, err)
			return
		}

		userDetail, err := opt.HandleAuthMiddleware(ctx, token)
		if err != nil {
			handleError(ctx, w, r, err)
			return
		}

		ctx = SetUserDetailToCtx(ctx, userDetail)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

	})
}

func (opt *Options) VerifyJWT(jwtToken string, claims any) (err error) {

	tokenParsed, err := jwt.Parse([]byte(jwtToken), opt.JwtOpt.jwtVerifier)
	if err != nil && errors.Is(err, jwt.ErrInvalidFormat) {
		return commonerr.SetNewBadRequest("authorization", err.Error())
	} else if err != nil {
		return err
	}

	marshalledToken, err := tokenParsed.Claims().MarshalJSON()
	if err != nil {
		return err
	}

	err = json.Unmarshal(marshalledToken, claims)
	if err != nil {
		return err
	}

	return nil
}

func (opt *Options) GetTokenClaims(token string) (claims *Claims, err error) {

	claims = &Claims{}

	err = opt.VerifyJWT(token, claims)
	if err != nil {
		return
	}

	return claims, nil
}

func (opt *Options) HandleAuthMiddleware(ctx context.Context, token string) (ret UserJWTPayload, err error) {

	claims, err := opt.GetTokenClaims(token)
	if err != nil {
		return
	}

	err = claims.Verify()
	if err != nil {
		return
	}

	return claims.Payload, nil
}

func GetBearerToken(token string) (string, error) {
	splitToken := strings.Split(token, "Bearer ")
	if len(splitToken) != 2 {
		return "", errors.New("invalid token")
	}

	return splitToken[1], nil
}

func SetUserDetailToCtx(ctx context.Context, data UserJWTPayload) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, userContextKey, data)
}

func GetUserDetailFromCtx(ctx context.Context) (UserJWTPayload, bool) {
	user, ok := ctx.Value(userContextKey).(UserJWTPayload)
	return user, ok
}

func handleError(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	commonwriter.SetError(ctx, w, err)
}
