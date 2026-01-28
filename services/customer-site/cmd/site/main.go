package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fedotovmax/envconfig"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/config"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/controller"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/keys"
	"github.com/fedotovmax/microservices-shop/customer-site/internal/middlewares"
	"github.com/fedotovmax/microservices-shop/customer-site/pkg/logger"
	"github.com/go-chi/chi/v5"
)

func Static(dir string) http.Handler {
	return http.StripPrefix("/public/", http.FileServer(http.Dir(dir)))
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

	r := chi.NewRouter()

	r.Use(middlewares.GzipMiddleware)

	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir(publicDir))))

	routeController := controller.New(r)

	routeController.RegisterPublic()
	routeController.RegisterProtected()

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
