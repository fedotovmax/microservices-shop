package controller

import (
	"context"
	"log/slog"
	"strings"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain/errs"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TgBotUsecase interface {
	TgBotStartCommand(ctx context.Context, chatID int64, userID string) error
}

type tgBotController struct {
	log     *slog.Logger
	usecase TgBotUsecase
}

func NewTgBotController(log *slog.Logger, usecase TgBotUsecase) *tgBotController {
	return &tgBotController{
		log:     log,
		usecase: usecase,
	}
}

// https://t.me/MicroservicesShopNotifyBot?start=userID-12345
func (tgc *tgBotController) Handler(ctx context.Context, b *bot.Bot, update *models.Update) {

	const op = "controller.tg_bot.Handler"

	l := tgc.log.With(slog.String("op", op))

	msg := update.Message

	if msg == nil || msg.Text == "" {
		return
	}

	cmd, args, err := tgc.parseCommand(msg.Text)

	if err != nil {
		l.Error(err.Error())
		return
	}

	switch cmd {
	case domain.Start:

		err := tgc.usecase.TgBotStartCommand(ctx, msg.Chat.ID, args[0])

		if err != nil {
			l.Error(err.Error())
			return
		}

		l.Info("send message tg", slog.Int64("CHAT_ID", msg.Chat.ID),
			slog.String("CMD", cmd), slog.Any("ARGS", args))

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Вы успешно подписаны на уведомления!",
		})
	case domain.Help:
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Запрос помощи",
		})
	default:
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Неизвестная команда!",
		})
	}

}

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

func (tgc *tgBotController) isCommand(cmd string) bool {
	return strings.HasPrefix(cmd, "/")
}

func (tgc *tgBotController) normalizeCmd(cmd string) string {

	if i := strings.Index(cmd, "@"); i != -1 {
		return cmd[:i]
	}
	return cmd
}
