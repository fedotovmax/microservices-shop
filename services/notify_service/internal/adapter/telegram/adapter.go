package telegram

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type tgAdapter struct {
	tgbot *bot.Bot
	ctx   context.Context
	stop  context.CancelFunc
}

type Config struct {
	Token   string
	Options []bot.Option
}

var ErrUnexpected = errors.New("unexpected error when set commands")

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

func New(cfg *Config) (*tgAdapter, error) {

	const op = "adapter.telegram.New"

	tgbot, err := bot.New(cfg.Token, cfg.Options...)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	setCmdCtx, cancelSetCmdCtx := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelSetCmdCtx()

	err = setCommands(setCmdCtx, tgbot)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &tgAdapter{
		tgbot: tgbot,
		ctx:   ctx,
		stop:  cancel,
	}, nil

}

func (tg *tgAdapter) Start() {
	go tg.tgbot.Start(tg.ctx)
}

func (tg *tgAdapter) Stop() {
	tg.stop()

}

func (tg *tgAdapter) RegisterCommand(cmdType bot.HandlerType, cmd domain.Cmd, f bot.HandlerFunc) error {

	const op = "adapter.telegram.RegisterCommand"

	err := cmd.Validate()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	tg.tgbot.RegisterHandler(cmdType, cmd.String(), bot.MatchTypePrefix, f)
	return nil
}

func (tg *tgAdapter) SendMessage(ctx context.Context, n *domain.TgNotification) error {

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
