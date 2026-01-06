package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fedotovmax/microservices-shop/customer-site/internal/client"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/client/users"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/components"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/middlewares"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/models"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/templates"
	"github.com/fedotovmax/microservices-shop/customer-site/pkg/htmx"
	"github.com/fedotovmax/microservices-shop/customer-site/pkg/utils"
	"github.com/go-chi/chi/v5"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

func Static(dir string) http.Handler {
	return http.StripPrefix("/public/", http.FileServer(http.Dir(dir)))
}

func main() {

	log := slog.Default()

	projectRoot, err := os.Getwd()

	publicDir := filepath.Join(projectRoot, "web", "public")

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	transport := httptransport.New("localhost:8081", "", nil)

	apiclient := client.New(transport, strfmt.Default)

	r := chi.NewRouter()

	r.Use(middlewares.GzipMiddleware)

	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir(publicDir))))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		err := utils.Render(w, r, templates.Homepage("Go HTMX+Templ"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	r.Get("/htmx", func(w http.ResponseWriter, r *http.Request) {
		log.Info("request to /htmx")
		if htmx.IsHTMX(r) {

			dto := &models.UserspbCreateUserRequest{
				Email:    "makc",
				Password: "123456",
			}
			params := users.NewPostCustomersUsersParams()
			params.SetDto(dto)
			params.SetContext(r.Context())

			response, err := apiclient.Users.PostCustomersUsers(params)

			if err != nil {

				log.Error(err.Error())
				return
			}

			log.Info("test request with apiclient is ok!", slog.Any("response", response))

			htmx.NewResponse().Retarget("#result").Reswap(htmx.SwapInnerHTML).MustRenderTempl(r.Context(), w, components.Test("htmx endpoint working!"))
			return
		}
	})

	s := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	err = s.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("error when start http server", slog.String("error", err.Error()))
		os.Exit(1)
	}

}
