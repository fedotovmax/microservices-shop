package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fedotovmax/envconfig"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/client"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/config"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/dom"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/keys"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/middlewares"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/router"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/templates/pages/home"
	"github.com/fedotovmax/microservices-shop/customer-site/pkg/logger"
	"github.com/fedotovmax/microservices-shop/customer-site/pkg/utils"
	"github.com/go-chi/chi/v5"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/starfederation/datastar-go/datastar"
)

func Static(dir string) http.Handler {
	return http.StripPrefix("/public/", http.FileServer(http.Dir(dir)))
}

type Store struct {
	Message string `json:"new_message"`
}

func setupLooger(env string) (*slog.Logger, error) {
	switch env {
	case keys.Development:
		return logger.NewDevelopmentHandler(slog.LevelDebug), nil
	case keys.Production:
		return logger.NewProductionHandler(slog.LevelWarn), nil
	default:
		return nil, envconfig.ErrInvalidAppEnv
	}
}

func main() {

	cfg, err := config.LoadAppConfig()

	if err != nil {
		logger.GetFallback().Error(err.Error())
		os.Exit(1)
	}

	log, err := setupLooger(cfg.Env)

	if err != nil {
		logger.GetFallback().Error(err.Error())
		os.Exit(1)
	}

	workdir, err := os.Getwd()

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	publicDir := filepath.Join(workdir, "web", "public")

	transport := httptransport.New(cfg.ApiGatewayAddr, "", nil)

	customersApi := client.New(transport, strfmt.Default).Customers
	//TODO: use for query to api gateway

	// params := customers.NewGetCustomersUsersIDParamsWithContext(context.Background())
	// locale := "ru"
	// params.SetXRequestLocale(&locale)
	// params.SetID("123")

	// customersApi.GetCustomersUsersID(params)
	_ = customersApi

	r := chi.NewRouter()

	r.Use(middlewares.GzipMiddleware)

	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir(publicDir))))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		csrf := utils.NewCSRF()
		csrfCookie := utils.CreateCSRFCookie(keys.CookieCsrf, csrf)
		http.SetCookie(w, csrfCookie)

		err := utils.Render(w, r, home.Home(&home.Props{CSRF: csrf}))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	r.Post(router.TOAST_ROUTE, func(w http.ResponseWriter, r *http.Request) {

		log.Info(r.Header.Get("X-Csrf-Token"))

		sse := datastar.NewSSE(w, r)

		err := sse.PatchElementTempl(home.TestNotification(), datastar.WithSelectorID(dom.RandomToastContainerID()), datastar.WithModeAppend())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})

	r.Get("/updates", func(w http.ResponseWriter, r *http.Request) {

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

	r.Post("/messages", func(w http.ResponseWriter, r *http.Request) {

		store := &Store{}

		err := datastar.ReadSignals(r, store)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sse := datastar.NewSSE(w, r)

		sse.PatchElements(fmt.Sprintf("<div>%s</div>", store.Message), datastar.WithSelectorID("messages"), datastar.WithModeAppend())

		store.Message = ""

		sse.MarshalAndPatchSignals(store)

	})

	port := fmt.Sprintf(":%d", cfg.Port)

	s := &http.Server{
		Addr:    port,
		Handler: r,
	}

	log.Info(fmt.Sprintf("🚀 Server starting at %s%s", "http://localhost", port))

	err = s.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("error when start http server", slog.String("error", err.Error()))
		os.Exit(1)
	}

}
