package customercontroller

import (
	"log/slog"
	"net/http"

	"github.com/fedotovmax/httputils"
	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	controllerPkg "github.com/fedotovmax/microservices-shop/api-gateway/internal/controller"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
	"google.golang.org/grpc/metadata"
)

// @Summary      Refresh session
// @Description  Refresh session by refresh token from headers
// @Router       /customers/session/refresh-session [post]
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param X-Request-Locale header string false "Locale"
// @Param X-Refresh-Token header string false "Refresh token"
// @Success      201  {object}  sessionspb.CreateSessionResponse
// @Failure      400  {object}  errdetails.BadRequest
// @Failure      500  {object}  httputils.ErrorResponse
func (c *controller) refreshSession(w http.ResponseWriter, r *http.Request) {
	const op = "controller.customer.refreshSession"

	l := c.log.With(slog.String("op", op))

	locale := r.Header.Get(keys.HeaderLocale)

	if locale == "" {
		locale = keys.FallbackLocale
	}

	userAgent := r.UserAgent()

	ip := controllerPkg.GetRealIP(r)

	refreshToken := r.Header.Get(keys.HeaderRefreshToken)

	if refreshToken == "" {
		msg, err := i18n.Local.Get(locale, keys.Unauthorized)
		if err != nil {
			l.Error(err.Error())
		}
		httputils.WriteJSON(w, http.StatusUnauthorized, httputils.NewError(msg))
	}

	userSessionActionReq := &sessionspb.RefreshSessionRequest{
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		Ip:           ip,
	}

	md := metadata.Pairs(
		keys.MetadataLocaleKey, locale,
	)

	ctx := metadata.NewOutgoingContext(r.Context(), md)

	response, err := c.sessions.RefreshSession(ctx, userSessionActionReq)

	if err != nil {
		httputils.HandleErrorFromGrpc(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusCreated, response)

}
