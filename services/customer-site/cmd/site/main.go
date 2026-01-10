package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fedotovmax/microservices-shop/customer-site/internal/client"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/middlewares"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/templates"
	"github.com/fedotovmax/microservices-shop/customer-site/pkg/htmx"
	"github.com/fedotovmax/microservices-shop/customer-site/pkg/utils"
	"github.com/go-chi/chi/v5"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

type Sender struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type Notify struct {
	Variant *string `json:"variant"`
	Title   *string `json:"title"`
	Message *string `json:"message"`
	Sender  *Sender `json:"sender"`
}

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
		err := utils.Render(w, r, templates.Homepage("Go HTMX+Templ"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	r.Get("/notify", func(w http.ResponseWriter, r *http.Request) {
		if htmx.IsHTMX(r) {
			time.Sleep(time.Second * 2)
			variant := "message"
			title := "Hello, all is working!"
			message := "Event sending from server!"
			sender := Sender{Name: "Fedotov Max", Avatar: "/public/images/gopher.jpg"}

			notify := Notify{Variant: &variant, Title: &title, Message: &message, Sender: &sender}

			if err := htmx.NewResponse().Reswap(htmx.SwapNone).AddTrigger(htmx.TriggerObject("notify", notify)).Write(w); err != nil {
				log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})

	port := ":3000"

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
