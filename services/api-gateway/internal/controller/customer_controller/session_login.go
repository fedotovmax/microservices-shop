package customercontroller

import (
	"log/slog"
	"net/http"

	"github.com/fedotovmax/httputils"
	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/domain"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
	"google.golang.org/grpc/metadata"
)

type handleSessionStatusParams struct {
	UserAgent        string
	IP               string
	Locale           string
	BypassCode       *string
	DeviceTrustToken *string
	Response         *userspb.UserSessionActionResponse
}

// @Summary      Login in account
// @Description  Login in account
// @Router       /customers/session/login [post]
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param dto body domain.LoginInput true "Dto for login in account"
// @Param X-Request-Locale header string false "Locale"
// @Success      201  {object}  sessionspb.SessionCreated
// @Failure      400  {object}  errdetails.BadRequest
// @Failure      401  {object}  httputils.ErrorResponse
// @Failure      403  {object}  domain.LoginErrorResponse
// @Failure      404  {object}  httputils.ErrorResponse
// @Failure      500  {object}  httputils.ErrorResponse
func (c *controller) sessionLogin(w http.ResponseWriter, r *http.Request) {
	const op = "controller.customer.sessionLogin"

	l := c.log.With(slog.String("op", op))

	locale := r.Header.Get(keys.HeaderLocale)

	if locale == "" {
		locale = keys.FallbackLocale
	}

	var loginInput domain.LoginInput

	err := httputils.DecodeJSON(r.Body, &loginInput)

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

	response, err := c.users.UserSessionAction(ctx, &userspb.UserSessionActionRequest{
		Email:    loginInput.Email,
		Password: loginInput.Password,
	})

	if err != nil {
		httputils.HandleErrorFromGrpc(w, err)
		return
	}

	c.handleUserSessionActionResponse(ctx, w, &handleSessionStatusParams{
		UserAgent:        loginInput.UserAgent,
		IP:               loginInput.Ip,
		Locale:           locale,
		BypassCode:       loginInput.BypassCode,
		DeviceTrustToken: loginInput.DeviceTrustToken,
		Response:         response,
	})

}
