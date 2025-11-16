package eventprocessor

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/fedotovmax/microservices-shop/user_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/user_service/internal/infra/logger"
	"github.com/fedotovmax/microservices-shop/user_service/internal/ports"
	"github.com/fedotovmax/microservices-shop/user_service/pkg/utils"
)

type EventService interface {
	ConfirmEvent(ctx context.Context, ev *domain.SuccessEvent) error
	ConfirmFailed(ctx context.Context, ev *domain.FailedEvent) error
	ReserveNewEvents(ctx context.Context, limit int, reserveDuration time.Duration) ([]*domain.Event, error)
}

type App struct {
	pa  ports.ProduceAdapter
	es  EventService
	log *slog.Logger
	cfg Config
	// flag for
	inProcess int32
}

type Config struct {
	// Limit of events to receive, min = 1, Max = 100
	Limit int
	// Min = 1, Max = 10
	Workers  int
	Interval time.Duration
	// Event reserve duration
	ReserveDuration time.Duration
	// Timeout for methods process, confirmSuccess, confirmFailed
	ProcessTimeout time.Duration
}

func New(l *slog.Logger, pa ports.ProduceAdapter, es EventService, cfg Config) *App {

	const minWorkers = 1
	const maxWorkers = 10

	const minLimit = 1
	const maxLimit = 100

	const minInterval = time.Second * 5

	const minReserve = time.Second * 25

	const minProcessTimeout = time.Millisecond * 1100

	if cfg.Workers < minWorkers {
		cfg.Workers = minWorkers
	}

	if cfg.Workers > maxWorkers {
		cfg.Workers = maxWorkers
	}

	if cfg.Limit < minLimit {
		cfg.Workers = minLimit
	}

	if cfg.Limit > maxLimit {
		cfg.Workers = maxLimit
	}

	if cfg.Interval < minInterval {
		cfg.Interval = minInterval
	}

	if cfg.ReserveDuration < minReserve {
		cfg.ReserveDuration = minReserve
	}

	if cfg.ProcessTimeout < minProcessTimeout {
		cfg.ProcessTimeout = minProcessTimeout
	}

	return &App{
		pa:  pa,
		log: l,
		es:  es,
		cfg: cfg,
	}
}

func (a *App) DispatchMonitoring(ctx context.Context) {

	const op = "usecase.event-processor.DispatchMonitoring"

	log := a.log.With(slog.String("op", op))

	eventsSuccesses := a.pa.GetSuccesses(ctx)
	eventsErrors := a.pa.GetErrors(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("dispatch monitoring [successes] stopped")
				return
			case success := <-eventsSuccesses:
				err := a.confirm(ctx, success)
				if err != nil {
					log.Error("error when confirm event, but event is sended",
						logger.Err(err))
					continue
				}
				log.Info("event sended", slog.String("event_id", success.ID))
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("dispatch monitoring [errors] stopped")
				return
			case err := <-eventsErrors:
				log.Error("event send failed",
					slog.String("event_id", err.ID), logger.Err(err.Error))

				confirmErr := a.es.ConfirmFailed(ctx, err)

				if confirmErr != nil {
					log.Error("error when confirm send fail", logger.Err(confirmErr))
				}
			}
		}
	}()

}

func (a *App) confirm(ctx context.Context, ev *domain.SuccessEvent) error {

	const op = "app.event-processor.confirm"

	queriesCtx, cancelQueriesCtx := context.WithTimeout(ctx, a.cfg.ProcessTimeout)
	defer cancelQueriesCtx()

	err := a.es.ConfirmEvent(queriesCtx, ev)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

func (a *App) ProcessingNewEvents(ctx context.Context) {
	const op = "usecase.event-processor.ProcessingNewEvents"

	log := slog.With(slog.String("op", op))

	go func() {
		ticker := time.NewTicker(a.cfg.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Info("event processing stopped")
				return
			case <-ticker.C:
				if !atomic.CompareAndSwapInt32(&a.inProcess, 0, 1) {
					continue
				}
				a.process(ctx)
				atomic.StoreInt32(&a.inProcess, 0)
			}
		}
	}()
}

func (a *App) process(ctx context.Context) {

	const op = "usecase.event-processor.process"

	log := slog.With(slog.String("op", op))

	queriesCtx, cancelQueriesCtx := context.WithTimeout(ctx, a.cfg.ProcessTimeout)
	defer cancelQueriesCtx()

	events, err := a.es.ReserveNewEvents(queriesCtx, a.cfg.Limit, a.cfg.ReserveDuration)

	if err != nil {
		log.Error("error in processing transaction", logger.Err(err))
		return
	}

	eventsCh := make(chan *domain.Event, len(events))
	for i := 0; i < len(events); i++ {
		eventsCh <- events[i]
	}
	close(eventsCh)

	workerPoolCtx, workerPoolCtxCancel := context.WithCancel(ctx)
	defer workerPoolCtxCancel()

	publishResults := utils.Workerpool(workerPoolCtx, eventsCh, a.cfg.Workers,
		func(ev *domain.Event) error {
			publishCtx, cancelPublishCtx := context.WithTimeout(workerPoolCtx, a.cfg.ProcessTimeout)
			defer cancelPublishCtx()
			return a.pa.Publish(publishCtx, ev)
		})

	for err := range publishResults {
		if err != nil {
			log.Info("publish error", logger.Err(err))
		}
	}
}
