package auth

import (
	"context"
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/faisalhardin/employee-payroll-system/pkg/common/commonerr"
	"github.com/pkg/errors"
)

type Options struct {
	Cfg *JWTConfig

	JwtOpt     JwtOpt
	ServerHost string
}

type JwtOpt struct {
	JWTPrivateKey string
	jwtSigner     *jwt.HSAlg
	jwtVerifier   *jwt.HSAlg
}

type Claims struct {
	jwt.RegisteredClaims
	Payload UserJWTPayload `json:"payload"`
}

func (claims Claims) Verify() (err error) {

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return commonerr.SetNewTokenExpiredError()
	}

	return nil
}

func New(cfg *JWTConfig) (*Options, error) {

	opt := &Options{
		Cfg: cfg,
	}
	opt.JwtOpt = JwtOpt{
		JWTPrivateKey: cfg.Credentials.Secret,
	}

	// Create signer
	signer, err := jwt.NewSignerHS(jwt.HS256, []byte(opt.JwtOpt.JWTPrivateKey))
	if err != nil {
		return opt, errors.Wrap(err, "NewAuthOpt")
	}
	opt.JwtOpt.jwtSigner = signer

	verifier, err := jwt.NewVerifierHS(jwt.HS256, []byte(opt.JwtOpt.JWTPrivateKey))
	if err != nil {
		return opt, errors.Wrap(err, "NewAuthOpt")
	}
	opt.JwtOpt.jwtVerifier = verifier

	return opt, nil
}

func (opt *Options) CreateJWTToken(ctx context.Context, payload UserJWTPayload, timeNow, timeExpired time.Time) (tokenStr string, err error) {

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    opt.Cfg.ServerHost,
			Audience:  jwt.Audience{opt.Cfg.ServerHost},
			ExpiresAt: jwt.NewNumericDate(timeExpired),
			IssuedAt:  jwt.NewNumericDate(timeNow),
		},
		Payload: payload,
	}

	return opt.generateToken(ctx, claims, timeNow, timeExpired)
}

func (opt *Options) generateToken(ctx context.Context, claims any, timeNow, timeExpired time.Time) (tokenStr string, err error) {

	// Build and sign token
	builder := jwt.NewBuilder(opt.JwtOpt.jwtSigner)
	token, err := builder.Build(&claims)
	if err != nil {

		return
	}

	return token.String(), nil
}
