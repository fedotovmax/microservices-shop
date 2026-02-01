package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain/errs"
)

func (u *usecases) FindChatByUser(ctx context.Context, userID string) (int64, error) {
	const op = "usecases.FindChatByUser"

	chatID, err := u.chatStorage.GetChatIDByUserID(ctx, userID)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return 0, fmt.Errorf("%s: %w: %v", op, errs.ErrChatIDNotFound, err)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return chatID, nil

}
