package customercontroller

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fedotovmax/httputils"
	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	controolerPkg "github.com/fedotovmax/microservices-shop/api-gateway/internal/controller"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
	"google.golang.org/grpc/metadata"
)

type handleSessionStatusParams struct {
	UserAgent string
	IP        string
	Locale    string
	Response  *userspb.UserSessionActionResponse
}

// @Summary      Login in account
// @Description  Login in account
// @Router       /customers/session/login [post]
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param dto body userspb.UserSessionActionRequest true "Dto for login in account"
// @Param X-Request-Locale header string false "Locale"
// @Success      201  {object}  userspb.CreateUserResponse
// @Failure      400  {object}  errdetails.BadRequest
// @Failure      401  {object}  httputils.ErrorResponse
// @Failure      403  {object}  userspb.UserSessionActionResponse
// @Failure      404  {object}  httputils.ErrorResponse
// @Failure      500  {object}  httputils.ErrorResponse
func (c *controller) sessionLogin(w http.ResponseWriter, r *http.Request) {
	const op = "controller.customer.sessionLogin"

	l := c.log.With(slog.String("op", op))

	locale := r.Header.Get(keys.HeaderLocale)

	if locale == "" {
		locale = keys.FallbackLocale
	}

	userAgent := r.UserAgent()

	ip := controolerPkg.GetRealIP(r)

	//TODO:remove
	if ip == "::1" {
		ip = "127.0.0.1"
	}

	bypassCode := r.URL.Query().Get("bypass_code")

	var userSessionActionReq userspb.UserSessionActionRequest

	err := httputils.DecodeJSON(r.Body, &userSessionActionReq)

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
		keys.MetadataSessionBypassCode, bypassCode,
	)

	ctx := metadata.NewOutgoingContext(r.Context(), md)

	response, err := c.users.UserSessionAction(ctx, &userSessionActionReq)

	if err != nil {
		httputils.HandleErrorFromGrpc(w, err)
		return
	}

	c.handleUserSessionActionStatus(ctx, w, &handleSessionStatusParams{
		UserAgent: userAgent,
		IP:        ip,
		Locale:    locale,
		Response:  response,
	})
}

func (c *controller) handleUserSessionActionStatus(ctx context.Context, w http.ResponseWriter, params *handleSessionStatusParams) {

	switch params.Response.UserSessionActionStatus {
	case userspb.UserSessionActionStatus_SESSION_STATUS_BAD_CREDENTIALS:
		msg, _ := i18n.Local.Get(params.Locale, keys.BadCredentials)
		httputils.WriteJSON(w, http.StatusUnauthorized, httputils.NewError(msg))
		return
	case userspb.UserSessionActionStatus_SESSION_STATUS_DELETED, userspb.UserSessionActionStatus_SESSION_STATUS_EMAIL_NOT_VERIFIED:
		httputils.WriteJSON(w, http.StatusForbidden, params.Response)
		return
	case userspb.UserSessionActionStatus_SESSION_STATUS_OK:
		if params.Response.UserId != nil && params.Response.Email != nil {
			res, err := c.sessions.CreateSession(ctx, &sessionspb.CreateSessionRequest{
				Issuer:    fmt.Sprintf("%s.customer_controller", keys.APP_NAME),
				Uid:       *params.Response.UserId,
				UserAgent: params.UserAgent,
				Ip:        params.IP,
			})

			if err != nil {
				httputils.HandleErrorFromGrpc(w, err)
				return
			}
			httputils.WriteJSON(w, http.StatusCreated, res)
			return
		}
		c.log.Error("unexpected grpc response from users client")
		httputils.WriteJSON(w, http.StatusInternalServerError, httputils.NewError("unexpected grpc response when get user information for prepare session"))
		return
	default:
		c.log.Error("unexpected session status")
		msg, _ := i18n.Local.Get(params.Locale, keys.UnexpectedSessionStatus)
		httputils.WriteJSON(w, http.StatusInternalServerError, httputils.NewError(msg))
		return
	}
}
