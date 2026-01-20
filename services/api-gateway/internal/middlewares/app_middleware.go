package middlewares

import (
	"log/slog"
	"net/http"
)

func NewTestMiddleware(
	log *slog.Logger,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//TODO

			next.ServeHTTP(w, r)
		})
	}
}
