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

func (c *controller) sessionLogin(w http.ResponseWriter, r *http.Request) {
	const op = "controller.customer.sessionLogin"

	l := c.log.With(slog.String("op", op))

	locale := r.Header.Get(keys.HeaderLocale)

	if locale == "" {
		locale = keys.FallbackLocale
	}

	var userSessionActionReq userspb.UserSessionActionRequest

	err := httputils.DecodeJSON(r.Body, &userSessionActionReq)

	if err != nil {
		msg, err := i18n.Local.Get(locale, keys.ValidationInvalidBody)
		if err != nil {
			l.Error(err.Error())
		}
		httputils.WriteJSON(w, http.StatusBadRequest, domain.NewError(msg))
		return
	}

	md := metadata.Pairs(
		keys.MetadataLocaleKey, locale,
	)

	ctx := metadata.NewOutgoingContext(r.Context(), md)

	response, err := c.users.UserSessionAction(ctx, &userSessionActionReq)

	if err != nil {
		httputils.HandleErrorFromGrpc(w, err)
		return
	}

	c.parseUserSessionActionStatus(w, locale, response)
}

func (c *controller) parseUserSessionActionStatus(w http.ResponseWriter, locale string, res *userspb.UserSessionActionResponse) {

	//TODO: login
	switch res.UserSessionActionStatus {
	case userspb.UserSessionActionStatus_SESSION_STATUS_BAD_CREDENTIALS:
		msg, _ := i18n.Local.Get(locale, keys.BadCredentials)
		httputils.WriteJSON(w, http.StatusBadRequest, domain.NewError(msg))
		return
	case userspb.UserSessionActionStatus_SESSION_STATUS_DELETED, userspb.UserSessionActionStatus_SESSION_STATUS_EMAIL_NOT_VERIFIED:
		httputils.WriteJSON(w, http.StatusForbidden, res)
		return
	case userspb.UserSessionActionStatus_SESSION_STATUS_OK:
		httputils.WriteJSON(w, http.StatusOK, res)
		return
	default:
		c.log.Error("unexpected session status")
		msg, _ := i18n.Local.Get(locale, keys.UnexpectedSessionStatus)
		httputils.WriteJSON(w, http.StatusInternalServerError, domain.NewError(msg))
		return
	}
}
