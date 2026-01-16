package customercontroller

import (
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

type controller struct {
	router   chi.Router
	users    userspb.UserServiceClient
	sessions sessionspb.SessionsServiceClient
	log      *slog.Logger
	issuer   string
}

func New(router chi.Router, log *slog.Logger, usersrpc userspb.UserServiceClient, sessionsrpc sessionspb.SessionsServiceClient) *controller {
	return &controller{router: router, users: usersrpc, sessions: sessionsrpc, log: log, issuer: fmt.Sprintf("%s.customer_controller", keys.APP_NAME)}
}

func (c *controller) Register() {

	accessTokenMiddleware := middlewares.NewAccessTokenMiddleware(c.users, c.sessions, c.log)

	c.router.Route("/customers", func(cr chi.Router) {

		cr.Route("/users", func(ur chi.Router) {

			ur.Post("/", c.createUser)
			//TODO: get ID from session
			ur.With(accessTokenMiddleware).Patch("/{id}", c.updateUserByID)
			ur.Get("/{id}", c.findUserByID)
		})

		cr.Route("/session", func(sr chi.Router) {
			sr.Post("/login", c.sessionLogin)
			sr.Post("/refresh-session", c.refreshSession)
		})

	})
}
