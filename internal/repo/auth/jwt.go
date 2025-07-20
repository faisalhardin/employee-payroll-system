package auth

import (
	"context"
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/faisalhardin/employee-payroll-system/internal/config"
	"github.com/faisalhardin/employee-payroll-system/internal/entity/model"
	"github.com/pkg/errors"
)

type Options struct {
	Cfg *config.Config

	JwtOpt JwtOpt
}

type JwtOpt struct {
	JWTPrivateKey string
	jwtSigner     *jwt.HSAlg
	jwtVerifier   *jwt.HSAlg
}

type Claims struct {
	jwt.RegisteredClaims
	Payload model.UserJWTPayload `json:"payload"`
}

func New(opt *Options) (*Options, error) {

	opt.JwtOpt = JwtOpt{
		JWTPrivateKey: opt.Cfg.JWTConfig.Credentials.Secret,
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

func (opt *Options) CreateJWTToken(ctx context.Context, payload model.UserJWTPayload, timeNow, timeExpired time.Time) (tokenStr string, err error) {

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    opt.Cfg.Server.Host,
			Audience:  jwt.Audience{opt.Cfg.Server.Host},
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
