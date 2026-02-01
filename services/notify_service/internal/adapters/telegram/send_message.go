package telegram

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain"
	"github.com/go-telegram/bot"
)

func (tg *telegram) SendMessage(ctx context.Context, n *domain.TgNotification) error {

	const op = "adapter.telegram.SendMessage"

	_, err := tg.tgbot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: n.ChatID,
		Text:   n.Text,
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
