package usecase

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain"
)

func (u *usecases) SendTgMessage(ctx context.Context, text string, userId string) error {

	//TODO::
	const op = "usecase.SendTgMessage"

	chatID, err := u.FindChatByUser(ctx, userId)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	n := &domain.TgNotification{
		ChatID: chatID,
		Text:   text,
	}

	err = u.tgSender.SendMessage(ctx, n)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
