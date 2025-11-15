package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/fedotovmax/microservices-shop/user_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/user_service/internal/ports"
	"github.com/fedotovmax/microservices-shop/user_service/pkg/utils"
	"github.com/fedotovmax/pgxtx"
)

type EventProcessorConfig struct {
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

type EventProcessorProps struct {
	ProduceAdapter     ports.ProduceAdapter
	EventAdapter       ports.EventAdapter
	TransactionManager pgxtx.Manager
	Config             EventProcessorConfig
}

type eventProcessor struct {
	pa        ports.ProduceAdapter
	ea        ports.EventAdapter
	txManager pgxtx.Manager
	cfg       EventProcessorConfig
	// flag for
	inProcess int32
}

func NewEventProcessorUsecase(props EventProcessorProps) *eventProcessor {

	const minWorkers = 1
	const maxWorkers = 10

	const minLimit = 1
	const maxLimit = 100

	const minInterval = time.Second * 5

	const minReserve = time.Second * 25

	const minProcessTimeout = time.Second * 800

	if props.Config.Workers < minWorkers {
		props.Config.Workers = minWorkers
	}

	if props.Config.Workers > maxWorkers {
		props.Config.Workers = maxWorkers
	}

	if props.Config.Limit < minLimit {
		props.Config.Workers = minLimit
	}

	if props.Config.Limit > maxLimit {
		props.Config.Workers = maxLimit
	}

	if props.Config.Interval < minInterval {
		props.Config.Interval = minInterval
	}

	if props.Config.ReserveDuration < minReserve {
		props.Config.ReserveDuration = minReserve
	}

	if props.Config.ProcessTimeout < minProcessTimeout {
		props.Config.ProcessTimeout = minProcessTimeout
	}

	return &eventProcessor{
		ea:        props.EventAdapter,
		pa:        props.ProduceAdapter,
		txManager: props.TransactionManager,
		cfg:       props.Config,
	}
}

func (e *eventProcessor) DispatchMonitoring(ctx context.Context) {

	const op = "usecase.event-processor.DispatchMonitoring"

	eventsSuccesses := e.pa.GetSuccesses(ctx)
	eventsErrors := e.pa.GetErrors(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Debug("dispatch monitoring [successes] stopped", slog.String("op", op))
				return
			case success := <-eventsSuccesses:
				err := e.confirmSuccess(ctx, success)
				if err != nil {
					slog.Error("error when confirm event, but event is sended", slog.String("op", op),
						slog.Any("error", err))
					continue
				}
				slog.Debug("event sended", slog.String("op", op), slog.String("event_id", success.ID))
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Debug("dispatch monitoring [errors] stopped", slog.String("op", op))

				return
			case err := <-eventsErrors:

				slog.Error("event send failed", slog.String("op", op),
					slog.String("event_id", err.ID), slog.Any("error", err.Error))

				confirmErr := e.confirmFailed(ctx, err)

				if confirmErr != nil {
					slog.Error("error when confirm send fail", slog.String("op", op), slog.Any("error", confirmErr))
				}
			}
		}
	}()

}

func (e *eventProcessor) ProcessingNewEvents(ctx context.Context) {
	const op = "usecase.event-processor.ProcessingNewEvents"

	go func() {
		ticker := time.NewTicker(e.cfg.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				slog.Debug("event processing stopped", slog.String("op", op))
				return
			case <-ticker.C:
				if !atomic.CompareAndSwapInt32(&e.inProcess, 0, 1) {
					continue
				}
				e.process(ctx)
				atomic.StoreInt32(&e.inProcess, 0)
			}
		}
	}()
}

func (e *eventProcessor) confirmFailed(ctx context.Context, ev *domain.FailedEvent) error {
	const op = "usecase.event-processor.confirmFailed"

	err := e.txManager.Wrap(ctx, func(txCtx context.Context) error {

		updateCtx, cancelUpdateCtx := context.WithTimeout(txCtx, time.Second*3)
		defer cancelUpdateCtx()

		err := e.ea.RemoveReserve(updateCtx, ev.ID)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (e *eventProcessor) confirmSuccess(ctx context.Context, ev *domain.SuccessEvent) error {

	const op = "usecase.event-processor.confirmSuccess"

	err := e.txManager.Wrap(ctx, func(txCtx context.Context) error {

		updateCtx, cancelUpdateCtx := context.WithTimeout(txCtx, time.Second*3)
		defer cancelUpdateCtx()

		err := e.ea.RemoveReserve(updateCtx, ev.ID)

		if err != nil {
			return err
		}

		err = e.ea.ChangeStatus(updateCtx, ev.ID)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (e *eventProcessor) process(ctx context.Context) {

	const op = "usecase.event-processor.process"

	var events []*domain.Event

	queriesCtx, cancelQueriesCtx := context.WithTimeout(ctx, e.cfg.ProcessTimeout)
	defer cancelQueriesCtx()

	err := e.txManager.Wrap(queriesCtx, func(txCtx context.Context) error {
		var err error
		events, err = e.ea.FindNewAndNotReserved(txCtx, e.cfg.Limit)

		if err != nil {
			return err
		}

		eventsIds := make([]string, len(events))

		for i := 0; i < len(events); i++ {
			eventsIds[i] = events[i].ID
		}

		err = e.ea.SetReservedToByIDs(txCtx, eventsIds, e.cfg.ReserveDuration)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		slog.Error("error in processing transaction", slog.String("op", op), slog.Any("error", err))
		return
	}

	eventsCh := make(chan *domain.Event, len(events))
	for i := 0; i < len(events); i++ {
		eventsCh <- events[i]
	}
	close(eventsCh)

	workerPoolCtx, workerPoolCtxCancel := context.WithCancel(ctx)
	defer workerPoolCtxCancel()

	publishResults := utils.Workerpool(workerPoolCtx, eventsCh, e.cfg.Workers,
		func(ev *domain.Event) error {
			publishCtx, cancelPublishCtx := context.WithTimeout(workerPoolCtx, e.cfg.ProcessTimeout)
			defer cancelPublishCtx()
			return e.pa.Publish(publishCtx, ev)
		})

	for err := range publishResults {
		if err != nil {
			slog.Debug("publish error", slog.String("op", op), slog.Any("error", err))
		}
	}

}
