package telegram

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/keys"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (tgc *tgBotController) handleStart(ctx context.Context, b *bot.Bot, u *models.Update) {

	const op = "controller.tg_bot.handleStart"

	l := tgc.log.With(slog.String("op", op))

	msg := u.Message

	if msg == nil || msg.Text == "" {
		return
	}

	_, args, err := tgc.parseCommand(msg.Text)

	if err != nil {
		l.Error(err.Error())
		return
	}

	locale := u.Message.From.LanguageCode

	if len(args) == 0 {
		responseText, err := i18n.Local.Get(locale, keys.UnableIdentifyUser)
		if err != nil {
			l.Warn(err.Error())
		}
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: msg.Chat.ID,
			Text:   responseText,
		})
		if err != nil {
			l.Error(err.Error())
			return
		}
	}

	err = tgc.saveChatUserPair.Execute(ctx, msg.Chat.ID, args[0])

	if err != nil {
		l.Error(err.Error())
		return
	}

	responseText, err := i18n.Local.Get(locale, keys.StartResponseText)

	if err != nil {
		l.Warn(err.Error())
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: msg.Chat.ID,
		Text:   responseText,
	})

	if err != nil {
		l.Error(err.Error())
		return
	}
}
