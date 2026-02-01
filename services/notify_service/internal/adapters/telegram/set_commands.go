package telegram

import (
	"context"
	"fmt"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func setCommands(ctx context.Context, b *bot.Bot) error {

	const op = "adapter.telegram.setCommands"

	locales, err := i18n.Local.GetSupportedLocales()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for locale := range locales {
		startDescription, _ := i18n.Local.Get(locale, domain.Start.String())
		helpDescription, _ := i18n.Local.Get(locale, domain.Help.String())

		ok, err := b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
			LanguageCode: locale,
			Commands: []models.BotCommand{
				{Command: domain.Start.String(), Description: startDescription},
				{Command: domain.Help.String(), Description: helpDescription},
			},
		})

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if !ok {
			return fmt.Errorf("%s: %w", op, ErrUnexpected)
		}
	}

	return nil

}
