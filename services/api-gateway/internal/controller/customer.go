package controller

import (
	"log/slog"
	"net/http"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/domain"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
	"github.com/fedotovmax/microservices-shop/api-gateway/pkg/utils/httphelper"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/metadata"
)

type customerController struct {
	router chi.Router
	users  userspb.UserServiceClient
	log    *slog.Logger
}

func NewCustomerController(router chi.Router, log *slog.Logger, rpc userspb.UserServiceClient) *customerController {
	return &customerController{router: router, users: rpc, log: log}
}

func (c *customerController) createUser(w http.ResponseWriter, r *http.Request) {

	locale := r.Header.Get(keys.HeaderLocale)

	if locale == "" {
		locale = keys.FallbackLocale
	}

	var createUserReq userspb.CreateUserRequest

	err := httphelper.DecodeJSON(r.Body, &createUserReq)

	if err != nil {

		msg := i18n.Manager.GetMessage(locale, keys.ValidationInvalidBody)

		httphelper.WriteJSON(w, http.StatusBadRequest, domain.NewError(msg))
		return
	}

	md := metadata.Pairs(
		keys.MetadataLocaleKey, locale,
	)

	ctx := metadata.NewOutgoingContext(r.Context(), md)

	response, err := c.users.CreateUser(ctx, &createUserReq)

	if err != nil {
		httphelper.HandleErrorFromGrpc(err, w)
		return
	}

	httphelper.WriteJSON(w, http.StatusCreated, response)

}

func (c *customerController) getUserById(w http.ResponseWriter, r *http.Request) {

}

func (c *customerController) Register() {

	c.router.Route("/customers", func(cr chi.Router) {

		cr.Route("/users", func(ur chi.Router) {

			ur.Post("/", c.createUser)
			ur.Get("/{id}", c.getUserById)
		})

	})
}
