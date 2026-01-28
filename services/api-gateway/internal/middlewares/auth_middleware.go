package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/fedotovmax/httputils"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/domain"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/utils"
	"github.com/fedotovmax/passport"
)

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

			accessToken, err := utils.ValidateAuthHeader(authHeader)

			if err != nil {
				httputils.WriteJSON(w, http.StatusUnauthorized, httputils.NewError(err.Error()))

			}

			sid, uid, err := passport.Verify(passport.VerifyParams{
				Token:  accessToken,
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
