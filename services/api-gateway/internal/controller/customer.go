package controller

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/domain"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/utils"
	"github.com/go-chi/chi/v5"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

	locale := r.Header.Get(headerLocale)

	if locale == "" {
		locale = headerFallbackLocale
	}

	var createUserReq userspb.CreateUserRequest

	err := utils.DecodeJSON(r.Body, &createUserReq)

	if err != nil {

		msg := i18n.Manager.GetMessage(locale, keys.ValidationInvalidBody)

		utils.WriteJSON(w, http.StatusBadRequest, domain.NewError(msg))
		return
	}

	md := metadata.Pairs(
		metadataLocaleKey, locale,
	)

	ctx := metadata.NewOutgoingContext(r.Context(), md)

	response, err := c.users.CreateUser(ctx, &createUserReq)

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			for _, d := range st.Details() {
				switch info := d.(type) {
				case *errdetails.BadRequest:
					utils.WriteJSON(w, http.StatusBadRequest, info)
					return
				}
			}
		}
		utils.WriteJSON(w, http.StatusBadRequest, domain.NewError(err.Error()))
		return
	}

	w.Write([]byte(fmt.Sprintf("CREATE USER: %s", response.GetId())))
}

func (c *customerController) Register() {
	c.router.Route("/customers", func(cr chi.Router) {

		cr.Route("/users", func(ur chi.Router) {

			ur.Post("/", c.createUser)
		})

	})
}
