package controller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/fedotovmax/microservices-shop/customer-site/internal/dom"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/keys"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/openapiclient"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/state"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/templates/pages/home"
	"github.com/fedotovmax/microservices-shop/customer-site/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/starfederation/datastar-go/datastar"
)

type controller struct {
	gateway   *openapiclient.APIClient
	base      *chi.Mux
	log       *slog.Logger
	refresher Refresher
}

type store struct {
	Message string `json:"new_message"`
}

func New(base *chi.Mux) *controller {

	gateway := openapiclient.NewAPIClient(&openapiclient.Configuration{
		Host: "",
	})

	refresher := func(ctx context.Context, dto openapiclient.SessionspbRefreshSessionRequest) (*openapiclient.SessionspbSessionCreated, *http.Response, error) {

		res, httpRes, err := gateway.CustomersAPI.CustomersSessionRefreshSessionPost(ctx).Dto(dto).Execute()

		return res, httpRes, err
	}

	return &controller{
		base:      base,
		gateway:   gateway,
		refresher: refresher,
	}
}

func (c *controller) RegisterProtected() {

	c.base.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		clientState := state.GetClientState(r)

		//TODO: logout request to API
		_ = clientState.RefreshToken
		ClearCookies(w)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	c.base.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
		clientState := state.GetClientState(r)

		sse := datastar.NewSSE(w, r)

		_ = sse

		ctx := CreateOpenapiAccessCtx(r.Context(), clientState.AccessToken)

		user, _, err := WithAuth(ctx, c.log, w, clientState, c.refresher, func(ctx context.Context) (*openapiclient.UserspbUser, *http.Response, error) {
			return c.gateway.CustomersAPI.CustomersUsersProfileGet(ctx).Execute()
		})

		if errors.Is(err, ErrUnauthorized) {
			http.Redirect(w, r, "/logout", http.StatusMovedPermanently)
			return
		}

		if err != nil {
			var genericOpenAPIError *openapiclient.GenericOpenAPIError
			if errors.As(err, &genericOpenAPIError) {
				switch t := genericOpenAPIError.Model().(type) {
				case openapiclient.ErrdetailsBadRequest:
					slog.Error("validation errors (from grpc)", slog.Any("violations", t.FieldViolations))
				case openapiclient.HttputilsErrorResponse:
					slog.Error("http error", slog.Any("error", t.GetMessage()))
				case openapiclient.GithubComFedotovmaxMicroservicesShopApiGatewayInternalDomainLoginErrorResponse:
					slog.Warn("login error response", slog.Any("current type", t))
				default:
					slog.Info("unknown http error", slog.Any("unknown error", err.Error()))
				}
			}

			//Handle default variant
		}

		_ = user

	})
}

func (c *controller) RegisterPublic() {

	c.base.Get("/", func(w http.ResponseWriter, r *http.Request) {

		csrf := utils.NewCSRF()
		csrfCookie := utils.CreateCSRFCookie(keys.CookieCsrf, csrf)
		http.SetCookie(w, csrfCookie)

		err := utils.Render(w, r, home.Home(&home.Props{CSRF: csrf}))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	c.base.Post(keys.ROUTE_TOAST, func(w http.ResponseWriter, r *http.Request) {

		c.log.Info(r.Header.Get("X-Csrf-Token"))

		sse := datastar.NewSSE(w, r)

		err := sse.PatchElementTempl(home.TestNotification(), datastar.WithSelectorID(dom.RandomToastContainerID()), datastar.WithModeAppend())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})

	c.base.Get("/updates", func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		patchedRequest := r.WithContext(ctx)

		sse := datastar.NewSSE(w, patchedRequest)

		ticker := time.NewTicker(time.Second * 3)
		defer ticker.Stop()

		for {
			select {
			case <-r.Context().Done():
				sse.PatchElements("<div>клиент вышел</div>", datastar.WithSelectorID("updates"), datastar.WithModeAppend())
				cancel()
				return
			case <-ticker.C:
				currentTime := time.Now().String()
				sse.PatchElements(fmt.Sprintf("<div>%s</div>", currentTime), datastar.WithSelectorID("updates"), datastar.WithModeAppend())
			}
		}

	})

	c.base.Post("/messages", func(w http.ResponseWriter, r *http.Request) {

		store := &store{}

		err := datastar.ReadSignals(r, store)
		if err != nil {
			c.log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sse := datastar.NewSSE(w, r)

		sse.PatchElements(fmt.Sprintf("<div>%s</div>", store.Message), datastar.WithSelectorID("messages"), datastar.WithModeAppend())

		store.Message = ""

		sse.MarshalAndPatchSignals(store)

	})

}
