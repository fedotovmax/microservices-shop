package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
)

func NewAccessTokenMiddleware(
	users userspb.UserServiceClient,
	sessions sessionspb.SessionsServiceClient,
	log *slog.Logger,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//TODO
			next.ServeHTTP(w, r)
		})
	}
}

/// =========== example

// ctx := context.WithValue(r.Context(), keys.UserSessionCtx, nil)

// 			authHeader := r.Header.Get(keys.AuthorizationHeader)

// 			splittedHeader := strings.Split(authHeader, " ")

// 			authFall := ex.NewErr(ex.Unauthorized, http.StatusUnauthorized)

// 			if len(splittedHeader) != 2 {
// 				utils.WriteJSON(w, authFall.Status(), authFall)
// 				return
// 			}

// 			accessToken := splittedHeader[1]
// 			claims, err := jwtService.Parse(accessToken, jwt.AccessToken)
// 			if err != nil {
// 				utils.WriteJSON(w, authFall.Status(), authFall)
// 				return
// 			}
// 			user, fall := userService.FindById(r.Context(), claims.UserId)
// 			if fall != nil {
// 				utils.WriteJSON(w, fall.Status(), fall)
// 				return
// 			}
// 			session, fall := sessionRepository.FindByAgentAndUserId(r.Context(), claims.UserAgent, user.UserId)

// 			if fall != nil {
// 				accessToken, refreshToken := utils.RemoveTokensCookie()
// 				http.SetCookie(w, accessToken)
// 				http.SetCookie(w, refreshToken)
// 				forbidden := ex.NewErr(ex.Forbidden, http.StatusForbidden)
// 				utils.WriteJSON(w, forbidden.Status(), forbidden)
// 				return
// 			}

// 			localSession := model.LocalSession{UserId: session.UserId, UserAgent: session.UserAgent, Email: user.Email, Ip: session.Ip}
// 			ctx = context.WithValue(ctx, keys.UserSessionCtx, localSession)
// 			next.ServeHTTP(w, r.WithContext(ctx))
