package tgbotcontroller

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/keys"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Usecases interface {
	SaveChatUserPair(ctx context.Context, chatID int64, userID string) error
}

type Telegram interface {
	RegisterCommand(cmdType bot.HandlerType, cmd domain.Cmd, f bot.HandlerFunc) error
}

type tgBotController struct {
	log     *slog.Logger
	usecase Usecases
	tg      Telegram
}

func NewTgBotController(log *slog.Logger, usecase Usecases, tg Telegram) *tgBotController {
	return &tgBotController{
		log:     log,
		usecase: usecase,
		tg:      tg,
	}
}

func (tgc *tgBotController) Register() error {

	const op = "controller.tg_bot.Register"

	err := tgc.tg.RegisterCommand(bot.HandlerTypeMessageText, domain.Start, tgc.handleStart)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = tgc.tg.RegisterCommand(bot.HandlerTypeMessageText, domain.Help, tgc.handleHelp)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

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

	err = tgc.usecase.SaveChatUserPair(ctx, msg.Chat.ID, args[0])

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

// https://t.me/MicroservicesShopNotifyBot?start=12345

// returning cmd:string, args:[]string, error
func (tgc *tgBotController) parseCommand(text string) (string, []string, error) {

	parts := strings.Fields(text)

	if len(parts) == 0 {
		return "", nil, errs.ErrInvalidCommand
	}

	cmd := tgc.normalizeCmd(parts[0])

	var args []string

	if len(parts) > 1 {
		args = parts[1:]
	}

	return cmd, args, nil

}

func (tgc *tgBotController) normalizeCmd(cmd string) string {

	if i := strings.Index(cmd, "@"); i != -1 {
		return cmd[:i]
	}
	return cmd
}
