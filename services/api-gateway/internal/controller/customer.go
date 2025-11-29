package controller

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/domain"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/utils"
	"github.com/go-chi/chi/v5"
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

	var createUserReq userspb.CreateUserRequest

	err := utils.DecodeJSON(r.Body, &createUserReq)

	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.NewError(domain.INVALID_BODY))
		return
	}

	response, err := c.users.CreateUser(r.Context(), &createUserReq)

	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.NewError(err.Error()))
		return
	}

	c.log.Info("create user response", slog.Any("response", response))

	w.Write([]byte(fmt.Sprintf("CREATE USER: %s", response.GetId())))
}

func (c *customerController) Register() {
	c.router.Route("/customers", func(cr chi.Router) {

		cr.Route("/users", func(ur chi.Router) {

			ur.Post("/", c.createUser)
		})

	})
}
