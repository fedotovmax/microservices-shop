package utils

import (
	"errors"
	"net/http"

	"github.com/fedotovmax/microservices-shop/api-gateway/internal/domain"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
)

var ErrUnauthorized = errors.New("unauthorized")

func GetLocalSession(r *http.Request) (*domain.LocalSession, error) {
	user, ok := r.Context().Value(keys.SessionCtxKey{}).(*domain.LocalSession)
	if !ok {
		return nil, ErrUnauthorized
	}
	return user, nil
}
