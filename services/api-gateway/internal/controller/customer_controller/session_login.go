package customercontroller

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/fedotovmax/httputils"
	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
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
// @Success      200  {object}  sessionspb.CreateSessionResponse
// @Failure      400  {object}  errdetails.BadRequest
// @Failure      401  {object}  httputils.ErrorResponse
// @Failure      403  {object}  userspb.UserSessionActionResponse
// @Failure      404  {object}  httputils.ErrorResponse
// @Failure      406  {object}  httputils.ErrorResponse
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

func (c *controller) handleUserSessionActionResponse(
	ctx context.Context,
	w http.ResponseWriter,
	params *handleSessionStatusParams,
) {

	const op = "controller.customers.handleUserSessionActionResponse"

	l := c.log.With(slog.String("op", op))

	switch t := params.Response.Payload.(type) {
	case *userspb.UserSessionActionResponse_Deleted:
		httputils.WriteJSON(w, http.StatusForbidden, t.Deleted)
		return
	case *userspb.UserSessionActionResponse_BadCredentials:
		httputils.WriteJSON(w, http.StatusBadRequest, t.BadCredentials)
		return
	case *userspb.UserSessionActionResponse_EmailNotVerified:
		httputils.WriteJSON(w, http.StatusForbidden, t.EmailNotVerified)
		return
	case *userspb.UserSessionActionResponse_Ok:
		sessionResponse, err := c.sessions.CreateSession(ctx, &sessionspb.CreateSessionRequest{
			Uid:              t.Ok.UserId,
			UserAgent:        params.UserAgent,
			Ip:               params.IP,
			BypassCode:       params.BypassCode,
			DeviceTrustToken: params.DeviceTrustToken,
		})

		if err != nil {
			httputils.HandleErrorFromGrpc(w, err)
			return
		}
		c.handleCreateSessionResponse(w, sessionResponse)
		return
	default:
		l.Error(
			"unknown response from grpc.users-service.UserSessionAction",
			slog.Any("response.payload", params.Response.Payload),
		)
		httputils.WriteJSON(
			w,
			http.StatusInternalServerError,
			httputils.NewError("unknown user session action response"),
		)
		return
	}
}

func (c *controller) handleCreateSessionResponse(w http.ResponseWriter, res *sessionspb.CreateSessionResponse) {
	const op = "controller.customers.handleCreateSessionResponse"

	l := c.log.With(slog.String("op", op))

	switch t := res.Payload.(type) {

	case *sessionspb.CreateSessionResponse_BadBypassCode:
		httputils.WriteJSON(w, http.StatusBadRequest, t.BadBypassCode)
		return
	case *sessionspb.CreateSessionResponse_LoginFromNewDevice:
		httputils.WriteJSON(w, http.StatusForbidden, t.LoginFromNewDevice)
		return
	case *sessionspb.CreateSessionResponse_SessionCreated:
		httputils.WriteJSON(w, http.StatusCreated, t.SessionCreated)
		return
	case *sessionspb.CreateSessionResponse_UserInBlacklist:
		httputils.WriteJSON(w, http.StatusForbidden, t.UserInBlacklist)
		return
	default:
		l.Error(
			"unknown response from grpc.sessions-service.CreateSession",
			slog.Any("response.payload", res.Payload),
		)
		httputils.WriteJSON(
			w,
			http.StatusInternalServerError,
			httputils.NewError("unknown create session response"),
		)
	}
}
