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

func (c *customerController) findUserByID(w http.ResponseWriter, r *http.Request) {

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

	user := domain.UserFromProto(locale, response)

	httphelper.WriteJSON(w, http.StatusOK, user)

}

func (c *customerController) createUser(w http.ResponseWriter, r *http.Request) {

	const op = "controller.customer.createUser"

	l := c.log.With(slog.String("op", op))

	locale := r.Header.Get(keys.HeaderLocale)

	if locale == "" {
		locale = keys.FallbackLocale
	}

	var createUserReq userspb.CreateUserRequest

	err := httphelper.DecodeJSON(r.Body, &createUserReq)

	if err != nil {

		msg, err := i18n.Local.Get(locale, keys.ValidationInvalidBody)

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

	response, err := c.users.CreateUser(ctx, &createUserReq)

	if err != nil {
		httphelper.HandleErrorFromGrpc(w, err)
		return
	}

	httphelper.WriteJSON(w, http.StatusCreated, response)

}

func (c *customerController) updateUserByID(w http.ResponseWriter, r *http.Request) {

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
		httphelper.WriteJSON(w, http.StatusBadRequest, domain.NewError(msg))
		return
	}

	var updateUserProfileReq userspb.UpdateUserProfileRequest

	err := httphelper.DecodeJSON(r.Body, &updateUserProfileReq)

	if err != nil {

		msg, err := i18n.Local.Get(locale, keys.ValidationInvalidBody)

		if err != nil {
			l.Error(err.Error())
		}

		httphelper.WriteJSON(w, http.StatusBadRequest, domain.NewError(msg))
		return
	}

	md := metadata.Pairs(
		keys.MetadataLocaleKey, locale,
		keys.MetadataUserIDKey, userId,
	)

	ctx := metadata.NewOutgoingContext(r.Context(), md)

	_, err = c.users.UpdateUserProfile(ctx, &updateUserProfileReq)

	if err != nil {
		httphelper.HandleErrorFromGrpc(w, err)
		return
	}

	httphelper.WriteJSON(w, http.StatusOK, domain.OK())
}

func (c *customerController) Register() {

	c.router.Route("/customers", func(cr chi.Router) {

		cr.Route("/users", func(ur chi.Router) {

			ur.Post("/", c.createUser)
			//TODO: get ID from session
			ur.Patch("/{id}", c.updateUserByID)
			ur.Get("/{id}", c.findUserByID)
		})

	})
}
