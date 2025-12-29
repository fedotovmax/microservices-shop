package jwtadapter

import (
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type jwtAdapter struct {
	config *Config
}

var ErrInvalidToken = errors.New("invalid token")

var ErrParseClaims = errors.New("error when parse claims")

func New(cfg *Config) *jwtAdapter {
	return &jwtAdapter{
		config: cfg,
	}
}

func (a *jwtAdapter) CreateAccessToken(issuer, uid, sid string) (*domain.NewAccessToken, error) {

	const op = "adapter.auth.jwt.CreateAccessToken"

	now := time.Now()

	accessExpTime := now.Add(a.config.AccessTokenExpDuration)

	accessClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(accessExpTime),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    issuer,
		ID:        sid,
		Subject:   uid,
	}

	accessTokenObject := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	accessToken, err := accessTokenObject.SignedString([]byte(a.config.AccessTokenSecret))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	token := &domain.NewAccessToken{AccessToken: accessToken, AccessExpTime: accessExpTime}

	return token, nil

}

func (a *jwtAdapter) ParseAccessToken(token string, issuer string) (jti string, uid string, perr error) {

	const op = "adapter.auth.jwt.Parse"

	opts := []jwt.ParserOption{
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithIssuer(issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	}

	result, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(a.config.AccessTokenSecret), nil
	}, opts...)

	if err != nil {
		perr = fmt.Errorf("%s: %w: %v", op, ErrParseClaims, err)
		return jti, uid, perr
	}

	if !result.Valid {
		perr = fmt.Errorf("%s: %w", op, ErrInvalidToken)
		return jti, uid, perr
	}

	claims, ok := result.Claims.(*jwt.RegisteredClaims)

	if !ok {
		perr = fmt.Errorf("%s: %w", op, ErrParseClaims)
		return jti, uid, perr
	}

	if claims.Subject == "" || claims.ID == "" {
		perr = fmt.Errorf("%s: %w", op, ErrInvalidToken)
		return "", "", perr
	}

	jti, uid = claims.ID, claims.Subject

	return jti, uid, perr
}
