package customercontroller

import (
	"net/http"

	"github.com/fedotovmax/httputils"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	controllerPkg "github.com/fedotovmax/microservices-shop/api-gateway/internal/controller"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
	"google.golang.org/grpc/metadata"
)

// @Summary      Get user profile
// @Description  Get user base profile
// @Router       /customers/users/profile [get]
// @Tags         customers
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param X-Request-Locale header string false "Locale"
// @Success      200  {object}  userspb.User
// @Failure      401  {object}  httputils.ErrorResponse
// @Failure      400  {object}  errdetails.BadRequest
// @Failure      404  {object}  httputils.ErrorResponse
// @Failure      500  {object}  httputils.ErrorResponse
func (c *controller) getMe(w http.ResponseWriter, r *http.Request) {

	const op = "controller.customer.getMe"

	locale := r.Header.Get(keys.HeaderLocale)

	if locale == "" {
		locale = keys.FallbackLocale
	}

	user, err := controllerPkg.GetLocalSession(r)

	if err != nil {
		httputils.WriteJSON(w, http.StatusUnauthorized, httputils.NewError(err.Error()))
		return
	}

	md := metadata.Pairs(
		keys.MetadataLocaleKey, locale,
	)

	ctx := metadata.NewOutgoingContext(r.Context(), md)

	response, err := c.users.FindUserByID(ctx, &userspb.FindUserByIDRequest{
		Id: user.UID,
	})

	if err != nil {
		httputils.HandleErrorFromGrpc(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, response)

}
