package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/ports"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/queries"
)

type SendTgMessageUsecase struct {
	log            *slog.Logger
	query          queries.Chat
	telegramSender ports.TelegramSender
}

func NewSendTgMessageUsecase(
	log *slog.Logger,
	query queries.Chat,
	telegramSender ports.TelegramSender,
) *SendTgMessageUsecase {
	return &SendTgMessageUsecase{
		log:            log,
		query:          query,
		telegramSender: telegramSender,
	}
}

func (u *SendTgMessageUsecase) Execute(ctx context.Context, text string, userId string) error {

	const op = "usecases.send_tg_message"

	chatID, err := u.query.FindByUserID(ctx, userId)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	n := &domain.TgNotification{
		ChatID: chatID,
		Text:   text,
	}

	err = u.telegramSender.SendMessage(ctx, n)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, errs.ErrSendTelegramMessage, err)
	}

	return nil
}
