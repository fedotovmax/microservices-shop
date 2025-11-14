package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/user_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/user_service/internal/ports"
	"github.com/fedotovmax/pgxtx"
)

type userAdapter interface {
	CreateUser(ctx context.Context, d domain.CreateUser) (string, error)
}

type userUsecase struct {
	ua  userAdapter
	ea  ports.EventAdapter
	txm pgxtx.Manager
}

func NewUserUsecase(ua userAdapter, ea ports.EventAdapter, txm pgxtx.Manager) *userUsecase {
	return &userUsecase{
		ua:  ua,
		ea:  ea,
		txm: txm,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context, d domain.CreateUser) (string, error) {
	const op = "usecase.user.CreateUser"

	var userId string

	err := u.txm.Wrap(ctx, func(txCtx context.Context) error {
		var err error
		userId, err = u.ua.CreateUser(txCtx, d)

		if err != nil {
			return err
		}

		payload := domain.UserCreatedEvent{ID: userId}

		b, err := json.Marshal(payload)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		_, err = u.ea.CreateEvent(txCtx, domain.CreateEvent{
			Topic:   events.USER_EVENTS,
			Type:    "user.created",
			Payload: b,
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil

}
