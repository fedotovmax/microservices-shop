package customercontroller

import (
	"log/slog"
	"net/http"

	"github.com/fedotovmax/httputils"
	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
	"google.golang.org/grpc/metadata"
)

// @Summary      Send new email verify link
// @Description  Send new email verification link on demand
// @Router       /customers/users/send-new-verify-email-link [patch]
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param X-Request-Locale header string false "Locale"
// @Param dto body userspb.SendNewEmailVerifyLinkRequest true "Send new email verification link with body dto"
// @Success      200  {object}  httputils.ErrorResponse
// @Failure      404  {object}  httputils.ErrorResponse
// @Failure      400  {object}  errdetails.BadRequest
// @Failure      500  {object}  httputils.ErrorResponse
func (c *controller) sendNewEmailVerifyLink(w http.ResponseWriter, r *http.Request) {
	const op = "controller.customer.sendNewEmailVerifyLink"

	l := c.log.With(slog.String("op", op))

	locale := r.Header.Get(keys.HeaderLocale)

	if locale == "" {
		locale = keys.FallbackLocale
	}

	var req userspb.SendNewEmailVerifyLinkRequest

	err := httputils.DecodeJSON(r.Body, &req)

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

	_, err = c.users.SendNewEmailVerifyLink(ctx, &req)

	if err != nil {
		httputils.HandleErrorFromGrpc(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, httputils.OK())

}
