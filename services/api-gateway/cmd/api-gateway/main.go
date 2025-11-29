package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/fedotovmax/microservices-shop/api-gateway/internal/app"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/infra/logger"
)

func main() {

	log := logger.NewDevelopmentHandler()

	var port uint16 = 8081

	userClientAddr := "localhost:5555"

	application, err := app.New(log, app.Config{
		HttpPort:      port,
		UsersGRPCAddr: userClientAddr,
	})

	if err != nil {
		panic(err)
	}

	sig, triggerSignal := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer triggerSignal()

	application.MustRun(triggerSignal)

	<-sig.Done()

	log.Info("Signal recieved, shutdown app")

	shutdownContext, shutdownCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer shutdownCancel()

	application.Stop(shutdownContext)

}
