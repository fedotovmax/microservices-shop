package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/user_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/user_service/internal/ports"
	"github.com/fedotovmax/outbox"
	"github.com/fedotovmax/pgxtx"
)

type UserUsecase struct {
	ua     ports.UserAdapter
	outbox *outbox.Outbox
	txm    pgxtx.Manager
}

func NewUserUsecase(ua ports.UserAdapter, txm pgxtx.Manager, o *outbox.Outbox) *UserUsecase {
	return &UserUsecase{
		ua:     ua,
		outbox: o,
		txm:    txm,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, d domain.CreateUser) (string, error) {
	const op = "usecase.user.CreateUser"

	var userId string

	err := u.txm.Wrap(ctx, func(txCtx context.Context) error {
		var err error
		userId, err = u.ua.Create(txCtx, d)

		if err != nil {
			return err
		}

		payload := domain.UserCreatedEvent{ID: userId}

		b, err := json.Marshal(payload)

		if err != nil {
			return err
		}

		u.outbox.AddNewEvent(txCtx, outbox.CreateEvent{
			AggregateID: userId,
			Topic:       events.USER_EVENTS,
			Type:        "user.created",
			Payload:     b,
		})

		slog.Info("TODO: add new event", slog.Any("b", b))
		//TODO: add new event

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil

}
