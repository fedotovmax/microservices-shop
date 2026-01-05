package controller

import (
	"net/http"
	"strings"
)

func GetRealIP(r *http.Request) string {

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	return r.Header.Get("X-Real-IP")
}
