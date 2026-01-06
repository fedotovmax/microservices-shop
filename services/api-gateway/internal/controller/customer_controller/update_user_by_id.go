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

func (c *controller) updateUserByID(w http.ResponseWriter, r *http.Request) {

	const op = "controller.customer.updateUserByID"

	l := c.log.With(slog.String("op", op))

	locale := r.Header.Get(keys.HeaderLocale)

	if locale == "" {
		locale = keys.FallbackLocale
	}

	//TODO: get from session
	userId := r.PathValue("id")

	if userId == "" {
		msg, err := i18n.Local.Get(locale, keys.ValidationEmptyID)
		if err != nil {
			l.Error(err.Error())
		}
		httputils.WriteJSON(w, http.StatusBadRequest, httputils.NewError(msg))
		return
	}

	var updateUserProfileData userspb.UpdateUserProfileData

	err := httputils.DecodeJSON(r.Body, &updateUserProfileData)

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
		UserId: userId,
		Data:   &updateUserProfileData,
	}

	_, err = c.users.UpdateUserProfile(ctx, updateUserProfileReq)

	if err != nil {
		httputils.HandleErrorFromGrpc(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, httputils.OK())
}
