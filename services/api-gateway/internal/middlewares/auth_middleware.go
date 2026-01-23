package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/fedotovmax/httputils"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/domain"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
	"github.com/fedotovmax/passport"
)

var emptyHeader = "empty authorization header"
var badHeaderFormat = "bad authorization header format"
var expiredOrBadSignature = "the token has expired or the signature has been changed"

func NewAuthMiddleware(
	log *slog.Logger,
	tokenSecret string,
	issuer string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//TODO

			authHeader := r.Header.Get(keys.HeaderAuthorization)

			if authHeader == "" {
				httputils.WriteJSON(w, http.StatusUnauthorized, httputils.NewError(emptyHeader))
				return
			}

			authHeaderParts := strings.Split(authHeader, " ")

			if len(authHeaderParts) != 2 {
				httputils.WriteJSON(w, http.StatusUnauthorized, httputils.NewError(badHeaderFormat))
				return
			}

			if authHeaderParts[0] != keys.BearerAuthorizationToken {
				httputils.WriteJSON(w, http.StatusUnauthorized, httputils.NewError(badHeaderFormat))
				return
			}

			sid, uid, err := passport.Verify(passport.VerifyParams{
				Token:  authHeaderParts[1],
				Issuer: issuer,
				Secret: tokenSecret,
			})

			if err != nil {
				httputils.WriteJSON(w, http.StatusUnauthorized, httputils.NewError(expiredOrBadSignature))
				return
			}

			ctx := context.WithValue(r.Context(), keys.SessionCtxKey{}, &domain.LocalSession{UID: uid, SID: sid})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
