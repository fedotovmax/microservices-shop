package customercontroller

import (
	"log/slog"
	"net/http"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/domain"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
	"github.com/fedotovmax/microservices-shop/api-gateway/pkg/utils/httphelper"
	"google.golang.org/grpc/metadata"
)

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
		httphelper.WriteJSON(w, http.StatusBadRequest, domain.NewError(msg))
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
		httphelper.HandleErrorFromGrpc(w, err)
		return
	}

	httphelper.WriteJSON(w, http.StatusOK, response)

}
