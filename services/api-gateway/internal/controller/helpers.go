package controller

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/fedotovmax/microservices-shop/api-gateway/internal/domain"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
)

func GetRealIP(r *http.Request) string {
	// 1. X-Forwarded-For (может быть список)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ip := strings.TrimSpace(strings.Split(xff, ",")[0])
		if ip != "" {
			return ip
		}
	}

	// 2. X-Real-IP
	if xrip := strings.TrimSpace(r.Header.Get("X-Real-IP")); xrip != "" {
		return xrip
	}

	// 3. RemoteAddr (ip:port)
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return ip
	}

	return r.RemoteAddr
}

var ErrUnauthorized = errors.New("unauthorized")

func GetLocalSession(r *http.Request) (*domain.LocalSession, error) {
	user, ok := r.Context().Value(keys.SessionCtxKey{}).(*domain.LocalSession)
	if !ok {
		return nil, ErrUnauthorized
	}
	return user, nil
}
