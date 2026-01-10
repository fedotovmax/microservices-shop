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

// @Summary      Get user base info by id
// @Description  Get user base info by id
// @Router       /customers/users/{id} [get]
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param id path string true "User id parameter"
// @Param X-Request-Locale header string false "Locale"
// @Success      200  {object}  userspb.User
// @Failure      400  {object}  errdetails.BadRequest
// @Failure      404  {object}  httputils.ErrorResponse
// @Failure      500  {object}  httputils.ErrorResponse
func (c *controller) findUserByID(w http.ResponseWriter, r *http.Request) {

	const op = "controller.customer.findUserByID"

	l := c.log.With(slog.String("op", op))

	locale := r.Header.Get(keys.HeaderLocale)

	if locale == "" {
		locale = keys.FallbackLocale
	}

	userId := r.PathValue("id")

	if userId == "" {
		msg, err := i18n.Local.Get(locale, keys.ValidationEmptyID)
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

	response, err := c.users.FindUserByID(ctx, &userspb.FindUserByIDRequest{
		Id: userId,
	})

	if err != nil {
		httputils.HandleErrorFromGrpc(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, response)

}
