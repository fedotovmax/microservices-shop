package telegram

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/keys"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (tgc *tgBotController) handleHelp(ctx context.Context, b *bot.Bot, u *models.Update) {

	const op = "controller.tg_bot.handleHelp"

	l := tgc.log.With(slog.String("op", op))

	locale := u.Message.From.LanguageCode

	responseText, err := i18n.Local.Get(locale, keys.StartResponseText)

	if err != nil {
		l.Warn(err.Error())
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: u.Message.Chat.ID,
		Text:   responseText,
	})

	if err != nil {
		l.Error(err.Error())
		return
	}
}
