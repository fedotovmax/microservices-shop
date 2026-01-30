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

func New(router chi.Router, log *slog.Logger, usersrpc userspb.UserServiceClient, sessionsrpc sessionspb.SessionsServiceClient, cfg *Config) *controller {
	return &controller{router: router, users: usersrpc, sessions: sessionsrpc, log: log, cfg: cfg}
}

func (c *controller) Register() {

	authMiddleware := middlewares.NewAuthMiddleware(
		c.log,
		c.cfg.SessionsTokenSecret,
		c.cfg.SessionsTokenIssuer,
	)

	c.router.Route("/customers", func(customersRouter chi.Router) {

		customersRouter.Route("/users", func(userRouter chi.Router) {

			userRouter.Post("/", c.createUser)
			userRouter.Get("/verify-email/{link}", c.verifyEmail)
			userRouter.Patch("/send-new-verify-email-link", c.sendNewEmailVerifyLink)

			userRouter.With(authMiddleware).Route("/profile", func(profileRouter chi.Router) {
				profileRouter.Get("/", c.getMe)
				profileRouter.Patch("/", c.updateUserProfile)
			})

		})

		customersRouter.Route("/session", func(sessionRouter chi.Router) {
			sessionRouter.Post("/login", c.sessionLogin)
			sessionRouter.Post("/refresh-session", c.refreshSession)
			sessionRouter.With(authMiddleware).Get("/check", c.checkSession)
		})
	})
}
