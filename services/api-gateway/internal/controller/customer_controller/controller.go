package customercontroller

import (
	"log/slog"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/go-chi/chi/v5"
)

type controller struct {
	router chi.Router
	users  userspb.UserServiceClient
	log    *slog.Logger
}

func New(router chi.Router, log *slog.Logger, rpc userspb.UserServiceClient) *controller {
	return &controller{router: router, users: rpc, log: log}
}

func (c *controller) Register() {

	c.router.Route("/customers", func(cr chi.Router) {

		cr.Route("/users", func(ur chi.Router) {

			ur.Post("/", c.createUser)
			//TODO: get ID from session
			ur.Patch("/{id}", c.updateUserByID)
			ur.Get("/{id}", c.findUserByID)
		})

		cr.Route("/session", func(sr chi.Router) {
			sr.Post("/login", c.sessionLogin)
		})

	})
}
