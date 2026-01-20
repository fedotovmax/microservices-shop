package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fedotovmax/microservices-shop/apps-auth/internal/app"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/config"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/utils"
	"github.com/fedotovmax/microservices-shop/apps-auth/pkg/logger"
)

func main() {

	cfg, err := config.LoadAuthAppConfig()

	if err != nil {
		logger.GetFallback().Error(err.Error())
		os.Exit(1)
	}

	log, err := utils.SetupLooger(cfg.Env)

	if err != nil {
		logger.GetFallback().Error(err.Error())
		os.Exit(1)
	}

	application, err := app.New(cfg, log)

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	sig, triggerSignal := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer triggerSignal()

	application.Run(triggerSignal)

	<-sig.Done()

	log.Info("Signal recieved, shutdown app")

	//TODO: real time out for shutdown
	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer shutdownCtxCancel()

	application.Stop(shutdownCtx)

}
