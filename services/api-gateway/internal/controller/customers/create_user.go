package customers

import (
	"log/slog"
	"net/http"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"

	"github.com/fedotovmax/httputils"
	"google.golang.org/grpc/metadata"
)

// @Summary      Create user account
// @Description  Create new user account
// @Router       /customers/users [post]
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param dto body userspb.CreateUserRequest true "Create user account with body dto"
// @Param X-Request-Locale header string false "Locale"
// @Success      201  {object}  userspb.CreateUserResponse
// @Failure      400  {object}  errdetails.BadRequest
// @Failure      403  {object}  httputils.ErrorResponse
// @Failure      500  {object}  httputils.ErrorResponse
func (c *controller) createUser(w http.ResponseWriter, r *http.Request) {

	const op = "controller.customer.createUser"

	l := c.log.With(slog.String("op", op))

	locale := r.Header.Get(keys.HeaderLocale)

	if locale == "" {
		locale = keys.FallbackLocale
	}

	var createUserReq userspb.CreateUserRequest

	err := httputils.DecodeJSON(r.Body, &createUserReq)

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

	response, err := c.users.CreateUser(ctx, &createUserReq)

	if err != nil {
		httputils.HandleErrorFromGrpc(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusCreated, response)

}
