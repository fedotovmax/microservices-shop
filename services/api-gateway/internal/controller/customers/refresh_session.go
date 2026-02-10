package customers

import (
	"log/slog"
	"net/http"

	"github.com/fedotovmax/httputils"
	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	_ "github.com/fedotovmax/microservices-shop/api-gateway/internal/domain"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
	"google.golang.org/grpc/metadata"
)

// @Summary      Refresh session
// @Description  Refresh session by refresh token from headers
// @Router       /customers/session/refresh-session [post]
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param dto body sessionspb.RefreshSessionRequest true "Refresh session dto"
// @Param X-Request-Locale header string false "Locale"
// @Success      201  {object}  sessionspb.SessionCreated
// @Failure      400  {object}  errdetails.BadRequest
// @Failure      401  {object}  httputils.ErrorResponse
// @Failure      403  {object} 	domain.LoginErrorResponse
// @Failure      404  {object}  httputils.ErrorResponse
// @Failure      500  {object}  httputils.ErrorResponse
func (c *controller) refreshSession(w http.ResponseWriter, r *http.Request) {
	const op = "controller.customer.refreshSession"

	l := c.log.With(slog.String("op", op))

	locale := r.Header.Get(keys.HeaderLocale)

	if locale == "" {
		locale = keys.FallbackLocale
	}

	var refreshReq sessionspb.RefreshSessionRequest

	err := httputils.DecodeJSON(r.Body, &refreshReq)

	if err != nil {

		msg, err := i18n.Local.Get(locale, keys.ValidationInvalidBody)

		if err != nil {
			l.Error(err.Error())
		}

		httputils.WriteJSON(w, http.StatusBadRequest, httputils.NewError(msg))
		return
	}

	md := metadata.Pairs(
		keys.MetadataLocaleKey, locale,
	)

	ctx := metadata.NewOutgoingContext(r.Context(), md)

	response, err := c.sessions.RefreshSession(ctx, &sessionspb.RefreshSessionRequest{
		RefreshToken: refreshReq.RefreshToken,
		UserAgent:    refreshReq.UserAgent,
		Ip:           refreshReq.Ip,
	})

	if err != nil {
		httputils.HandleErrorFromGrpc(w, err)
		return
	}

	c.handleCreateSessionResponse(w, response)

}
