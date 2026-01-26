package customercontroller

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/fedotovmax/httputils"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/domain"
)

func (c *controller) handleUserSessionActionResponse(
	ctx context.Context,
	w http.ResponseWriter,
	params *handleSessionStatusParams,
) {

	const op = "controller.customers.handleUserSessionActionResponse"

	l := c.log.With(slog.String("op", op))

	switch t := params.Response.Payload.(type) {
	case *userspb.UserSessionActionResponse_Deleted:
		r := domain.NewLoginErrorResponse(
			domain.LoginErrorResponseTypeUserDeleted,
			t.Deleted.GetMessage(),
		)
		httputils.WriteJSON(w, http.StatusForbidden, r)
		return
	case *userspb.UserSessionActionResponse_BadCredentials:
		r := domain.NewLoginErrorResponse(
			domain.LoginErrorResponseTypeBadCredentials,
			t.BadCredentials.GetMessage(),
		)
		httputils.WriteJSON(w, http.StatusForbidden, r)
		return
	case *userspb.UserSessionActionResponse_EmailNotVerified:
		r := domain.NewLoginErrorResponse(
			domain.LoginErrorResponseTypeEmailNotVerified,
			t.EmailNotVerified.GetMessage(),
		)
		httputils.WriteJSON(w, http.StatusForbidden, r)
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
		r := domain.NewLoginErrorResponse(
			domain.LoginErrorResponseTypeBadBypassCode,
			t.BadBypassCode.GetMessage(),
		)
		httputils.WriteJSON(w, http.StatusForbidden, r)
		return
	case *sessionspb.CreateSessionResponse_LoginFromNewDevice:
		r := domain.NewLoginErrorResponse(
			domain.LoginErrorResponseTypeLoginFromNewDevice,
			t.LoginFromNewDevice.GetMessage(),
		)
		httputils.WriteJSON(w, http.StatusForbidden, r)
		return
	case *sessionspb.CreateSessionResponse_SessionCreated:
		httputils.WriteJSON(w, http.StatusCreated, t.SessionCreated)
		return
	case *sessionspb.CreateSessionResponse_UserInBlacklist:
		r := domain.NewLoginErrorResponse(
			domain.LoginErrorResponseTypeUserInBlacklist,
			t.UserInBlacklist.GetMessage(),
		)
		httputils.WriteJSON(w, http.StatusForbidden, r)
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
