package customercontroller

import (
	"log/slog"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	SessionsTokenIssuer     string
	ApplicationsTokenIssuer string
	SessionsTokenSecret     string
	ApplicationsTokenSecret string
}

type controller struct {
	router   chi.Router
	users    userspb.UserServiceClient
	sessions sessionspb.SessionsServiceClient
	log      *slog.Logger
	cfg      *Config
}

//issuer: fmt.Sprintf("%s.customer_controller", keys.APP_NAME)

func New(router chi.Router, log *slog.Logger, usersrpc userspb.UserServiceClient, sessionsrpc sessionspb.SessionsServiceClient, cfg *Config) *controller {
	return &controller{router: router, users: usersrpc, sessions: sessionsrpc, log: log, cfg: cfg}
}

func (c *controller) Register() {

	//TODO
	testMiddleware := middlewares.NewTestMiddleware(c.log)

	c.router.With(testMiddleware).Route("/customers", func(cr chi.Router) {

		cr.Route("/users", func(ur chi.Router) {

			ur.Post("/", c.createUser)
			//TODO: get ID from session
			ur.Patch("/{id}", c.updateUserByID)
			ur.Get("/{id}", c.findUserByID)
		})

		cr.Route("/session", func(sr chi.Router) {
			sr.Post("/login", c.sessionLogin)
			sr.Post("/refresh-session", c.refreshSession)
		})
	})
}
