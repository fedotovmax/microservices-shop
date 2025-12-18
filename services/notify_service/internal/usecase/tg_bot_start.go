package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain/errs"
)

func (u *usecases) TgBotStartCommand(ctx context.Context, chatID int64, userID string) error {
	const op = "usecase.TgBotStartCommand"

	err := u.storage.SaveChatIDByUserID(ctx, chatID, userID)
	if err != nil {
		if errors.Is(err, adapter.ErrAlreadyExists) {
			return fmt.Errorf("%s: %w: %v", op, errs.ErrUserIDAlreadyExists, err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	err = u.storage.SaveUserIDByChatID(ctx, chatID, userID)
	if err != nil {
		if errors.Is(err, adapter.ErrAlreadyExists) {
			return fmt.Errorf("%s: %w: %v", op, errs.ErrChatIDAlreadyExists, err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
