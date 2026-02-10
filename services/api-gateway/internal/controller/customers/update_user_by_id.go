package customers

import (
	"log/slog"
	"net/http"

	"github.com/fedotovmax/httputils"
	"github.com/fedotovmax/i18n"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/utils"
	"google.golang.org/grpc/metadata"
)

// @Summary      Update user profile
// @Description  Update user profile
// @Router       /customers/users/profile [patch]
// @Tags         customers
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param dto body userspb.UpdateUserProfileData true "Update user profile with body dto"
// @Param X-Request-Locale header string false "Locale"
// @Success      200  {object}  httputils.ErrorResponse
// @Failure      401  {object}  httputils.ErrorResponse
// @Failure      400  {object}  errdetails.BadRequest
// @Failure      500  {object}  httputils.ErrorResponse
func (c *controller) updateUserProfile(w http.ResponseWriter, r *http.Request) {

	const op = "controller.customer.updateUserProfile"

	l := c.log.With(slog.String("op", op))

	locale := r.Header.Get(keys.HeaderLocale)

	if locale == "" {
		locale = keys.FallbackLocale
	}

	user, err := utils.GetLocalSession(r)

	if err != nil {
		httputils.WriteJSON(w, http.StatusUnauthorized, httputils.NewError(err.Error()))
		return
	}

	var updateUserProfileData userspb.UpdateUserProfileData

	err = httputils.DecodeJSON(r.Body, &updateUserProfileData)

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

	updateUserProfileReq := &userspb.UpdateUserProfileRequest{
		UserId: user.UID,
		Data:   &updateUserProfileData,
	}

	_, err = c.users.UpdateUserProfile(ctx, updateUserProfileReq)

	if err != nil {
		httputils.HandleErrorFromGrpc(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, httputils.OK())
}
