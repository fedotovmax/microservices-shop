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
		r := domain.LoginErrorResponse{
			Type:    domain.LoginErrorResponseTypeUserDeleted,
			Message: t.Deleted.GetMessage(),
		}
		httputils.WriteJSON(w, http.StatusForbidden, r)
		return
	case *userspb.UserSessionActionResponse_BadCredentials:
		r := domain.LoginErrorResponse{
			Type:    domain.LoginErrorResponseTypeBadCredentials,
			Message: t.BadCredentials.GetMessage(),
		}
		httputils.WriteJSON(w, http.StatusForbidden, r)
		return
	case *userspb.UserSessionActionResponse_EmailNotVerified:

		r := domain.LoginErrorResponse{
			Type:    domain.LoginErrorResponseTypeEmailNotVerified,
			Message: t.EmailNotVerified.GetMessage(),
			UserId:  &t.EmailNotVerified.UserId,
		}
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
		r := domain.LoginErrorResponse{
			Type:    domain.LoginErrorResponseTypeBadBypassCode,
			Message: t.BadBypassCode.GetMessage(),
		}
		httputils.WriteJSON(w, http.StatusForbidden, r)
		return
	case *sessionspb.CreateSessionResponse_LoginFromNewDevice:
		r := domain.LoginErrorResponse{
			Type:    domain.LoginErrorResponseTypeLoginFromNewDevice,
			Message: t.LoginFromNewDevice.GetMessage(),
		}
		httputils.WriteJSON(w, http.StatusForbidden, r)
		return
	case *sessionspb.CreateSessionResponse_SessionCreated:
		httputils.WriteJSON(w, http.StatusCreated, t.SessionCreated)
		return
	case *sessionspb.CreateSessionResponse_UserInBlacklist:
		r := domain.LoginErrorResponse{
			Type:    domain.LoginErrorResponseTypeUserInBlacklist,
			Message: t.UserInBlacklist.GetMessage(),
		}
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
		return
	}
}
